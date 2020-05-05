package model

import (
	"errors"
	"fmt"
	"github.com/jinzhu/gorm"
	"graduate_registrator/utils/dbutil"
)

const (
	UserCapitalTableName = "user_capital"
)

type UserCapitalModel struct {
	Id     int64   `gorm:"primary_key;column:id" json:"id"`
	Email  string  `gorm:"column:email" json:"email"`
	Coin   int     `gorm:"column:coin" json:"coin"`
	Coupon int     `gorm:"column:coupon" json:"coupon"`
	State  int     `gorm:"column:state" json:"state"`
	ts     gorm.DB `json:"-"`
}

func (ucm *UserCapitalModel) GetDb() *gorm.DB {
	return dbutil.RegistratorDBPool.Table(UserCapitalTableName)
}
func (ucm *UserCapitalModel) Create(rechargeNum int, email string) error {
	ucm.Email = email
	ucm.Coupon = rechargeNum
	return dbutil.RegistratorDBPool.Table(UserCapitalTableName).Create(ucm).Error
}

//查询信息
func (ucm *UserCapitalModel) GetCouponByEmail(email string) (err error) {
	return dbutil.RegistratorDBPool.Table(UserCapitalTableName).Where("email=?", email).First(ucm).Error
}

//充值
func (ucm *UserCapitalModel) RechargeByEmail(rechargeNum int, email string) (err error) {
	//先查询再更新
	err = ucm.GetCouponByEmail(email)
	if err != nil && err != gorm.ErrRecordNotFound {
		fmt.Printf("GetCouponByEmail err ,err=%v", err)
		return err
	} else if err != nil && err == gorm.ErrRecordNotFound {
		err = ucm.Create(rechargeNum, email)
		if err != nil {
			fmt.Printf("RechargeByEmail [err=%v]", err)
			return err
		}
	}
	oldNum := ucm.Coupon
	updateData := make(map[string]int, 0)
	updateData["coupon"] = oldNum + rechargeNum
	return dbutil.RegistratorDBPool.Table(UserCapitalTableName).Update(updateData).Error
}

//点券扣减
func (ucm *UserCapitalModel) SubstanceCouponByEmail(substanceNum int, email string) (err error) {
	//先查询再更新
	err = ucm.GetCouponByEmail(email)
	if err != nil {
		fmt.Printf("SubstanceCouponByEmail err ,err=%v", err)
		return err
	}
	oldNum := ucm.Coupon
	curNum := 0
	if curNum = oldNum - substanceNum; curNum < 0 {
		fmt.Println("余额不足")
		err = errors.New("余额不足")
		return err
	}
	updateData := make(map[string]int, 0)
	updateData["coupon"] = curNum
	return dbutil.RegistratorDBPool.Table(UserCapitalTableName).Update(updateData).Error
}
