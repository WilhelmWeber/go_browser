package main

import (
	"github.com/WilhelmWeber/go_browser/src/cssparser"
	"github.com/WilhelmWeber/go_browser/src/filereader"
	"github.com/WilhelmWeber/go_browser/src/htmlparser"
	"github.com/WilhelmWeber/go_browser/src/layout"
	"github.com/WilhelmWeber/go_browser/src/styling"
	"github.com/rivo/tview"
)

func main() {
	html := filereader.Reader("public/test.html")
	css := filereader.Reader("public/test.css")

	dom := htmlparser.MainParser(html)
	cssom := cssparser.MainParser(css)

	rendering_tree, _ := styling.ToStyledNode(dom, cssom)
	layoutBox := layout.ToLayoutBox(rendering_tree)
	painting := layout.Painting(layoutBox)

	app := tview.NewApplication()
	if err := app.SetRoot(painting, true).Run(); err != nil {
		panic(err)
	}
}
