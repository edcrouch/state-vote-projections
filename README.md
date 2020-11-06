# state-vote-projections

Creates VERY simple projections for each state / county based on the current
percentage of total votes in and the current number of votes each candidate has.

Run like this: `go run main.go nevada` or `go run main.go new-mexico`.

There are also prebuilt binaries available that accept the state name as an argument.

You can run those like this: `./state-vote-projections-linux nevada`.

Uses the (unpublished?) API that MSNBC / NBC News uses to update their websites.
