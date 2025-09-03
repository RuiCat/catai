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
	Name         string        `json:"name"`         // 函数名
	Args         []FunctionArg `json:"args"`         // 参数
	Return       []FunctionArg `json:"return"`       // 输出
	Introduction string        `json:"introduction"` // 介绍
}

// FunctionArg 函数参数
type FunctionArg struct {
	Name         string `json:"name"`         // 参数
	Type         string `json:"type"`         // 参数类型
	Introduction string `json:"introduction"` // 介绍
}
