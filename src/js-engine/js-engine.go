package jsengine

import (
	"github.com/WilhelmWeber/go_browser/src/htmlparser"
	"github.com/dop251/goja"
)

type JsRuntime struct {
	Runtime *goja.Runtime
	Program *goja.Program
	State   goja.Value       //保持される状態
	DocRef  *htmlparser.Node //DOM構造（DOM操作に使用するために保持）
}

/*js-engine初期化*/
func New() *JsRuntime {
	vm := goja.New()
	js := &JsRuntime{Runtime: vm}
	//TODO DOM操作関数を定義して、vm.Set()する
	vm.Set("SetInnerTextById", js.SetInnerTextById)
	return js
}

/*スクリプトの実行*/
func (js *JsRuntime) Execute() {
	val, _ := js.Runtime.RunProgram(js.Program)
	js.State = val
}

/*スクリプトのコンパイル*/
func (js *JsRuntime) Compile(script string) {
	program, _ := goja.Compile("program", script, false)
	js.Program = program
}

/*idでElementを探してinnerTextをセットする関数の定義*/
func (js *JsRuntime) SetInnerTextById(id string, text string) {
	dom := js.DocRef
	searched := SearchById(id, dom)
	if searched.Type == htmlparser.Nil {
		return
	}
	elm := htmlparser.Element{Text: text}
	searched.Children[0] = &htmlparser.Node{Type: htmlparser.Tex, Elem: elm}
}

/*idでDOM要素を探す*/
func SearchById(id string, dom *htmlparser.Node) *htmlparser.Node {
	var node *htmlparser.Node
	//自分のidが探しているidであったら自らを返す
	if dom.Elem.Attributes["id"] == id {
		node = dom
		return node
	}
	//子要素を探す
	for _, child := range dom.Children {
		searched := SearchById(id, child)
		//Nilでないものが返ってきたらそのnodeをバケツリレー式に返す
		if searched.Type != htmlparser.Nil {
			node = searched
			return node
		}
	}
	//それでもなかったらNilNodeを返す
	return &htmlparser.Node{Type: htmlparser.Nil}
}
