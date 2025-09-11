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
	return tool
}
