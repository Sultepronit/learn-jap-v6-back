package sync

import (
	"fmt"
	"japv6/db"
	"japv6/models"
)

func filterUpdates(table string, group string, inputC []models.Card, isFresh bool) ([]models.Card, error) {
	ids := make([]int, len(inputC))
	for i := range inputC {
		ids[i] = inputC[i].ID
	}

	storedC, err := db.SelectMetaCardsByIds(table, group, ids)
	if err != nil {
		return nil, err
	}

	m := make(map[int]models.CardMeta)
	for _, c := range storedC {
		m[c.ID] = c
	}
	// fmt.Println("stored:", m)

	re := make([]models.Card, 0, len(inputC))
	for _, c := range inputC {
		sc := m[c.ID]
		// fmt.Println("syncV:", c.SyncV, sc.SyncV)
		// fmt.Println("V:", c.V, sc.V)
		if c.SyncV == sc.SyncV {
			re = append(re, c)
		} else if c.V > sc.V && isFresh {
			re = append(re, c)
		}
	}

	return re, nil
}

func update(table string, group string, updates []models.Card, lastV int, clientV int) ([]models.CardMeta, int, error) {
	filtered, err := filterUpdates(table, group, updates, (clientV+100 > lastV))
	if err != nil {
		return nil, 0, err
	}
	// fmt.Println("filtered:", filtered)

	if len(filtered) == 0 {
		return nil, lastV, nil
	}

	return db.UpdateCards(filtered, lastV, table, group)
}

var typeToArgs = map[string][]string{
	"wordCards": {"word", "card"},
	"wordProgs": {"word", "prog"},
}

func sync(report models.Msg) (result *models.Msg, err error) {
	result = &report

	args := typeToArgs[report.Type]

	clientV := report.V
	lastV, err := db.GetVersion(fmt.Sprintf("%s_%ss", args[0], args[1]))
	if err != nil {
		return
	}

	newV := lastV
	if report.Updated != nil {
		report.Accepted, newV, err = update(args[0], args[1], report.Updated, lastV, clientV)
		if err != nil {
			return nil, err
		}
		// filtered, err := filterUpdates(args[0], args[1], report.Updated, (clientV + 100 > lastV))
		// if err != nil {
		// 	return nil, err
		// }
		// // fmt.Println("filtered:", filtered)

		// if len(filtered) != 0 {
		// 	report.Accepted, newV, err = db.UpdateCards(filtered, lastV, args[0], args[1])
		// 	if err != nil {
		// 		return nil, err
		// 	}
		// }
	}

	if report.Created != nil {
		err = db.TempFillCards(report.Created, args[0], args[1])

	}

	if clientV == lastV && len(report.Accepted) == 0 {
		return nil, nil
	}
	// fmt.Println(report.Accepted)

	report.V = newV
	report.Updated, err = db.SelectCardsSyncRange(args[0], args[1], clientV+1, lastV)
	if err != nil {
		return
	}

	return result, nil
}

func Do(inputR []models.Msg) ([]*models.Msg, error) {
	outputR := make([]*models.Msg, 0, len(inputR))
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
