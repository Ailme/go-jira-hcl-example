package main

import (
	"fmt"
	"github.com/andygrunwald/go-jira"
	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/gohcl"
	"github.com/hashicorp/hcl/v2/hclparse"
	"github.com/kr/pretty"
	"log"
	"os"
)

type Root struct {
	Create []config `hcl:"create,block"`
}

type config struct {
	Type        string   `hcl:"type,label"`
	Project     string   `hcl:"project"`
	Summary     string   `hcl:"summary"`
	Assignee    string   `hcl:"assignee"`
	Description string   `hcl:"description,optional"`
	Labels      []string `hcl:"labels,optional"`
}

func renderDiags(diags hcl.Diagnostics, files map[string]*hcl.File) {
	wr := hcl.NewDiagnosticTextWriter(
		os.Stdout, // writer to send messages to
		files,     // the parser's file cache, for source snippets
		78,        // wrapping width
		true,      // generate colored/highlighted output
	)
	_ = wr.WriteDiagnostics(diags)
}

func main() {
	filename := "example.hcl"

	parser := hclparse.NewParser()
	f, diags := parser.ParseHCLFile(filename)
	if diags.HasErrors() {
		renderDiags(diags, parser.Files())

		log.Fatal(diags)
	}

	var root Root
	diags = gohcl.DecodeBody(f.Body, nil, &root)
	if diags.HasErrors() {
		renderDiags(diags, parser.Files())

		log.Fatal(diags)
	}

	_, _ = pretty.Println(root)

	basicAuth := jira.BasicAuthTransport{
		Username: os.Getenv("JIRA_USERNAME"),
		Password: os.Getenv("JIRA_PASSWORD"),
	}
	jiraClient, _ := jira.NewClient(basicAuth.Client(), os.Getenv("JIRA_URL"))

	for _, create := range root.Create {
		i := jira.Issue{
			Fields: &jira.IssueFields{
				Description: create.Description,
				Type:        jira.IssueType{Name: create.Type},
				Project:     jira.Project{Key: create.Project},
				Summary:     create.Summary,
				Labels:      create.Labels,
			},
		}

		issue, _, err := jiraClient.Issue.Create(&i)
		if err != nil {
			log.Fatal(err)
		}

		fmt.Println(issue.Key)
	}
}
