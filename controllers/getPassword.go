package controllers

import (
	"context"
	"errors"
	"fmt"
	"github.com/jinzhu/gorm"
	basic "graduate_registrator/const"
	"graduate_registrator/model"
	"graduate_registrator/utils/email"
	"graduate_registrator/utils/redisUtil"
)

func GetPasswordByEmail(ctx context.Context, mail, validateCode, password string) (err error) {
	isEmail := email.EmailValidate(mail)
	if !isEmail {
		fmt.Println("邮箱不合法")
		return errors.New("邮箱不合法")
	}

	if len(validateCode) != 6 {
		fmt.Println("验证码不合法")
		return errors.New("验证码不合法")
	}
	if len(password) < 8 {
		fmt.Println("密码不能小于8位")
		return errors.New("密码不能小于8位")
	}
	var key = mail
	//查询验证码是否过期
	res, err := redisUtil.Get(key)
	//if !strings.Contains(err.Error(), "nil returned") && err != nil  {
	//	fmt.Println("redis 查询错误")
	//	return
	//} else if err != nil && strings.Contains(err.Error(), "nil returned") && res == "" {
	//	fmt.Println("redis 查旬空")
	//	return
	//}
	if res.(string) == "" && err != nil {
		return errors.New("验证码或已失效")
	}
	fmt.Println(res.(string))
	validateC := Substring(res.(string), 18, len(res.(string)))
	fmt.Println(validateC)
	if validateCode != validateC {
		return errors.New("验证码或已失效")
	}
	var userInfo model.UserInfoModel
	err = userInfo.GetUserByEmail(mail)
	if err != nil && err != gorm.ErrRecordNotFound {
		fmt.Println("数据库查询失败")
		return
	} else if err != nil && err == gorm.ErrRecordNotFound {
		fmt.Println("用户不存在")
		return errors.New("用户不存在")
	}
	//更新用户信息
	err = userInfo.UpdateInfo(mail, password)
	if err != nil {
		fmt.Println("数据库插入失败")
		return
	}
	return nil
}

func UserRegistrator(mail, phone, password, nickname, validateCode string, age int) error {
	var userInfo model.UserInfoModel
	//验证邮箱是否存在
	err := userInfo.GetUserByEmail(mail)
	if err != nil && err != gorm.ErrRecordNotFound {
		fmt.Printf("数据库异常%s", mail)
		return err
	} else if err == nil {
		fmt.Printf("邮箱已被注册%s", mail)
		return errors.New(basic.EmailHasBeenExisted)
	}
	//验证phone
	if phone != "" {
		userInfo.GetUserByPhone(phone)
		if err != nil && err != gorm.ErrRecordNotFound {
			fmt.Printf("数据库异常%s", phone)
			return err
		} else if err == nil {
			fmt.Printf("手机号已被注册%s", phone)
			return errors.New(basic.PhoneHasBeenExisted)
		}
	}
	//验证游戏昵称
	err = userInfo.GetUserByNick(nickname)
	if err != nil && err != gorm.ErrRecordNotFound {
		fmt.Printf("数据库异常%s", nickname)
		return err
	} else if err == nil {
		fmt.Printf("用户名已存在%s", nickname)
		return errors.New(basic.NicknameHasBeenExisted)
	}
	//查询验证码是否过期
	_, err = redisUtil.Get(mail)
	fmt.Println(err)
	if err != nil {
		fmt.Printf("验证码或已过期%s", mail)
		return errors.New(basic.ValidationCodeExp)
	}

	//数据正常,进行注册
	userInfo.Email = mail
	userInfo.Phone = phone
	userInfo.Age = age
	userInfo.NickName = nickname
	userInfo.Password = password
	err = userInfo.CreateUser()
	if err != nil {
		fmt.Println(err)
		return err
	}
	return nil
}
func Substring(source string, start int, end int) string {
	var r = []rune(source)
	length := len(r)

	if start < 0 || end > length || start > end {
		return ""
	}

	if start == 0 && end == length {
		return source
	}

	return string(r[start:end])
}
