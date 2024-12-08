package main

import (
	"context"
	_ "embed"
	"fmt"
	"github.com/urfave/cli/v3"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
	"log"
	"os"
	"text/template"
)

//go:embed resource.tmpl
var tmpl string

type Data struct {
	PackageName      string
	ResourceName     string
	ResourceTypeName string
}

func main() {

	cmd := &cli.Command{
		Name:  "boom",
		Usage: "make an explosive entrance",
		Action: func(ctx context.Context, cmd *cli.Command) error {
			terraformType := cmd.Args().Get(0)
			resourceName := cmd.Args().Get(1)
			cwd, _ := os.Getwd()
			fmt.Printf("Getting Resource Type of: %q\n", terraformType)
			fmt.Printf("Current Working Directory: %q\n", cwd)
			formattedResourceName := cases.Title(language.English, cases.NoLower).String(resourceName)
			templ, _ := template.New("resource").Parse(tmpl)
			templ.Execute(os.Stdout, Data{
				ResourceName: formattedResourceName,
			})
			return nil
		},
	}

	if err := cmd.Run(context.Background(), os.Args); err != nil {
		log.Fatal(err)
	}
}
