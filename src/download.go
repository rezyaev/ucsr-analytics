package main

import (
	"compress/gzip"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
)

var client = &http.Client{}

func FetchDemos() error {
	history, error := fetchHistory()
	if error != nil {
		return error
	}
	fmt.Printf("fetched history \n %+v \n", len(history.Items))

	for _, historyItem := range history.Items {
		match, error := fetchMatch(historyItem.Match_id)
		if error != nil {
			return error
		}
		fmt.Printf("fetched match \n %+v \n", match)

		fmt.Printf("fetching demo... \n %+v \n", match.Demo_url[0])
		fetchDemo(match)
	}

	return nil
}

func fetchDemo(match Match) error {
	response, error := http.Get(match.Demo_url[0])
	if error != nil {
		return error
	}
	defer response.Body.Close()

	out, error := os.Create("../demos/" + match.Voting.Map.Pick[0] + "/" + match.Match_id + ".dem")
	if error != nil {
		return error
	}
	defer out.Close()

	gzreader, error := gzip.NewReader(response.Body)
	if error != nil {
		return error
	}
	defer gzreader.Close()

	_, error = io.Copy(out, gzreader)
	if error != nil {
		return error
	}
	return nil
}

type Match struct {
	Match_id string
	Demo_url []string
	Voting   struct{ Map struct{ Pick []string } }
}

func fetchMatch(matchId string) (Match, error) {
	request, error := http.NewRequest("GET", "https://open.faceit.com/data/v4/matches/"+matchId, nil)
	if error != nil {
		return Match{}, nil
	}

	request.Header.Add("Authorization", "Bearer "+os.Getenv("FACEIT_API_TOKEN"))
	response, error := client.Do(request)
	if error != nil {
		return Match{}, nil
	}
	defer response.Body.Close()

	var match Match
	if error := json.NewDecoder(response.Body).Decode(&match); error != nil {
		return Match{}, nil
	}

	return match, nil
}

type History struct {
	Items []struct{ Match_id string }
}

func fetchHistory() (History, error) {
	request, error := http.NewRequest("GET", "https://open.faceit.com/data/v4/players/54dfd14c-c2e3-4a2d-acb0-286ffac366aa/history?game=csgo&offset=0&limit=250", nil)
	if error != nil {
		return History{}, nil
	}
	request.Header.Add("Authorization", "Bearer "+os.Getenv("FACEIT_API_TOKEN"))

	response, error := client.Do(request)
	if error != nil {
		return History{}, nil
	}
	defer response.Body.Close()

	var history History
	if error := json.NewDecoder(response.Body).Decode(&history); error != nil {
		return History{}, nil
	}

	return history, nil
}
