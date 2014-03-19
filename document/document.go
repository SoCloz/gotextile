package document

import(
	"fmt"
	"html"
	"strings"
)

type D struct {
	Tag string
	Text string
	Attr map[string]string
	Children []*D
}

func New(tag string) *D {
	return &D{Tag: tag, Children: make([]*D, 0, 128)}
}

func (d *D) AddChild(c *D) {
	d.Children = append(d.Children, c)
}

func (d *D) ToHtml() string {
	var text string
	if len(d.Children) > 0 {
		for _, c := range d.Children {
			text += c.ToHtml()
		}
	} else {
		text = html.EscapeString(d.Text)
	}
	if d.Tag == "" {
		return text
	} else {
		if text == "" {
			return fmt.Sprintf("<%s>", d.Tag)
		} else {
			format := "%s"
			list := strings.Split(d.Tag, "+")
			for _, tag := range list {
				format = fmt.Sprintf(format, fmt.Sprintf("<%s>%s</%s>", tag, "%s", tag))
			}
			return fmt.Sprintf(format, text)
		}
	}
}

func (d *D) ToText() string {
	return ""
}
