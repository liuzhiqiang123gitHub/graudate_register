package main

import (
	"fmt"
	"graduate_registrator/routers"
	config "graduate_registrator/utils/conf"
	"graduate_registrator/utils/dbutil"
	etcdIni "graduate_registrator/utils/etcd"
	"graduate_registrator/utils/kafkautil"
	"graduate_registrator/utils/redisUtil"
	"runtime"
	"time"
)

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())
	//初始化etcd
	etcdIni.InitEtcd()
	//从etcd加载配置
	err := dbutil.InitDb(config.Conf)
	if err != nil {
		fmt.Println(err)
		return
	}
	err = redisUtil.InitRedis(config.Conf)
	if err != nil {
		fmt.Println(err)
		return
	}
	//连接kakfa
	go kafkautil.InitConsumer()
	go kafkautil.InitProducer()
	time.Sleep(time.Second)
	if kafkautil.G_Consumer != nil && kafkautil.G_SyncProducer != nil {
		go kafkautil.ConsumeMsg()
	}
	routers.StartHttpServer(9370)

	//res, err := redisUtil.Get("123")
	//if res == "" && err != nil {
	//	fmt.Println(err)
	//}
	//fmt.Print(63 & 1)
}
