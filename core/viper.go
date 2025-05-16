package core

import (
	"fmt"
	"github.com/fsnotify/fsnotify"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"ops-server/core/internal"
	"ops-server/global"
	"os"
)

// Viper 读取配置文件, 并将数据解析到每个对应的结构体
func Viper(config string) *viper.Viper {

	if config == "" { // 判断命令行参数是否为空
		if configEnv := os.Getenv(internal.ConfigEnv); configEnv == "" { // 判断 internal.ConfigEnv 常量存储的环境变量是否为空
			switch gin.Mode() {
			case gin.DebugMode:
				config = internal.ConfigDefaultFile
				fmt.Printf("您正在使用gin模式的%s环境名称,config的路径为%s\n", gin.Mode(), internal.ConfigDefaultFile)
			case gin.ReleaseMode:
				config = internal.ConfigReleaseFile
				fmt.Printf("您正在使用gin模式的%s环境名称,config的路径为%s\n", gin.Mode(), internal.ConfigReleaseFile)
			case gin.TestMode:
				config = internal.ConfigTestFile
				fmt.Printf("您正在使用gin模式的%s环境名称,config的路径为%s\n", gin.Mode(), internal.ConfigTestFile)
			}
		} else { // internal.ConfigEnv 常量存储的环境变量不为空 将值赋值于config
			config = configEnv
			fmt.Printf("您正在使用%s环境变量,config的路径为%s\n", internal.ConfigEnv, config)
		}
	} else { // 命令行参数不为空 将值赋值于config
		fmt.Printf("您正在使用命令行的-c参数传递的值,config的路径为%s\n", config)
	}

	v := viper.New()
	v.SetConfigFile(config)
	v.SetConfigType("yaml")
	err := v.ReadInConfig()

	if err != nil {
		panic(fmt.Errorf("Fatal error config file: %s \n", err))
	}

	// 监听配置是否变动, 变动后则调用OnConfigChange重新加载配置
	v.WatchConfig()

	v.OnConfigChange(func(e fsnotify.Event) {
		// 重新解析配置到global.GVA_CONFIG
		fmt.Println("config file changed:", e.Name)
		if err = v.Unmarshal(&global.OPS_CONFIG); err != nil {
			fmt.Println(err)
		}
	})

	// 调用Unmarshal解析配置文件到global.GVA_CONFIG
	if err = v.Unmarshal(&global.OPS_CONFIG); err != nil {
		panic(err)
	}

	return v

}
