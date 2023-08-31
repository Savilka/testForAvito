package test

import (
	"bytes"
	"context"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/ilyakaznacheev/cleanenv"
	"github.com/jackc/pgx/v5"
	"github.com/stretchr/testify/assert"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	api "testForAvito/internal/api/v1/handlers"
	"testForAvito/internal/app"
	"testForAvito/internal/config"
	"testForAvito/internal/storage/postgres"
	"testing"
)

var App app.App

func load() *config.Config {
	configPath := os.Getenv("CONFIG_PATH_FOR_TEST")
	if configPath == "" {
		log.Fatal("CONFIG_PATH is not set")
	}

	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		log.Fatalf("config file does not exist: %s", configPath)
	}

	var cfg config.Config

	if err := cleanenv.ReadConfig(configPath, &cfg); err != nil {
		log.Fatalf("cannot read config: %s", err)
	}

	return &cfg
}

func setUpDatabase() *pgx.Conn {
	cfg := load()
	con, _ := pgx.Connect(context.Background(), cfg.DataBaseUrl)
	con.Exec(context.Background(), "create table if not exists \"user\"\n(\n    id serial not null\n        constraint user_pk\n            primary key\n);\n\ncreate table if not exists segment\n(\n    id   serial not null\n        constraint segment_pk\n            primary key,\n    name varchar(200) unique\n);\n\ncreate table if not exists user_segment\n(\n    user_id    serial not null\n        constraint user_segment_user_id_fk\n            references \"user\"\n            on delete cascade,\n    segment_name varchar(200) not null\n        constraint user_segment_segment_name_fk\n            references segment (name)\n            on delete cascade,\n    constraint user_segment_pk\n        primary key (user_id, segment_name)\n);")
	return con
}

func cleanDataBase() {
	cfg := load()
	con, _ := pgx.Connect(context.Background(), cfg.DataBaseUrl)
	con.Exec(context.Background(), "DROP TABLE \"user\", segment, user_segment")
	con.Close(context.Background())
}

func TestMain(m *testing.M) {
	cfg := load()
	con, err := postgres.Connect(cfg.DataBaseUrl)
	if err != nil {
		panic(err)
	}
	App = app.App{
		Db:     con,
		Router: gin.Default(),
	}
	if err != nil {
		panic(err)
	}

	v1 := App.Router.Group("/v1")
	{
		userGroup := v1.Group("/user")
		{
			userGroup.GET("/add", api.AddUser(App.Db))
			userGroup.GET("/:id/getSegments", api.GetUserSegments(App.Db))
			userGroup.PUT("/:id/addSegments", api.AddUserToSegments(App.Db))
		}

		segmentGroup := v1.Group("/segment")
		{
			segmentGroup.POST("/add", api.AddSegment(App.Db))
			segmentGroup.DELETE("/delete", api.DeleteSegment(App.Db))
		}
	}

	m.Run()
}

func TestAddUser(t *testing.T) {
	setUpDatabase()

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/v1/user/add", nil)
	App.Router.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)
	assert.Equal(t, `{"user_id":"1"}`, w.Body.String())

	cleanDataBase()
}

func TestSegmentAdd(t *testing.T) {
	setUpDatabase()

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/v1/user/add", nil)
	App.Router.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)
	assert.Equal(t, `{"user_id":"1"}`, w.Body.String())

	w = httptest.NewRecorder()
	req, _ = http.NewRequest("GET", "/v1/user/add", nil)
	App.Router.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)
	assert.Equal(t, `{"user_id":"2"}`, w.Body.String())

	tests := []struct {
		body         string
		expected     string
		expectedCode int
	}{
		{
			`{"slug": "test1"}`,
			`{"segment_id":"1"}`,
			200,
		},
		{
			`{"slug": "test2","percent":-1}`,
			`{"error":"Key: 'AddSegmentRequest.Percent' Error:Field validation for 'Percent' failed on the 'gte' tag"}`,
			500,
		},
		{
			`{"slug": "test3","percent":101}`,
			`{"error":"Key: 'AddSegmentRequest.Percent' Error:Field validation for 'Percent' failed on the 'lte' tag"}`,
			500,
		},
		{
			`{"slug": "test4","percent":100}`,
			`{"segment_id":"2"}`,
			200,
		},
		{
			`{"slug": "test1"}`,
			`{"error":"Segment already exist"}`,
			500,
		},
	}

	for _, test := range tests {
		w = httptest.NewRecorder()
		jsonBody := []byte(test.body)
		bodyReader := bytes.NewReader(jsonBody)
		req, _ = http.NewRequest("POST", "/v1/segment/add", bodyReader)
		App.Router.ServeHTTP(w, req)

		assert.Equal(t, test.expectedCode, w.Code)
		assert.Equal(t, test.expected, w.Body.String())
	}

	cleanDataBase()
}

