package main

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	"github.com/cloudfoundry/cli/plugin"
)

// StackSearchResults represents top level attributes of JSON response from Cloud Foundry API
type StackSearchResults struct {
	TotalResults int                    `json:"total_results"`
	TotalPages   int                    `json:"total_pages"`
	Resources    []StackSearchResources `json:"resources"`
}

// StackSearchResources represents resources attribute of JSON response from Cloud Foundry API
type StackSearchResources struct {
	Entity   StackSearchEntity `json:"entity"`
	Metadata Metadata          `json:"metadata"`
}

// StackSearchEntity represents entity attribute of resources attribute within JSON response from Cloud Foundry API
type StackSearchEntity struct {
	Name             string `json:"name"`
	StackDescription string `json:"description"`
}

// GetStackData requests all of the Application data from Cloud Foundry
func (c Sleuth) GetStacks(cli plugin.CliConnection) map[string]StackSearchEntity {
	var data map[string]StackSearchEntity
	data = make(map[string]StackSearchEntity)
	stacks := c.GetStackData(cli)

	for _, val := range stacks.Resources {
		data[val.Metadata.GUID] = val.Entity
	}

	return data
}

// GetStackData requests all of the Application data from Cloud Foundry
func (c Sleuth) GetStackData(cli plugin.CliConnection) StackSearchResults {
	var res StackSearchResults
	res = c.UnmarshallStackSearchResults("/v2/stacks?order-direction=asc&results-per-page=100", cli)

	if res.TotalPages > 1 {
		for i := 2; i <= res.TotalPages; i++ {
			apiUrl := fmt.Sprintf("/v2/stacks?order-direction=asc&page=%v&results-per-page=100", strconv.Itoa(i))
			tRes := c.UnmarshallStackSearchResults(apiUrl, cli)
			res.Resources = append(res.Resources, tRes.Resources...)
		}
	}

	return res
}

func (c Sleuth) UnmarshallStackSearchResults(apiUrl string, cli plugin.CliConnection) StackSearchResults {
	var tRes StackSearchResults
	cmd := []string{"curl", apiUrl}
	output, _ := cli.CliCommandWithoutTerminalOutput(cmd...)
	json.Unmarshal([]byte(strings.Join(output, "")), &tRes)

	return tRes
}
