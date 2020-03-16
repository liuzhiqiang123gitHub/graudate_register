package routers

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"graduate_registrator/views"
)

const ServerAdmin = "http://localhost:4420"

func StartHttpServer(port int) {
	router := gin.New()
	//router.Use(httputil.ReqData2Form())
	regis := router.Group("/v2")
	{
		//用户注册
		regis.POST("/registrator", views.UserRegistrator)
		//用户找回密码,暂时使用邮箱
		regis.POST("/get_password_by_email", views.GetPasswordByEmail)
		//用于请求验证码
		regis.POST("/get_validate_code", views.GetValidateCode)

	}
	//err := router.Run(fmt.Sprintf("%s:%d"),"0.0.0.0", port)
	err := router.Run(fmt.Sprintf("%s:%d", "0.0.0.0", port))
	fmt.Println(err)
}
