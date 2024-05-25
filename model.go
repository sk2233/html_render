/*
@author: sk
@date: 2024/5/16
*/
package main

import (
	"fmt"
	"image/color"
)

type Token struct {
	Type TokenType
	Raw  string
}

func NewToken(Type TokenType, raw string) *Token {
	return &Token{Type: Type, Raw: raw}
}

type Node struct {
	// 通用属性
	Type  NodeType
	Nodes []*Node
	// NodeText
	Content string
	// NodeElement
	TagName string
	Attrs   map[string]string
	// 相关属性与其优先级
	Styles map[DeclarationType]*Declaration
	Orders map[DeclarationType]int
}

func NewTextNode(content string) *Node {
	return &Node{Type: NodeText, Nodes: make([]*Node, 0), Content: content}
}

func NewElementNode(tagName string, attrs map[string]string, nodes []*Node) *Node {
	return &Node{Type: NodeElement, Nodes: nodes, TagName: tagName, Attrs: attrs,
		Styles: make(map[DeclarationType]*Declaration), Orders: make(map[DeclarationType]int)}
}

type Selector struct {
	Type    SelectorType // id tag class  同属性优先级 id > class > tag 相同属性间互相覆盖
	Content string       // 具体值
}

func (s *Selector) Match(node *Node) bool {
	switch s.Type {
	case SelectorID:
		return node.Attrs["id"] == s.Content
	case SelectorTag:
		return node.TagName == s.Content
	case SelectorClass:
		return node.Attrs["class"] == s.Content
	case SelectorAll:
		return true
	default:
		panic(fmt.Sprintf("unknown selector type: %v", s.Type))
	}
}

type Declaration struct {
	Type DeclarationType // 修改那些属性
	// 具体属性值，不同属性支持的值不同
	Str string
	Num float64
	Clr color.Color
}

type Rule struct {
	Selectors    []*Selector
	Declarations []*Declaration
}

func (r *Rule) Match(node *Node) (int, []*Declaration) {
	order := 0
	for _, item := range r.Selectors {
		if item.Type.Order() > order && item.Match(node) {
			order = item.Type.Order()
		}
	}
	if order > 0 {
		return order, r.Declarations
	}
	return order, make([]*Declaration, 0)
}

type StyleSheet struct {
	Rules []*Rule
}

func NewStyleSheet() *StyleSheet {
	return &StyleSheet{Rules: make([]*Rule, 0)}
}

//func NewStyleSheet(sheets ...StyleSheet) *StyleSheet {
//	res := &StyleSheet{Rules: make([]*Rule, 0)}
//	for _, sheet := range sheets {
//		res.Rules = append(res.Rules, sheet.Rules...)
//	}
//	return res
//}
