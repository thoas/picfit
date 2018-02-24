package testconv

import (
	"bytes"
	"errors"
	"fmt"
	"go/doc"
	"go/format"
	"go/parser"
	"go/token"
	"io"
	"io/ioutil"
	"sort"
	"strings"
	"testing"
	"text/template"
)

func RunReadmeTest(t *testing.T, srcs ...string) {
	fset := token.NewFileSet()

	var examples []*doc.Example
	for _, src := range srcs {
		astFile, err := parser.ParseFile(fset, src, nil, parser.ParseComments)
		if err != nil {
			t.Fatal(err)
		}
		examples = append(examples, doc.Examples(astFile)...)
	}
	generateReadme(t, fset, examples)
}

func generateReadme(t *testing.T, fset *token.FileSet, examples []*doc.Example) {
	sg := NewSrcGen("README.md")
	readmeExamples := make(ReadmeExamples, len(examples))

	var buf bytes.Buffer
	for i, example := range examples {
		buf.Reset()

		var code string
		if example.Play != nil {
			format.Node(&buf, fset, example.Play)

			play, search := buf.String(), "func main() {"
			idx := strings.Index(play, search)
			if idx == -1 {
				t.Fatalf("bad formatting in example %v, could not find main() func", example.Name)
			}
			code = play[idx+len(search) : len(play)]
		} else {
			format.Node(&buf, fset, example.Code)
			code = buf.String()
		}

		code = strings.Trim(code, "\t\r\n{}")
		code = rewrap(&buf, "\n  > ", code)
		if 0 == len(code) {
			t.Fatalf("bad formatting in example %v, had no code", example.Name)
		}

		output := rewrap(&buf, "\n  > ", example.Output)
		if 0 == len(output) {
			t.Fatalf("bad formatting in example %v, had no output", example.Name)
		}

		title := strings.Title(strings.TrimLeft(example.Name, "_"))
		if 0 == len(title) {
			title = "Overview"
		}

		// the header example has no summary
		summary := rewrap(&buf, "\n  ", example.Doc)
		if 0 == len(summary) && title != "Package" {
			t.Fatalf("bad formatting in example %v, had no summary", example.Name)
		}

		readmeExamples[i] = ReadmeExample{
			Example: example,
			Title:   title,
			Summary: summary,
			Code:    code,
			Output:  output,
		}
	}

	sort.Sort(readmeExamples)
	sg.FuncMap["Examples"] = func() []ReadmeExample {
		return readmeExamples
	}

	if err := sg.Run(); err != nil {
		t.Fatal(err)
	}
}

type ReadmeExample struct {
	Example  *doc.Example
	Complete bool
	Title    string
	Summary  string
	Code     string
	Output   string
}

type ReadmeExamples []ReadmeExample

func (r ReadmeExamples) Len() int           { return len(r) }
func (r ReadmeExamples) Swap(i, j int)      { r[i], r[j] = r[j], r[i] }
func (r ReadmeExamples) Less(i, j int) bool { return r[i].Example.Order < r[j].Example.Order }

func rewrap(buf *bytes.Buffer, with, s string) string {
	buf.Reset()
	if len(s) == 0 {
		return ``
	}
	for i := 1; i < len(s); i++ {
		if s[i-1] == '\n' {
			if s[i] == '\t' {
				i++
			}
			buf.WriteString(with)
			continue
		}
		buf.WriteByte(s[i-1])
	}
	if end := s[len(s)-1]; end != '\n' {
		buf.WriteByte(end)
	}
	return buf.String()
}

type SrcGen struct {
	t           *testing.T
	Disabled    bool
	Data        interface{}
	FuncMap     template.FuncMap
	Name        string
	SrcPath     string
	SrcBytes    []byte
	TplPath     string
	TplBytes    []byte // Actual template bytes
	TplGenBytes []byte // Bytes produced from executed template
}

func NewSrcGen(name string) *SrcGen {
	funcMap := make(template.FuncMap)
	funcMap["args"] = func(s ...interface{}) interface{} {
		return s
	}
	funcMap["TrimSpace"] = strings.TrimSpace
	return &SrcGen{Name: "README.md", FuncMap: funcMap}
}

func (g *SrcGen) Run() error {
	if g.Disabled {
		return fmt.Errorf(`error: run failed because Disabled field is set for "%s"`, g.Name)
	}
	firstErr := func(funcs ...func() error) (err error) {
		for _, f := range funcs {
			err = f()
			if err != nil {
				return
			}
		}
		return
	}
	return firstErr(g.Check, g.Load, g.Generate, g.Format, g.Commit)
}

func (g *SrcGen) Check() error {
	g.Name = strings.TrimSpace(g.Name)
	g.TplPath = strings.TrimSpace(g.TplPath)
	g.SrcPath = strings.TrimSpace(g.SrcPath)
	if len(g.Name) == 0 {
		return errors.New("error: check for Name field failed because it was empty")
	}
	if len(g.TplPath) == 0 {
		g.TplPath = fmt.Sprintf(`internal/testdata/%s.tpl`, g.Name)
	}
	if len(g.SrcPath) == 0 {
		g.SrcPath = fmt.Sprintf(`%s`, g.Name)
	}
	return nil
}

func (g *SrcGen) Load() error {
	var err error
	if g.TplBytes, err = ioutil.ReadFile(g.TplPath); err != nil {
		return fmt.Errorf(`error: load io error "%s" reading TplPath "%s"`, err, g.TplPath)
	}
	if g.SrcBytes, err = ioutil.ReadFile(g.SrcPath); err != nil {
		return fmt.Errorf(`error: load io error "%s" reading SrcPath "%s"`, err, g.SrcPath)
	}
	return nil
}

func (g *SrcGen) Generate() error {
	tpl, err := template.New(g.Name).Funcs(g.FuncMap).Parse(string(g.TplBytes))
	if err != nil {
		return fmt.Errorf(`error: generate error "%s" parsing TplPath "%s"`, err, g.TplPath)
	}

	var buf bytes.Buffer
	if err = tpl.Execute(&buf, g.Data); err != nil {
		return fmt.Errorf(`error: generate error "%s" executing TplPath "%s"`, err, g.TplPath)
	}
	g.TplGenBytes = buf.Bytes()
	return nil
}

func (g *SrcGen) Format() error {

	// Only run gofmt for .go source code.
	if !strings.HasSuffix(g.SrcPath, ".go") {
		return nil
	}

	fmtBytes, err := format.Source(g.TplGenBytes)
	if err != nil {
		return fmt.Errorf(`error: format error "%s" executing TplPath "%s"`, err, g.TplPath)
	}
	g.TplGenBytes = fmtBytes
	return err
}

func (g *SrcGen) IsStale() bool {
	return !bytes.Equal(g.SrcBytes, g.TplGenBytes)
}

func (g *SrcGen) Dump(w io.Writer) string {
	sep := strings.Repeat("-", 80)
	fmt.Fprintf(w, "%[1]s\n  TplBytes:\n%[1]s\n%s\n%[1]s\n", sep, g.TplBytes)
	fmt.Fprintf(w, "  SrcBytes:\n%[1]s\n%s\n%[1]s\n", sep, g.SrcBytes)
	fmt.Fprintf(w, "  TplGenBytes (IsStale: %v):\n%s\n%[3]s\n%[2]s\n",
		g.IsStale(), sep, g.TplGenBytes)
	return g.Name
}

func (g *SrcGen) String() string {
	return g.Name
}

func (g *SrcGen) Commit() error {
	if !g.IsStale() {
		return nil
	}
	return ioutil.WriteFile(g.SrcPath, g.TplGenBytes, 0644)
}
