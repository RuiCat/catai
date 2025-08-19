package main

import (
	"bytes"
	"catai"
	"catai/sytem"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"gorm.io/gorm"
)

// Prompt 提示词
type Prompt struct {
	Prompt         string            `json:"prompt"`
	NegativePrompt string            `json:"negative_prompt"`
	Cat            *catai.Chat       `json:"-"`
	Prompts        map[string]any    `json:"-"` // 提示词信息
	ServerAddress  string            `json:"-"` // 服务器地址
	ClientID       string            `json:"-"` // 客户端ID
	PromptKey      map[string]string `json:"-"` // 提示词输入位置
}

func NewComfyUI(db *gorm.DB) *Prompt {
	config := &Config{Value: sytem.Comfyui}
	db.Where("Key = ?", "ComfyUI").Take(config)
	prompt := &Prompt{
		Cat: catai.NewChat(""),
	}
	prompt.Cat.IsSingleInvocation = true
	prompt.Cat.Messages = []*catai.Message{
		{MessageData: catai.MessageData{Role: "user", Content: "要求扩充提示词,性行为描与动作,环境与人物等描写细腻,根据聊天内容将对话场景生成为ComfyUI正向提示词,要求严格满足相关语法与英文格式."}},
		{MessageData: catai.MessageData{Role: "user", Content: config.Value}},
		{MessageData: catai.MessageData{Role: "user", Content: ""}},
		{MessageData: catai.MessageData{Role: "user", Content: ""}},
	}
	// 打开连接
	if prompt.ServerAddress == "" {
		prompt.ServerAddress = "127.0.0.1:8188"
	}
	if prompt.ClientID == "" {
		prompt.ClientID = uuid.New().String()
	}
	config = &Config{}
	db.Where("Key = ?", "Prompt").Take(config)
	if err := json.Unmarshal([]byte(config.Value), &prompt.Prompts); err != nil {
		panic(err)
	}
	config = &Config{}
	db.Where("Key = ?", "PromptKey").Take(config)
	if err := json.Unmarshal([]byte(config.Value), &prompt.PromptKey); err != nil {
		panic(err)
	}
	return prompt
}

// Chat 解析
func (prompt *Prompt) Chat(cat *Cat, str string) (images []string, err error) {
	if !cat.IsComfyUI {
		return nil, nil
	}
	state, _ := json.Marshal(cat.CatState)
	prompt.Cat.Messages[2].Content = fmt.Sprintf("历史提示词\n```\nPrompt:%s\nNegativePrompt:%s\n```\n", prompt.Prompt, prompt.NegativePrompt)
	prompt.Cat.Messages[3].Content = fmt.Sprintf("人设信息:\n'''\n%s\n'''\n人物状态:\n'''\n%s\n'''\n对话内容:\n'''\n%s\n'''\n", cat.CatState.System, string(state), str)
	// 获取数据
	r, err := prompt.Cat.Chat(nil)
	if err != nil {
		return nil, err
	}
	if err = r.Choices[0].Message.Json(prompt); err != nil {
		return nil, err
	}
	// 获取绘制图片
	images, err = prompt.GetImages()
	if err != nil {
		cat.IsComfyUI = false
	}
	return images, err
}

type PromptResponse struct {
	PromptID string `json:"prompt_id"`
}

func (prompt *Prompt) QueuePrompt() (string, error) {
	if prompt.Prompt == "" {
		return "", nil
	}
	prompt.Prompts[prompt.PromptKey["prompt"]].(map[string]any)["inputs"].(map[string]any)["text"] = prompt.Prompt
	prompt.Prompts[prompt.PromptKey["negative_prompt"]].(map[string]any)["inputs"].(map[string]any)["text"] = prompt.NegativePrompt
	requestBody := map[string]any{
		"prompt":    prompt.Prompts,
		"client_id": prompt.ClientID,
	}
	jsonData, err := json.Marshal(requestBody)
	if err != nil {
		return "", err
	}
	resp, err := http.Post(
		fmt.Sprintf("http://%s/prompt", prompt.ServerAddress),
		"application/json",
		bytes.NewBuffer(jsonData),
	)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	var promptResp PromptResponse
	if err := json.Unmarshal(body, &promptResp); err != nil {
		return "", err
	}
	return promptResp.PromptID, nil
}

func (prompt *Prompt) GetImages() (images []string, err error) {
	// 打开连接
	var conn *websocket.Conn
	if conn, _, err = websocket.DefaultDialer.Dial(
		fmt.Sprintf("ws://%s/ws?clientId=%s", prompt.ServerAddress, prompt.ClientID),
		nil,
	); err != nil {
		return nil, err
	}
	defer conn.Close()
	// 传递工作流
	if _, err = prompt.QueuePrompt(); err != nil {
		return nil, err
	}
	// 处理回复
	var msg map[string]any
	for {
		_, message, err := conn.ReadMessage()
		if err != nil {
			return images, err
		}
		if err := json.Unmarshal(message, &msg); err != nil {
			return images, err
		}
		switch msg["type"] {
		case "executed":
			if output, ok := msg["data"].(map[string]any); ok {
				if img, ok := output["output"].(map[string]any); ok {
					for _, image := range img["images"].([]any) {
						images = append(images, fmt.Sprintf("http://%s/api/view?filename=%s&type=output&subfolder=", prompt.ServerAddress, image.(map[string]any)["filename"].(string)))
					}
					return images, err
				}
			}
		}
	}
}
