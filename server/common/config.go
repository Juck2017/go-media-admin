package common

import (
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"suitbim.com/go-media-admin/utils"
)
//func GetGlobal() *GlobalConfig {
//	return &g
//}

type application struct {
	Port int `mapstructure:"port"`
	Mode string `mapstructure:"mode"`
	Log  struct {
		File  string `mapstructure:"file"`
		Level string `mapstructure:"level"`
	} `mapstructure:"log"`
	Db struct {
		Path string `mapstructure:"path"`
		DdlAuto string `mapstructure:"ddl-auto"`
	} `mapstructure:"db"`
}

//type SHP struct {
//	Schema string
//	Host   string
//	Port   int
//}

type dispatcher struct {
	Stream struct {
		Schema string `mapstructure:"schema"`
		Host   string `mapstructure:"host"`
		Port   int `mapstructure:"port"`
	} `mapstructure:"stream"`
	Signal struct {
		Schema string `mapstructure:"schema"`
		Host   string `mapstructure:"host"`
		Port   int `mapstructure:"port"`
	} `mapstructure:"signal"`
	Viewer struct {
		Schema string `mapstructure:"schema"`
		Host   string `mapstructure:"host"`
		Port   int `mapstructure:"port"`
	} `mapstructure:"viewer"`
}

// 定义一个全局配置的结构体,编译viper进行解析
type Config struct {
	// 应用程序实例
	App application `mapstructure:"app"`
	// 远端流分发服务器相关配置实例
	Dispatcher dispatcher `mapstructure:"dispatcher"`
	// Zap日志实例
	Logger *zap.SugaredLogger
}

// 结构体转json字符串
func (g *Config) String() string {
	return utils.Struct2JsonString(g)
}

// 加载viper配置
func LoadViperConfig(config interface{}, filePath string) {
	v := viper.New()
	v.SetConfigFile(filePath)
	handError(v.ReadInConfig())
	handError(v.Unmarshal(config))
}

// 错误处理
func handError(err error) {
	if err != nil {
		panic(err)
	}
}

// Init 初始化配置文件解析
func init() {
	// 加载配置文件
	LoadViperConfig(&GlobalConf, "./config.yaml")

	// 初始化zap日志实例
	GlobalConf.Logger = utils.InitLogger(
		GlobalConf.App.Log.Level,
		GlobalConf.App.Log.File,
		GlobalConf.App.Mode == "dev")

	GlobalConf.Logger.Infof("\n%s", GetBannerString())
	GlobalConf.Logger.Info("==== Global init viper config success ====")
	GlobalConf.Logger.Info("==== Global init zap log success ====")
	// 初始化数据库
	InitDb()

	// 全局配置初始化完成
	GlobalConf.Logger.Info("==== Global all configuration initialled ====")
	// g.Log.Debugf("==== Global configuration value:\n%s\n", g.String())
}
