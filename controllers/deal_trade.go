package controllers

import (
	"fmt"
	"graduate_registrator/model"
)

func Recharge(rechargeNum int, email string) error {
	uCM := model.UserCapitalModel{}
	err := uCM.RechargeByEmail(rechargeNum, email)
	if err != nil {
		fmt.Printf("充值失败 err=%v", err)
		return err
	}
	return nil
}
func BuyEquipment(couponNum int, weaponId int, email string) error {
	uWM := model.UserWeaponModel{}
	err := uWM.BuyWeaponByEmail(couponNum, weaponId, email)
	if err != nil {
		fmt.Printf("购买装备失败")
		return err
	}
	return nil
}
func GetAllEquipments(email string) (error, model.WeaponModelList) {
	wML := model.WeaponModelList{}
	err := wML.GetAllWeapons(email)
	if err != nil {
		fmt.Println("获得装备失败")
		return err, wML
	}
	return err, wML
}
func GetCoupon(email string) (error, model.UserCapitalModel) {
	uCM := model.UserCapitalModel{}
	err := uCM.GetCouponByEmail(email)
	if err != nil {
		fmt.Printf("用户获得信息失败，err=%v", err)
		return err, uCM
	}
	return nil, uCM
}
func GetEquipments(email string) (error, model.UserWeaponModelList) {
	uWML := model.UserWeaponModelList{}
	err := uWML.GetEquipments(email)
	if err != nil {
		if err != nil {
			fmt.Printf("获取已有的装备失败，err=%v", err)
			return err, uWML
		}
	}
	return nil, uWML
}
func GetWithoutEquipments(email string) (error, model.WeaponModelList) {
	wmL := model.WeaponModelList{}
	err := wmL.GetWithoutEquipments(email)
	if err != nil {
		if err != nil {
			fmt.Printf("获取未有的装备失败，err=%v", err)
			return err, wmL
		}
	}
	return nil, wmL
}
