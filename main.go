package main

import (
	"fmt"
	"graduate_registrator/routers"
	"graduate_registrator/utils/dbutil"
	"graduate_registrator/utils/redisUtil"
)

func main() {
	err := dbutil.InitDb()
	if err != nil {
		fmt.Println(err)
		return
	}
	err = redisUtil.InitRedis("")
	if err != nil {
		fmt.Println(err)
		return
	}
	routers.StartHttpServer(9370)
	//res, err := redisUtil.Get("123")
	//if res == "" && err != nil {
	//	fmt.Println(err)
	//}
}
