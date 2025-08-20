package sytem

import (
	"catai"
	"database/sql/driver"
	_ "embed"
	"encoding/json"
	"fmt"
)

//go:embed extract.system
var Extract string

//go:embed cat.system
var Cat string

//go:embed 基础工作流.json
var Comfyui string

// CatState 猫娘状态结构体
// 记录猫娘在对话过程中的动态状态
type CatState struct {
	Self     Maps           `json:"自身状态"` // 猫娘属性字典(名字、年龄等)
	Feeling  string         `json:"当前感情"` // 当前情感状态(开心、生气等)
	Scene    string         `json:"当前场景"` // 当前所处场景(卧室、厨房等)
	Behavior []string       `json:"当前行为"` // 当前行为列表
	Message  *catai.Message `json:"-"`    // 当前消息对象(内部使用)
	System   string         `json:"-"`    // 人设信息
}

func (cat *CatState) String() (str string) {
	str = "当前属性:\n"
	for key, v := range cat.Self {
		str += " " + key + "\n"
		for k, s := range v {
			str += "  " + k + ":" + s + "\n"
		}
	}
	str += fmt.Sprintf("当前感情:\n%s\n", cat.Feeling)
	str += fmt.Sprintf("当前场景:\n%s\n当前行为:\n", cat.Scene)
	for _, v := range cat.Behavior {
		str += fmt.Sprintf(" %s\n", v)
	}
	return str
}

// Strings 自定义字符串切片类型
// 实现数据库扫描和值转换接口
type Strings []string

func (t *Strings) Scan(value interface{}) error {
	bytesValue, _ := value.([]byte)
	return json.Unmarshal(bytesValue, t)
}

func (t Strings) Value() (driver.Value, error) {
	return json.Marshal(t)
}

// Maps 自定义映射类型
// 实现数据库扫描和值转换接口
type Maps map[string]map[string]string

func (t *Maps) Scan(value interface{}) error {
	bytesValue, _ := value.([]byte)
	return json.Unmarshal(bytesValue, t)
}

func (t Maps) Value() (driver.Value, error) {
	return json.Marshal(t)
}

// CatEntities 对话实体信息
// CatEntities 表示对话中的实体信息
// 包含实体、关联、限制条件和执行动作
type CatEntities struct {
	ID            uint    `json:"-" gorm:"primaryKey;autoIncrement"`
	CatEntitiesID uint    `json:"-"`
	Entity        string  `json:"行为对象" gorm:"type:text"`
	Relations     Strings `json:"关联信息" gorm:"type:text"`
	Conditions    Strings `json:"限制信息" gorm:"type:text"`
	Actions       Strings `json:"执行内容" gorm:"type:text"`
}

// CatResponse 对话响应结构体
// 包含猫娘对用户输入的完整响应信息
type CatResponse struct {
	ID          uint          `json:"-" gorm:"primaryKey;autoIncrement"` // 数据库主键
	Answer      string        `json:"回答" gorm:"type:text"`               // 回答内容
	Self        Maps          `json:"自身状态" gorm:"type:text"`             // 自身属性更新
	Feeling     string        `json:"当前感情" gorm:"type:text"`             // 更新后的感情状态
	Scene       string        `json:"当前场景" gorm:"type:text"`             // 更新后的场景
	Behavior    Strings       `json:"当前行为" gorm:"type:text"`             // 当前行为列表
	NextOptions Strings       `json:"接下来可选行为" gorm:"type:text"`          // 用户下一步可选行为
	CatEntities []CatEntities `json:"关联" gorm:"type:text"`               // 关联的实体信息
}
