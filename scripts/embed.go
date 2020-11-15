// +build ignore

package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/lobre/doodle/pkg/embeds/htmldir"
	"github.com/lobre/doodle/pkg/embeds/staticdir"
	"github.com/shurcooL/vfsgen"
)

func main() {
	pkgs := map[string]http.FileSystem{
		"staticdir": staticdir.FS,
		"htmldir":   htmldir.FS,
	}

	for pkg, fs := range pkgs {
		filename := fmt.Sprintf("pkg/embeds/%s/vfsdata.go", pkg)
		fmt.Printf("Generating %s...\n", filename)

		err := vfsgen.Generate(fs, vfsgen.Options{
			Filename:     filename,
			PackageName:  pkg,
			BuildTags:    "embeds",
			VariableName: "FS",
		})

		if err != nil {
			log.Fatalln(err)
		}
	}
}
