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

	return r
}
