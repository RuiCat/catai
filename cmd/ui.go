package main

import (
	"image"
	"image/color"

	"gioui.org/io/event"
	"gioui.org/io/pointer"
	"gioui.org/layout"
	"gioui.org/op/clip"
	"gioui.org/op/paint"
	"gioui.org/unit"
	"gioui.org/widget"
	"gioui.org/widget/material"
)

// Editor 输入框实现
type Editor struct {
	Rect  int // 圆角
	Width int // 线框
	Theme *material.Theme
	// 按钮
	Button      *widget.Clickable
	ButtonStyle material.ButtonStyle
	// 输入框
	Editor      widget.Editor
	EditorStyle material.EditorStyle
	// 模型返回
	CatEditor      widget.Editor
	CatEditorStyle material.EditorStyle
	// 属性
	AttributeEditor      widget.Editor
	AttributeEditorStyle material.EditorStyle
	// 显示
	ImgOp paint.ImageOp
	// 显示
	DefaultColor color.NRGBA // 默认状态颜色
	ActiveColor  color.NRGBA // 点燃状态颜色
	IsActive     bool
	// 回调
	Call func(str string)
}

func (editor *Editor) New(th *material.Theme, hint string) {
	editor.Theme = th
	editor.EditorStyle = material.Editor(th, &editor.Editor, hint)
	editor.CatEditor.ReadOnly = true
	editor.CatEditor.SetText("内部测试软件")
	editor.CatEditorStyle = material.Editor(th, &editor.CatEditor, hint)
	editor.AttributeEditor.ReadOnly = true
	editor.AttributeEditor.SetText("-")
	editor.ImgOp = paint.NewImageOp(image.NewAlpha(image.Rect(0, 0, 255, 255)))
	editor.AttributeEditorStyle = material.Editor(th, &editor.AttributeEditor, hint)
	editor.Button = new(widget.Clickable)
	editor.ButtonStyle = material.Button(editor.Theme, editor.Button, "发送")
}

func (editor *Editor) Layout(gtx layout.Context) layout.Dimensions {
	// 边框
	paint.FillShape(gtx.Ops, editor.DefaultColor, clip.Stroke{
		Path:  clip.RRect{Rect: image.Rect(0, 0, gtx.Constraints.Min.X, gtx.Constraints.Min.Y), SE: editor.Rect, SW: editor.Rect, NW: editor.Rect, NE: editor.Rect}.Path(gtx.Ops),
		Width: 1,
	}.Op(),
	)
	// 布局
	return layout.Flex{Axis: layout.Vertical}.Layout(gtx, layout.Flexed(0.9, func(gtx layout.Context) layout.Dimensions {
		return layout.Flex{Axis: layout.Horizontal}.Layout(gtx,
			layout.Flexed(0.7,
				func(gtx layout.Context) layout.Dimensions {
					// 图片
					return layout.Flex{
						Axis:      layout.Vertical,
						Alignment: layout.Middle,
						Spacing:   layout.SpaceEnd,
					}.Layout(gtx,
						layout.Rigid(func(gtx layout.Context) layout.Dimensions {
							gtx.Constraints.Min.X = 1
							return layout.UniformInset(50).Layout(gtx, widget.Image{Src: editor.ImgOp}.Layout)
						}),
						layout.Flexed(0.3, func(gtx layout.Context) layout.Dimensions {
							return layout.UniformInset(unit.Dp(editor.Rect)).Layout(gtx, editor.CatEditorStyle.Layout)
						}),
					)
				}),
			layout.Flexed(0.3, func(gtx layout.Context) layout.Dimensions {
				return layout.UniformInset(unit.Dp(editor.Rect)).Layout(gtx, func(gtx layout.Context) layout.Dimensions {
					return widget.Border{Color: editor.DefaultColor, CornerRadius: unit.Dp(editor.Rect), Width: unit.Dp(editor.Width)}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
						return layout.UniformInset(unit.Dp(editor.Rect)).Layout(gtx, func(gtx layout.Context) layout.Dimensions {
							return editor.AttributeEditorStyle.Layout(gtx)
						})
					})
				})
			}),
		)
	}),
		layout.Flexed(0.1, func(gtx layout.Context) layout.Dimensions {
			return layout.UniformInset(unit.Dp(editor.Rect)).Layout(gtx, func(gtx layout.Context) layout.Dimensions {
				return layout.Flex{}.Layout(gtx, layout.Flexed(0.9, func(gtx layout.Context) layout.Dimensions {
					// 事件
					rect := clip.Rect(image.Rectangle{Max: gtx.Constraints.Min})
					childArea := rect.Push(gtx.Ops)
					defer childArea.Pop()
					if ev, ok := gtx.Event(pointer.Filter{
						Target: editor,
						Kinds:  pointer.Enter | pointer.Leave,
					}); ok {
						if x, ok := ev.(pointer.Event); ok {
							switch x.Kind {
							case pointer.Enter:
								editor.IsActive = true
							case pointer.Leave:
								editor.IsActive = false
							}
						}
					}
					event.Op(gtx.Ops, editor)
					// 输入框处理
					c := editor.DefaultColor
					if editor.IsActive {
						c = editor.ActiveColor
					}
					return widget.Border{Color: c, CornerRadius: unit.Dp(editor.Rect), Width: unit.Dp(editor.Width)}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
						return layout.UniformInset(unit.Dp(editor.Rect)).Layout(gtx, editor.EditorStyle.Layout)
					})
				}), layout.Rigid(func(gtx layout.Context) layout.Dimensions {
					// 发送按钮
					gtx.Constraints.Min.X = 60
					if editor.Button.Clicked(gtx) && editor.Call != nil {
						text := editor.Editor.Text()
						if text != "" {
							go editor.Call(text)
							editor.CatEditor.SetText("猫娘思考中....")
							editor.Editor.SetText("")
						}
					}
					return layout.UniformInset(unit.Dp(editor.Rect)).Layout(gtx, editor.ButtonStyle.Layout)
				}),
				)
			})
		}),
	)
}