func TestDeleteSegment(t *testing.T) {
	setUpDatabase()

	w := httptest.NewRecorder()
	jsonBody := []byte(`{"slug": "test1"}`)
	bodyReader := bytes.NewReader(jsonBody)
	req, _ := http.NewRequest("POST", "/v1/segment/add", bodyReader)
	App.Router.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)
	assert.Equal(t, `{"segment_id":"1"}`, w.Body.String())

	w = httptest.NewRecorder()
	jsonBody = []byte(`{"slug": "test1"}`)
	bodyReader = bytes.NewReader(jsonBody)
	req, _ = http.NewRequest("DELETE", "/v1/segment/delete", bodyReader)
	App.Router.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)

	w = httptest.NewRecorder()
	jsonBody = []byte(`{}`)
	bodyReader = bytes.NewReader(jsonBody)
	req, _ = http.NewRequest("DELETE", "/v1/segment/delete", bodyReader)
	App.Router.ServeHTTP(w, req)

	assert.Equal(t, 500, w.Code)

	cleanDataBase()
}

func TestAddUserToSegments(t *testing.T) {
	setUpDatabase()

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/v1/user/add", nil)
	App.Router.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)
	assert.Equal(t, `{"user_id":"1"}`, w.Body.String())

	w = httptest.NewRecorder()
	jsonBody := []byte(`{"slug": "1"}`)
	bodyReader := bytes.NewReader(jsonBody)
	req, _ = http.NewRequest("POST", "/v1/segment/add", bodyReader)
	App.Router.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)
	assert.Equal(t, `{"segment_id":"1"}`, w.Body.String())

	w = httptest.NewRecorder()
	jsonBody = []byte(`{"slug": "2"}`)
	bodyReader = bytes.NewReader(jsonBody)
	req, _ = http.NewRequest("POST", "/v1/segment/add", bodyReader)
	App.Router.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)
	assert.Equal(t, `{"segment_id":"2"}`, w.Body.String())

	w = httptest.NewRecorder()
	jsonBody = []byte(`{"slug": "3"}`)
	bodyReader = bytes.NewReader(jsonBody)
	req, _ = http.NewRequest("POST", "/v1/segment/add", bodyReader)
	App.Router.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)
	assert.Equal(t, `{"segment_id":"3"}`, w.Body.String())

	w = httptest.NewRecorder()
	jsonBody = []byte(`{"new_segments": ["1", "2", "3"],"old_segments": []}`)
	bodyReader = bytes.NewReader(jsonBody)
	req, _ = http.NewRequest("PUT", "/v1/user/1/addSegments", bodyReader)
	App.Router.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)

	w = httptest.NewRecorder()
	req, _ = http.NewRequest("GET", "/v1/user/1/getSegments", nil)
	App.Router.ServeHTTP(w, req)
	res := struct {
		Segments []string `json:"segments"`
	}{
		Segments: nil,
	}

	err := json.Unmarshal(w.Body.Bytes(), &res)
	if err != nil {
		panic(err)
	}
	_ = string(w.Body.Bytes()[:])
	assert.Equal(t, 200, w.Code)
	assert.Equal(t, []string{"1", "2", "3"}, res.Segments)
	cleanDataBase()
}
