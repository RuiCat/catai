package sytem

import _ "embed"

//go:embed extract.system
var Extract string

//go:embed cat.system
var Cat string

//go:embed 基础工作流.json
var Comfyui string
