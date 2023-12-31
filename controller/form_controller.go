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

func (c *FormController) DuplicateForm(ctx *gin.Context) {
	userID, _ := ctx.Get("userId")

	id, err := strconv.ParseInt(ctx.Param("id"), 10, 64)
	if err != nil {
		returnFormIdInvalid(ctx)
		return
	}

	var form model.Form
	form.ID = id
	form.OwnerID = userID.(int64)

	fs := service.FormService{}

	exist, err := fs.CheckFormExist(form)
	if err != nil {
		returnMySQLError(ctx, err)
		return
	}
	if !exist {
		returnFormNotFound(ctx)
		return
	}

	newID, err := fs.DuplicateForm(form)
	if err != nil {
		returnMySQLError(ctx, err)
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"code": "0",
		"msg":  "Success.",
		"data": gin.H{
			"id": newID,
		},
	})
}

func (c *FormController) GetFormList(ctx *gin.Context) {
	userID, _ := ctx.Get("userId")

	page, err := strconv.Atoi(ctx.DefaultQuery("page", "1"))
	if err != nil || page < 1 {
		page = 1
	}
	size, err := strconv.Atoi(ctx.DefaultQuery("size", "10"))
	if err != nil || size < 1 {
		size = 20
	}
	if size > 100 {
		size = 100
	}

	fs := service.FormService{}

	total, err := fs.GetFormListCount(userID.(int64))
	if err != nil {
		returnMySQLError(ctx, err)
		return
	}

	var forms []model.FormBase
	err = fs.GetFormList(userID.(int64), page, size, &forms)
	if err != nil {
		returnMySQLError(ctx, err)
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"code": "0",
		"msg":  "Success.",
		"data": gin.H{
			"totalPage": total/size + func() int {
				if total%size > 0 {
					return 1
				}
				return 0
			}(),
			"totalCount": total,
			"results":    forms,
		},
	})
}

func (c *FormController) GetFormDetail(ctx *gin.Context) {
	userID, t := ctx.Get("userId")

	id, err := strconv.ParseInt(ctx.Param("id"), 10, 64)
	if err != nil {
		returnFormIdInvalid(ctx)
		return
	}

	var form model.Form
	form.ID = id
	if t {
		form.OwnerID = userID.(int64)
	}

	fs := service.FormService{}

	exist, err := fs.CheckFormExist(form)
	if err != nil {
		returnMySQLError(ctx, err)
		return
	}
	if !exist {
		returnFormNotFound(ctx)
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

func (c *FormController) UpdateForm(ctx *gin.Context) {
	userID, _ := ctx.Get("userId")

	id, err := strconv.ParseInt(ctx.Param("id"), 10, 64)
	if err != nil {
		returnFormIdInvalid(ctx)
		return
	}

	var form model.Form
	form.ID = id
	form.OwnerID = userID.(int64)

	fs := service.FormService{}

	exist, err := fs.CheckFormExist(form)
	if err != nil {
		returnMySQLError(ctx, err)
		return
	}
	if !exist {
		returnFormNotFound(ctx)
		return
	}

	if err := fs.GetFormDetail(&form); err != nil {
		returnMySQLError(ctx, err)
		return
	}

	if form.Status == 2 {
		ctx.JSON(http.StatusOK, gin.H{
			"code": "2",
			"msg":  "Form is published.",
		})
		return
	}

	form.IsPublish = -1

	if err := ctx.BindJSON(&form); err != nil {
		ctx.JSON(http.StatusOK, gin.H{
			"code": "1",
			"msg":  "Invalid form data.",
			"data": err.Error(),
		})
		return
	}

	if len(form.Title) > 255 {
		ctx.JSON(http.StatusOK, gin.H{
			"code": "1",
			"msg":  "Form title too long.",
		})
		return
	}
	if form.IsPublish == 0 || form.IsPublish == 1 {
		form.Status = int(form.IsPublish) + 1
	} else {
		if form.IsPublish != -1 {
			ctx.JSON(http.StatusOK, gin.H{
				"code": "1",
				"msg":  "Invalid form status.",
			})
			return
		}
	}

	if err := fs.UpdateForm(&form); err != nil {
		returnMySQLError(ctx, err)
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"code": "0",
		"msg":  "Success.",
		"data": form,
	})
}

func (c *FormController) DeleteForm(ctx *gin.Context) {
	userID, _ := ctx.Get("userId")

	id, err := strconv.ParseInt(ctx.Param("id"), 10, 64)
	if err != nil {
		returnFormIdInvalid(ctx)
		return
	}

	var form model.Form
	form.ID = id
	form.OwnerID = userID.(int64)

	fs := service.FormService{}

	exist, err := fs.CheckFormExist(form)
	if err != nil {
		returnMySQLError(ctx, err)
		return
	}
	if !exist {
		returnFormNotFound(ctx)
		return
	}

	if err := fs.DeleteForm(form.ID); err != nil {
		returnMySQLError(ctx, err)
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"code": "0",
		"msg":  "Success.",
	})
}

func (c *FormController) GenerateFormBody(ctx *gin.Context) {
	userID, _ := ctx.Get("userId")

	id, err := strconv.ParseInt(ctx.Param("id"), 10, 64)
	if err != nil {
		returnFormIdInvalid(ctx)
		return
	}

	var req struct {
		Question string `json:"question" binding:"required"`
	}
	if err := ctx.BindJSON(&req); err != nil {
		ctx.JSON(http.StatusOK, gin.H{
			"code": "1",
			"msg":  "Invalid data.",
			"data": err.Error(),
		})
		return
	}

	var form model.Form
	form.ID = id
	form.OwnerID = userID.(int64)

	fs := service.FormService{}

	exist, err := fs.CheckFormExist(form)
	if err != nil {
		returnMySQLError(ctx, err)
		return
	}
	if !exist {
		returnFormNotFound(ctx)
		return
	}

	if err := fs.GetFormDetail(&form); err != nil {
		returnMySQLError(ctx, err)
		return
	}

	if form.Status == 2 {
		ctx.JSON(http.StatusOK, gin.H{
			"code": "2",
			"msg":  "Form is published.",
		})
		return
	}

	form.Components, err = fs.OpenAIGenerate(req.Question)
	if err != nil {
		ctx.JSON(http.StatusOK, gin.H{
			"code": "3",
			"msg":  "Generate form body failed.",
			"data": err.Error(),
		})
		return
	}

	err = fs.UpdateForm(&form)
	if err != nil {
		returnMySQLError(ctx, err)
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"code": "0",
		"msg":  "Success.",
	})
}
