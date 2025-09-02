package main

import (
	"catai"
	"catai/sytem"
	"encoding/json"
	"fmt"

	"gorm.io/gorm"
)

// Cat 猫娘核心结构体
// 包含猫娘的所有运行时状态和功能组件
type Cat struct {
	Db             *gorm.DB        // 数据库连接，用于持久化数据
	CatState       *sytem.CatState // 当前状态(感情、场景、行为等)
	CatChat        *catai.Chat     // 主聊天处理器，处理用户对话
	CatChatExtract *catai.Chat     // 回答提取处理器，用于生成规则实体
	IsEntities     bool            // 转化为规则
	IsComfyUI      bool            // 绘制图片
}

// NewCat 创建并初始化一个新的猫娘实例
// 返回值: *Cat 初始化完成的猫娘指针
func NewCat(db *gorm.DB) *Cat {
	// 自动迁移数据库表结构
	db.AutoMigrate(&sytem.CatEntities{})
	db.AutoMigrate(&sytem.CatResponse{})
	// 初始化猫娘实例
	cat := &Cat{
		Db: db,
		CatState: &sytem.CatState{
			Self: map[string]map[string]string{
				"服装": {
					"上衣": "", "下衣": "",
					"内衣": "", "内裤": "",
					"胸罩": "", "鞋子": "",
					"袜子": "", "装饰物": ""},
				"身体": {
					"头": "", "脸": "", "眼": "", "鼻": "", "嘴": "", "舌头": "", "耳朵": "", "头发": "", "胳膊": "", "手": "",
					"胸": "", "腰": "", "肚子": "", "尾巴": "", "小穴": "", "屁股": "", "子宫": "", "阴道": "", "尿道": "", "膀胱": "", "肠道": "", "胃": "",
					"输卵管": "", "卵巢": "", "肛门": "", "大腿": "", "小腿": "", "脚": "", "后背": "", "锁骨": "", "脉搏": "", "腋窝": "", "阴唇": "", "宫颈口(子宫口)": "", "生命状态": ""},
				"数值状态": {"好感度": "?/100", "饱食度": "?/100", "体力": "?/100", "敏捷": "?/100", "智力": "?/100", "魅力": "?/100", "性欲": "?/100"},
			},
			Message: &catai.Message{MessageData: catai.MessageData{Role: "assistant"}},
		}}
	// 初始化主聊天处理器
	cat.CatChat = catai.NewChat("")
	// 设置默认猫娘提示词
	config := &Config{Value: sytem.Cat}
	db.Where("Key = ?", "SystemCat").Take(config)
	cat.CatChat.ChatSystem(config.Value) // 默认猫娘提示词
	cat.CatState.System = config.Value
	cat.CatChat.Messages = append(cat.CatChat.Messages, cat.CatState.Message) // 绑定猫娘状态
	data, _ := json.Marshal(cat.CatState)
	cat.CatState.Message.MessageData.Content = string(data)
	// 初始化提取处理器，用于将回答转化为事件
	cat.CatChatExtract = catai.NewChat("") // 将猫娘回答转化为事件
	config = &Config{Value: sytem.Extract}
	db.Where("Key = ?", "SystemExtract").Take(config)
	cat.CatChatExtract.ChatSystem(config.Value)
	cat.CatChatExtract.IsSingleInvocation = true
	config = &Config{}
	db.Where("Key = ?", "IsEntities").Take(config)
	if config.Value == "true" {
		cat.IsEntities = true
	}
	config = &Config{}
	db.Where("Key = ?", "IsComfyUI").Take(config)
	if config.Value == "true" {
		cat.IsComfyUI = true
	}
	return cat
}

// Chat 猫娘核心聊天方法
// 处理用户输入并生成响应，同时更新猫娘状态和数据库记录
//
// 参数:
//
//	str string - 用户输入的聊天内容
//	cr *CatResponse - 用于存储响应数据的结构体指针
//
// 返回值:
//
//	*catai.ChatReturn - 包含AI生成的消息内容
//	error - 错误对象，处理成功时为nil
//
// 处理流程:
//  1. 调用主聊天处理器处理用户输入
//  2. 解析响应并更新数据库
//  3. 更新猫娘内部状态
//  4. 使用提取处理器生成规则实体
func (cat *Cat) Chat(str string, cr *sytem.CatResponse) (*catai.ChatReturn, error) {
	// 调用聊天处理器处理用户输入
	return cat.CatChat.ChatUserCall(str, func(md catai.MessageData) catai.MessageData {
		// 解析聊天结果到响应对象
		cr.ID = 0
		err := md.Json(cr, cat.CatState.Message.Content)
		if err == nil {
			// 将响应记录到数据库
			cat.Db.Create(&cr)
			// 根据响应更新猫娘状态
			for key, v := range cr.Self {
				for k, s := range v {
					if cat.CatState.Self[key] == nil {
						cat.CatState.Self[key] = map[string]string{}
					}
					cat.CatState.Self[key][k] = s
				}
			}
			cat.CatState.Feeling = cr.Feeling
			cat.CatState.Scene = cr.Scene
			cat.CatState.Behavior = cr.Behavior
			// 更新猫娘状态消息
			data, _ := json.Marshal(cat.CatState)
			cat.CatState.Message.MessageData.Content = string(data)
			cat.CatState.Message.Data, _ = json.Marshal(cat.CatState.Message.MessageData)
			md.Content = cr.Answer
			// 格式化回答内容，添加选项编号
			for i, v := range cr.NextOptions {
				md.Content += fmt.Sprintf("\n %d.'%s'", i, v)
			}
		} else {
			fmt.Printf("[回答内容解析错误]:%s\n[错误信息]:%s", md.String(), err)
		}
		// 使用提取处理器将回答转化为规则实体
		if cat.IsEntities {
			cat.CatChatExtract.ChatUserCall(md.Content, func(md catai.MessageData) catai.MessageData {
				r := make([]sytem.CatEntities, 0)
				if md.Json(&r) == nil {
					for i := range r {
						r[i].CatEntitiesID = cr.ID
						cat.Db.Create(&r[i])
					}
				} else {
					fmt.Println(md.String())
				}
				return md
			})
		}
		return md
	})
}
