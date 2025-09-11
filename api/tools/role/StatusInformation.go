package role

import "catai/api/tools"

// StatusInformation 状态信息
func StatusInformation() *tools.Information {
	tool := tools.NewInformation("状态信息")
	tool.Value["背包信息"] = map[string]string{"": "实体/素材名称:相关的说明 #每一条信息都是自身背包或者身上的素材或者实体"}
	tool.Value["角色状态"] = map[string]string{
		"":      "状态:角色当前状态 #在每一次与用户对话中更新相关内容",
		"能量值":   "int",
		"饱食度":   "int",
		"清洁度":   "int",
		"魅力值":   "int",
		"信任度":   "int",
		"经验值":   "int",
		"性反应":   "int",
		"饥饿度":   "int",
		"口渴度":   "int",
		"情绪强度":  "int",
		"排泄需求":  "string",
		"当前情绪":  "string",
		"意识水平":  "string",
		"认知能力":  "string",
		"精神症状":  "string",
		"情绪稳定性": "string",
	}
	tool.Value["身体状态"] = map[string]string{
		"心脏":   "string",
		"血管":   "string",
		"肺部":   "string",
		"气道":   "string",
		"肾脏":   "string",
		"膀胱":   "string",
		"尿道":   "string",
		"子宫":   "string",
		"卵巢":   "string",
		"体温":   "string",
		"疲劳度":  "string",
		"疼痛指数": "string",
		"意识水平": "string",
		"上消化道": "string",
		"下消化道": "string",
		"附属器官": "string",
		"外生殖器": "string",
		"内生殖器": "string",
		"妊娠状态": "string",
		"腺体功能": "string",
		"激素水平": "string",
		"中枢神经": "string",
		"周围神经": "string",
		"运动能力": "string",
		"感官功能": "string",
		"代谢状态": "string",
	}
	tool.Value["社会环境状态"] = map[string]string{
		"社会角色": "string",
		"经济状况": "string",
		"文化背景": "string",
		"环境适应": "string",
	}
	tool.Value["交互状态"] = map[string]string{
		"": `器官:状态 # 器官名称:{"插入状态": map[string:器官位置]any:{"插入对象": "string","插入深度": "string","组织反应": "string"},"器官变形":map[string:对象]any,"体液分布":map[string:对象]any}`,
	}
	return tool
}
