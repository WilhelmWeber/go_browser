package layout

import (
	"github.com/WilhelmWeber/go_browser/src/htmlparser"
	"github.com/WilhelmWeber/go_browser/src/styling"
	"github.com/rivo/tview"
)

type BoxType int

const (
	BlockBox BoxType = iota
	InlineBox
	AnonymousBox
	NilBox //lastTaker関数でchildrenが空であったときに帰ってくるBoxのType
)

type LayoutBox struct {
	Type       BoxType
	Children   []LayoutBox
	Elm        htmlparser.Element
	NodeType   htmlparser.NodeType
	Properties map[string]string
}

func ToLayoutBox(snode *styling.StyledNode) LayoutBox {
	var boxtype BoxType
	nodetype := snode.Type
	properties := snode.Properties
	elm := snode.Elem
	var children []LayoutBox

	switch properties["display"] {
	case "block":
		boxtype = BlockBox
	case "inline":
		boxtype = InlineBox
	}

	for _, child := range snode.Children {
		switch child.Properties["display"] {
		case "block":
			//ブロックであればそのままpush
			child_box := ToLayoutBox(child)
			children = append(children, child_box)
		case "inline":
			switch lastTaker(children).Type {
			case AnonymousBox:
				break
			default:
				children = append(children, LayoutBox{Type: AnonymousBox})
			}
			last := lastTaker(children)
			children[len(children)-1].Children = append(last.Children, ToLayoutBox(child))
		}
	}

	return LayoutBox{Type: boxtype, Children: children, NodeType: nodetype, Properties: properties, Elm: elm}
}

func Painting(box LayoutBox) *tview.Flex {
	switch box.Type {
	case BlockBox, InlineBox:
		switch box.NodeType {
		case htmlparser.Elem:
			flex := tview.NewFlex()
			flex.SetTitle(box.Elm.TagName)
			flex.SetBorder(true)
			flex.SetDirection(tview.FlexRow)
			for _, child := range box.Children {
				childflex := Painting(child)
				flex.AddItem(childflex, 0, 1, false)
			}
			return flex
		case htmlparser.Tex:
			flex := tview.NewFlex()
			flex.SetBorder(true)
			text := tview.NewTextView()
			text.SetText(box.Elm.Text)
			flex.AddItem(text, 0, 1, false)
			return flex
		}
	case AnonymousBox:
		flex := tview.NewFlex()
		flex.SetBorder(true)
		flex.SetDirection(tview.FlexColumn)
		for _, child := range box.Children {
			childflex := Painting(child)
			flex.AddItem(childflex, 0, 1, false)
		}
		return flex
	}
	panic("Unreachable")
}

func lastTaker(boxes []LayoutBox) LayoutBox {
	if len(boxes) == 0 {
		return LayoutBox{Type: NilBox}
	} else {
		return boxes[len(boxes)-1]
	}
}
