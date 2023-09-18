package pkg

import "time"

type GlobalGrepRequest struct {
	ServerList []string      `json:"server_list"`
	ServerPort string        `json:"server_port"`
	Command    string        `json:"command"`
	Timeout    time.Duration `json:"timeout"`
}

type GlobalGrepResponse struct {
}

type LocalGrepRequest struct {
	Command string `json:"command"`
}

type LocalGrepResponse struct {
	Command     string   `json:"command"`
	Count       int      `json:"count"`
	MatchedText []string `json:"matched_text"`
	OnlyCount   bool     `json:"only_count"`
}
