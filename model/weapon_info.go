package model

import (
	"github.com/jinzhu/gorm"
	"graduate_registrator/utils/dbutil"
)

const (
	WeaponInfoTableName = "weapon"
)

type WeaponModel struct {
	Id          int64   `gorm:"primary_key;column:id" json:"id"`
	Name        string  `gorm:"name;column:name" json:"name"`
	Description string  `gorm:"description" json:"description"`
	state       int     `json:"-"`
	ts          gorm.DB `json:"-"`
}
type WeaponModelList []WeaponModel

//获得所有武器
func (wML *WeaponModelList) GetAllWeapons(email string) error {
	return dbutil.RegistratorDBPool.Table(WeaponInfoTableName).Find(wML).Error
}
