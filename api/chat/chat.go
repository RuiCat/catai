package chat

import (
	"catai/api/config"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"reflect"
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
		case map[string][]any:
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
			fmt.Println("->", reflect.TypeOf(m).Kind())
			return nil, false
		}
	}
	return r, true
}

// MessagesFace 上下文接口
type MessagesFace interface {
	Get() io.Reader
	GetMessages() *Messages
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
	return ret, json.NewDecoder(res.Body).Decode(&ret)
}
