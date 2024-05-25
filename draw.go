/*
@author: sk
@date: 2024/5/25
*/
package main

import (
	"image/color"
	"math"

	"github.com/fogleman/gg"
	"golang.org/x/image/colornames"
)

type IDraw interface {
	Draw(ctx *gg.Context)
	Measure(width float64)
	GetWidth() float64
	GetHeight() float64
	Layout(x, y float64)
}

type BaseDraw struct {
	X, Y          float64
	Width, Height float64
	Draws         []IDraw
}

func (b *BaseDraw) GetWidth() float64 {
	return b.Width
}

func (b *BaseDraw) GetHeight() float64 {
	return b.Height
}

func TransDraw(ctx *gg.Context, node *Node) IDraw {
	if node.Type == NodeText {
		return NewTextDraw(ctx, node.Content)
	} else {
		margin := GetNum(node.Styles, DeclarationMargin)
		color0 := GetClr(node.Styles, DeclarationColor)
		display := GetStr(node.Styles, DeclarationDisplay)
		padding := GetNum(node.Styles, DeclarationPadding)
		width := GetNum(node.Styles, DeclarationWidth)
		height := GetNum(node.Styles, DeclarationHeight)
		return NewRectDraw(ctx, margin, padding, width, height, color0, display, node.Nodes)
	}
}

func GetStr(styles map[DeclarationType]*Declaration, type0 DeclarationType) string {
	if res, ok := styles[type0]; ok {
		return res.Str
	}
	return DisplayBlock
}

func GetClr(styles map[DeclarationType]*Declaration, type0 DeclarationType) color.Color {
	if res, ok := styles[type0]; ok {
		return res.Clr
	}
	return colornames.White
}

func GetNum(styles map[DeclarationType]*Declaration, type0 DeclarationType) float64 {
	if res, ok := styles[type0]; ok {
		return res.Num
	}
	return 0
}

//=================文本不接受任何属性，字体大小颜色都是固定的inline设置==================

type TextDraw struct {
	*BaseDraw
	Content string
}

func (t *TextDraw) Layout(x, y float64) {
	t.X, t.Y = x, y // 自己布局一下就行了，没有其他孩子需要布局
}

func (t *TextDraw) Measure(width float64) {
	// 什么都不用做，文本的宽度,高度本来就是固定的
}

func NewTextDraw(ctx *gg.Context, content string) *TextDraw {
	w, h := ctx.MeasureString(content)
	return &TextDraw{Content: content, BaseDraw: &BaseDraw{Width: w, Height: h}}
}

func (t *TextDraw) Draw(ctx *gg.Context) {
	ctx.SetColor(colornames.Black)
	ctx.DrawStringAnchored(t.Content, t.X, t.Y, 0, 1)
}

//======================矩形有各种属性用于绘制=======================

type RectDraw struct {
	*BaseDraw
	Margin, Padding float64
	Color           color.Color
	Display         string
}

func (r *RectDraw) Layout(x, y float64) {
	r.X, r.Y = x, y
	width := r.Width - r.Padding*2 - r.Margin*2
	x += r.Padding + r.Margin
	y += r.Padding + r.Margin
	currW := float64(0)
	maxH := float64(0)
	for _, draw := range r.Draws {
		if currW+draw.GetWidth() > width {
			x = r.X + r.Padding + r.Margin
			y += maxH
			currW = draw.GetWidth()
			maxH = draw.GetHeight()
		} else {
			currW += draw.GetWidth()
			maxH = math.Max(maxH, draw.GetHeight())
		}
		draw.Layout(x, y)
		x += draw.GetWidth()
	}
}

func (r *RectDraw) Measure(width float64) {
	// 先测量宽度
	width -= r.Padding*2 + r.Margin*2 // 孩子能用的部分
	if r.Display == DisplayBlock {
		r.Width = width + r.Padding*2 + r.Margin*2 // 这时设置宽度无效，固定占用全部的
		for _, draw := range r.Draws {
			draw.Measure(width)
		}
	} else { // inline 模式
		if r.Width > 0 { // 指定了宽度直接使用
			r.Width += r.Margin * 2 // 用户指定的宽度是不算外边距的
			width = r.Width - r.Padding*2 - r.Margin*2
			for _, draw := range r.Draws {
				draw.Measure(width)
			}
		} else { // 否则需要进行计算
			currW := float64(0)
			maxW := float64(0)
			for _, draw := range r.Draws {
				draw.Measure(width)
				if currW+draw.GetWidth() > width {
					maxW = math.Max(maxW, currW)
					currW = draw.GetWidth()
				} else {
					currW += draw.GetWidth()
				}
			}
			r.Width = math.Max(maxW, currW) + r.Padding*2 + r.Margin*2 // 宽度由子类控制
		}
	}
	// 计算高度
	if r.Height > 0 {
		r.Height = r.Height + r.Margin*2 // 用户指定的高度是不算外边距的
		return                           // 指定了高度，无需计算
	}
	width = r.Width - r.Padding*2 - r.Margin*2
	currH, currW := float64(0), float64(0)
	maxH := float64(0)
	for _, draw := range r.Draws {
		if currW+draw.GetWidth() > width {
			currH += maxH
			currW = draw.GetWidth()
			maxH = draw.GetHeight()
		} else {
			currW += draw.GetWidth()
			maxH = math.Max(maxH, draw.GetHeight())
		}
	}
	r.Height = currH + maxH + r.Padding*2 + r.Margin*2
}

func NewRectDraw(ctx *gg.Context, margin, padding, width, height float64, color0 color.Color, display string, nodes []*Node) *RectDraw {
	draws := make([]IDraw, 0)
	for _, node := range nodes {
		draws = append(draws, TransDraw(ctx, node))
	}
	return &RectDraw{Margin: margin, Padding: padding, Color: color0, Display: display,
		BaseDraw: &BaseDraw{Width: width, Height: height, Draws: draws}}
}

func (r *RectDraw) Draw(ctx *gg.Context) {
	ctx.SetColor(r.Color) // 先绘制自己
	ctx.DrawRectangle(r.X+r.Margin, r.Y+r.Margin, r.Width-r.Margin*2, r.Height-r.Margin*2)
	ctx.Fill()
	for _, draw := range r.Draws { // 再绘制儿子
		draw.Draw(ctx)
	}
}
