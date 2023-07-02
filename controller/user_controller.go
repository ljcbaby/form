package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/ljcbaby/form/model"
	"github.com/ljcbaby/form/service"
	"github.com/ljcbaby/form/util"
)

type UserController struct{}

func (c *UserController) Login(ctx *gin.Context) {
	var user model.User
	if err := ctx.ShouldBindJSON(&user); err != nil {
		ctx.JSON(http.StatusOK, gin.H{
			"code": "1",
			"msg":  "Login failed.",
			"data": err.Error(),
		})
		return
	}
	if len(user.Username) < 1 || len(user.Username) > 32 || !util.IsAscii(user.Username) ||
		len(user.Password) != 32 || !util.IsHex(user.Password) {
		ctx.JSON(http.StatusOK, gin.H{
			"code": "1",
			"msg":  "Login failed.",
		})
		return
	}

	us := service.UserService{}

	jwt, err := us.Login(user)
	if err != nil {
		ctx.JSON(http.StatusOK, gin.H{
			"code": "1",
			"msg":  "Login failed.",
			"data": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"code": "0",
		"msg":  "Success.",
		"data": gin.H{
			"jwt": jwt,
		},
	})
}

func (c *UserController) Register(ctx *gin.Context) {
	var user model.User
	if err := ctx.ShouldBindJSON(&user); err != nil {
		ctx.JSON(http.StatusOK, gin.H{
			"code": "1",
			"msg":  "Invalid input.",
			"data": err.Error(),
		})
		return
	}
	if len(user.Username) < 1 || len(user.Username) > 32 || !util.IsAscii(user.Username) ||
		len(user.Password) != 32 || !util.IsHex(user.Password) ||
		len(user.Nickname) < 1 || len(user.Nickname) > 32 {
		ctx.JSON(http.StatusOK, gin.H{
			"code": "1",
			"msg":  "Content illegal.",
		})
		return
	}

	us := service.UserService{}

	if exist, err := us.CheckUsernameExist(user.Username); err != nil {
		returnMySQLError(ctx, err)
		return
	} else if exist {
		ctx.JSON(http.StatusOK, gin.H{
			"code": "2",
			"msg":  "Username already exists.",
		})
		return
	}

	user.Salt = util.GenerateSalt()
	user.Password = util.MD5(user.Password + user.Salt)

	if err := us.Register(&user); err != nil {
		returnMySQLError(ctx, err)
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"code": "0",
		"msg":  "Success.",
	})
}

func (c *UserController) GetUserInfo(ctx *gin.Context) {
	userID, _ := ctx.Get("userId")

	var user model.User
	user.ID = userID.(int64)

	us := service.UserService{}

	if err := us.GetUserByID(&user); err != nil {
		if err.Error() == "USER_NOT_FOUND" {
			ctx.JSON(http.StatusOK, gin.H{
				"code": "1",
				"msg":  "User not found.",
			})
			return
		}
		returnMySQLError(ctx, err)
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"code": "0",
		"msg":  "Success.",
		"data": gin.H{
			"username": user.Username,
			"nickname": user.Nickname,
		},
	})
}
