package sync

import (
	"japv6/db"
	"japv6/models"
)

func SyncWordCards(inputC []models.Card) ([]models.CardMeta, error) {
	// v, err := db.GetVersion("word_cards")
	// if err != nil {
	// 	return nil, err
	// }

	ids := make([]int, len(inputC))
	for i := range inputC {
		ids[i] = inputC[i].ID
	}
	// log.Println(ids)

	storedC, err := db.SelectMetaCardsByIds(ids)
	if err != nil {
		return nil, err
	}
	// log.Println(storedC)

	m := make(map[int]models.CardMeta)
	for _, c := range storedC {
		m[c.ID] = c
	}
	// log.Println(m)

	filtered := make([]models.Card, 0, len(inputC))
	for _, c := range inputC {
		sc := m[c.ID]
		if c.SyncV == sc.SyncV {
			filtered = append(filtered, c)
		} else if c.V > sc.V {
			filtered = append(filtered, c)
		}
	}

	return db.UpdateWordCards(filtered)
}
