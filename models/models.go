package models

import "encoding/json"

type CardMeta struct {
	ID    int `json:"id"`
	V     int `json:"v"`
	SyncV int `json:"syncV"`
}

type Card struct {
	CardMeta
	Data json.RawMessage `json:"data"`
}

// SyncBlock?
type SyncBlock struct {
	Type     string     `json:"type"`
	V        int        `json:"v"`
	Updated  []Card     `json:"updated,omitempty"`
	Accepted []CardMeta `json:"accepted,omitempty"`
}

type Message struct {
	Standard     []*SyncBlock `json:"standard,omitempty"`
	DeletedWords []int        `json:"deletedWords,omitempty"`
}
