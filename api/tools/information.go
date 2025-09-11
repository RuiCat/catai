package tools

import (
	"catai/api/chat"
	"fmt"
)

// Information 信息工具
type Information struct {
	Name     string                       `json:"Name"`  // 信息状态
	Mess     *chat.Data                   `json:"-"`     // 绑定对话
	Value    map[string]map[string]string `json:"Value"` // 信息值
	IsUpdate bool                         `json:"-"`     // 需要更新
}

// NewInformation创建一个信息块
func NewInformation(name string) *Information {
	return &Information{
		Name:  name,
		Mess:  &chat.Data{Role: "system", Enable: true},
		Value: map[string]map[string]string{},
	}
}

// Bind 绑定数据
func (infor *Information) Bind(mess chat.MessagesFace) {
	// 绑定工具结果
	mess.BindData(infor.Mess)
	mess.AddTool(chat.Function{
		Name: infor.Name,
		Parameters: []chat.Parameter{
			{Type: "object", Properties: map[string]chat.Location{
				"操作": chat.NewLocation("string", fmt.Sprintf("对 %s 添加一个分类,选择一个 添加分类,删除分类,删除记录,修改记录 相关操作", infor.Name)),
			}, Required: []string{"操作"}},
			{Type: "object", Properties: map[string]chat.Location{
				"分类": chat.NewLocation("string", "需要修改值的 分类 值索引"),
			}, Required: []string{"分类"}},
			{Type: "object", Properties: map[string]chat.Location{
				"索引": chat.NewLocation("string", "需要修改值的 Key 值索引"),
			}, Required: []string{"索引"}},
			{Type: "object", Properties: map[string]chat.Location{
				"新值": chat.NewLocation("string", "如果为修改操作则传递本参数"),
			}, Required: []string{"新值"}},
		},
		Description: `严格按照如下格式进行工具调用:\n 添加分类: {"操作": "添加分类", "分类": "需要操作的分类名称", "索引": "", "新值": ""}" \n 删除分类: {"操作": "删除分类", "分类": "需要操作的分类名称", "索引": "", "新值": ""}"\n 删除记录: {"操作": "删除记录", "分类": "需要操作的分类名称", "索引": "分类下的记录", "新值": ""}"\n 修改记录: {"操作": "修改记录", "分类": "需要操作的分类名称", "索引": "分类下的记录", "新值": "设置的新值"}"`,
		Call: chat.Call{
			Call: func(mes *chat.Messages, tool *chat.Tool, val map[string]any) {
				defer func() {
					recover()
				}()
				fl, sy := val["分类"].(string), val["索引"].(string)
				switch val["操作"] {
				case "添加分类":
					if _, ok := infor.Value[fl]; ok {
						return
					}
					infor.Value[fl] = map[string]string{}
				case "删除分类":
					if _, ok := infor.Value[fl]; !ok {
						return
					}
					delete(infor.Value, fl)
				case "删除记录":
					if _, ok := infor.Value[fl]; !ok {
						return
					}
					if _, ok := infor.Value[fl]["索引"]; !ok {
						return
					}
					delete(infor.Value[fl], sy)
				case "修改记录":
					infor.Value[fl][sy] = val["新值"].(string)
				default:
					return
				}
				infor.IsUpdate = true
			},
			CallUpdate: func(tool *chat.Tool, ret *chat.ChatRet) {
				if infor.IsUpdate {
					infor.IsUpdate = false
					infor.Update()
				}
			},
		},
	})
	// 更新值
	infor.Update()
}

// Update 更新
func (infor *Information) Update() error {
	infor.Mess.Content = fmt.Sprintf("[ 数据块: %s ]\n", infor.Name)
	for key, v := range infor.Value {
		infor.Mess.Content += "分类: " + key + "\n"
		if _, ok := v[""]; ok {
			infor.Mess.Content += "  参考:\n	" + v[""] + "\n"
		}
		for k, s := range v {
			if k != "" {
				infor.Mess.Content += "  " + k + ":" + s + "\n"
			}
		}
	}
	return infor.Mess.Update()
}
