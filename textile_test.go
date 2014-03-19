package gotextile

import (
	"log"
	"testing"

	"launchpad.net/gocheck"
)

func Test(t *testing.T) { gocheck.TestingT(t) }

type TestSuite struct {}

var _ = gocheck.Suite(&TestSuite{})

type TestCase struct {
	Source string
	Expected string
}

var tests = []TestCase{
	TestCase{
		`p. paragraph`,
		`<p>paragraph</p>`,
	},
	TestCase{
		`p. paragraph
next line`,
		`<p>paragraph<br>next line</p>`,
	},
	TestCase{
		`paragraph`,
		`<p>paragraph</p>`,
	},
	TestCase{
		`with entities < " ' >`,
		`<p>with entities &lt; &#34; &#39; &gt;</p>`,
	},
	TestCase{
		`_emphasis_ __italic__ *strong* **bold** ^sup^ ~sub~`,
		`<p><em>emphasis</em> <i>italic</i> <strong>strong</strong> <b>bold</b> <sup>sup</sup> <sub>sub</sub></p>`,
	},
	TestCase{
		`e[_m_]phasis i[__t__]alic s[*t*]rong b[**o**]ld s[^u^]p s[~u~]b`,
		`<p>e<em>m</em>phasis i<i>t</i>alic s<strong>t</strong>rong b<b>o</b>ld s<sup>u</sup>p s<sub>u</sub>b</p>`,
	},
	TestCase{
		`e_m_phasis i__t__alic s*t*rong b**old** s^u^p s~u~b`,
		`<p>e_m_phasis i__t__alic s*t*rong b**old** s^u^p s~u~b</p>`,
	},
	TestCase{
		`*strong strong* *stron*g strong* strong*`,
		`<p><strong>strong strong</strong> <strong>stron*g strong</strong> strong*</p>`,
	},
	TestCase{
		`_emphasis emphasis_ _emphasi_s emphasis_ emphasis_`,
		`<p><em>emphasis emphasis</em> <em>emphasi_s emphasis</em> emphasis_</p>`,
	},
	TestCase{
		`@code code@ c@o@de @cod@e code@ code@`,
		`<p><code>code code</code> c@o@de <code>cod@e code</code> code@</p>`,
	},
	TestCase{
		`h1. paragraph`,
		`<h1>paragraph</h1>`,
	},
	TestCase{
		`h1. paragraph
next line`,
		`<h1>paragraph<br>next line</h1>`,
	},
	TestCase{
		`* list`,
		`<ul><li>list</li></ul>`,
	},
	TestCase{
		`* list1
* list2`,
		`<ul><li>list1</li><li>list2</li></ul>`,
	},
	TestCase{
		`* list1
** list1.1
** list1.2
* list2`,
		`<ul><li>list1<ul><li>list1.1</li><li>list1.2</li></ul></li><li>list2</li></ul>`,
	},
	TestCase{
		`** list1.1
** list1.2
* list2`,
		`<ul><li><ul><li>list1.1</li><li>list1.2</li></ul></li><li>list2</li></ul>`,
	},
	TestCase{
		`# list`,
		`<ol><li>list</li></ol>`,
	},
	TestCase{
		`# list1
# list2`,
		`<ol><li>list1</li><li>list2</li></ol>`,
	},
	TestCase{
		`# list1
## list1.1
## list1.2
# list2`,
		`<ol><li>list1<ol><li>list1.1</li><li>list1.2</li></ol></li><li>list2</li></ol>`,
	},
	TestCase{
		`## list1.1
## list1.2
# list2`,
		`<ol><li><ol><li>list1.1</li><li>list1.2</li></ol></li><li>list2</li></ol>`,
	},
	TestCase{
		`- name := definition`,
		`<dl><dt>name</dt><dd>definition</dd></dl>`,
	},
	TestCase{
		`- name1 := definition1
- name2 := definition2`,
		`<dl><dt>name1</dt><dd>definition1</dd><dt>name2</dt><dd>definition2</dd></dl>`,
	},
	TestCase{
		`- name1 := definition1
- name2 := definition2
line2 =:`,
		`<dl><dt>name1</dt><dd>definition1</dd><dt>name2</dt><dd>definition2<br>line2</dd></dl>`,
	},
	TestCase{
		`bc. line1
line2`,
		`<pre><code>line1
line2</code></pre>`,
	},
	TestCase{
		`pre. line1
line2`,
		`<pre>line1
line2</pre>`,
	},
	TestCase{
		`###. line1
line2`,
		``,
	},
	TestCase{
		`??cite??`,
		`<p><cite>cite</cite></p>`,
	},
	TestCase{
		`??cite1??
??cite2??
-- Author`,
		`<p><cite>cite1</cite><br><cite>cite2</cite><br>-- Author</p>`,
	},
	TestCase{
		`|cell1.1|cell1.2|cell1.3|
|cell2.1|cell2.2|cell2.3|`,
		`<table><tr><td>cell1.1</td><td>cell1.2</td><td>cell1.3</td></tr><tr><td>cell2.1</td><td>cell2.2</td><td>cell2.3</td></tr></table>`,
	},
}

func (s *TestSuite) TestTextile(c *gocheck.C) {
	for _, t := range tests {
		log.Printf("==== %#v ====", t.Source)
		txt, _ := TextileToHtml(t.Source)
		c.Assert(txt, gocheck.Equals, t.Expected, gocheck.Commentf(t.Source))
	}
}
