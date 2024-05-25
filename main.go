/*
@author: sk
@date: 2024/5/16
*/
package main

import "github.com/fogleman/gg"

// https://limpet.net/mbrubeck/2014/08/08/toy-layout-engine-1.html

func main() {
	tokens := ParseToken("/Users/bytedance/Documents/go/my_browser/res/test.html")
	node := ParseNode(tokens)
	style := ParseCss(node)
	ApplyCss(node, style)
	ctx := gg.NewContext(1280, 720)
	err := ctx.LoadFontFace("/Users/bytedance/Documents/go/my_browser/res/fusion-pixel-12px-monospaced-zh_hans.ttf", 36)
	HandleErr(err)
	draw := TransDraw(ctx, node)
	draw.Measure(1280)
	draw.Layout(0, 0)
	draw.Draw(ctx)
	err = ctx.SavePNG("/Users/bytedance/Documents/go/my_browser/res/test.png")
	HandleErr(err)
}
