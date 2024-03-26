package styling

//スタイリング情報をDOMに追加（レンダリングツリーを作る）

import (
	"fmt"

	"github.com/WilhelmWeber/go_browser/src/cssparser"
	"github.com/WilhelmWeber/go_browser/src/htmlparser"
)

type StyledNode struct {
	Type       htmlparser.NodeType
	Children   []*StyledNode
	Elem       htmlparser.Element
	Text       string
	Properties map[string]string
}

func Matches(sels []cssparser.Selector, node *htmlparser.Node) bool {
	for _, sel := range sels {
		switch sel.Type {
		case cssparser.Univ:
			return true
		case cssparser.TagType:
			if sel.TagName == node.Elem.TagName {
				return true
			} else {
				continue
			}
		case cssparser.Class:
			if sel.ClassName == node.Elem.Attributes["class"] {
				return true
			} else {
				continue
			}
		case cssparser.Attr:
			if sel.TagName == node.Elem.TagName && sel.AttrValue == node.Elem.Attributes[sel.AttrName] {
				return true
			} else {
				continue
			}
		}
	}
	//何も合致しなかったら返す
	return false
}

func ToStyledNode(node *htmlparser.Node, stylesheet cssparser.Stylesheet) (*StyledNode, error) {
	properites := make(map[string]string)
	var children []*StyledNode

	for _, rule := range stylesheet.Rules {
		if Matches(rule.Selectors, node) {
			for _, dec := range rule.Declarations {
				properites[dec.Name] = dec.Value
			}
		}
	}

	//displayが設定されてなかったら初期値でblockを格納する
	if _, ok := properites["display"]; !ok && node.Type == htmlparser.Elem {
		properites["display"] = "block"
	}
	//エラーを返すことでDOMに含めないようにする
	if properites["display"] == "none" {
		return &StyledNode{}, fmt.Errorf("none")
	}

	//ChildNodeの処理
	for _, child := range node.Children {
		styled, err := ToStyledNode(child, stylesheet)
		if err != nil {
			//errが返ってくるということはdisplayがnoneなので含めない
			continue
		}
		children = append(children, styled)
	}
	stylednode := &StyledNode{Type: node.Type, Children: children, Elem: node.Elem, Properties: properites}
	return stylednode, nil
}
