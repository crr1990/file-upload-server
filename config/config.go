package config
import (
	"github.com/fsnotify/fsnotify"
	"github.com/lexkong/log"
	"github.com/spf13/viper"
	"strings"
)

type Config struct {
	Name string
}

//初始化 配置文件
func Init(cfg string) error {
	config := Config{Name: cfg}
	if err := config.initConfig(); err != nil {
		return err
	}


	//初始话日志文件
	config.initLog()

	//监控配置文件并热加载程序
	config.watchConfig()
	return nil
}

//初始化配置文件
func (c *Config) initConfig() error {
	if c.Name != "" {
		viper.SetConfigFile(c.Name) //如果制定了配置文件 ，那么就解析这个指定的
	} else {
		// ./conf/config.yaml
		viper.AddConfigPath("./config") //添加需要解析配置文件的路径
		viper.SetConfigName("config") //添加需要解析的文件
	}
	viper.SetConfigType("yaml") //设置配置文件的格式为： yaml

	// 通过指定配置文件可以很方便地连接不同的环境（开发环境、测试环境）并加载不同的配置，方便开发和测试。
	viper.AutomaticEnv()            //读取匹配的环境变量
	viper.SetEnvPrefix("APISERVER") // 读取环境变量的前缀为APISERVER
	replacer := strings.NewReplacer(".", "_")
	viper.SetEnvKeyReplacer(replacer)

	if err := viper.ReadInConfig(); err != nil { // viper解析配置文件
		log.Infof("错误是", err)
		return err
	}
	return nil
}

//初始化Log 日志包
func (c *Config) initLog() {
	//读取配置文件信息
	passLagerCfg := log.PassLagerCfg{
		Writers:        viper.GetString("log.writers"),
		LoggerLevel:    viper.GetString("log.logger_level"),
		LoggerFile:     viper.GetString("log.logger_file"),
		LogFormatText:  viper.GetBool("log.log_format_text"),
		RollingPolicy:  viper.GetString("log.rollingPolicy"),
		LogRotateDate:  viper.GetInt("log.log_rotate_date"),
		LogRotateSize:  viper.GetInt("log.log_rotate_size"),
		LogBackupCount: viper.GetInt("log.log_backup_count"),
	}
	log.InitWithConfig(&passLagerCfg)
}

// 监控配置文件变化并热加载程序
func (c *Config) watchConfig() {
	viper.WatchConfig()
	viper.OnConfigChange(func(e fsnotify.Event) {
		log.Infof("Config file changed: %s", e.Name)
	})
}
