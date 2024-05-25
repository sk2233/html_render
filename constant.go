/*
@author: sk
@date: 2024/5/16
*/
package main

import "fmt"

type NodeType string

const (
	NodeText    NodeType = "NodeText"
	NodeElement NodeType = "NodeElement"
)

type TokenType string

const (
	TokenStart TokenType = "TokenStart"
	TokenEnd   TokenType = "TokenEnd"
	TokenText  TokenType = "TokenText"
)

type SelectorType string

const (
	SelectorID    SelectorType = "SelectorID"
	SelectorTag   SelectorType = "SelectorTag"
	SelectorClass SelectorType = "SelectorClass"
	SelectorAll   SelectorType = "SelectorAll"
)

func (s SelectorType) Order() int {
	switch s {
	case SelectorID:
		return 99
	case SelectorClass:
		return 98
	case SelectorTag:
		return 97
	case SelectorAll:
		return 1
	default:
		panic(fmt.Sprintf("invalid selector type: %s", s))
	}
}

type DeclarationType string

const (
	// 默认 0
	DeclarationMargin DeclarationType = "margin"
	// 默认白色
	DeclarationColor DeclarationType = "color"
	// 默认 block
	DeclarationDisplay DeclarationType = "display"
	// 默认为 0
	DeclarationPadding DeclarationType = "padding"
	// 默认取其内部最小值
	DeclarationWidth DeclarationType = "width"
	// 默认取其内部最小值
	DeclarationHeight DeclarationType = "height"
)

const (
	// inline 在宽度无限的情况下堆到一行，若超出宽度向下 warp 若是单个宽度超出最大宽度直接报错
	DisplayInline = "inline"
	// block 搞成一行，若是当前已经存在部分元素了换行
	DisplayBlock = "block"
)
