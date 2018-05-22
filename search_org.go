package main

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	"github.com/cloudfoundry/cli/plugin"
)

// OrgSearchResults represents top level attributes of JSON response from Cloud Foundry API
type OrgSearchResults struct {
	TotalResults int                  `json:"total_results"`
	TotalPages   int                  `json:"total_pages"`
	Resources    []OrgSearchResources `json:"resources"`
}

// OrgSearchResources represents resources attribute of JSON response from Cloud Foundry API
type OrgSearchResources struct {
	Entity   OrgSearchEntity `json:"entity"`
	Metadata Metadata        `json:"metadata"`
}

// OrgSearchEntity represents entity attribute of resources attribute within JSON response from Cloud Foundry API
type OrgSearchEntity struct {
	Name string `json:"name"`
}

func (c Sleuth) GetOrgs(cli plugin.CliConnection) map[string]string {
	var data map[string]string
	data = make(map[string]string)
	orgs := c.GetOrgData(cli)

	for _, val := range orgs.Resources {
		data[val.Metadata.GUID] = val.Entity.Name
	}

	return data
}

// GetOrgData requests all of the Application data from Cloud Foundry
func (c Sleuth) GetOrgData(cli plugin.CliConnection) OrgSearchResults {
	var res OrgSearchResults
	res = c.UnmarshallOrgSearchResults("/v2/organizations?order-direction=asc&results-per-page=100", cli)

	if res.TotalPages > 1 {
		for i := 2; i <= res.TotalPages; i++ {
			apiUrl := fmt.Sprintf("/v2/organizations?order-direction=asc&page=%v&results-per-page=100", strconv.Itoa(i))
			tRes := c.UnmarshallOrgSearchResults(apiUrl, cli)
			res.Resources = append(res.Resources, tRes.Resources...)
		}
	}

	return res
}

func (c Sleuth) UnmarshallOrgSearchResults(apiUrl string, cli plugin.CliConnection) OrgSearchResults {
	var tRes OrgSearchResults
	cmd := []string{"curl", apiUrl}
	output, _ := cli.CliCommandWithoutTerminalOutput(cmd...)
	json.Unmarshal([]byte(strings.Join(output, "")), &tRes)

	return tRes
}
