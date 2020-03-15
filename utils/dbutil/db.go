package dbutil

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
	"graduate_registrator/config"
)

var (
	RegistratorDBPool *gorm.DB
)

//初始化go-sql-driver/mysql 连接池
func InitDbPool(config *config.MysqlConfig) (*sql.DB, error) {

	dbPool, err := sql.Open("mysql", config.MysqlConn)
	if nil != err {
		return nil, err
	}
	dbPool.SetMaxOpenConns(config.MysqlConnectPoolSize)
	dbPool.SetMaxIdleConns(config.MysqlConnectPoolSize / 2)

	err = dbPool.Ping()
	if err != nil {
		return nil, err
	}

	fmt.Println("init db pool OK")
	return dbPool, nil
}

//初始化gorm 连接池
func InitGormDbPool(config *config.MysqlConfig, setLog bool) (err error) {

	RegistratorDBPool, err = gorm.Open("mysql", config.MysqlConn)
	if err != nil {
		fmt.Println("init db err : ", config, err)
		return err
	}

	RegistratorDBPool.DB().SetMaxOpenConns(config.MysqlConnectPoolSize)
	RegistratorDBPool.DB().SetMaxIdleConns(config.MysqlConnectPoolSize / 2)
	if setLog {
		RegistratorDBPool.LogMode(true)
		//db.SetLogger(clog.Logger)
	}
	RegistratorDBPool.SingularTable(true)

	err = RegistratorDBPool.DB().Ping()
	if err != nil {
		return err
	}
	//	fmt.Println("init db pool OK")

	return nil
}
func InitDb() error {
	mysqlConf := &config.MysqlConfig{
		MysqlConn:            fmt.Sprintf("%s:%s@(%s:%d)/%s?charset=utf8&parseTime=True&loc=Local", "root", "Liuzhi19972123", "148.70.248.33", 3306, "graduate_project"),
		MysqlConnectPoolSize: 10,
	}
	err := InitGormDbPool(mysqlConf, true)
	if err != nil {
		fmt.Println(err)
		return err
	}
	return nil
}
