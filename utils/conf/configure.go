package config

var Conf Configure

type MysqlConfig struct {
	MysqlConn            string `json:"mysql_conn"`
	MysqlConnectPoolSize int    `json:"mysql_connect_pool_size"`
}

type RedisConfig struct {
	RedisConn      string
	RedisPasswd    string
	ReadTimeout    int
	ConnectTimeout int
	WriteTimeout   int
	IdleTimeout    int
	MinIdle        int
	MaxIdle        int
	MaxActive      int
	RedisDb        string
	MaxRetries     int
}
type RPCSetting struct {
	Addr string `json:"addr"`
	Net  string `json:"net"`
}

type CeleryQueue struct {
	Url   string `json:"url"`
	Queue string `json:"queue"`
}

type OssConfig struct {
	AccessKeyId     string `json:"access_key_id"`
	AccessKeySecret string `json:"access_key_secret"`
	Region          string `json:"region"`
	Bucket          string `json:"bucket"`
}

type Configure struct {
	MysqlSetting map[string]MysqlConfig `json:"mysql_setting"`
	//PostgresSetting map[string]PostgresConfig
	RedisSetting map[string]RedisConfig `json:"redis_setting"`
	//RpcSetting      map[string]RPCSetting
	//CelerySetting   map[string]CeleryQueue
	//OssSetting      map[string]OssConfig
	LogDir   string //不推荐
	LogFile  string //不推荐
	LogLevel string //不推荐
	//LogSetting      LogConfig //推荐
	Listen        string
	RpcListen     string
	External      map[string]string
	ExternalInt64 map[string]int64
	GormDebug     bool   //sql 输出开关
	StaticDir     string //静态文件目录设置
	Environment   string //环境变量区分
}
