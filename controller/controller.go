package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type Controller struct {
	User *UserController
	Form *FormController
}

func (c *Controller) Index(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, gin.H{
		"code": "0",
		"msg":  "Ljcbaby's form backend.",
	})
}

func returnMySQLError(ctx *gin.Context, err error) {
	ctx.JSON(http.StatusOK, gin.H{
		"code": "100",
		"msg":  "MySQL error.",
		"data": err.Error(),
	})
}
