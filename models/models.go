package models

import "encoding/json"

type CardMeta struct {
	// ID    int `json:"id"`
	ID    any `json:"id"`
	V     int `json:"v"`
	SyncV int `json:"syncV"`
}

type Card struct {
	CardMeta
	Data json.RawMessage `json:"data"`
}

// for TempFillCards only?
// type AnyCard struct {
// 	ID    any `json:"id"`
// 	V     int `json:"v"`
// 	SyncV int `json:"syncV"`
// 	Data json.RawMessage `json:"data"`
// }

type SyncBlock struct {
	Type     string     `json:"type"`
	V        int        `json:"v"`
	ForcedUpdate bool `json:"forcedUpdate,omitempty"`
	Updated  []Card     `json:"updated,omitempty"`
	Accepted []CardMeta `json:"accepted,omitempty"`
}

type Message struct {
	Standard     []*SyncBlock `json:"standard,omitempty"`
	DeletedWords []int        `json:"deletedWords,omitempty"`
}
