package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
)

type GitHubEvent struct {
	Type    string  `json:"type"`
	Repo    Repo    `json:"repo"`
	Payload Payload `json:"payload"`
}

type Repo struct {
	Name string `json:"name"`
}

type Payload struct {
	Ref     string   `json:"ref"`
	RefType string   `json: ref_type`
	Action  string   `json:"action"`
	Commits []Commit `json:"commits"`
}
type Commit struct {
	Message string `json:"message"`
}

func GetGitHubEvents(username string) ([]GitHubEvent, error) {
	res, err := http.Get(fmt.Sprintf("https://api.github.com/users/%s/events", username))
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	if res.StatusCode == 404 {
		return nil, fmt.Errorf("User not found")
	}

	if res.StatusCode != http.StatusOK {

		return nil, fmt.Errorf("Error fetching events:", res.StatusCode)
	}

	var events []GitHubEvent
	err = json.NewDecoder(res.Body).Decode(&events)
	if err != nil {
		return nil, fmt.Errorf("Error decoding JSON:", err)
	}

	return events, nil
}

func PrintEvents(events []GitHubEvent) {
	for _, event := range events {
		var res string
		switch event.Type {
		case "PushEvent":
			commits := len(event.Payload.Commits)
			res = fmt.Sprintf("Pushed %d commits to %s", commits, event.Repo)

		case "CreateEvent":
			res = fmt.Sprintf("Created %sin %s", event.Payload.RefType, event.Repo.Name)
		case "WatchEvent":
			res = fmt.Sprintf("Starred %s", event.Repo.Name)
		case "ForkEvent":
			res = fmt.Sprintf("Forked %s", event.Repo.Name)
		default:
			res = fmt.Sprintf("%s in %s", event.Type, event.Repo.Name)
		}
		fmt.Printf("- %s\n", res)
	}

}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Please provide a github username.")
		return
	}
	username := os.Args[1]
	events, err := GetGitHubEvents(username)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	PrintEvents(events)

}
