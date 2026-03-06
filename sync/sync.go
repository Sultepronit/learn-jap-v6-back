package sync

import (
	"japv6/db"
	"japv6/models"
	"fmt"
)

func filterUpdates(inputC []models.Card, isFresh bool) ([]models.Card, error) {
	ids := make([]int, len(inputC))
	for i := range inputC {
		ids[i] = inputC[i].ID
	}

	storedC, err := db.SelectMetaCardsByIds(ids)
	if err != nil {
		return nil, err
	}

	m := make(map[int]models.CardMeta)
	for _, c := range storedC {
		m[c.ID] = c
	}

	re := make([]models.Card, 0, len(inputC))
	for _, c := range inputC {
		sc := m[c.ID]
		if c.SyncV == sc.SyncV {
			re = append(re, c)
		} else if c.V > sc.V && isFresh {
			re = append(re, c)
		}
	}

	return re, nil
}

// func implementUpdates(updates []models.Card, v int) ([]models.CardMeta, int, error) {
// 	filtered, err := filterUpdates(updates)
// 	if err != nil {
// 		return nil, 0, err
// 	}

// 	// return db.UpdateWordCards(filtered)
// 	return db.UpdateCards(filtered, v)
// }

var typeToArgs = map[string][]string{
	"wordCards": {"word", "card"},
	"wordProgs": {"word", "prog"},
}

func sync(report models.Report) (result *models.Report, err error) {
	result = &report

	args := typeToArgs[report.Type]

	clientV := report.V
	lastV, err := db.GetVersion(fmt.Sprintf("%s_%ss", args[0], args[1]))
	if err != nil {
		return
	}

	newV := lastV
	if (report.Updated != nil) {
		filtered, err := filterUpdates(report.Updated, (clientV + 100 > lastV))
		if err != nil {
			return nil, err
		}

		report.Accepted, newV, err = db.UpdateCards(filtered, lastV, args[0], args[1])
		if err != nil {
			return nil, err
		}
	} else if clientV == lastV {
		return nil, nil
	}

	report.V = newV
	report.Updated, err = db.SelectCardsSyncRange(args[0], args[1], clientV+1, lastV)
	if err != nil {
		return
	}

	return result, nil
}

func Do(inputR []models.Report) ([]*models.Report, error) {
	outputR := make([]*models.Report, 0, len(inputR))
	for _, r := range inputR {
		rs, err := sync(r)
		if err != nil {
			return nil, err
		}
		if rs != nil {
			outputR = append(outputR, rs)
		}
	}

	return outputR, nil
}

// func SyncWordCards(inputC []models.Card) ([]models.CardMeta, error) {
// 	// v, err := db.GetVersion("word_cards")
// 	// if err != nil {
// 	// 	return nil, err
// 	// }

// 	ids := make([]int, len(inputC))
// 	for i := range inputC {
// 		ids[i] = inputC[i].ID
// 	}
// 	// log.Println(ids)

// 	storedC, err := db.SelectMetaCardsByIds(ids)
// 	if err != nil {
// 		return nil, err
// 	}
// 	// log.Println(storedC)

// 	m := make(map[int]models.CardMeta)
// 	for _, c := range storedC {
// 		m[c.ID] = c
// 	}
// 	// log.Println(m)

// 	filtered := make([]models.Card, 0, len(inputC))
// 	for _, c := range inputC {
// 		sc := m[c.ID]
// 		if c.SyncV == sc.SyncV {
// 			filtered = append(filtered, c)
// 		} else if c.V > sc.V {
// 			filtered = append(filtered, c)
// 		}
// 	}

// 	return db.UpdateWordCards(filtered)
// }
