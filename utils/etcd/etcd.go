package etcdIni

import (
	"context"
	"encoding/json"
	"github.com/coreos/etcd/clientv3"
	config "graduate_registrator/utils/conf"
	"time"
)

func InitEtcd() {
	etcdConfig := clientv3.Config{
		Endpoints:   []string{"148.70.248.33:2379"},
		DialTimeout: 10 * time.Second,
	}
	client, err := clientv3.New(etcdConfig)
	if err != nil {
		panic(err)
	}
	kv := clientv3.NewKV(client)
	ctx, _ := context.WithTimeout(context.Background(), 3*time.Second)
	getRsp, err := kv.Get(ctx, "/zhiqiang/login_service/config", clientv3.WithPrefix())
	for _, v := range getRsp.Kvs {
		//fmt.Printf("你疯了吗%+v", string(v.Value))
		json.Unmarshal(v.Value, &config.Conf)
		//fmt.Println(config.Conf)
		//time.Sleep(time.Second * 100)
	}
}
