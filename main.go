package main

import (
	"github.com/hashicorp/hcl/v2/gohcl"
	"github.com/hashicorp/hcl/v2/hclparse"
	"github.com/kr/pretty"
	"log"
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
		log.Fatal(diags)
	}

	var root Root
	diags = gohcl.DecodeBody(f.Body, nil, &root)
	if diags.HasErrors() {
		log.Fatal(diags)
	}

	_, _ = pretty.Println(root)
}
