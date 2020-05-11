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
	financial := router.Group("/v2/financial")
	{
		//充值
		financial.POST("/recharge", views.Recharge)
		//购买装备
		financial.POST("/buy_equipment", views.BuyEquipment)
		//查询用户拥有的的装备
		financial.POST("/get_equipments_by_email", views.GetEquipments)
		//查询用户未拥有的装备
		financial.POST("/get_without_equipments_by_email", views.GetWithoutEquipments)
		//查询所有装备
		financial.POST("/get_equipments", views.GetAllEquipments)
		//查询用户金融信息,比如点券
		financial.POST("/get_coupon_by_email", views.GetCoupon)
		//脚本，批量刷装备信息.注册时就有四样基础装备
		financial.POST("/batch_insert_equipments", views.BatchInsertEquipment)
	}
	err := router.Run(fmt.Sprintf("%s:%d", "0.0.0.0", port))
	fmt.Println(err)
}
