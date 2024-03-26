package htmlparser

import (
	"regexp"
	"strings"
)

type NodeType int

const (
	Elem NodeType = iota
	Tex
	Nil //DOM探査用のtype(DOM探査してなかったら返す)
)

type Node struct {
	Type     NodeType
	Children []*Node
	Elem     Element
}

type Element struct {
	TagName    string
	Attributes map[string]string
	Text       string
}

type Parser struct {
	Tags  []string
	Index int
}

// HTMLパースメイン関数
func MainParser(html string) *Node {
	var sents []string
	_sents := TagParser(html)
	//微調整
	for _, e := range _sents {
		sent := strings.TrimSpace(e)
		if sent == "" {
			continue
		}
		sents = append(sents, sent)
	}
	parser := Parser{Tags: sents, Index: 0}
	dom := parser.NodeParser()
	return dom
}

func TagParser(html string) []string {
	var bodies []string
	rhtml := []rune(html)
	i := 0
	for i < len(rhtml) {
		if rhtml[i] == '<' {
			var tag []rune
			tag = append(tag, rhtml[i])
			for rhtml[i] != '>' {
				i++
				tag = append(tag, rhtml[i])
			}
			bodies = append(bodies, string(tag))
			i++
			continue
		} else {
			var text []rune
			for rhtml[i] != '<' {
				text = append(text, rhtml[i])
				i++
			}
			bodies = append(bodies, string(text))
		}
	}
	return bodies
}

func OpenTagName(tag string) (string, map[string]string) {
	tr := []rune(tag)
	var tagname []rune
	attribution := make(map[string]string)
	i := 0
	if tr[i] == '<' {
		i++
	} else {
		panic("Expected: <")
	}
	for tr[i] != ' ' && tr[i] != '>' {
		tagname = append(tagname, tr[i])
		i++
	}
	for tr[i] != '>' {
		if isSpace(tr[i]) {
			i++
			continue
		}
		attrName, attrValue, next_i := Attribution(tr, i)
		attribution[attrName] = attrValue
		i = next_i
		i++
	}
	return string(tagname), attribution
}

func Attribution(tr []rune, i int) (string, string, int) {
	var attrName []rune
	var attrValue []rune
	for {
		if isLetter(tr[i]) {
			attrName = append(attrName, tr[i])
			i++
			continue
		} else {
			break
		}
	}
	for {
		if isSpace(tr[i]) {
			i++
			continue
		} else {
			break
		}
	}
	if tr[i] == '=' {
		i++
	} else {
		panic("Expected: =")
	}
	for {
		if isSpace(tr[i]) {
			i++
			continue
		} else {
			break
		}
	}
	if tr[i] == '"' {
		i++
	} else {
		panic("Expected: \"")
	}
	for {
		if tr[i] == '"' {
			break
		}
		attrValue = append(attrValue, tr[i])
		i++
	}
	return string(attrName), string(attrValue), i
}

func (parser *Parser) NodeParser() *Node {
	otp := regexp.MustCompile(`<(".*?"|'.*?'|[^'"])*?>`) //開始タグ

	if otp.MatchString(parser.Tags[parser.Index]) {
		tagname, attribution := OpenTagName(parser.Tags[parser.Index])
		ctp := regexp.MustCompile(`^</` + tagname + `\s*>$`)
		parser.Index++
		var nodes []*Node
		for !ctp.MatchString(parser.Tags[parser.Index]) {
			//終了タグになるまでnodeを呼び出す
			node := parser.NodeParser()
			nodes = append(nodes, node)
			parser.Index++
		}
		elem := Element{TagName: tagname, Attributes: attribution}
		return &Node{Type: Elem, Children: nodes, Elem: elem}
	} else {
		//テキストであったらそのままノードを返す
		elm := Element{Text: parser.Tags[parser.Index]}
		return &Node{Type: Tex, Elem: elm}
	}
}

func isLetter(x rune) bool {
	return ('a' <= x && x <= 'z') || ('A' <= x && x <= 'Z') || (x == '-')
}

func isSpace(x rune) bool {
	return x == ' '
}
