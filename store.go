package main

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"sync"
)

// Simple file-based persistent store for completed games and leaderboard

type FileStore struct {
	gamesPath string
	lbPath    string
	mu        sync.Mutex
}

func NewFileStore(gamesPath, lbPath string) *FileStore {
	// ensure files exist
	os.MkdirAll("data", 0755)
	if _, err := os.Stat(gamesPath); os.IsNotExist(err) {
		ioutil.WriteFile(gamesPath, []byte("[]"), 0644)
	}
	if _, err := os.Stat(lbPath); os.IsNotExist(err) {
		ioutil.WriteFile(lbPath, []byte("{}"), 0644)
	}
	return &FileStore{gamesPath: gamesPath, lbPath: lbPath}
}

func (s *FileStore) AppendGame(rec GameRecord) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	bs, _ := ioutil.ReadFile(s.gamesPath)
	var arr []GameRecord
	json.Unmarshal(bs, &arr)
	arr = append(arr, rec)
	b2, _ := json.MarshalIndent(arr, "", "  ")
	return ioutil.WriteFile(s.gamesPath, b2, 0644)
}

func (s *FileStore) LoadLeaderboard() Leaderboard {
	s.mu.Lock()
	defer s.mu.Unlock()
	bs, _ := ioutil.ReadFile(s.lbPath)
	var lb Leaderboard
	json.Unmarshal(bs, &lb)
	if lb == nil { lb = Leaderboard{} }
	return lb
}

func (s *FileStore) IncrementWinner(username string) error {
	if username == "draw" || username == "" { return nil }
	s.mu.Lock()
	defer s.mu.Unlock()
	bs, _ := ioutil.ReadFile(s.lbPath)
	var lb Leaderboard
	json.Unmarshal(bs, &lb)
	if lb == nil { lb = Leaderboard{} }
	lb[username] = lb[username] + 1
	b2, _ := json.MarshalIndent(lb, "", "  ")
	return ioutil.WriteFile(s.lbPath, b2, 0644)
}
