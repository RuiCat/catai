package role

import "catai/api/chat"

// RolePlaying 角色扮演
type RolePlaying struct {
	Chat                chat.MessagesFace `json:"-"`                   // 绑定对话
	BodyInformation     string            `json:"BodyInformation"`     // 角色身份信息
	StatusInformation   string            `json:"StatusInformation"`   // 状态信息
	ShortTermGoal       string            `json:"ShortTermGoal"`       // 短期目标
	LongTermGoal        string            `json:"LongTermGoal"`        // 长期目标
	CurrentPlan         string            `json:"CurrentPlan"`         // 当前规划
	CurrentBehavior     string            `json:"CurrentBehavior"`     // 当前行为
	ShortTermMemory     string            `json:"ShortTermMemory"`     // 短期记忆
	LongTermMemoryIndex string            `json:"LongTermMemoryIndex"` // 长期记忆
}

// Bind 绑定
func (role *RolePlaying) Bind() {
	mes := role.Chat.GetMessages()
	mes.AddTool(chat.Function{})
}
