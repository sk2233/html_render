/*
@author: sk
@date: 2024/5/20
*/
package main

import (
	"fmt"
	"regexp"
	"strings"
)

func ParseCss(root *Node) *StyleSheet {
	res := NewStyleSheet()
	parseCss(root, res)
	return res
}

func parseCss(node *Node, res *StyleSheet) bool {
	if node.Type == NodeText {
		return false
	}
	if node.TagName == "style" {
		if len(node.Nodes) == 1 && node.Nodes[0].Type == NodeText {
			res.Rules = append(res.Rules, parseRules(node.Nodes[0].Content)...)
		}
		return true
	}
	nodes := make([]*Node, 0)
	for _, item := range node.Nodes {
		if !parseCss(item, res) {
			nodes = append(nodes, item)
		}
	} // 解析并移出 style
	node.Nodes = nodes
	return false
}

var (
	ruleRe        = regexp.MustCompile("[{}]")
	declarationRe = regexp.MustCompile("[:;]")
)

func parseRules(content string) []*Rule {
	res := make([]*Rule, 0)
	items := ruleRe.Split(content, -1)
	for i := 1; i < len(items); i += 2 {
		res = append(res, &Rule{
			Selectors:    parseSelectors(items[i-1]),
			Declarations: parseDeclarations(items[i]),
		})
	}
	return res
}

func parseDeclarations(content string) []*Declaration {
	res := make([]*Declaration, 0)
	items := declarationRe.Split(content, -1)
	for i := 1; i < len(items); i += 2 {
		items[i-1] = strings.TrimSpace(items[i-1])
		items[i] = strings.TrimSpace(items[i])
		type0 := DeclarationType(items[i-1])
		switch type0 {
		case DeclarationMargin, DeclarationPadding, DeclarationWidth, DeclarationHeight:
			l := len(items[i])
			res = append(res, &Declaration{
				Type: type0,
				Num:  MustFloat(items[i][:l-2]),
			})
		case DeclarationColor:
			res = append(res, &Declaration{
				Type: DeclarationColor,
				Clr:  MustClr(items[i][1:]),
			})
		case DeclarationDisplay:
			res = append(res, &Declaration{
				Type: DeclarationDisplay,
				Str:  items[i],
			})
		default:
			panic(fmt.Errorf("unknown declaration type %q", items[i]))
		}
	}
	return res
}

func parseSelectors(content string) []*Selector {
	res := make([]*Selector, 0)
	items := strings.Split(content, ",")
	for _, item := range items {
		item = strings.TrimSpace(item)
		if item == "*" {
			res = append(res, &Selector{
				Type: SelectorAll,
			})
		} else if strings.HasPrefix(item, "#") {
			res = append(res, &Selector{
				Type:    SelectorID,
				Content: item[1:],
			})
		} else if strings.HasPrefix(item, ".") {
			res = append(res, &Selector{
				Type:    SelectorClass,
				Content: item[1:],
			})
		} else {
			res = append(res, &Selector{
				Type:    SelectorTag,
				Content: item,
			})
		}
	}
	return res
}
