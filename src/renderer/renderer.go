package renderer

import (
	"github.com/WilhelmWeber/go_browser/src/cssparser"
	"github.com/WilhelmWeber/go_browser/src/htmlparser"
	jsengine "github.com/WilhelmWeber/go_browser/src/js-engine"
	"github.com/WilhelmWeber/go_browser/src/layout"
	"github.com/WilhelmWeber/go_browser/src/styling"
	"github.com/rivo/tview"
)

type Renderer struct {
	View       *tview.Flex
	DocElem    *htmlparser.Node
	JsInstance *jsengine.JsRuntime
}

/*Rendere構造体の初期化*/
func New(html string, css string, js string) *Renderer {
	dom := htmlparser.MainParser(html)
	stylesheet := cssparser.MainParser(css)

	styled_node, _ := styling.ToStyledNode(dom, stylesheet)
	box := layout.ToLayoutBox(styled_node)
	view := layout.Painting(box)

	//JSruntimeの初期化とスクリプトのコンパイル
	jr := jsengine.New()
	jr.Compile(js)
	jr.DocRef = dom

	renderer := &Renderer{View: view, DocElem: dom, JsInstance: jr}
	return renderer
}

/*スクリプトの実行*/
func (r *Renderer) ExecuteScripts() {
	r.JsInstance.Execute()
}
