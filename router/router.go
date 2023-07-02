package router

import (
	"github.com/ljcbaby/form/controller"
	"github.com/ljcbaby/form/middleware"

	"github.com/gin-gonic/gin"
)

// SetupRouter 配置路由
func SetupRouter() *gin.Engine {
	// gin.SetMode(gin.ReleaseMode)

	r := gin.Default()

	// 创建控制器实例
	Controller := &controller.Controller{}

	// 中间件
	r.Use(middleware.CORS())

	// 路由配置
	r.GET("/", Controller.Index)

	// User
	r.POST("/users/login", Controller.User.Login)
	r.POST("/users/register", Controller.User.Register)
	r.GET("/users", middleware.Auth(), Controller.User.GetUserInfo)

	// Form
	r.POST("/forms", middleware.Auth(), Controller.Form.CreateForm)
	r.POST("/forms/:id/duplicate", middleware.Auth(), Controller.Form.DuplicateForm)
	r.GET("/forms", middleware.Auth(), Controller.Form.GetFormList)
	r.GET("/forms/:id", middleware.Auth(), Controller.Form.GetFormDetail)
	r.PATCH("/forms/:id", middleware.Auth(), Controller.Form.UpdateForm)
	r.DELETE("/forms/:id", middleware.Auth(), Controller.Form.DeleteForm)
	r.POST("/forms/:id/submit", middleware.Auth(), Controller.Form.SubmitForm)
	r.GET("/forms/:id/results", middleware.Auth(), Controller.Form.GetFormResults)

	return r
}
