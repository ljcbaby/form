package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/ljcbaby/form/service"
)

type FormController struct{}

func (c *FormController) CreateForm(ctx *gin.Context) {
	userId, _ := ctx.Get("userId")

	fs := service.FormService{}

	id, err := fs.CreateForm(userId.(int64))
	if err != nil {
		ctx.JSON(http.StatusOK, gin.H{
			"code": "1",
			"msg":  "Create form failed.",
			"data": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"code": "0",
		"msg":  "Success.",
		"data": gin.H{
			"id": id,
		},
	})
}

func (c *FormController) DuplicateForm(ctx *gin.Context) {}

func (c *FormController) GetFormList(ctx *gin.Context) {}

func (c *FormController) GetFormDetail(ctx *gin.Context) {}

func (c *FormController) UpdateForm(ctx *gin.Context) {}

func (c *FormController) DeleteForm(ctx *gin.Context) {}

func (c *FormController) SubmitForm(ctx *gin.Context) {}

func (c *FormController) GetFormResults(ctx *gin.Context) {}
