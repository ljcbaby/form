package controller

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/ljcbaby/form/model"
	"github.com/ljcbaby/form/service"
)

type ResultController struct{}

func (c *ResultController) SubmitForm(ctx *gin.Context) {
	id, err := strconv.ParseInt(ctx.Param("id"), 10, 64)
	if err != nil {
		returnFormIdInvalid(ctx)
		return
	}

	fs := service.FormService{}

	exist, err := fs.CheckFormExist(model.Form{ID: id})
	if err != nil {
		returnMySQLError(ctx, err)
		return
	}
	if !exist {
		returnFormNotFound(ctx)
		return
	}

	var result model.ResultRequest

	if err := ctx.BindJSON(&result); err != nil {
		ctx.JSON(http.StatusOK, gin.H{
			"code": "1",
			"msg":  "Invalid form data.",
			"data": err.Error(),
		})
		return
	}

	rs := service.ResultService{}

	err = rs.SubmitForm(id, result)
	if err != nil {
		returnMySQLError(ctx, err)
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"code": "0",
		"msg":  "Success.",
	})
}

func (c *ResultController) GetFormResult(ctx *gin.Context) {
	userID, _ := ctx.Get("userId")

	id, err := strconv.ParseInt(ctx.Param("id"), 10, 64)
	if err != nil {
		returnFormIdInvalid(ctx)
		return
	}

	rs := service.ResultService{}

	exist, err := rs.CheckFormResultExist(model.Form{ID: id, OwnerID: userID.(int64)})
	if err != nil {
		returnMySQLError(ctx, err)
		return
	}
	if !exist {
		returnFormResultInvalid(ctx)
		return
	}

	var res model.ResultResponse

	err = rs.GetFormResult(id, &res)
	if err != nil {
		returnMySQLError(ctx, err)
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"code": "0",
		"msg":  "Success.",
		"data": res,
	})
}

func (c *ResultController) GetFormResultsList(ctx *gin.Context) {
	userID, _ := ctx.Get("userId")

	id, err := strconv.ParseInt(ctx.Param("id"), 10, 64)
	if err != nil {
		returnFormIdInvalid(ctx)
		return
	}

	fe_id := ctx.DefaultQuery("fe_id", "")
	if fe_id == "" {
		ctx.JSON(http.StatusOK, gin.H{
			"code": "1",
			"msg":  "Invalid Query.",
		})
		return
	}
	page, err := strconv.Atoi(ctx.DefaultQuery("page", "1"))
	if err != nil || page < 1 {
		page = 1
	}
	size, err := strconv.Atoi(ctx.DefaultQuery("size", "10"))
	if err != nil || size < 1 {
		size = 10
	}
	if size > 50 {
		size = 50
	}

	rs := service.ResultService{}

	exist, err := rs.CheckFormResultExist(model.Form{ID: id, OwnerID: userID.(int64)})
	if err != nil {
		returnMySQLError(ctx, err)
		return
	}
	if !exist {
		returnFormResultInvalid(ctx)
		return
	}

	total, err := rs.GetFormResultsCount(id)
	if err != nil {
		returnMySQLError(ctx, err)
		return
	}

	var results []model.ResultList

	err = rs.GetFormResultsList(id, fe_id, page, size, &results)
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
			"results":    results,
		},
	})
}

func (c *ResultController) GetFormResultsDetail(ctx *gin.Context) {
	userID, _ := ctx.Get("userId")

	id, err := strconv.ParseInt(ctx.Param("id"), 10, 64)
	if err != nil {
		returnFormIdInvalid(ctx)
		return
	}
	rid, err := strconv.ParseInt(ctx.Param("rid"), 10, 64)
	if err != nil {
		ctx.JSON(http.StatusOK, gin.H{
			"code": "1",
			"msg":  "Invalid result id.",
		})
		return
	}

	rs := service.ResultService{}

	exist, err := rs.CheckFormResultExist(model.Form{ID: id, OwnerID: userID.(int64)})
	if err != nil {
		returnMySQLError(ctx, err)
		return
	}
	if !exist {
		returnFormResultInvalid(ctx)
		return
	}

	var res []model.Component
	err = rs.GetFormResultsDetail(id, rid, &res)
	if err != nil {
		returnMySQLError(ctx, err)
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"code": "0",
		"msg":  "Success.",
		"data": gin.H{
			"components": res,
		},
	})
}

func (c *ResultController) GetFormResultsFile(ctx *gin.Context) {
	userID, _ := ctx.Get("userId")

	id, err := strconv.ParseInt(ctx.Param("id"), 10, 64)
	if err != nil {
		returnFormIdInvalid(ctx)
		return
	}

	rs := service.ResultService{}

	exist, err := rs.CheckFormResultExist(model.Form{ID: id, OwnerID: userID.(int64)})
	if err != nil {
		returnMySQLError(ctx, err)
		return
	}
	if !exist {
		returnFormResultInvalid(ctx)
		return
	}

	file, err := rs.GetFormResultsFile(id)
	if err != nil {
		returnMySQLError(ctx, err)
		return
	}

	ctx.Header("Content-Transfer-Encoding", "binary")
	ctx.Data(http.StatusOK, "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet", file.Bytes())
}
