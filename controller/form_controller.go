package controller

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/ljcbaby/form/model"
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

func (c *FormController) GetFormDetail(ctx *gin.Context) {
	userID, _ := ctx.Get("userId")

	id, err := strconv.ParseInt(ctx.Param("id"), 10, 64)
	if err != nil {
		ctx.JSON(http.StatusOK, gin.H{
			"code": "1",
			"msg":  "Invalid form id.",
			"data": err.Error(),
		})
		return
	}

	var form model.Form
	form.ID = id
	form.OwnerID = userID.(int64)

	fs := service.FormService{}

	exist, err := fs.CheckFormNameExist(form)
	if err != nil {
		returnMySQLError(ctx, err)
		return
	}
	if !exist {
		ctx.JSON(http.StatusOK, gin.H{
			"code": "10",
			"msg":  "Form not found.",
			"data": err.Error(),
		})
		return
	}

	if err := fs.GetFormDetail(&form); err != nil {
		returnMySQLError(ctx, err)
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"code": "0",
		"msg":  "Success.",
		"data": form,
	})
}

func (c *FormController) UpdateForm(ctx *gin.Context) {}

func (c *FormController) DeleteForm(ctx *gin.Context) {}

func (c *FormController) SubmitForm(ctx *gin.Context) {}

func (c *FormController) GetFormResults(ctx *gin.Context) {}
