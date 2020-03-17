package model

import (
	"graduate_registrator/utils/dbutil"
)

const (
	UserInfoModelName = "user_info"
)

type UserInfoModel struct {
	//	CREATE TABLE `user_info` (
	//	`id` bigint(20) NOT NULL AUTO_INCREMENT,
	//	`email` varchar(36) NOT NULL COMMENT '邮箱',
	//	`phone` int NOT NULL COMMENT '手机号',
	//	`password` varchar(20) NOT NULL COMMENT '密码',
	//	`nickname` varchar(20) NOT NULL COMMENT '绰号',
	//	`state` int(3) NOT NULL COMMENT '用户状态',
	//	PRIMARY KEY (`id`),
	//	UNIQUE KEY `user_mail` (`email`),
	//	UNIQUE KEY `user_phone` (`phone`),
	//	UNIQUE KEY `user_nick` (`nickname`),
	//	KEY `idx_email` (`email`),
	//	KEY `idx_nick` (`nickname`)
	//) ENGINE=InnoDB  DEFAULT CHARSET=utf8mb4 COMMENT='用户信息';
	Id       int64  `gorm:"primary_key;column:id" json:"id"`
	Email    string `gorm:"column:email" json:"email"`
	Phone    string `gorm:"column:phone" json:"phone"`
	Password string `gorm:"column:password" json:"password"`
	NickName string `gorm:"column:nickname" json:"nickname"`
	Age      int    `gorm:"column:age" json:"age"`
	State    int    `gorm:"column:state" json:"state"`
}

func GetName() string {
	return UserInfoModelName
}

//插入用户
func (userInfo *UserInfoModel) CreateUser() error {
	return dbutil.RegistratorDBPool.Table(GetName()).Create(&userInfo).Error
}

//根据邮箱查询
func (userInfo *UserInfoModel) GetUserByEmail(email string) error {
	return dbutil.RegistratorDBPool.Table(GetName()).Where("email=?", email).Last(&userInfo).Error
}

//根据昵称查询
func (userInfo *UserInfoModel) GetUserByNick(nick string) error {
	return dbutil.RegistratorDBPool.Table(GetName()).Where("nickname=?", nick).Last(&userInfo).Error
}

//根据手机号查询
func (userInfo *UserInfoModel) GetUserByPhone(phone string) error {
	return dbutil.RegistratorDBPool.Table(GetName()).Where("phone=?", phone).Last(&userInfo).Error
}

//更新密码
func (userInfo *UserInfoModel) UpdateInfo(email, pass string) error {
	//data := make(map[string]interface{})
	//data["password"]=pass
	userInfo.Password = pass
	return dbutil.RegistratorDBPool.Table(GetName()).Where("email=? ", email).Update(&userInfo).Error
}
