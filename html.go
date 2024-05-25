/*
@author: sk
@date: 2024/5/16
*/
package main

import (
	"bytes"
	"fmt"
	"os"
	"strings"
)

var (
	tokens0 []*Token
	index0  int
)

func ParseNode(tokens []*Token) *Node {
	tokens0 = tokens
	index0 = 0
	return parseNode()
}

func parseNode() *Node {
	token := readToken()
	if token.Type == TokenText {
		return NewTextNode(token.Raw)
	}
	if token.Type == TokenEnd {
		panic(fmt.Sprintf("invalid token %v", token))
	}
	tagName, attrs := parseTagAndAttr(token.Raw)
	nodes := make([]*Node, 0)
	for peekToken().Type != TokenEnd {
		nodes = append(nodes, parseNode())
	}
	token = readToken()
	if token.Raw != tagName {
		panic(fmt.Sprintf("invalid end tag %v tagName is %v", token, tagName))
	}
	return NewElementNode(tagName, attrs, nodes)
}

func peekToken() *Token {
	return tokens0[index0]
}

func parseTagAndAttr(raw string) (string, map[string]string) {
	items := strings.Split(raw, " ")
	attr := make(map[string]string)
	for i := 1; i < len(items); i++ {
		index := strings.IndexByte(items[i], '=')
		attr[items[i][:index]] = strings.Trim(items[i][index+1:], "\"")
	}
	return items[0], attr
}

func readToken() *Token {
	index0++
	return tokens0[index0-1]
}

func ParseToken(file string) []*Token {
	bs, err := os.ReadFile(file)
	HandleErr(err)
	res := make([]*Token, 0)
	index := 0
	for index < len(bs) {
		switch bs[index] {
		case ' ', '\t', '\n':
			// 忽略空格
			index++
		case '<':
			index++
			type0 := TokenStart
			if bs[index] == '/' {
				index++
				type0 = TokenEnd
			}
			buff := bytes.Buffer{}
			for bs[index] != '>' {
				buff.WriteByte(bs[index])
				index++
			}
			index++
			res = append(res, NewToken(type0, buff.String()))
		default:
			buff := bytes.Buffer{}
			for bs[index] != '<' {
				buff.WriteByte(bs[index])
				index++
			}
			res = append(res, NewToken(TokenText, strings.TrimSpace(buff.String())))
		}
	}
	return res
}

func ApplyCss(node *Node, style *StyleSheet) {
	if node.Type == NodeText {
		return
	}
	applyRules(node, style.Rules)
	for _, item := range node.Nodes {
		ApplyCss(item, style)
	}
}

func applyRules(node *Node, rules []*Rule) {
	for _, rule := range rules {
		order, styles := rule.Match(node)
		for _, style := range styles {
			if order > node.Orders[style.Type] {
				node.Orders[style.Type] = order
				node.Styles[style.Type] = style
			}
		}
	}
}
