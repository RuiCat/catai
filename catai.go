package catai

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	_ "embed"
)

// 默认密匙
var Key = ""

// Message 表示完整的消息对象，包含原始数据和序列化后的字节
// MessageData: 消息内容
// Data: 序列化后的JSON数据，缓存用
type Message struct {
	MessageData
	Data []byte `json:"-"`
}

// MarshalJSON 编码
func (message *Message) MarshalJSON() ([]byte, error) {
	if len(message.Data) == 0 {
		message.Data, _ = json.Marshal(message.MessageData)
	}
	return message.Data, nil
}

// ChatReturn 表示API返回的完整响应结构
// Choices: 包含返回的消息选择
// Created: 时间戳
// Id: 请求ID
// Model: 使用的模型名称
// Object: 对象类型
// SystemFingerprint: 系统指纹
// Usage: token使用统计
type ChatReturn struct {
	Choices           []ChatChoice `json:"choices"`
	Created           int          `json:"created"`
	Id                string       `json:"id"`
	Model             string       `json:"model"`
	Object            string       `json:"object"`
	SystemFingerprint string       `json:"system_fingerprint"`
	Usage             ChatUsage    `json:"usage"`
}

// String 格式和输出
func (chat *ChatReturn) String() (str string) {
	return chat.Choices[0].Message.String()
}

// ChatChoice 表示API返回的单个消息选择
// FinishReason: 完成原因
// Index: 选择索引
// Message: 实际消息内容
type ChatChoice struct {
	FinishReason string      `json:"finish_reason"`
	Index        int         `json:"index"`
	Message      MessageData `json:"message"`
}

// ChatUsage 记录API调用的token使用情况
// CompletionTokens: 生成使用的token数
// PromptTokens: 提示使用的token数
// TotalTokens: 总token数
type ChatUsage struct {
	CompletionTokens int `json:"completion_tokens"`
	PromptTokens     int `json:"prompt_tokens"`
	TotalTokens      int `json:"total_tokens"`
}

// Chat 管理整个聊天会话
// Messages: 消息历史记录
// Buffer: 用于构建HTTP请求的缓冲区
// IsSingleInvocation: 设置为单次调用
type Chat struct {
	Messages           Messages
	Buffer             Buffer
	IsSingleInvocation bool
}

// NewChat 初始化
func NewChat(model string) (cat *Chat) {
	if model == "" {
		model = "deepseek-ai/DeepSeek-V3"
	}
	// 构建
	cat = &Chat{Messages: []*Message{}, Buffer: make(Buffer, 5)}
	cat.Buffer[0] = new([]byte)
	cat.Buffer[1] = new([]byte)
	cat.Buffer[2] = new([]byte)
	cat.Buffer[3] = new([]byte)
	cat.Buffer[4] = new([]byte)
	(*cat.Buffer[0]) = []byte(`{"model": "`)
	(*cat.Buffer[1]) = []byte(model)
	(*cat.Buffer[2]) = []byte(`","messages":`)
	(*cat.Buffer[3]) = make([]byte, 0, 255)
	(*cat.Buffer[4]) = []byte(`}`)
	return cat
}

// String 格式和输出
func (messages *Chat) String() (str string) {
	for _, message := range messages.Messages {
		str += fmt.Sprintln(message)
	}
	return str
}

// ChatSystem 添加系统设置
func (messages *Chat) ChatSystem(str string) {
	messages.Messages = append([]*Message{{MessageData: MessageData{Role: "system", Content: str}}}, messages.Messages...)
}

// ChatUser 处理用户输入并获取AI回复
// str: 用户输入的消息内容
// 返回值: API响应和可能的错误
func (messages *Chat) ChatUser(str string) (*ChatReturn, error) {
	return messages.ChatUserCall(str, nil)
}

// ChatUserCall 处理支持回调的用户输入
func (messages *Chat) ChatUserCall(str string, call func(MessageData) MessageData) (*ChatReturn, error) {
	if messages.IsSingleInvocation {
		mes := messages.Messages[len(messages.Messages)-1]
		if mes.Role == "user" {
			mes.Content = str
		} else {
			messages.Messages = append(messages.Messages, &Message{MessageData: MessageData{Role: "user", Content: str}})
		}
	} else {
		messages.Messages = append(messages.Messages, &Message{MessageData: MessageData{Role: "user", Content: str}})
	}
	return messages.Chat(call)
}

// Chat 执行实际的API调用并处理响应
// 返回值: API响应和可能的错误
func (messages *Chat) Chat(call func(MessageData) MessageData) (*ChatReturn, error) {
	// 写入上下文
	if err := messages.Buffer.ReadJson(3, messages.Messages); err != nil {
		return nil, err
	}
	// 开始通信

	// data, _ := io.ReadAll(messages.Buffer.Get())
	// fmt.Println(string(data))

	body, err := ChatPost(messages.Buffer.Get())
	if err != nil {
		return nil, err
	}
	// 处理获取的内容
	chat := &ChatReturn{}
	err = json.NewDecoder(body.Body).Decode(chat)
	if err != nil {
		return nil, err
	}
	if len(chat.Choices) == 0 {
		return nil, nil
	}
	// 处理输出
	if call != nil {
		chat.Choices[0].Message = call(chat.Choices[0].Message)
	}
	// 单次调用直接返回
	if messages.IsSingleInvocation {
		return chat, nil
	}
	// 记录上下文
	messages.Messages = append(messages.Messages, &Message{MessageData: chat.Choices[0].Message})
	return chat, nil
}

// ChatPost 发送数据
func ChatPost(body io.Reader) (*http.Response, error) {
	req, err := http.NewRequest("POST", "https://api.siliconflow.cn/v1/chat/completions", body)
	if err != nil {
		return nil, err
	}
	req.Header.Set("content-type", "application/json")
	req.Header.Set("authorization", Key)
	req.Header.Set("accept", "application/json")
	res, err := (&http.Client{}).Do(req)
	if err != nil {
		return nil, err
	}
	return res, err
}
