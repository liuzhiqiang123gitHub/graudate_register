package model

import (
	"fmt"
	"github.com/jinzhu/gorm"
	"graduate_registrator/utils/dbutil"
)

const (
	UserWeaponTableName = "user_weapon"
)

type UserWeaponModel struct {
	Id          int64   `gorm:"primary_key;column:id" json:"id"`
	Email       string  `email:"name;column:email" json:"email"`
	WeaponId    int     `gorm:"weapon_id" json:"weapon_id"`
	Description string  `gorm:"description" json:"description"`
	State       int     `gorm:"column:state" json:"state"`
	ts          gorm.DB `json:"-"`
}
type UserWeaponModelList []UserWeaponModel

//购买武器
func (uWM *UserWeaponModel) BuyWeaponByEmail(couponNum int, weaponId int64, email string) error {
	uWM.Email = email
	uWM.Id = weaponId
	//点券扣减
	uCM := UserCapitalModel{}
	err := uCM.SubstanceCouponByEmail(couponNum, email)
	if err != nil {
		fmt.Printf(" uCM.SubstanceCouponByEmail err ,err= %v", err)
		return err
	}
	//插入数据
	return dbutil.RegistratorDBPool.Table(WeaponInfoTableName).Create(uWM).Error
}

//查询拥有的装备
func (uWML *UserWeaponModelList) GetEquipments(email string) error {
	return dbutil.RegistratorDBPool.Table(WeaponInfoTableName).Where("email=?", email).Find(uWML).Error
}

//查询未拥有的
func (uWML *UserWeaponModelList) GetWithoutEquipments(email string) error {
	return dbutil.RegistratorDBPool.Table(WeaponInfoTableName).Where("email !=?", email).Find(uWML).Error
}
