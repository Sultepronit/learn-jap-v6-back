package models

import "encoding/json"

type WordCard struct {
	ID    int             `json:"id"`
	V     int             `json:"v"`
	SyncV int             `json:"syncV"`
	Data  json.RawMessage `json:"data"`
}
