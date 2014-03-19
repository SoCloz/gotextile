package parser

import(
	"fmt"
	"regexp"
	"strings"

	"github.com/SoCloz/gotextile/document"
)

type Parser struct {
	Re *regexp.Regexp
	Tag string
	Type string
	AllowedChildren []string
	OffsetEnd int
	Parse func(string, string) (*document.D, error)
	ParseAttributes func(string)
}

func re(str string) *regexp.Regexp {
	return regexp.MustCompile("(?Us)"+str)
}

var parsers []Parser

func init() {
	parsers = []Parser{
		Parser{
			Re: re("p. (.*)\n*$"),
			Tag: "p",
			Type: "block",
			AllowedChildren: []string{"phrase"},
		},
		Parser{
			Re: re("h1. (.*)\n*$"),
			Tag: "h1",
			Type: "block",
			AllowedChildren: []string{"phrase"},
		},
		Parser{
			Re: re("(\\*+ .*)\n*$"),
			Tag: "ul",
			Type: "list",
			Parse: parseList,
		},
		Parser{
			Re: re("(\\#+ .*)\n*$"),
			Tag: "ol",
			Type: "list",
			Parse: parseList,
		},
		Parser{
			Re: re("(- .* := .*)\n*$"),
			Tag: "dl",
			Type: "list",
			AllowedChildren: []string{"definition-term", "definition-definition"},
		},
		Parser{
			Re: re("- (.*) :="),
			Tag: "dt",
			Type: "definition-term",
			AllowedChildren: []string{"phrase"},
			OffsetEnd: -2,
		},
		Parser{
			Re: re(":= ([^\n]*)\n-"),
			Tag: "dd",
			Type: "definition-definition",
			AllowedChildren: []string{"phrase"},
			OffsetEnd: -1,
		},
		Parser{
			Re: re(":= ([^\n]*)\n$"),
			Tag: "dd",
			Type: "definition-definition",
			AllowedChildren: []string{"phrase"},
		},
		Parser{
			Re: re(":= (.*) =:\n"),
			Tag: "dd",
			Type: "definition-definition",
			AllowedChildren: []string{"phrase"},
		},
		Parser{
			Re: re("bc. (.*)\n*$"),
			Tag: "pre+code",
			Type: "block",
		},
		Parser{
			Re: re("pre. (.*)\n*$"),
			Tag: "pre",
			Type: "block",
		},
		Parser{
			Re: re("###. .*$"),
			Tag: "",
			Type: "block",
		},
		Parser{
			Re: re("(\\|.*\\|)\n*$"),
			Tag: "table",
			Type: "block",
			Parse: parseTable,
		},
		Parser{
			Re: re("(.*)\n*$"),
			Tag: "p",
			Type: "block",
			AllowedChildren: []string{"phrase"},
		},
		Parser{
			Re: re("\n"),
			Tag: "br",
			Type: "phrase",
		},
		Parser{
			Re: re("\\b\"([^\"]+)\":[^\\s]+\\b"),
			Tag: "a",
			Type: "phrase",
		},
		Parser{
			Re: re("\\b![^\\s!]+!\\b"),
			Tag: "img",
			Type: "phrase",
		},
		Parser{
			Re: re("\\B\\*\\*(.+)\\*\\*\\B"),
			Tag: "b",
			Type: "phrase",
			AllowedChildren: []string{"phrase"},
		},
		Parser{
			Re: re("\\[\\*\\*(.+)\\*\\*\\]"),
			Tag: "b",
			Type: "phrase",
			AllowedChildren: []string{"phrase"},
		},
		Parser{
			// far from perfect...
			// - add space at the end, otherwise it matches foo**bar**baz
			// - add [^\*], otherwise would catch foo**bar**
			Re: re("\\B\\*(.+[^\\*])\\*\\s"),
			Tag: "strong",
			Type: "phrase",
			AllowedChildren: []string{"phrase"},
			OffsetEnd: -1,
		},
		Parser{
			Re: re("\\[\\*(.+)\\*\\]"),
			Tag: "strong",
			Type: "phrase",
			AllowedChildren: []string{"phrase"},
		},
		Parser{
			Re: re("\\b__(.+)__\\b"),
			Tag: "i",
			Type: "phrase",
			AllowedChildren: []string{"phrase"},
		},
		Parser{
			Re: re("\\[__(.+)__\\]"),
			Tag: "i",
			Type: "phrase",
			AllowedChildren: []string{"phrase"},
		},
		Parser{
			Re: re("\\b_(.+)_\\b"),
			Tag: "em",
			Type: "phrase",
			AllowedChildren: []string{"phrase"},
		},
		Parser{
			Re: re("\\[_(.+)_\\]"),
			Tag: "em",
			Type: "phrase",
			AllowedChildren: []string{"phrase"},
		},
		Parser{
			Re: re("\\B~(.+)~\\B"),
			Tag: "sub",
			Type: "phrase",
			AllowedChildren: []string{"phrase"},
		},
		Parser{
			Re: re("\\[~(.+)~\\]"),
			Tag: "sub",
			Type: "phrase",
			AllowedChildren: []string{"phrase"},
		},
		Parser{
			Re: re("\\B\\^(.+)\\^\\B"),
			Tag: "sup",
			Type: "phrase",
			AllowedChildren: []string{"phrase"},
		},
		Parser{
			Re: re("\\[\\^(.+)\\^\\]"),
			Tag: "sup",
			Type: "phrase",
			AllowedChildren: []string{"phrase"},
		},
		Parser{
			Re: re("\\B\\+(.+)\\+\\B"),
			Tag: "ins",
			Type: "phrase",
			AllowedChildren: []string{"phrase"},
		},
		Parser{
			Re: re("\\[\\+(.+)\\+\\]"),
			Tag: "ins",
			Type: "phrase",
			AllowedChildren: []string{"phrase"},
		},
		Parser{
			Re: re("\\B-(.+)-\\B"),
			Tag: "del",
			Type: "phrase",
			AllowedChildren: []string{"phrase"},
		},
		Parser{
			Re: re("\\[-(.+)-\\]"),
			Tag: "del",
			Type: "phrase",
			AllowedChildren: []string{"phrase"},
		},
		Parser{
			Re: re("\\B@(.+)@\\B"),
			Tag: "code",
			Type: "phrase",
			AllowedChildren: []string{"phrase"},
		},
		Parser{
			Re: re("\\?\\?(.+)\\?\\?"),
			Tag: "cite",
			Type: "phrase",
			AllowedChildren: []string{"phrase"},
		},
	}
}

