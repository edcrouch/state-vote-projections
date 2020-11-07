package getresults

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strings"
)

var host = "https://www.nbcnews.com"
var endpoint = "/politics/2020-elections/%state%-president-results?format=json"

// Retrieve ...
func Retrieve(state string) ElectionData {
	resp, err := http.Get(host + strings.Replace(endpoint, "%state%", state, 1))

	if err != nil {
		panic(err)
	}

	defer resp.Body.Close()
	body, closeErr := ioutil.ReadAll(resp.Body)

	if closeErr != nil {
		panic(closeErr)
	}

	test := &WebResponse{}

	jsonErr := json.Unmarshal(body, test)

	if jsonErr != nil {
		panic(closeErr)
	}

	return formatData(test)
}

// FormatData ...
func formatData(data *WebResponse) ElectionData {
	var bidenData map[string]CountyResults
	var trumpData map[string]CountyResults
	results := make(map[string]ResultRow)

	for _, v := range data.CandidateCountyResults {
		if v.FullName == "Joe Biden" {
			bidenData = createCandidateResults(v.CountyResults)
		} else if v.FullName == "Donald Trump" {
			trumpData = createCandidateResults(v.CountyResults)
		}
	}
	for _, v := range data.CountiesPercentIn {
		results[v.Name] = ResultRow{
			Name:      v.Name,
			PercentIn: v.PercentIn,
			Biden:     bidenData[v.Name],
			Trump:     trumpData[v.Name],
		}
	}

	return ElectionData{
		Counties: results,
		State: StateData{
			PercentIn: data.PercentOfExpectedVoteReceived,
		},
	}
}

func createCandidateResults(results []CountyResults) map[string]CountyResults {
	data := make(map[string]CountyResults)
	for _, v := range results {
		data[v.CountyName] = v
	}

	return data
}

// ElectionData ...
type ElectionData struct {
	Counties map[string]ResultRow
	State    StateData
}

// ResultRow ...
type ResultRow struct {
	Name      string
	PercentIn float32
	Biden     CountyResults
	Trump     CountyResults
}

// StateData ...
type StateData struct {
	PercentIn int
}

// WebResponse ...
type WebResponse struct {
	CandidateCountyResults        []CandidateCountyResult `json:"candidateCountyResults"`
	CountiesPercentIn             []CountyPercentIn       `json:"countiesPercentIn"`
	PercentOfExpectedVoteReceived int                     `json:"percentOfExpectedVoteReceived"`
}

// CandidateCountyResult ...
type CandidateCountyResult struct {
	DeclaredWinner       bool            `json:"declaredWinner"`
	FullName             string          `json:"fullName"`
	TotalPercentReceived float32         `json:"totalPercentReceived"`
	TotalVotesReceived   int             `json:"totalVotesReceived"`
	CountyResults        []CountyResults `json:"countyResults"`
}

// CountyPercentIn ...
type CountyPercentIn struct {
	Name      string  `json:"name"`
	PercentIn float32 `json:"percentIn"`
}

// CountyResults ...
type CountyResults struct {
	CountyName    string  `json:"countyName"`
	PercentOfVote float32 `json:"percentOfVote"`
	VotesReceived int     `json:"votesReceived"`
	IsLeading     bool    `json:"isLeading"`
}
