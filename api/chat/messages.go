package chat

import (
	"bytes"
	"catai/api/buffer"
	"catai/api/config"
	"encoding/json"
	"io"
)

// Messages 上下文
type Messages struct {
	Buffer buffer.Buffer
	Model  string
	Tools  []*Tool
	Datas  []*Data
	IsJson bool
	System *Data
}

// NewMessages 创建
func NewMessages() *Messages {
	return &Messages{
		Model:  config.Get("Model", "qwen-plus"),
		Buffer: make(buffer.Buffer, 0),
		Tools:  make([]*Tool, 0),
		Datas:  make([]*Data, 0),
		System: &Data{Role: "system"},
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
	buffer := bytes.NewBuffer(data.data)
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
	mes.Buffer = mes.Buffer[:0]
	return nil
}

// AddTool 添加工具
func (mes *Messages) AddTool(function Function) error {
	tool := &Tool{Type: "function", Function: function}
	if err := tool.Update(); err != nil {
		return err
	}
	mes.Tools = append(mes.Tools, tool)
	mes.Buffer = mes.Buffer[:0]
	return nil
}

// Get 默认读取
func (mes *Messages) Get() io.Reader {
	if len(mes.Buffer) > 0 {
		return mes.Buffer.Get()
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
			mes.Buffer.AddDyte([]byte(`]`))
			break
		}
	}
	// 工具列表
	if n := len(mes.Tools) - 1; n >= 0 {
		mes.Buffer.AddDyte([]byte(`,"tools":[`))
		for i := 0; ; i++ {
			mes.Buffer.AddPtr(&mes.Tools[i].data)
			if i < n {
				mes.Buffer.AddDyte([]byte(`,`))
			} else {
				mes.Buffer.AddDyte([]byte(`]`))
				break
			}
		}
	}
	if mes.IsJson {
		mes.Buffer.AddDyte([]byte(`,"response_format":{"type": "json_object"}`))
	}
	mes.Buffer.AddDyte([]byte(`}`))
	return mes.Buffer.Get()
}

// GetMessages 获取自身
func (mes *Messages) GetMessages() *Messages { return mes }
