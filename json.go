package catai

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"reflect"
	"strings"

	"github.com/kaptinlin/jsonrepair"
)

// MessageData 表示单条对话消息的数据结构
// Role: 消息角色，如"user"或"assistant"
// Content: 消息内容
type MessageData struct {
	Role    string `json:"role"`    // 消息发送者角色，如"user"(用户)或"assistant"(AI助手)
	Content string `json:"content"` // 消息的具体内容
}

// String 将MessageData格式化为字符串输出
// 返回值: 格式化后的字符串，格式为"[Role:    user]: 消息内容"
func (data *MessageData) String() (str string) {
	return fmt.Sprintf("[Role:%10s]: %s", data.Role, data.Content)
}

// Json 将MessageData中的JSON内容解析到指定的结构体r中
// 参数r: 目标结构体指针，用于存储解析后的JSON数据
// 返回值: 错误信息，如果解析成功则为nil
// 功能: 自动修复JSON格式错误并解析内容
func (data *MessageData) Json(r any, v ...any) (err error) {
	jsonData := GetBetweenStr(data.Content, "```json", "```")
	// 解析
	repaired, err := jsonrepair.JSONRepair(jsonData)
	if err == nil {
		err = json.NewDecoder(bytes.NewBufferString(repaired)).Decode(r)
		if err == nil {
			return nil
		}
	}
	// 连续修复
	typ := ReflectFormat(r)
	// 处理错误
	jsonChat := NewChat("")
	jsonChat.IsSingleInvocation = true
	jsonChat.Messages = []*Message{
		{MessageData: MessageData{Role: "system", Content: "根据历史数据修正新内容的json格式错误同时去除格式中表情与符号."}},
		{MessageData: MessageData{Role: "system", Content: fmt.Sprintf("历史参考数据:\n'''\n%v\n'''", v)}},
		{MessageData: MessageData{Role: "system", Content: fmt.Sprintf("要求返回格式:\n'''\n%s\n'''", typ)}},
	}
	for i := 0; i < 5; i++ {
		// 发送修复请求
		r, er := jsonChat.ChatUser(fmt.Sprintf("回答内容:\n'''\n%s\n'''\n错误信息:\n'''\n%s\n'''\n", jsonData, err))
		if er != nil {
			// 重新尝试
			repaired, er = jsonrepair.JSONRepair(GetBetweenStr(r.Choices[0].Message.Content, "```json", "```"))
			if er == nil {
				er = json.NewDecoder(bytes.NewBufferString(repaired)).Decode(r)
				// 修复完成
				if er == nil {
					return er
				}
			} else {
				err = errors.Join(err, er)
			}
		} else {
			err = errors.Join(err, er)
		}
	}
	return err
}

// GetBetweenStr 从字符串str中提取位于start和end之间的子字符串
// 参数str: 源字符串
// 参数start: 起始标记字符串
// 参数end: 结束标记字符串
// 返回值: 提取的子字符串，如果找不到标记则返回空字符串
func GetBetweenStr(str, start, end string) string {
	n := strings.Index(str, start)
	if n == -1 {
		n = 0
	} else {
		n = n + len(start) // 增加了else，不加的会把start带上
	}
	str = string([]byte(str)[n:])
	m := strings.Index(str, end)
	if m == -1 {
		m = len(str)
	}
	str = string([]byte(str)[:m])
	return str
}

// ReflectFormat 通过反射获取类型的格式化字符串表示
// 参数s: 任意类型的值
// 返回值: 类型的格式化字符串表示
// 支持的类型: 结构体、切片、指针、map和接口
func ReflectFormat(s any) string {
	var v reflect.Value
	var ok bool
	if v, ok = s.(reflect.Value); !ok {
		v = reflect.ValueOf(s)
	}
	t := v.Type()
	switch t.Kind() {
	case reflect.Struct:
		var result strings.Builder
		result.WriteString("{\n")
		for i := 0; i < v.NumField(); i++ {
			field := t.Field(i).Tag.Get("json")
			if field != "" && field != "-" {
				result.WriteString(fmt.Sprintf("  %s: %v\n", field, ReflectFormat(v.Field(i))))
			}
		}
		result.WriteString("}")
		return result.String()
	case reflect.Slice:
		return fmt.Sprintf("[]%s", ReflectFormat(reflect.New(t.Elem())))
	case reflect.Pointer:
		return ReflectFormat(v.Elem())
	case reflect.Map:
		return fmt.Sprintf("map[%s]%s\n", ReflectFormat(reflect.New(t.Key())), ReflectFormat(reflect.New(t.Elem())))
	case reflect.Interface:
		return "any"
	}
	return t.Kind().String()
}
