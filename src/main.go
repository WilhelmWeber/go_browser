package main

import (
	"fmt"

	"github.com/WilhelmWeber/go_browser/src/cssparser"
	"github.com/WilhelmWeber/go_browser/src/filereader"
	"github.com/WilhelmWeber/go_browser/src/htmlparser"
	"github.com/WilhelmWeber/go_browser/src/styling"
)

func main() {

	html := filereader.Reader("public/test.html")
	css := filereader.Reader("public/test.css")

	dom := htmlparser.MainParser(html)
	cssom := cssparser.MainParser(css)

	rendering_tree, _ := styling.ToStyledNode(dom, cssom)

	fmt.Println(rendering_tree.Children[0].Children[0])
}
