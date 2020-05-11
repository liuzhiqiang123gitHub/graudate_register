package views

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"graduate_registrator/controllers"
	"graduate_registrator/model"
	"graduate_registrator/utils"
	"graduate_registrator/utils/email"
	"graduate_registrator/utils/httputils"
	"graduate_registrator/utils/redisUtil"
)

type GetRegistratorReq struct {
	Email        string `form:"email" json:"email" binding:"required"`
	Phone        string `form:"phone" json:"phone"`
	Password     string `form:"password" json:"password" binding:"required"`
	NickName     string `form:"nickname" json:"nickname" binding:"required"`
	Age          int    `form:"age" json:"age" binding:"required"`
	ValidateCode string `form:"validate_code" json:"validate_code" binding:"required"`
}

type GetRegistratorRsp struct {
	Status      string      `json:"status"`
	Description string      `json:"description"`
	Data        interface{} `json:"data"`
}

func UserRegistrator(c *gin.Context) {
	req := &GetRegistratorReq{}
	//rsp := GetRegistratorRsp{}
	if err := c.Bind(req); err != nil {
		fmt.Printf("%+v", req)
		err := errors.New("invalid params")
		//clog.Logger.Warning("LoginController failed to %v", err.Error())
		fmt.Printf("UserRegistrator failed to %v", err.Error())
		httputils.ResponseError(c, "", err.Error())
		return
	}
	//验证邮箱
	isEmail := email.EmailValidate(req.Email)
	if !isEmail {
		fmt.Println("邮箱格式不合法")
		err := errors.New("邮箱不合法")
		httputils.ResponseError(c, "", err.Error())
		return
	} else if len(req.Email) < 8 {
		fmt.Println("该邮箱不能注册")
		err := errors.New("该邮箱不能注册")
		httputils.ResponseError(c, "", err.Error())
		return
	}
	//验证phone
	if req.Phone != "" {
		isMobile := email.VerifyMobileFormat(req.Phone)
		if !isMobile {
			fmt.Println("手机号不合法")
			err := errors.New("手机号不合法")
			httputils.ResponseError(c, "", err.Error())
			return
		}
	}
	//验证验证码格式是否有效
	if len(req.ValidateCode) != 6 {
		fmt.Println("验证码格式无效")
		err := errors.New("验证码格式无效")
		httputils.ResponseError(c, "", err.Error())
		return
	}
	//验证密码
	if len(req.Password) < 8 || len(req.Password) > 18 {
		fmt.Println("密码少于8位")
		err := errors.New("密码长度8-18位")
		httputils.ResponseError(c, "", err.Error())
		return
	}
	//验证年龄
	if req.Age < 15 || req.Age > 100 {
		fmt.Println("根据您的年龄,不能注册")
		err := errors.New("该年龄段不能注册")
		httputils.ResponseError(c, "", err.Error())
		return
	}
	err := controllers.UserRegistrator(req.Email, req.Phone, req.Password, req.NickName, req.ValidateCode, req.Age)
	if err != nil {
		fmt.Println(err)
		httputils.ResponseError(c, "", err.Error())
		return
	}
	httputils.ResponseOk(c, "", "")
	return

}

type GetPasswordReq struct {
	Mail         string `form:"mail" json:"mail" binding:"required"`
	ValidateCode string `form:"validate_code" json:"validate_code" binding:"required"`
	Password     string `form:"password" json:"password" binding:"required"`
}
type GetPasswordRsp struct {
	Status      string      `json:"status"`
	Description string      `json:"description"`
	Data        interface{} `json:"data"`
}

func GetPasswordByEmail(c *gin.Context) {
	req := &GetPasswordReq{}
	//rsp := GetPasswordRsp{}
	if err := c.Bind(req); err != nil {
		fmt.Printf("%+v", req)
		err := errors.New("invalid params")
		//clog.Logger.Warning("LoginController failed to %v", err.Error())
		fmt.Printf("GetPasswordByEmail failed to %v", err.Error())
		httputils.ResponseError(c, "", err.Error())
		return
	}
	err := controllers.GetPasswordByEmail(nil, req.Mail, req.ValidateCode, req.Password)
	if err != nil {
		httputils.ResponseError(c, "", err.Error())
		return
	}
	httputils.ResponseOk(c, "", "")
	fmt.Printf("LoginController req=%+v ", req)
	return
}

