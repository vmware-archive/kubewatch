// Yannotated generates an annotated yaml config boilerplate from a Go structure tree.
package main

import (
	"flag"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"io"
	"log"
	"os"
	"strings"

	"github.com/fatih/structtag"
	"github.com/segmentio/textio"
)

const (
	YAMLFormat = "yaml"
	GoFormat   = "go"
)

type Flags struct {
	// Dir is the directory containing the source code.
	Dir string
	// Package of the root type name.
	Package string
	// Type is the name of the root type of the config tree.
	Type string
	// Output is the file name of the generated output.
	Output string
	// Format is the format of the output.
	Format string
}

func (f *Flags) Bind(fs *flag.FlagSet) {
	if fs == nil {
		fs = flag.CommandLine
	}
	fs.StringVar(&f.Dir, "dir", ".", "Directory of the Go source code")
	fs.StringVar(&f.Package, "package", "", "Name of the root struct type")
	fs.StringVar(&f.Type, "type", "", "Name of the root struct type")
	fs.StringVar(&f.Output, "o", "", "Filename of the generated output")
	fs.StringVar(&f.Format, "format", "", "Output format: Yaml, Go")
}

func mainE(flags Flags) error {
	var fset token.FileSet
	pkgs, err := parser.ParseDir(&fset, flags.Dir, nil, parser.ParseComments)
	if err != nil {
		return err
	}
	pkg, found := pkgs[flags.Package]
	if !found {
		return fmt.Errorf("cannot find package %q in %v", flags.Package, pkgs)
	}

	// trim all unexported symbols
	for _, f := range pkg.Files {
		ast.FileExports(f)
	}

	types := collectTypes(pkg)

	root, found := types[flags.Type]
	if !found {
		return fmt.Errorf("cannot find root type %q in %v", flags.Type, types)
	}
	w, err := os.Create(flags.Output)
	if err != nil {
		return err
	}
	defer w.Close()

	if flags.Format == GoFormat {
		fmt.Fprintf(w, "package %s\n\n", flags.Package)
		fmt.Fprintf(w, "var yannotated = `")
		defer fmt.Fprintf(w, "`\n")
	}

	return emit(w, types, root)
}

func emit(w io.Writer, types map[string]*ast.StructType, node *ast.StructType) error {
	for _, field := range node.Fields.List {
		name, err := fieldName(field)
		if err != nil {
			return err
		}
		if field.Doc != nil {
			lines := strings.Split(strings.TrimSpace(field.Doc.Text()), "\n")
			for _, l := range lines {
				fmt.Fprintf(w, "# %s\n", l)
			}
		}

		fmt.Fprintf(w, "%s:", name)

		switch typ := field.Type.(type) {
		case *ast.Ident:

			switch name := typ.Name; name {
			case "string":
				fmt.Fprintln(w, ` ""`)
			case "int":
				fmt.Fprintln(w, " 0")
			case "bool":
				fmt.Fprintln(w, " false")
			default:
				t, found := types[name]
				if !found {
					return fmt.Errorf("cannot find type %q", name)
				}
				fmt.Fprintf(w, "\n")
				iw := textio.NewPrefixWriter(w, "  ")
				if err := emit(iw, types, t); err != nil {
					return err
				}
			}
		case *ast.MapType:
			fmt.Fprintf(w, " {}\n")
		default:
			return fmt.Errorf("unsupported field type: %T (%s)", field.Type, field.Type)
		}
	}
	return nil
}

func fieldName(field *ast.Field) (string, error) {
	if field.Tag != nil {
		// remove backticks
		clean := field.Tag.Value[1 : len(field.Tag.Value)-1]
		tags, err := structtag.Parse(clean)
		if err != nil {
			return "", fmt.Errorf("while parsing %q: %w", clean, err)
		}

		var yamlName, jsonName string
		for _, tag := range tags.Tags() {
			switch tag.Key {
			case "json":
				jsonName = tag.Name
			case "yaml":
				yamlName = tag.Name
			}
		}
		if yamlName != "" {
			return yamlName, nil
		}
		if jsonName != "" {
			return jsonName, nil
		}
	}
	if got, want := len(field.Names), 1; got != want {
		return "", fmt.Errorf("unsupported number of struct field names, got: %d, want: %d", got, want)
	}

	return strings.ToLower(field.Names[0].Name), nil
}

func collectTypes(n ast.Node) map[string]*ast.StructType {
	v := typeCollectingVisitor(map[string]*ast.StructType{})
	ast.Walk(v, n)
	return v
}

type typeCollectingVisitor map[string]*ast.StructType

func (v typeCollectingVisitor) Visit(node ast.Node) (w ast.Visitor) {
	switch n := node.(type) {
	case *ast.TypeSpec:
		if t, ok := n.Type.(*ast.StructType); ok {
			v[n.Name.Name] = t
		}
	case *ast.Package:
		return v
	case *ast.File:
		return v
	case *ast.GenDecl:
		if n.Tok == token.TYPE {
			return v
		}
	}
	return nil
}

func main() {
	var flags Flags
	flags.Bind(nil)
	flag.Parse()

	if err := mainE(flags); err != nil {
		log.Fatal(err)
	}
}
