package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type Controller struct {
	User   *UserController
	Form   *FormController
	Result *ResultController
}

func (c *Controller) Index(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, gin.H{
		"code": "0",
		"msg":  "Ljcbaby's form backend.",
	})
}

func returnFormNotFound(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, gin.H{
		"code": "10",
		"msg":  "Form not found.",
	})
}

func returnFormIdInvalid(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, gin.H{
		"code": "11",
		"msg":  "Invalid form id.",
	})
}

func returnFormResultInvalid(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, gin.H{
		"code": "12",
		"msg":  "Invalid form result.",
	})
}

func returnMySQLError(ctx *gin.Context, err error) {
	ctx.JSON(http.StatusOK, gin.H{
		"code": "100",
		"msg":  "MySQL error.",
		"data": err.Error(),
	})
}
