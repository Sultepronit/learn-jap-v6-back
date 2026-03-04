package models

import "encoding/json"

type CardMeta struct {
	ID    int `json:"id"`
	V     int `json:"v"`
	SyncV int `json:"syncV"`
}

type Card struct {
	// ID    int             `json:"id"`
	// V     int             `json:"v"`
	// SyncV int             `json:"syncV"`
	CardMeta
	Data json.RawMessage `json:"data"`
}
