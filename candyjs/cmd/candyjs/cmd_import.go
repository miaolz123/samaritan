package main

import (
	"bytes"
	"fmt"
	"go/ast"
	"go/format"
	"go/parser"
	"go/token"
	"io/ioutil"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"text/template"
)

// CmdImport generates a new candyjs.PackagePusher function for the given
// Package, the package can be any builtin or any other third party one.
type CmdImport struct {
	Output string `short:"" long:"output" description:"output file name" default:"pkg_%s.go"`
	Debug  bool   `short:"" long:"debug" description:"active debug messages"`
	Args   struct {
		Package string `positional-arg-name:"package" description:"package to import"`
	} `positional-args:"yes" required:"true"`

	curPkgName, fullPkgName, pkgName string
}

// Execute run the CmdImport, follows the go-flags interface
func (c *CmdImport) Execute(args []string) error {
	c.fullPkgName = c.Args.Package
	fmt.Printf("Processing %q\n", c.Args.Package)

	objects, err := c.getObjects()
	if err != nil {
		return err
	}

	c.getCurrentPckgName()
	return c.render(objects)
}

func (c *CmdImport) getCurrentPckgName() {
	pkgs, _ := parser.ParseDir(token.NewFileSet(), ".", nil, 0)
	for pkgName := range pkgs {
		c.curPkgName = pkgName
	}
}

func (c *CmdImport) getObjects() (map[string]*ast.Object, error) {
	pkgs, err := c.parserPackage()
	if err != nil {
		return nil, err
	}

	var objects map[string]*ast.Object
	for _, pkg := range pkgs {
		pkgObjs := c.getPackageObjects(pkg)
		if pkg.Name == "main" || len(pkgObjs) == 0 {
			continue
		}

		c.pkgName = pkg.Name
		objects = pkgObjs
	}

	return objects, nil
}

func (c *CmdImport) parserPackage() (map[string]*ast.Package, error) {
	dir, err := c.getPackagePath(c.fullPkgName)
	if err != nil {
		return nil, err
	}

	return parser.ParseDir(token.NewFileSet(), dir, nil, 0)
}

func (c *CmdImport) getPackageObjects(pkg *ast.Package) map[string]*ast.Object {
	objects := make(map[string]*ast.Object)

	for filename, f := range pkg.Files {
		if strings.HasSuffix(filepath.Base(filename), "_test.go") {
			continue
		}

		for name, object := range f.Scope.Objects {
			if ast.IsExported(name) {
				objects[name] = object
			}
		}

		if c.Debug {
			fmt.Printf("Processed package file %q\n", filename)
		}
	}

	return objects
}

func (c *CmdImport) getPackagePath(pkgName string) (string, error) {
	for _, base := range []string{os.Getenv("GOPATH"), runtime.GOROOT()} {
		dir := filepath.Join(base, "src", pkgName)
		_, err := os.Stat(dir)
		if err == nil {
			return dir, nil
		}
	}

	return "", fmt.Errorf("package %q not found", pkgName)
}

func (c *CmdImport) render(objs map[string]*ast.Object) error {
	t := template.New("tmpl")
	t.Funcs(template.FuncMap{
		"isFunc":           isFunc,
		"isVar":            isVar,
		"isConst":          isConst,
		"isStruct":         isStruct,
		"nameToJavaScript": nameToJavaScript,
	})

	_, err := t.Parse(formatTemplateNewLines(tmpl))
	if err != nil {
		return err
	}

	buf := bytes.NewBuffer(nil)
	err = t.Execute(buf, struct {
		FullPkgName, CurPkgName, PkgName, Objs interface{}
	}{
		FullPkgName: c.fullPkgName,
		CurPkgName:  c.curPkgName,
		PkgName:     c.pkgName,
		Objs:        objs,
	})

	if err != nil {
		return err
	}

	output, err := format.Source(buf.Bytes())
	if err != nil {
		return err
	}

	file := fmt.Sprintf(c.Output, c.pkgName)
	fmt.Printf("File generated %q\n", file)

	return ioutil.WriteFile(file, output, 0644)
}

func formatTemplateNewLines(tmpl string) string {
	return strings.Replace(tmpl, "\\\n", " ", -1)
}

func isFunc(obj *ast.Object) bool {
	return obj.Kind == ast.Fun
}

func isVar(obj *ast.Object) bool {
	return obj.Kind == ast.Var
}

func isConst(obj *ast.Object) bool {
	return obj.Kind == ast.Con
}

func isStruct(obj *ast.Object) bool {
	if obj.Kind != ast.Typ {
		return false
	}

	_, isStruct := obj.Decl.(*ast.TypeSpec).Type.(*ast.StructType)
	return isStruct
}

//TODO: is copy pasted from the main package
func nameToJavaScript(name string) string {
	var toLower, keep string
	for _, c := range name {
		if c >= 'A' && c <= 'Z' && len(keep) == 0 {
			toLower += string(c)
		} else {
			keep += string(c)
		}
	}

	lc := len(toLower)
	if lc > 1 && lc != len(name) {
		keep = toLower[lc-1:] + keep
		toLower = toLower[:lc-1]

	}

	return strings.ToLower(toLower) + keep
}

const tmpl = `
{{$fullPkg := .FullPkgName}}
{{$pkg := .PkgName}}
package {{.CurPkgName}}

import (
	"{{$fullPkg}}"

	"github.com/mcuadros/go-candyjs"
)

func init() {
	candyjs.RegisterPackagePusher("{{$fullPkg}}", func(ctx *candyjs.Context) {
		ctx.PushObject()
		{{range .Objs}} \
		{{if isFunc .}} \
			ctx.PushGoFunction({{$pkg}}.{{.Name}})
			ctx.PutPropString(-2, "{{nameToJavaScript .Name}}")
		{{else if isStruct .}} \
			ctx.PushType({{$pkg}}.{{.Name}}{})
			ctx.PutPropString(-2, "{{.Name}}")
		{{else if isVar .}} \
			ctx.PushProxy({{$pkg}}.{{.Name}})
			ctx.PutPropString(-2, "{{.Name}}")
		{{else if isConst .}} \
			ctx.PushInterface({{$pkg}}.{{.Name}})
			ctx.PutPropString(-2, "{{.Name}}")
		{{else}} \
			//ignored {{.Name}} - {{.Kind}}
		{{end}} \
		{{end}} \
	})
}
`