type GetValidateReq struct {
	Mail string `form:"mail" json:"mail" binding:"required"`
}
type GetValidateData struct {
}
type GetValidateRsp struct {
	Status      string `json:"status"`
	Description string `json:"description"`
	Data        string `json:"data"`
}

func GetValidateCode(c *gin.Context) {
	req := &GetValidateReq{}
	//rsp := GetValidateRsp{}
	if err := c.Bind(req); err != nil {
		fmt.Printf("%+v", req)
		err := errors.New("invalid params")
		//clog.Logger.Warning("LoginController failed to %v", err.Error())
		fmt.Printf("GetPasswordByEmail failed to %v", err.Error())
		httputils.ResponseError(c, "", err.Error())
		return
	}
	fmt.Printf("LoginController req=%v ", req)
	//验证邮箱
	isEmail := email.EmailValidate(req.Mail)
	if !isEmail {
		fmt.Println("邮箱不合法")
		err := errors.New("邮箱不合法")
		httputils.ResponseError(c, "", err.Error())
	}
	strCode := utils.GenValidateCode(6)
	//向该用户发送邮件
	body := fmt.Sprintf("您的验证码为:%s,5分钟之后过期,如不是本人操作请忽略。", strCode)
	err := email.StartSendEmail(req.Mail, "字节飞舞计算机系统有限公司", body)
	if err != nil {
		fmt.Printf("发送邮件失败%s", req.Mail)
		return
	}
	//删除key
	err = redisUtil.Delete(req.Mail)
	if err != nil {
		fmt.Println("删除key失败")
		httputils.ResponseError(c, "", err.Error())
	}
	var value = fmt.Sprintf("%s_%s", req.Mail, strCode)
	err = redisUtil.Set(req.Mail, value, 300)
	if err != nil {
		fmt.Println("写入redis失败")
		return
	}
	httputils.ResponseOk(c, "", "")
	return
}

type GetRechargeReq struct {
	Email       string `form:"email" json:"email" binding:"required"`
	RechargeNum int    `form:"recharge_num" json:"recharge_num" binding:"required"`
}

type GetRechargeRsp struct {
	Status      string `json:"status"`
	Description string `json:"description"`
	Data        string `json:"data"`
}

func Recharge(c *gin.Context) {
	req := &GetRechargeReq{}
	//rsp := GetValidateRsp{}
	if err := c.Bind(req); err != nil {
		fmt.Printf("%+v", req)
		err := errors.New("invalid params")
		//clog.Logger.Warning("LoginController failed to %v", err.Error())
		fmt.Printf("Recharge failed to %v", err.Error())
		httputils.ResponseError(c, "", err.Error())
		return
	}
	fmt.Printf("recharge req=%v ", req)
	err := controllers.Recharge(req.RechargeNum, req.Email)
	if err != nil {
		err := errors.New("充值失败")
		//clog.Logger.Warning("LoginController failed to %v", err.Error())
		fmt.Printf("Recharge failed to %v", err.Error())
		httputils.ResponseError(c, "", err.Error())
		return
	}
	httputils.ResponseOk(c, "", "")
	return
}

type GetBuyEquipmentReq struct {
	Email     string `form:"email" json:"email" binding:"required"`
	CouponNum int    `form:"coupon_num" json:"coupon_num" binding:"required"`
	WeaponId  int    `form:"weapon_id" json:"weapon_id" binding:"required"`
}

type GetBuyEquipmentRsp struct {
	Status      string `json:"status"`
	Description string `json:"description"`
	Data        string `json:"data"`
}

