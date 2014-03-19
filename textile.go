package gotextile

import(
	"github.com/SoCloz/gotextile/parser"
)

func TextileToHtml(text string) (string, error) {
	doc, err := parser.ParseDocument(text)
	return doc.ToHtml(), err
}
