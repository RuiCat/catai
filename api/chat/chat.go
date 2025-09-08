package chat

import (
	"catai/api/config"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

// ChatRet 对话结果
type ChatRet map[string]any

// Get 得到数据
func (ret ChatRet) Get(Key ...any) (r any, ok bool) {
	r = ret
	for _, key := range Key {
		switch m := r.(type) {
		case ChatRet:
			if r, ok = m[key.(string)]; !ok {
				return nil, false
			}
		case map[string]any:
			if r, ok = m[key.(string)]; !ok {
				return nil, false
			}
		case []any:
			if i, ok := key.(int); !ok {
				return nil, false
			} else {
				r = m[i]
			}
		default:
			return nil, false
		}
	}
	return r, true
}

// MessagesFace 上下文接口
type MessagesFace interface {
	Get() io.Reader
	AddTool(function Function) error
	BindData(data *Data)
	GetMessages() *Messages
	SetUser(content string) error
}

// ChatPost 发送数据
func ChatPost(mes MessagesFace) (ret ChatRet, _ error) {
	req, err := http.NewRequest("POST", config.Get("URL", ""), mes.Get())
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", config.Get("Key", ""))
	res, err := (&http.Client{}).Do(req)
	if err != nil {
		return nil, err
	}
	err = json.NewDecoder(res.Body).Decode(&ret)
	// 处理工具
	mess := mes.GetMessages()
	if calls, ok := ret.Get("choices", 0, "message", "tool_calls"); ok {
		for _, a := range calls.([]any) {
			call := (ChatRet)(a.(map[string]any))
			name, ok := call.Get("function", "name")
			if !ok {
				fmt.Println("???", a)
				break
			}
			arguments, ok := call.Get("function", "arguments")
			if !ok {
				fmt.Println("???", a)
				break
			}
			val := map[string]any{}
			if json.Unmarshal([]byte(arguments.(string)), &val) != nil {
				fmt.Println("???!!!", a)
				break
			}
			if i, ok := mess.ToolsMap[name.(string)]; ok && mess.Tools[i].Function.Call != nil {
				mess.Tools[i].Function.Call(mess, mess.Tools[i], val)
			}
		}
	}
	for _, tool := range mess.Tools {
		if tool.Function.CallUpdate != nil {
			tool.Function.CallUpdate(&ret)
		}
	}
	return ret, err
}

// ChatUser 发送数据
func ChatUser(mes MessagesFace, uesr string) (ret ChatRet, _ error) {
	mes.SetUser(uesr)
	return ChatPost(mes)
}
