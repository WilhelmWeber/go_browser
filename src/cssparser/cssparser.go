package cssparser

import (
	"strings"
)

type SelectorType int
type AttrSelectorOp int

// セレクターの種類のenum
const (
	Univ    SelectorType = iota //*
	TagType                     //htmlタグの名称で指定
	Attr                        //attributionで指定
	Class                       //class名での指定
)

const (
	None    AttrSelectorOp = iota //未設定の時は0で代入されるので
	Eq                            //=
	Contain                       //~=
)

// stylesheetは複数のルールからなる
type Stylesheet struct {
	Rules []Rule
}

type Rule struct {
	Selectors    []Selector
	Declarations []Declaration
}

type Selector struct {
	Type      SelectorType   //セレクターの種類
	TagName   string         //タグ名
	AttrName  string         //属性名
	Op        AttrSelectorOp //属性指定のときの演算子
	AttrValue string         //属性値
	ClassName string         //クラス値
}

type Declaration struct {
	Name  string
	Value string
}

// 疑似クラス化
type CssParser struct {
	Css   []rune
	Index int
}

// CSSパースメイン関数
func MainParser(css string) Stylesheet {
	parser := CssParser{Css: []rune(css), Index: 0}
	stylesheet := parser.Parser()
	return stylesheet
}

func (parser *CssParser) Parser() Stylesheet {
	var rules []Rule
	for parser.Index < len(parser.Css) {
		if isSpace(parser.Css[parser.Index]) {
			parser.Index++
			continue
		}
		rule := parser.RuleParser()
		rules = append(rules, rule)
		//}まで読み込みが終了した時点でruleが返ってくるので一つ進める
		parser.Index++
	}
	return Stylesheet{Rules: rules}
}

// CssParser本元から呼び出される
// 文字列開始までの半角スペースの判定についてはCssParserで行っているためすぐに文字列が来るものとして処理を開始する
func (parser *CssParser) RuleParser() Rule {
	selectors := parser.SlcsParser()
	//次の文字列が来るまで半角スペースを飛ばす
	for isSpace(parser.Css[parser.Index]) {
		parser.Index++
	}
	parser.Index++
	declarations := parser.DecsParser()
	return Rule{Selectors: selectors, Declarations: declarations}
}

func (parser *CssParser) DecParser() Declaration {
	var name []rune
	var val []rune
	for parser.Css[parser.Index] != ':' {
		if isSpace(parser.Css[parser.Index]) {
			parser.Index++
			continue
		}
		name = append(name, parser.Css[parser.Index])
		parser.Index++
	}
	for isSpace(parser.Css[parser.Index]) {
		parser.Index++
	}
	for parser.Css[parser.Index] != ';' {
		val = append(val, parser.Css[parser.Index])
		parser.Index++
	}
	return Declaration{Name: string(name), Value: string(val)}
}

func (parser *CssParser) DecsParser() []Declaration {
	var declarations []Declaration
	for parser.Css[parser.Index] != '}' {
		if isSpace(parser.Css[parser.Index]) {
			parser.Index++
			continue
		}
		declaration := parser.DecParser()
		declarations = append(declarations, declaration)
		parser.Index++
	}
	return declarations
}

func (parser *CssParser) SlcsParser() []Selector {
	var sentences []rune
	var selectors []Selector

	for parser.Css[parser.Index] != '{' {
		sentences = append(sentences, parser.Css[parser.Index])
		parser.Index++
	}
	ss := strings.Split(string(sentences), ",")
	for _, s := range ss {
		selector := SlcParser(s)
		selectors = append(selectors, selector)
	}
	return selectors
}

/*
方針
'{'が出てくるまで文字列をそのままSlcsParserが文字列を抜き出す
,があればカンマで分割し、その分だけSlcParser()を呼び出す
後ろのスペースだけを取り除く
後はパターンマッチングで適切なselectorを返す
*/
func SlcParser(s string) Selector {
	selSent := strings.TrimSpace(s)
	var selector Selector

	if selSent == "*" {
		selector = Selector{Type: Univ}
	} else if selRune := []rune(selSent); selRune[0] == '.' {
		selector = Selector{Type: Class, ClassName: string(selRune[1:])}
	} else if strings.Contains(selSent, "[") {
		var tagname []rune
		var attrname []rune
		var attrvalue []rune
		var op AttrSelectorOp
		i := 0
		for selRune[i] != '[' {
			if isSpace(selRune[i]) {
				i++
				continue
			}
			tagname = append(tagname, selRune[i])
			i++
		}
		i++
		for selRune[i] != '=' && selRune[i] != '~' {
			if isSpace(selRune[i]) {
				i++
				continue
			}
			attrname = append(attrname, selRune[i])
			i++
		}
		if selRune[i] == '=' {
			op = Eq
			i++
		} else {
			if selRune[i+1] != '=' {
				panic("Expected: =")
			}
			op = Contain
			i += 2
		}
		for selRune[i] != '"' {
			i++
		}
		i++
		for selRune[i] != '"' {
			attrvalue = append(attrvalue, selRune[i])
			i++
		}
		selector = Selector{Type: Attr, TagName: string(tagname), AttrName: string(attrname), Op: op, AttrValue: string(attrvalue)}
	} else {
		//さもなければtag名指定だろう
		selector = Selector{Type: TagType, TagName: selSent}
	}
	return selector
}

func isSpace(x rune) bool {
	return x == ' '
}