func BuyEquipment(c *gin.Context) {
	req := &GetBuyEquipmentReq{}
	//rsp := GetValidateRsp{}
	if err := c.Bind(req); err != nil {
		fmt.Printf("%+v", req)
		err := errors.New("invalid params")
		//clog.Logger.Warning("LoginController failed to %v", err.Error())
		fmt.Printf("BuyEquipment failed to %v", err.Error())
		httputils.ResponseError(c, "", err.Error())
		return
	}
	fmt.Printf("BuyEquipment req=%v ", req)
	err := controllers.BuyEquipment(req.CouponNum, req.WeaponId, req.Email)
	if err != nil {
		//clog.Logger.Warning("LoginController failed to %v", err.Error())
		fmt.Printf("BuyEquipment failed to %v", err.Error())
		httputils.ResponseError(c, "", err.Error())
		return
	}
	httputils.ResponseOk(c, "", "")
	return
}

type GetAllEquipmentsReq struct {
	Email string `form:"email" json:"email" binding:"required"`
}

type GetAllEquipmentsRsp struct {
	Status      string                `json:"status"`
	Description string                `json:"description"`
	Data        model.WeaponModelList `json:"data"`
}

func GetAllEquipments(c *gin.Context) {
	req := &GetAllEquipmentsReq{}
	//rsp := GetValidateRsp{}
	if err := c.Bind(req); err != nil {
		fmt.Printf("%+v", req)
		err := errors.New("invalid params")
		//clog.Logger.Warning("LoginController failed to %v", err.Error())
		fmt.Printf("GetAllEquipments failed to %v", err.Error())
		httputils.ResponseError(c, "", err.Error())
		return
	}
	fmt.Printf("BuyEquipment req=%v ", req)
	data, err := controllers.GetAllEquipments(req.Email)
	if err != nil {
		err := errors.New("获得装备失败，请稍后重试")
		//clog.Logger.Warning("LoginController failed to %v", err.Error())
		fmt.Printf("GetAllEquipments failed to %v", err.Error())
		httputils.ResponseError(c, "", err.Error())
		return
	}
	httputils.ResponseOk(c, data, "")
	return
}

type GetCouponByEmailReq struct {
	Email string `form:"email" json:"email" binding:"required"`
}

type GetCouponByEmailRsp struct {
	Status      string                 `json:"status"`
	Description string                 `json:"description"`
	Data        model.UserCapitalModel `json:"data"`
}

//获得点券信息
func GetCoupon(c *gin.Context) {
	req := &GetCouponByEmailReq{}
	//rsp := GetValidateRsp{}
	if err := c.Bind(req); err != nil {
		fmt.Printf("%+v", req)
		err := errors.New("invalid params")
		//clog.Logger.Warning("LoginController failed to %v", err.Error())
		fmt.Printf("GetCoupon failed to %v", err.Error())
		httputils.ResponseError(c, "", err.Error())
		return
	}
	fmt.Printf("BuyEquipment req=%v ", req)
	err, data := controllers.GetCoupon(req.Email)
	if err != nil {
		err := errors.New("获得用户金融账户失败，请稍后重试")
		//clog.Logger.Warning("LoginController failed to %v", err.Error())
		fmt.Printf("GetAllEquipments failed to %v", err.Error())
		httputils.ResponseError(c, "", err.Error())
		return
	}
	httputils.ResponseOk(c, data, "")
	return
}

type GetEquipmentsReq struct {
	Email string `form:"email" json:"email" binding:"required"`
}

type GetEquipmentsRsp struct {
	Status      string                    `json:"status"`
	Description string                    `json:"description"`
	Data        model.UserWeaponModelList `json:"data"`
}

//获得已有装备信息
func GetEquipments(c *gin.Context) {
	req := &GetCouponByEmailReq{}
	//rsp := GetValidateRsp{}
	if err := c.Bind(req); err != nil {
		fmt.Printf("%+v", req)
		err := errors.New("invalid params")
		//clog.Logger.Warning("LoginController failed to %v", err.Error())
		fmt.Printf("GetEquipments failed to %v", err.Error())
		httputils.ResponseError(c, "", err.Error())
		return
	}
	fmt.Printf("BuyEquipment req=%v ", req)
	err, data := controllers.GetEquipments(req.Email)
	if err != nil {
		err := errors.New("获得已有装备失败，请稍后重试")
		//clog.Logger.Warning("LoginController failed to %v", err.Error())
		fmt.Printf("GetAllEquipments failed to %v", err.Error())
		httputils.ResponseError(c, "", err.Error())
		return
	}
	httputils.ResponseOk(c, data, "")
	return
}

