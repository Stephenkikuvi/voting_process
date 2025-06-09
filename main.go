package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sync"
)

type Candidate struct {
	Name  string `json:"name"`
	Votes int    `json:"votes"`
}

type Voter struct {
	ID string `json:"id"`
}

var (
	candidates = map[string]*Candidate{
		"1": {Name: "Candidate A", Votes: 0},
		"2": {Name: "Candidate B", Votes: 0},
		"3": {Name: "Candidate C", Votes: 0},
	}
	votedIDs = make(map[string]bool)
	mutex    sync.Mutex
)

func castVoteHandler(w http.ResponseWriter, r *http.Request) {
	var voter Voter
	err := json.NewDecoder(r.Body).Decode(&voter)
	if err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	candidateID := r.URL.Query().Get("candidate")
	if _, exists := candidates[candidateID]; !exists {
		http.Error(w, "Candidate not found", http.StatusNotFound)
		return
	}

	mutex.Lock()
	defer mutex.Unlock()

	if votedIDs[voter.ID] {
		http.Error(w, "Voter has already voted", http.StatusForbidden)
		return
	}

	candidates[candidateID].Votes++
	votedIDs[voter.ID] = true

	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "Vote cast successfully for %s", candidates[candidateID].Name)
}

func resultsHandler(w http.ResponseWriter, r *http.Request) {
	mutex.Lock()
	defer mutex.Unlock()

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(candidates)
}

func main() {
	http.HandleFunc("/cast_vote", castVoteHandler)
	http.HandleFunc("/results", resultsHandler)

	fmt.Println("Server is running on port 8080...")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
