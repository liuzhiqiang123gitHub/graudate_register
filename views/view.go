package views

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"graduate_registrator/controllers"
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
	if len(req.Password) < 8 {
		fmt.Println("密码少于8位")
		err := errors.New("密码少于8位")
		httputils.ResponseError(c, "", err.Error())
		return
	}
	//验证年龄
	if req.Age < 15 {
		fmt.Println("根据您的年龄,不能注册")
		err := errors.New("根据您的年龄,不能注册")
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
	err := email.StartSendEmail(req.Mail, "至强科技", body)
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
