package role

import "catai/api/tools"

// BodyInformation 角色信息
func BodyInformation() *tools.Information {
	tool := tools.NewInformation("角色信息")
	tool.Value["基础信息"] = map[string]string{
		"姓名":   "string",
		"种族":   "string",
		"外表":   "string",
		"性格":   "string",
		"身世":   "string",
		"体型":   "string",
		"年龄":   "int",
		"身高":   "int",
		"体重":   "int",
		"体脂率":  "int",
		"骨架类型": "string",
	}
	tool.Value["外貌特征"] = map[string]string{
		"肤色":   "string",
		"肤质":   "string",
		"疤痕":   "string",
		"发型":   "string",
		"发色":   "string",
		"腋毛":   "string",
		"胡须":   "string",
		"脸型":   "string",
		"眼型":   "string",
		"鼻型":   "string",
		"唇型":   "string",
		"纹身":   "string",
		"穿孔":   "string",
		"化妆":   "string",
		"乳房发育": "string",
		"发量密度": "string",
		"阴毛分布": "string",
		"痣/胎记": "string",
		"皮肤病变": "string",
	}
	tool.Value["穿着状态"] = map[string]string{
		"贴身层": "string",
		"中间层": "string",
		"外层":  "string",
		"下装":  "string",
		"整洁度": "string",
		"功能性": "string",
		"适配度": "string",
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
	tool.Value["心理状态"] = map[string]string{
		"情绪强度":  "int",
		"当前情绪":  "string",
		"意识水平":  "string",
		"认知能力":  "string",
		"精神症状":  "string",
		"情绪稳定性": "string",
	}
	tool.Value["生理需求"] = map[string]string{
		"能量值":  "int",
		"饱食度":  "int",
		"清洁度":  "int",
		"魅力值":  "int",
		"信任度":  "int",
		"经验值":  "int",
		"性反应":  "int",
		"饥饿度":  "int",
		"口渴度":  "int",
		"排泄需求": "string",
	}
	tool.Value["社会环境状态"] = map[string]string{
		"社会角色": "string",
		"经济状况": "string",
		"文化背景": "string",
		"环境适应": "string",
	}
	tool.Value["交互状态"] = map[string]string{
		"状态名称": `参考结构: {"插入状态": map[string:器官位置]any:{"插入对象": "string","插入深度": "string","组织反应": "string"},"器官变形":map[string:对象]any,"体液分布":map[string:对象]any}`,
	}
	return tool
}
