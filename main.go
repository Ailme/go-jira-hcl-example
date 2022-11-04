package main

import (
	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/gohcl"
	"github.com/hashicorp/hcl/v2/hclparse"
	"github.com/kr/pretty"
	"log"
	"os"
)

type Root struct {
	Create config `hcl:"create,block"`
}

type config struct {
	Type        string   `hcl:"type,label"`
	Project     string   `hcl:"project"`
	Summary     string   `hcl:"summary"`
	Assignee    string   `hcl:"assignee"`
	Description string   `hcl:"description,optional"`
	Labels      []string `hcl:"labels,optional"`
}

func main() {
	filename := "example.hcl"

	parser := hclparse.NewParser()
	f, diags := parser.ParseHCLFile(filename)
	if diags.HasErrors() {
		wr := hcl.NewDiagnosticTextWriter(
			os.Stdout,      // writer to send messages to
			parser.Files(), // the parser's file cache, for source snippets
			78,             // wrapping width
			true,           // generate colored/highlighted output
		)
		_ = wr.WriteDiagnostics(diags)

		log.Fatal(diags)
	}

	var root Root
	diags = gohcl.DecodeBody(f.Body, nil, &root)
	if diags.HasErrors() {
		wr := hcl.NewDiagnosticTextWriter(
			os.Stdout,      // writer to send messages to
			parser.Files(), // the parser's file cache, for source snippets
			78,             // wrapping width
			true,           // generate colored/highlighted output
		)
		_ = wr.WriteDiagnostics(diags)

		log.Fatal(diags)
	}

	_, _ = pretty.Println(root)
}
