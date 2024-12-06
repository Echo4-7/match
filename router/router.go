package router

import (
	api "Fire/api/v1"
	"Fire/middleware"
	"github.com/gin-gonic/gin"
	"net/http"
)

func NewRouter() *gin.Engine {

	r := gin.Default()
	r.Use(middleware.Cors())
	r.StaticFS("/static", http.Dir("./static")) // 加载静态文件

	v1 := r.Group("api/v1")
	{
		//测试
		v1.GET("ping", func(c *gin.Context) {
			c.String(http.StatusOK, "pong")
		})

		// 用户操作
		v1.POST("user/register", api.UserRegister)
		v1.POST("user/login", api.UserLogin)
		// 发送验证码
		v1.POST("user/send_code", api.SendCheckCode)
		// 校验验证码
		v1.POST("user/check_code", api.CheckCode)
		// 忘记密码
		v1.PUT("user/findPwd", api.FindPwd)

		// 轮播图
		v1.GET("carousels", api.ListCarousel)

		auth := v1.Group("/") // 需要登陆保护  api/v1
		auth.Use(middleware.JWT())
		{
			// 更新操作
			auth.PUT("user/update", api.UserUpdate)
			// 上传头像
			auth.POST("avatar", api.UploadAvatar)
			// 获取用户详细信息
			auth.GET("user/info", api.UserInfo)
		}
	}

	//r.NoRoute(func(c *gin.Context) {
	//	c.JSON(http.StatusOK, gin.H{
	//		"message": "404",
	//	})
	//})
	return r
}
