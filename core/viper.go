package core

import (
	"flag"
	"fmt"
	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
	"ops-server/core/internal"
	"ops-server/global"
	"os"
)

// Viper 读取配置文件, 并将数据解析到每个对应的结构体
func Viper(path ...string) *viper.Viper {
	var config string

	if len(path) == 0 {
		flag.StringVar(&config, "c", "", "choose config file.")
		flag.Parse()

		if config == "" {
			if configEnv := os.Getenv(internal.ConfigEnv); configEnv == "" {
				panic("config is not exits")
			} else {
				fmt.Printf("您正在使用%s环境变量,config的路径为%s\n", internal.ConfigEnv, config)
				config = configEnv
			}
		} else {
			fmt.Printf("您正在使用命令行的-c参数传递的值,config的路径为%s\n", config)
		}
	} else { // 函数传递的可变参数的第一个值赋值于config
		config = path[0]
		fmt.Printf("您正在使用func Viper()传递的值,config的路径为%s\n", config)
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
