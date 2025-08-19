// main 包包含猫娘聊天机器人的主程序
package main

import (
	"fmt"
	"image"
	"image/color"
	"net/http"
	"time"

	_ "image/png"

	"gioui.org/app"
	"gioui.org/io/event"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/op/paint"
	"gioui.org/widget/material"
)

func main() {
	// 创建app
	window := new(app.Window)
	// 创建猫娘实例
	db := Init()
	cat := NewCat(db)
	// 提示词
	prompt := NewComfyUI(db)
	// 默认显示
	theme := material.NewTheme()
	editor := &Editor{
		Rect:         5,
		Width:        1,
		DefaultColor: color.NRGBA{102, 139, 139, 255},
		ActiveColor:  color.NRGBA{150, 205, 205, 255},
	}
	editor.New(theme, "用户输入")
	callChat := make(chan string)
	defer close(callChat)
	// 处理聊天
	editor.Call = func(str string) {
		callChat <- str
	}
	go func() {
		for str := range callChat {
			// 创建响应对象
			r := &CatResponse{}
			// 调用猫娘聊天方法处理输入
			l, err := cat.Chat(str, r)
			if l == nil || len(l.Choices) == 0 {
				panic("无法读取数据")
			}
			if err != nil {
				fmt.Println("错误:", err)
				cat.CatChat.Messages = cat.CatChat.Messages[:len(cat.CatChat.Messages)-1]
				r.Answer = l.Choices[0].Message.Content
			} else {
				// 记录当前回答
				editor.AttributeEditor.SetText(cat.CatState.String())
				// 回答
				fmt.Println(r.NextOptions)
			}
			// 处理对话
			editor.CatEditor.SetText(r.Answer)
			// 处理绘画
			images, err := prompt.Chat(cat, fmt.Sprintf("提问:\n```\n%s\n```\n回答:\n```\n%s\n```\n", str, l.Choices[0].Message.Content))
			if err != nil {
				fmt.Println("错误:", err)
			}
			if len(images) > 0 {
				res, err := http.Get(images[0])
				if err != nil {
					fmt.Println(err)
				}
				img, _, err := image.Decode(res.Body)
				if err != nil {
					fmt.Println(err)
				}
				editor.ImgOp = paint.NewImageOp(img)
				res.Body.Close()
			}
		}
	}()
	// 周期
	Ticker := time.NewTicker(time.Millisecond * 500)
	// 事件处理
	events := make(chan event.Event)
	acks := make(chan struct{})
	defer close(acks)
	go func() {
		// 创建界面
		for range acks {
			events <- window.Event()
		}
	}()
	// 创建界面
	var ops op.Ops
	acks <- struct{}{}
	for {
		select {
		case eve := <-events:
			switch e := eve.(type) {
			case app.DestroyEvent:
				fmt.Println(e.Err)
				return
			case app.FrameEvent:
				gtx := app.NewContext(&ops, e)
				layout.UniformInset(10).Layout(gtx, func(gtx layout.Context) layout.Dimensions {
					return editor.Layout(gtx)
				})
				e.Frame(gtx.Ops)
			}
			acks <- struct{}{}
		case <-Ticker.C:
			window.Invalidate()
		}
	}
}
