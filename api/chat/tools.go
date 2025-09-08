package chat

import (
	"bytes"
	"encoding/json"
)

// Tool 工具定义
type Tool struct {
	data     []byte   `json:"-"` // 缓存数据
	Type     string   `json:"type"`
	Function Function `json:"function"`
}

// GetName gg工具名称
func (tool *Tool) GetName() string {
	return tool.Function.Name
}

// Update 更新缓冲
func (tool *Tool) Update() error {
	buffer := bytes.NewBuffer(tool.data)
	err := json.NewEncoder(buffer).Encode(tool)
	if err != nil {
		return err
	}
	tool.data = buffer.Bytes()
	return err
}

// Function 函数信息
type Function struct {
	Name        string                                               `json:"name"`        // 工具名
	Parameters  []FunctionParameter                                  `json:"parameters"`  // 工具参数
	Description string                                               `json:"description"` // 工具描述
	Call        func(mes *Messages, tool *Tool, args map[string]any) `json:"-"`           // 回调
	CallUpdate  func(ret *ChatRet)                                   `json:"-"`           // 回调
}

// FunctionParameter 函数参数
type FunctionParameter struct {
	Type       string            `json:"type"`       // 类型
	Properties map[string]string `json:"properties"` // 参数描述
	Required   []string          `json:"required"`   // 必填项
}
