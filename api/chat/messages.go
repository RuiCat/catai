package chat

import (
	"bytes"
	"catai/api/buffer"
	"catai/api/config"
	"encoding/json"
	"fmt"
	"io"
)

// Messages 上下文
type Messages struct {
	Buffer   buffer.Buffer  `json:"-"`
	Model    string         `json:"-"`
	User     *Data          `json:"-"`
	Tools    []*Tool        `json:"-"`
	Datas    []*Data        `json:"-"`
	System   *Data          `json:"-"`
	IsJson   bool           `json:"-"`
	ToolsMap map[string]int `json:"-"`
}

// NewMessages 创建
func NewMessages() *Messages {
	return &Messages{
		Model:    config.Get("Model", "qwen-plus"),
		Buffer:   make(buffer.Buffer, 0),
		Tools:    make([]*Tool, 0),
		ToolsMap: map[string]int{},
		Datas:    make([]*Data, 0),
		User:     &Data{Role: "user"},
		System:   &Data{Role: "system"},
	}
}

// Data 上下文数据
type Data struct {
	data    []byte `json:"-"`       // 缓存数据
	Role    string `json:"role"`    // 消息发送者角色
	Content string `json:"content"` // 消息的具体内容
}

// Update 更新缓冲
func (data *Data) Update() error {
	buffer := bytes.NewBuffer(data.data[:0])
	err := json.NewEncoder(buffer).Encode(data)
	if err != nil {
		return err
	}
	data.data = buffer.Bytes()
	return err
}

// AddData 添加上下文数据
func (mes *Messages) AddMessage(role string, content string) error {
	switch role {
	case "user", "assistant":
		data := &Data{Role: role, Content: content}
		if err := data.Update(); err != nil {
			return err
		}
		mes.Datas = append(mes.Datas, data)
	case "system":
		mes.System.Content = content
		if err := mes.System.Update(); err != nil {
			return err
		}
	}
	return nil
}

// AddDialogue 添加对话
func (mes *Messages) AddDialogue(assistant string, user string) error {
	if err := mes.AddMessage("assistant", assistant); err != nil {
		return err
	}
	if err := mes.AddMessage("user", user); err != nil {
		mes.Datas = mes.Datas[:len(mes.Datas)-1]
		return err
	}
	return nil
}

// SetUser 设置提问不添加到上下文
func (mes *Messages) SetUser(content string) error {
	mes.User.Content = content
	return mes.User.Update()
}

// AddTool 添加工具
func (mes *Messages) AddTool(function Function) error {
	if _, ok := mes.ToolsMap[function.Name]; ok {
		return fmt.Errorf("重复工具定义: %s", function.Name)
	}
	tool := &Tool{Type: "function", Enable: true, Function: function}
	if err := tool.Update(); err != nil {
		return err
	}
	mes.ToolsMap[function.Name] = len(mes.Tools)
	mes.Tools = append(mes.Tools, tool)
	return nil
}

// BindData 绑定数据
func (mes *Messages) BindData(data *Data) {
	if data != nil {
		mes.Datas = append(mes.Datas, data)
	}
}

// Get 默认读取
func (mes *Messages) Get() io.Reader {
	// 更新缓冲区
	if len(mes.Buffer) > 0 {
		mes.Buffer = mes.Buffer[:0]
	}
	// 初始化
	mes.Buffer.AddDyte([]byte(`{"model": "`))
	mes.Buffer.AddString(mes.Model)
	mes.Buffer.AddDyte([]byte(`","messages":[`))
	// 构建上下文
	mes.Buffer.AddPtr(&mes.System.data)
	mes.Buffer.AddDyte([]byte(`,`))
	for i, n := 0, len(mes.Datas)-1; ; i++ {
		mes.Buffer.AddPtr(&mes.Datas[i].data)
		if i < n {
			mes.Buffer.AddDyte([]byte(`,`))
		} else {
			break
		}
	}
	if mes.User.Content != "" {
		mes.User.Content = ""
		mes.Buffer.AddDyte([]byte(`,`))
		mes.Buffer.AddPtr(&mes.User.data)
	}
	mes.Buffer.AddDyte([]byte(`]`))
	// 工具列表
	if n := len(mes.Tools) - 1; n >= 0 {
		mes.Buffer.AddDyte([]byte(`,"tools":[`))
		for i := 0; ; i++ {
			if !mes.Tools[i].Enable {
				continue
			}
			mes.Buffer.AddPtr(&mes.Tools[i].data)
			if i < n {
				mes.Buffer.AddDyte([]byte(`,`))
			} else {
				mes.Buffer.AddDyte([]byte(`]`))
				break
			}
		}
		mes.Buffer.AddDyte([]byte(`,"parallel_tool_calls":true`))
	}
	if mes.IsJson {
		mes.Buffer.AddDyte([]byte(`,"response_format":{"type": "json_object"}`))
	}
	mes.Buffer.AddDyte([]byte(`}`))
	return mes.Buffer.Get()
}

// GetMessages 获取自身
func (mes *Messages) GetMessages() *Messages { return mes }