type GetWithoutEquipmentsReq struct {
	Email string `form:"email" json:"email" binding:"required"`
}

type GetWithoutEquipmentsRsp struct {
	Status      string                    `json:"status"`
	Description string                    `json:"description"`
	Data        model.UserWeaponModelList `json:"data"`
}

//获得已有装备信息
func GetWithoutEquipments(c *gin.Context) {
	req := &GetWithoutEquipmentsReq{}
	//rsp := GetValidateRsp{}
	if err := c.Bind(req); err != nil {
		fmt.Printf("%+v", req)
		err := errors.New("invalid params")
		//clog.Logger.Warning("LoginController failed to %v", err.Error())
		fmt.Printf("GetWithoutEquipments failed to %v", err.Error())
		httputils.ResponseError(c, "", err.Error())
		return
	}
	fmt.Printf("BuyEquipment req=%v ", req)
	err, data := controllers.GetWithoutEquipments(req.Email)
	if err != nil {
		err := errors.New("获得已有装备失败，请稍后重试")
		//clog.Logger.Warning("LoginController failed to %v", err.Error())
		fmt.Printf("GetAllEquipments failed to %v", err.Error())
		httputils.ResponseError(c, "", err.Error())
		return
	}
	httputils.ResponseOk(c, data, "")
	return
}

type GetBatchInsertEquipmentReq struct {
}

type GetBatchInsertEquipmentRsp struct {
	Status      string `json:"status"`
	Description string `json:"description"`
	Data        string `json:"data"`
}

//获得已有装备信息
func BatchInsertEquipment(c *gin.Context) {
	req := &GetBatchInsertEquipmentReq{}
	//rsp := GetValidateRsp{}
	if err := c.Bind(req); err != nil {
		fmt.Printf("%+v", req)
		err := errors.New("invalid params")
		//clog.Logger.Warning("LoginController failed to %v", err.Error())
		fmt.Printf("GetWithoutEquipments failed to %v", err.Error())
		httputils.ResponseError(c, "", err.Error())
		return
	}
	//开启事务
	usM := model.UserWeaponModel{}
	ts := usM.GetDB().Begin()
	commit := false
	defer func() {
		if commit {
			ts.Commit()
		} else {
			ts.Rollback()
		}
	}()
	//查询用户信息表，批量导入
	userInfoList := model.UserInfoModelList{}
	err := userInfoList.GetAllUsers()
	if err != nil {
		fmt.Printf("userInfoList.GetAllUsers() err [err=%v]", err)
		err := errors.New("批量导入失败")
		//clog.Logger.Warning("LoginController failed to %v", err.Error())
		fmt.Printf("GetWithoutEquipments failed to %v", err.Error())
		httputils.ResponseError(c, "", err.Error())
		return
	}
	fmt.Printf("vip %+v", userInfoList)
	userWeapon := model.UserWeaponModel{}
	for _, userInfo := range userInfoList {
		//查询user_weapon表是否有该账户，如果有则跳过
		err = userWeapon.GetByEmail(userInfo.Email)
		if err != nil && err != gorm.ErrRecordNotFound {
			fmt.Printf("GetByEmail [err=%v]", err)
			err := errors.New("批量导入失败")
			//clog.Logger.Warning("LoginController failed to %v", err.Error())
			fmt.Printf("GetWithoutEquipments failed to %v", err.Error())
			httputils.ResponseError(c, "", err.Error())
			return
		} else if err != nil && err == gorm.ErrRecordNotFound {
			//导入3条基础数据
			//0，1，5
			userWeapon.Create(userInfo.Email, 1)
			userWeapon = model.UserWeaponModel{}
			userWeapon.Create(userInfo.Email, 2)
			userWeapon = model.UserWeaponModel{}
			userWeapon.Create(userInfo.Email, 5)
			userWeapon = model.UserWeaponModel{}
		} else {
			userWeapon = model.UserWeaponModel{}
			continue
		}
	}
	commit = true
	httputils.ResponseOk(c, "", "")
	return
}
