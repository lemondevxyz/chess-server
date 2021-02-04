package rest

import (
	"github.com/gin-gonic/gin"
	"gopkg.in/go-playground/validator.v9"
)

var v = validator.New()

func Init(r *gin.RouterGroup) {
	{
		r.POST("/cmd", UserWall(cmdHandler))
	}
	{
		r.POST("/game", gameHandler)
	}
}
