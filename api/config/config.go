package config

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// ConfigType 配置项
type ConfigType struct {
	Key   string `gorm:"primaryKey;type:text"`
	Value string `gorm:"type:text"`
}

// ConfigList 配置列表
var ConfigList map[string]string = map[string]string{}

// Quit 退出
var Quit chan os.Signal = make(chan os.Signal, 1)

func init() {
	// 打开SQLite数据库连接
	db, err := gorm.Open(sqlite.Open("Config.db"), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}
	// 读取配置
	var configList []ConfigType
	config := db.Table("configs")
	config.Find(&configList)
	for _, v := range configList {
		ConfigList[v.Key] = v.Value
	}
	// 拦截信号
	signal.Notify(Quit, syscall.SIGABRT, syscall.SIGALRM, syscall.SIGBUS, syscall.SIGCHLD, syscall.SIGCONT, syscall.SIGFPE, syscall.SIGHUP, syscall.SIGILL, syscall.SIGINT, syscall.SIGKILL, syscall.SIGPIPE, syscall.SIGPOLL, syscall.SIGPROF, syscall.SIGPWR, syscall.SIGQUIT, syscall.SIGSEGV, syscall.SIGSTKFLT, syscall.SIGSTOP, syscall.SIGTERM, syscall.SIGTRAP, syscall.SIGTSTP, syscall.SIGTTIN, syscall.SIGTTOU, syscall.SIGUNUSED, syscall.SIGUSR1, syscall.SIGUSR2, syscall.SIGVTALRM, syscall.SIGWINCH, syscall.SIGXCPU, syscall.SIGXFSZ)
	// 关闭更新
	go func() {
		<-Quit
		// 处理写入
		configList = configList[:0]
		for k, v := range ConfigList {
			configList = append(configList, ConfigType{Key: k, Value: v})
		}
		config.Save(configList)
		// 使应用程序正常退出
		os.Exit(0)
	}()
}

// Close 关闭并写出配置
func Close() {
	Quit <- nil
	<-make(chan struct{})
}

// Get 得到配置项
func Get(key string, Value string) string {
	if v, ok := ConfigList[key]; ok {
		return v
	}
	if Value == "" {
		panic(fmt.Errorf("配置项 %s 默认值不存在", key))
	}
	ConfigList[key] = Value
	return Value
}
