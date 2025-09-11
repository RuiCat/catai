package tools

import (
	"catai/api/chat"
	"errors"
	"fmt"
)

// TaskList 任务列表
type TaskList struct {
	*chat.Messages `json:"-"`           // 底层绑定
	Name           string               `json:"-"` // 名称
	Mess           *chat.Data           `json:"-"` // 绑定对话
	Tools          map[string]chat.Call `json:"-"` // 工具列表
}

// NewTaskList 创建任务
func NewTaskList(name string, mes *chat.Messages) (task *TaskList) {
	task = &TaskList{
		Messages: mes,
		Name:     name,
		Mess:     &chat.Data{Role: "system", Enable: true},
		Tools:    map[string]chat.Call{},
	}
	mes.BindData(task.Mess)
	mes.AddTool(chat.Function{
		Name: task.Name,
		Parameters: []chat.Parameter{
			{Type: "object", Properties: map[string]chat.Location{
				"工具名称": chat.NewLocation("string", fmt.Sprintf("对 %s 块内部工具的名称", task.Name)),
			}, Required: []string{"工具名称"}},
			{Type: "object", Properties: map[string]chat.Location{
				"操作类型": chat.NewLocation("string", "操作分为: 启用工具,停用工具"),
			}, Required: []string{"操作类型"}},
		},
		Description: fmt.Sprintf("对 %s 块内部工具的管理,可以通过此工具将指定的工具禁用或者启用", task.Name),
		Call: chat.Call{
			Call:       task.Call,
			CallUpdate: task.CallUpdate,
		},
	})
	task.Update()
	return task
}

// AddTool 添加工具
func (task *TaskList) AddTool(function chat.Function) error {
	task.Tools[function.Name] = function.Call
	function.Call.Call = task.Call
	function.Call.CallUpdate = task.CallUpdate
	return task.Messages.AddTool(function)
}

// Call 回调
func (task *TaskList) Call(mes *chat.Messages, tool *chat.Tool, args map[string]any) {
	if tool.GetName() == task.Name {

	} else if call := task.Tools[tool.GetName()].Call; call != nil {
		call(mes, tool, args)
	}
}

// Update 更新
func (task *TaskList) Update() (err error) {
	defer func() {
		if r := recover(); r != nil {
			err = errors.Join(err, fmt.Errorf("%v", r))
		}
	}()
	task.Mess.Content = fmt.Sprintf("[ 工具列表: %s ]\n", task.Name)
	var s, b string
	for str := range task.Tools {
		tool := task.Messages.Tools[task.Messages.ToolsMap[str]]
		c := fmt.Sprintf("	工具: %s\n	工具说明: %s\n", str, tool.Function.Description)
		if tool.Enable {
			s += c
		} else {
			b += c
		}
	}
	task.Mess.Content += fmt.Sprintf("[已启用工具列表]\n%s\n[已禁用工具列表]\n%s\n", s, b)
	return task.Mess.Update()
}

// CallUpdate 回调
func (task *TaskList) CallUpdate(tool *chat.Tool, ret *chat.ChatRet) {
	if tool.GetName() == task.Name {
		task.Update()
	} else if call := task.Tools[tool.GetName()].CallUpdate; call != nil {
		call(tool, ret)
	}
}
