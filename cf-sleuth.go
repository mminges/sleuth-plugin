package main

import (
	"fmt"
	"os"
	"strconv"
//	"sort"
//	"flag"
	"github.com/cloudfoundry/cli/plugin"
	"github.com/cloudfoundry/cli/cf/terminal"
	"github.com/cloudfoundry/cli/cf/trace"
	"github.com/fatih/color"
)


// Sleuth represents Buildpack Usage CLI interface
type Sleuth struct{
	UI terminal.UI
}


// OutputResults represents the filtered event results for the input args
type OutputResults struct {
	Comment      string  `json:"comment"`
	Resources    []AppSearchResources `json:"resources"`
}

// Metadata is the data retrieved from the response json
type Metadata struct {
	GUID string `json:"guid"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}

// GetMetadata provides the Cloud Foundry CLI with metadata to provide user about how to use `sleuth` command
func (c *Sleuth) GetMetadata() plugin.PluginMetadata {
	return plugin.PluginMetadata{
		Name: "sleuth",
		Version: plugin.VersionType{
			Major: 0,
			Minor: 1,
			Build: 0,
		},
		Commands: []plugin.Command{
			{
				Name:     "instances",
				HelpText: "Sleuth CF foundation (by michael.minges@cgi.com)",
				UsageDetails: plugin.Usage {
					Usage: UsageText(),
				},
			},
			{
				Name:     "singletons",
				HelpText: "Sleuth CF foundation (by michael.minges@cgi.com)",
				UsageDetails: plugin.Usage {
					Usage: UsageText(),
				},
			},
		},
	}
}

func main() {
	plugin.Start(new(Sleuth))
}

// Run is what is executed by the Cloud Foundry CLI when the sleuth command is specified
func (c Sleuth) Run(cli plugin.CliConnection, args []string) {
	orgs := c.GetOrgs(cli)
	spaces := c.GetSpaces(cli)
	apps := c.GetAppData(cli)
	results := c.FilterResults(cli, orgs, spaces, apps)

	if args[0] == "instances" {
		c.AllApps(results)
	} else if args[0] == "singletons" {
		c.SingletonApps(results)
	}

}

func Usage(code int) {
	fmt.Println("\nUsage: ", UsageText())
	os.Exit(code)
}

func UsageText() (string) {
	usage := "cf {instances|singletons}"
	return usage
}

func SayOK() {
	c := color.New(color.FgGreen).Add(color.Bold)
	c.Println("OK\n")
}

func (c Sleuth) SingletonApps(results OutputResults) {
	Writer := color.Output

	ui := terminal.NewUI(
		os.Stdin,
		Writer,
		terminal.NewTeePrinter(Writer),
		trace.NewLogger(Writer, false, "false", ""),
	)

	ui.Say("Getting apps ...\n")
	SayOK()

	headers := []string{
		"Org",
		"Space",
		"App",
	}

	t := ui.Table(headers)

	for _, val := range results.Resources  {
		if val.Entity.Instances == 1 {
			t.Add(val.Entity.Org, val.Entity.Space, val.Entity.Name)
		}
	}

	t.Print()
}

func (c Sleuth) AllApps(results OutputResults) {
	Writer := color.Output

	ui := terminal.NewUI(
		os.Stdin,
		Writer,
		terminal.NewTeePrinter(Writer),
		trace.NewLogger(Writer, false, "false", ""),
	)

	ui.Say("Getting apps ...\n")
	SayOK()

	headers := []string{
		"Org",
		"Space",
		"App",
		"Instances",
	}

	t := ui.Table(headers)

	for _, val := range results.Resources  {
		instances := strconv.Itoa(val.Entity.Instances)
		t.Add(val.Entity.Org, val.Entity.Space, val.Entity.Name, instances)
	}

	t.Print()
}
