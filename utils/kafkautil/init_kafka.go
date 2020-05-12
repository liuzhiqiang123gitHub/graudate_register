package kafkautil

import (
	"fmt"
	"github.com/Shopify/sarama"
	"graduate_registrator/model"
	"strings"
	"time"
)

var G_SyncProducer *SyncProducer
var G_Consumer *Consumer

func InitProducer() {
	defer func() {
		fmt.Println("生产者初始化成功！！")
	}()
	G_SyncProducer, _ = NewSyncProducer([]string{"148.70.248.33:9092"}, sarama.WaitForLocal)
}
func InitConsumer() {
	defer func() {
		fmt.Println("消费者初始化成功！！")
	}()
	G_Consumer, _ = NewConsumer([]string{"148.70.248.33:9092"}, "register-1", []string{"register"}, sarama.OffsetOldest, true)
}
func ProduceMsg(email string) {
	key := fmt.Sprintf("email+%s", time.Now())
	_, _, err := G_SyncProducer.SendJSON("register", key, &email)
	if err != nil {
		fmt.Println(err)
	}
	//fmt.Println(G_Consumer)
	fmt.Printf("生产者msg, [%s:%s]\n", key, email)
	//time.Sleep(time.Second)

}
func ConsumeMsg() {
	for {
		msg, err := G_Consumer.Receive(NeverTimeout)
		if err != nil {
			fmt.Println(err)
		}
		fmt.Printf("消费者topic:%v, partition:%v/%v, %v\n", msg.Topic, msg.Partition, msg.Offset, string(msg.Value))
		email := strings.Trim(string(msg.Value), "\"")
		//接收到注册的用户自动绑定绑定基础装备
		userWeapon := model.UserWeaponModel{}
		userWeapon.Create(email, 1)
		userWeapon = model.UserWeaponModel{}
		userWeapon.Create(email, 2)
		userWeapon = model.UserWeaponModel{}
		userWeapon.Create(email, 6)
		userWeapon = model.UserWeaponModel{}
		//开金融账户
		userCap := model.UserCapitalModel{}
		userCap.Create(4000, email)
	}
}
