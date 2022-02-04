package main

import (
	"flag"
	"fmt"
	"go/ast"
	"go/format"
	"go/parser"
	"go/token"
	"os"

	"golang.org/x/tools/go/ast/astutil"
)

const MAX_CHANNELS = 9

func main() {

	flag.Parse()
	args := flag.Args()
	if len(args) != 1 {
		fmt.Println(args)
		panic("takes one arg")
	}
	out := args[0]

	fset := token.NewFileSet()

	f2, err := parser.ParseFile(fset, "template", template,
		parser.AllErrors|parser.ParseComments,
	)
	if err != nil {
		panic(err)
	}

	f2.Comments = nil

	var done bool
	astutil.Apply(f2, nil, func(c *astutil.Cursor) bool {
		n := c.Node()
		switch n.(type) {
		case *ast.FuncDecl:
			if done {
				return true
			}
			done = true

			c.Delete()
			for i := 2; i <= MAX_CHANNELS; i++ {
				c.InsertBefore(genSelectCall(c, i, false, false))
				c.InsertBefore(genSelectCall(c, i, true, false))
				c.InsertBefore(genSelectCall(c, i, false, true))
				c.InsertBefore(genSelectCall(c, i, true, true))
			}
		}

		return true
	})

	outStream, err := os.Create(out)
	if err != nil {
		panic(err)
	}
	defer outStream.Close()

	fmt.Fprint(outStream, "// Code generated by a tool. DO NOT EDIT.\n\n") // easier than AST
	err = format.Node(outStream, fset, f2)
	if err != nil {
		panic(err)
	}
}

const template = `
package oselect

func init() {
	panic("this should never be included")
}
`