package config
type MysqlConfig struct {
	MysqlConn            string
	MysqlConnectPoolSize int
}

type PostgresConfig struct {
	PostgresConn            string
	PostgresConnectPoolSize int
}

type LogConfig struct {
	LogDir      string
	LogFile     string
	LogLevel    string
	LogFormat   string
	ProcessName string
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
type Configure struct {
	MysqlSetting    map[string]MysqlConfig
	PostgresSetting map[string]PostgresConfig
	RedisSetting    map[string]RedisConfig
	//RpcSetting      map[string]RPCSetting
	//CelerySetting   map[string]CeleryQueue
	//OssSetting      map[string]OssConfig
	LogDir          string    //不推荐
	LogFile         string    //不推荐
	LogLevel        string    //不推荐
	LogSetting      LogConfig //推荐
	Listen          string
	RpcListen       string
	External        map[string]string
	ExternalInt64   map[string]int64
	GormDebug       bool   //sql 输出开关
	StaticDir       string //静态文件目录设置
	Environment     string //环境变量区分
}

var Config *Configure

