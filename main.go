package main

import (
	"math"
	"os"
	"smorty/electoral-college/getresults"
	"sort"

	"golang.org/x/text/language"
	"golang.org/x/text/message"

	"github.com/olekukonko/tablewriter"
)

func main() {

	var state = os.Args[1]
	results := getresults.Retrieve(state)
	createData(results, state)
}

func createData(results getresults.ElectionData, state string) {
	data := make([][]string, 0)

	var totalVotesBiden int
	var totalVotesTrump int
	var percentage float32
	var totalCounties int
	var totalProjectedBiden int
	var totalProjectedTrump int

	p := message.NewPrinter(language.English)

	for i, v := range results.Counties {
		projectedBiden := generateProjection(v.Biden.VotesReceived, v.PercentIn)
		projectedTrump := generateProjection(v.Trump.VotesReceived, v.PercentIn)

		trumpLead := v.Trump.VotesReceived - v.Biden.VotesReceived
		bidenLead := v.Biden.VotesReceived - v.Trump.VotesReceived
		trumpLeadDisplay := ""
		bidenLeadDisplay := ""

		if trumpLead > 0 {
			trumpLeadDisplay = p.Sprintf("+%d", trumpLead)
		}

		if bidenLead > 0 {
			bidenLeadDisplay = p.Sprintf("+%d", bidenLead)
		}

		data = append(data, []string{
			i,
			p.Sprintf("%.2f%%", v.PercentIn),
			"Biden",
			p.Sprintf("%d", v.Biden.VotesReceived),
			p.Sprintf("%d", projectedBiden),
			bidenLeadDisplay,
		})

		data = append(data, []string{
			i,
			"",
			"Trump",
			p.Sprintf("%d", v.Trump.VotesReceived),
			p.Sprintf("%d", projectedTrump),
			trumpLeadDisplay,
		})

		totalVotesBiden += v.Biden.VotesReceived
		totalVotesTrump += v.Trump.VotesReceived
		totalProjectedBiden += projectedBiden
		totalProjectedTrump += projectedTrump
		percentage += v.PercentIn
		totalCounties++
	}

	sort.SliceStable(data, func(i, j int) bool {
		return data[i][0] < data[j][0]
	})

	table := tablewriter.NewWriter(os.Stdout)

	for _, row := range data {
		color := tablewriter.FgBlueColor
		if row[2] == "Trump" {
			color = tablewriter.FgRedColor
		}
		appendRow(table, row, color)
	}

	table.SetHeader([]string{"County", "Percent In", "Candidate", "Current", "Projected", "Projected Difference"})
	table.SetFooter([]string{"County", "Percent In", "Candidate", "Current", "Projected", "Projected Difference"})

	table.SetAutoMergeCellsByColumnIndex([]int{0})
	table.SetRowLine(true)

	trumpLead := totalProjectedTrump - totalProjectedBiden
	bidenLead := totalProjectedBiden - totalProjectedTrump
	trumpLeadDisplay := ""
	bidenLeadDisplay := ""

	if trumpLead > 0 {
		trumpLeadDisplay = p.Sprintf("+%d", trumpLead)
	}

	if bidenLead > 0 {
		bidenLeadDisplay = p.Sprintf("+%d", bidenLead)
	}

	appendRow(table, []string{
		"Total",
		p.Sprintf("%.d%%", results.State.PercentIn),
		"Biden",
		p.Sprintf("%d", totalVotesBiden),
		p.Sprintf("%d", totalProjectedBiden),
		bidenLeadDisplay,
	}, tablewriter.FgBlueColor)

	appendRow(table, []string{
		"Total",
		"",
		"Trump",
		p.Sprintf("%d", totalVotesTrump),
		p.Sprintf("%d", totalProjectedTrump),
		trumpLeadDisplay,
	}, tablewriter.FgRedColor)

	table.SetCaption(true, state)
	table.Render()

}

func generateProjection(votes int, percent float32) int {
	projection := float32(votes) / (percent / 100)
	result := math.Floor(float64(projection))
	return int(result)
}

// func appendRow(table, [], rowCount, color)
func appendRow(table *tablewriter.Table, row []string, color int) {
	backgroundColor := 0
	table.Rich(row, []tablewriter.Colors{
		tablewriter.Colors{tablewriter.Normal, backgroundColor},
		tablewriter.Colors{tablewriter.Normal, backgroundColor},
		tablewriter.Colors{tablewriter.Normal, backgroundColor},
		tablewriter.Colors{tablewriter.Normal, backgroundColor},
		tablewriter.Colors{tablewriter.Normal, backgroundColor},
		tablewriter.Colors{tablewriter.Normal, color},
	})
}
