// +build ignore

package main

import (
	"log"

	"github.com/shurcooL/vfsgen"
	"github.com/straightdave/trunks/lib/plot"
)

func main() {
	err := vfsgen.Generate(plot.Assets, vfsgen.Options{
		PackageName:  "plot",
		BuildTags:    "!dev",
		VariableName: "Assets",
	})

	if err != nil {
		log.Fatalln(err)
	}
}
