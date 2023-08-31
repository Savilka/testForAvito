package api

import (
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"
	"net/http"
	"testForAvito/internal/api/v1/models"
	"testForAvito/internal/storage/postgres"
)

func AddUser(db *postgres.Storage) gin.HandlerFunc {
	return func(c *gin.Context) {
		id, err := db.AddNewUser()
		if err != nil {
			c.Error(err)
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{"user_id": id})
	}
}

func AddSegment(db *postgres.Storage) gin.HandlerFunc {
	return func(c *gin.Context) {
		var data api.AddSegmentRequest

		err := c.ShouldBindJSON(&data)
		if err != nil {
			c.Error(err)
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		id, err := db.AddSegment(data.Slug, data.Percent)
		if err != nil && err.Error() == "storage.postgres.AddSegment: "+pgx.ErrNoRows.Error() {
			c.Error(err)
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Segment already exist"})
			return
		}
		if err != nil {
			c.Error(err)
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{"segment_id": id})
	}
}

func DeleteSegment(db *postgres.Storage) gin.HandlerFunc {
	return func(c *gin.Context) {
		var data api.DeleteSegmentRequest

		err := c.ShouldBindJSON(&data)
		if err != nil {
			c.Error(err)
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		err = db.DeleteSegment(data.Slug)
		if err != nil {
			c.Error(err)
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{})
	}
}

func GetUserSegments(db *postgres.Storage) gin.HandlerFunc {
	return func(c *gin.Context) {
		names, err := db.GetUserSegments(c.Param("id"))
		if err != nil {
			c.Error(err)
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{"segments": names})
	}
}

func AddUserToSegments(db *postgres.Storage) gin.HandlerFunc {
	return func(c *gin.Context) {
		var data api.AddUserSegmentRequest
		err := c.ShouldBindJSON(&data)
		if err != nil {
			c.Error(err)
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		err = db.AddUserToSegments(data.NewSegments, data.OldSegments, c.Param("id"))
		if err != nil {
			c.Error(err)
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{})
	}
}
