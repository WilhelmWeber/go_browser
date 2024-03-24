package main

import (
	"fmt"

	"github.com/WilhelmWeber/go_browser/src/cssparser"
	"github.com/WilhelmWeber/go_browser/src/filereader"
	"github.com/WilhelmWeber/go_browser/src/htmlparser"
)

func main() {

	html := filereader.Reader("public/test.html")
	css := filereader.Reader("public/test.css")

	dom := htmlparser.MainParser(html)
	cssom := cssparser.MainParser(css)

	fmt.Println(dom, cssom)
}