func ParseDocument(text string) (*document.D, error) {
	d := document.New("")
	blocks := strings.Split(text, "\n\n")
	for _, text := range blocks {
		c, err := parse(text, "", []string{"block","list","phrase"})
		if err == nil {
			d.AddChild(c)
		}
	}
	return d, nil
}

func canAddChild(t string, allowedChildren []string) bool {
	for _, v := range allowedChildren {
		if v == t {
			return true
		}
	}
	return false
}

func parse(text, tag string, allowedChildren []string) (*document.D, error) {
	d := document.New(tag)
	parsed := 0
Iteration:
	for parsed < len(text) {
		bestPos := [2]int{len(text)+1, 0}
		var bestParser *Parser
		for i, p := range parsers {
			if canAddChild(p.Type, allowedChildren) {
				loc := p.Re.FindStringIndex(text[parsed:]+"\n")
				if loc != nil && loc[0] < len(text)-parsed && loc[0] < bestPos[0] {
					bestPos[0] = loc[0]
					bestPos[1] = loc[1]
					bestParser = &parsers[i]
				}
			}
		}
		if bestParser != nil {
			if bestPos[0] > 0 {
				d.AddChild(&document.D{Text: text[parsed:parsed+bestPos[0]]})
				parsed += bestPos[0]
			}
			matches := bestParser.Re.FindStringSubmatch(text[parsed:]+"\n")
			if len(matches) > 0 {
				if len(matches) > 1 {
					var c *document.D
					var err error
					if bestParser.Parse != nil {
						c, err = bestParser.Parse(matches[1], bestParser.Tag)
					} else {
						c, err = parse(matches[1], bestParser.Tag, bestParser.AllowedChildren)
					}
					if err == nil {
						d.AddChild(c)
					}
				} else {
					d.AddChild(&document.D{Tag: bestParser.Tag})
				}
				parsed += len(matches[0])+bestParser.OffsetEnd
				continue Iteration
			}
		}
		d.AddChild(&document.D{Text: text[parsed:]})
		parsed = len(text)
	}
	return d, nil
}


func parseList(text, tag string) (*document.D, error) {
	d, _, err := parseListAtLevel(text, tag, 1)
	return d, err
}

func parseListAtLevel(text, tag string, level int) (*document.D, int, error) {
	prefix := regexp.QuoteMeta(text[0:1])
	re := regexp.MustCompile(fmt.Sprintf("(?s)((%s{%d}) ([^%s]+)\n*)?(%s+)?", prefix, level, prefix, prefix))
	d := document.New(tag)
	parsed := 0
	for parsed < len(text) {
		matches := re.FindStringSubmatch(text[parsed:])
		if len(matches) > 0 {
			c, err := parse(strings.Trim(matches[3], "\n"), "li", []string{"phrase"})
			if err != nil {
				return d, parsed, err
			}
			nextLevel := len(matches[4])
			parsed += len(matches[1])
			if nextLevel > level {
				cc, offset, err := parseListAtLevel(text[parsed:], tag, level+1)
				if err != nil {
					return d, parsed, err
				}
				c.AddChild(cc)
				parsed += offset
			}
			d.AddChild(c)
			if nextLevel < level {
				return d, parsed, nil
			}
		} else {
			return d, parsed, nil
		}
	}
	return d, parsed, nil
}

func parseTable(text, tag string) (*document.D, error) {
	re := regexp.MustCompile("(?Us)\\|(.*)\\|\n")
	d := document.New(tag)
	parsed := 0
	for parsed < len(text) {
		matches := re.FindStringSubmatch(text[parsed:]+"\n")
		if len(matches) > 0 {
			tr := document.New("tr")
			d.AddChild(tr)
			for _, cell := range strings.Split(matches[1], "|") {
				td, err := parse(cell, "td", []string{"phrase"})
				if err != nil {
					return d, err
				}
				tr.AddChild(td)
			}
			parsed += len(matches[0])
		} else {
			return d, nil
		}
	}
	return d, nil
}
