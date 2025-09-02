package main

import (
	"catai"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func Init() *gorm.DB {
	// 打开SQLite数据库连接
	db, err := gorm.Open(sqlite.Open("CatEntities.db"), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}
	// 配置信息
	db.AutoMigrate(&Config{})
	config := &Config{}
	db.Where("Key = ?", "Key").Take(config)
	catai.ApiKey = config.Value
	config = &Config{}
	db.Where("Key = ?", "Model").Take(config)
	catai.BaseModel = config.Value
	config = &Config{}
	db.Where("Key = ?", "URL").Take(config)
	catai.BaseURL = config.Value
	return db
}

// Config 全局配置
type Config struct {
	Key   string `gorm:"primaryKey;type:text"`
	Value string `gorm:"type:text"`
}
