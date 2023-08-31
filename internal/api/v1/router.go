package api

import (
	"github.com/gin-gonic/gin"
	api "testForAvito/internal/api/v1/handlers"
	"testForAvito/internal/storage/postgres"
)

func InitRouter(db *postgres.Storage) *gin.Engine {
	r := gin.Default()
	initHandlers(r, db)
	return r
}

func initHandlers(r *gin.Engine, db *postgres.Storage) {
	v1 := r.Group("/v1")
	{
		userGroup := v1.Group("/user")
		{
			userGroup.GET("/add", api.AddUser(db))
			userGroup.GET("/:id/getSegments", api.GetUserSegments(db))
			userGroup.PUT("/:id/addSegments", api.AddUserToSegments(db))
		}

		segmentGroup := v1.Group("/segment")
		{
			segmentGroup.POST("/add", api.AddSegment(db))
			segmentGroup.DELETE("/delete", api.DeleteSegment(db))
		}
	}
}
