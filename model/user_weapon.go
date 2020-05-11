package model

import (
	"errors"
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

func (uWM *UserWeaponModel) GetDB() *gorm.DB {
	return dbutil.RegistratorDBPool.Table(UserWeaponTableName)
}

//购买武器
func (uWM *UserWeaponModel) BuyWeaponByEmail(couponNum int, weaponId int, email string) error {
	//查询该装备是否存在
	err := uWM.GetByEmailAndWeaponId(email, weaponId)
	if err != nil && err != gorm.ErrRecordNotFound {
		fmt.Printf("系统开小差了，请稍后再试,err=%v", err)
		err = errors.New("系统开小差了，请稍后再试")
		return err
	} else if err == nil {
		fmt.Printf("装备已存在")
		err = errors.New("装备已存在")
		return err
	}
	uWM.Email = email
	uWM.WeaponId = weaponId
	//点券扣减
	uCM := UserCapitalModel{}
	err = uCM.SubstanceCouponByEmail(couponNum, email)
	if err != nil {
		fmt.Printf(" uCM.SubstanceCouponByEmail err ,err= %v", err)
		return err
	}
	//插入数据
	return dbutil.RegistratorDBPool.Table(UserWeaponTableName).Create(uWM).Error
}

//查询拥有的装备
func (uWML *UserWeaponModelList) GetEquipments(email string) error {
	return dbutil.RegistratorDBPool.Table(UserWeaponTableName).Where("email=?", email).Find(uWML).Error
}

//查询未拥有的
func (uWML *WeaponModelList) GetWithoutEquipments(email string) error {
	sql := "select b.* from user_weapon a right join weapon b on a.weapon_id=b.id where a.email != ?"
	return dbutil.RegistratorDBPool.Table(UserWeaponTableName).Raw(sql, email).Find(uWML).Error
}

//查询一条记录
func (uWM *UserWeaponModel) GetByEmail(email string) error {
	return dbutil.RegistratorDBPool.Table(UserWeaponTableName).Where("email=?", email).First(uWM).Error
}

//插入一条记录
func (uWM *UserWeaponModel) Create(email string, weaponId int) error {
	uWM.Email = email
	uWM.WeaponId = weaponId
	return dbutil.RegistratorDBPool.Table(UserWeaponTableName).Create(uWM).Error
}

//根据email和weapon id查询
func (uWM *UserWeaponModel) GetByEmailAndWeaponId(email string, weaponId int) error {
	uWM.Email = email
	uWM.WeaponId = weaponId
	return dbutil.RegistratorDBPool.Table(UserWeaponTableName).Where("email=? and weapon_id=?", email, weaponId).First(uWM).Error
}
