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
	r.GET("/users", middleware.Auth(true), Controller.User.GetUserInfo)

	// Form
	r.POST("/forms", middleware.Auth(true), Controller.Form.CreateForm)
	r.POST("/forms/:id/generateFormBody", middleware.Auth(true), Controller.Form.GenerateFormBody)
	r.POST("/forms/:id/duplicate", middleware.Auth(true), Controller.Form.DuplicateForm)
	r.GET("/forms", middleware.Auth(true), Controller.Form.GetFormList)
	r.GET("/forms/:id", middleware.Auth(false), Controller.Form.GetFormDetail)
	r.PATCH("/forms/:id", middleware.Auth(true), Controller.Form.UpdateForm)
	r.DELETE("/forms/:id", middleware.Auth(true), Controller.Form.DeleteForm)

	// Result
	r.POST("/forms/:id/submit", Controller.Result.SubmitForm)
	r.GET("/forms/:id/result", middleware.Auth(true), Controller.Result.GetFormResult)
	r.GET("/forms/:id/results", middleware.Auth(true), Controller.Result.GetFormResultsList)
	r.GET("/forms/:id/results/:rid", middleware.Auth(true), Controller.Result.GetFormResultsDetail)

	return r
}
