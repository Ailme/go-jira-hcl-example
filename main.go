package main

import (
	"fmt"
	"github.com/andygrunwald/go-jira"
	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/gohcl"
	"github.com/hashicorp/hcl/v2/hclparse"
	"github.com/kr/pretty"
	"github.com/zclconf/go-cty/cty"
	"github.com/zclconf/go-cty/cty/function"
	"log"
	"os"
)

type VariablesBlock struct {
	Variables variables `hcl:"variables,block"`
	Remains   hcl.Body  `hcl:",remain"`
}

type variables struct {
	Remains hcl.Body `hcl:",remain"`
}

type CreateBlocks struct {
	Create []createConfig `hcl:"create,block"`
}

type createConfig struct {
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
	Remains         hcl.Body `hcl:",remain"`
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

var EnvFunc = function.New(&function.Spec{
	Params: []function.Parameter{
		{
			Name:             "env",
			Type:             cty.String,
			AllowDynamicType: true,
		},
	},
	Type: function.StaticReturnType(cty.String),
	Impl: func(args []cty.Value, retType cty.Type) (cty.Value, error) {
		in := args[0].AsString()
		out := os.Getenv(in)
		return cty.StringVal(out), nil
	},
})

func parse(filename string) (*CreateBlocks, error) {
	parser := hclparse.NewParser()
	f, diags := parser.ParseHCLFile(filename)
	if diags.HasErrors() {
		renderDiags(diags, parser.Files())

		return nil, diags
	}

	ctx := &hcl.EvalContext{
		Variables: map[string]cty.Value{
			"iter": cty.EmptyObjectVal,
		},
		Functions: map[string]function.Function{
			"env": EnvFunc,
		},
	}

	var variablesBlock VariablesBlock
	_ = gohcl.DecodeBody(f.Body, ctx, &variablesBlock)

	if variablesBlock.Variables.Remains != nil {
		variables, diags := variablesBlock.Variables.Remains.JustAttributes()
		if diags.HasErrors() {
			renderDiags(diags, parser.Files())

			return nil, diags
		}

		for key, variable := range variables {
			var value cty.Value

			diags := gohcl.DecodeExpression(variable.Expr, nil, &value)
			if diags.HasErrors() {
				return nil, diags
			}

			ctx.Variables[key] = value
		}
	}

	var createBlocks CreateBlocks
	schema, _ := gohcl.ImpliedBodySchema(&createBlocks)
	bc, _, diags := variablesBlock.Remains.PartialContent(schema)
	if diags.HasErrors() {
		renderDiags(diags, parser.Files())

		return nil, diags
	}

	for _, block := range bc.Blocks {
		var config createConfig
		attr, diags := block.Body.JustAttributes()
		if diags.HasErrors() {
			renderDiags(diags, parser.Files())

			return nil, diags
		}

		forEach, found := attr["for_each"]

		if found {
			var forEachValue cty.Value

			diags := gohcl.DecodeExpression(forEach.Expr, ctx, &forEachValue)
			if diags.HasErrors() {
				return nil, diags
			}

			forEachValue.ForEachElement(func(key cty.Value, val cty.Value) (stop bool) {
				ctx.Variables["iter"] = val

				diags = gohcl.DecodeBody(block.Body, ctx, &config)
				if diags.HasErrors() {
					renderDiags(diags, parser.Files())

					return true
				}
				config.Type = block.Labels[0]
				config.Remains = nil
				createBlocks.Create = append(createBlocks.Create, config)

				return false
			})
		} else {
			diags = gohcl.DecodeBody(block.Body, ctx, &config)
			if diags.HasErrors() {
				renderDiags(diags, parser.Files())

				return nil, diags
			}
			config.Type = block.Labels[0]
			createBlocks.Create = append(createBlocks.Create, config)
		}
	}

	return &createBlocks, nil
}

func authJira() (*jira.Client, error) {
	basicAuth := jira.BasicAuthTransport{
		Username: os.Getenv("JIRA_USERNAME"),
		Password: os.Getenv("JIRA_PASSWORD"),
	}
	return jira.NewClient(basicAuth.Client(), os.Getenv("JIRA_URL"))
}

func processCreate(createBlocks *CreateBlocks, jiraClient *jira.Client) error {
	for _, config := range createBlocks.Create {
		i := jira.Issue{
			Fields: &jira.IssueFields{
				Description: config.Description,
				Type:        jira.IssueType{Name: config.Type},
				Project:     jira.Project{Key: config.Project},
				Summary:     config.Summary,
				Labels:      config.Labels,
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
