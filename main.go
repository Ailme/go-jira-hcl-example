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
	Type            string   `hcl:"type,label"`
	Project         string   `hcl:"project"`
	Summary         string   `hcl:"summary"`
	Description     string   `hcl:"description,optional"`
	AppLayer        string   `hcl:"app_layer,optional"`
	Components      []string `hcl:"components,optional"`
	SprintId        int      `hcl:"sprint,optional"`
	Epic            string   `hcl:"epic,optional"`
	Labels          []string `hcl:"labels,optional"`
	StoryPoint      int      `hcl:"story_point,optional"`
	QaStoryPoint    int      `hcl:"qa_story_point,optional"`
	Assignee        string   `hcl:"assignee,optional"`
	Developer       string   `hcl:"developer,optional"`
	TeamLead        string   `hcl:"team_lead,optional"`
	TechLead        string   `hcl:"tech_lead,optional"`
	ReleaseEngineer string   `hcl:"release_engineer,optional"`
	Tester          string   `hcl:"tester,optional"`
	Parent          string   `hcl:"parent,optional"`
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

func parse(filename string) (*Root, error) {
	parser := hclparse.NewParser()
	f, diags := parser.ParseHCLFile(filename)
	if diags.HasErrors() {
		renderDiags(diags, parser.Files())

		return nil, diags
	}

	var root Root
	diags = gohcl.DecodeBody(f.Body, nil, &root)
	if diags.HasErrors() {
		renderDiags(diags, parser.Files())

		return nil, diags
	}

	return &root, nil
}

func authJira() (*jira.Client, error) {
	basicAuth := jira.BasicAuthTransport{
		Username: os.Getenv("JIRA_USERNAME"),
		Password: os.Getenv("JIRA_PASSWORD"),
	}
	return jira.NewClient(basicAuth.Client(), os.Getenv("JIRA_URL"))
}

func processCreate(root *Root, jiraClient *jira.Client) error {
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
			return err
		}

		fmt.Println(issue.Key)
	}

	return nil
}

func main() {
	filename := "example.hcl"

	root, err := parse(filename)
	if err != nil {
		log.Fatal(err)
	}
	_, _ = pretty.Println(root)

	// костыль, чтобы отключить запросы в Jira
	if false {
		jiraClient, err := authJira()
		if err != nil {
			log.Fatal(err)
		}

		err = processCreate(root, jiraClient)
		if err != nil {
			log.Fatal(err)
		}
	}
}
