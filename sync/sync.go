package sync

import (
	"fmt"
	"japv6/db"
	"japv6/models"
)

var typeToArgs = map[string][]string{
	"wordCards": {"word", "card"},
	"wordProgs": {"word", "prog"},
}

func sync2(msg models.Msg) (*models.Msg, error) {
	var res models.Msg
	res.Type = msg.Type

	args := typeToArgs[msg.Type]

	clientV := msg.V
	lastV, err := db.GetVersion(fmt.Sprintf("%s_%ss", args[0], args[1]))
	if err != nil {
		return nil, err
	}

	newV := lastV
	if msg.Updated != nil {
		res.Accepted, newV, err = db.UpsertCards(msg.Updated, lastV, args[0], args[1])
		if err != nil {
			return nil, err
		}
	}

	if clientV == lastV && len(res.Accepted) == 0 {
		return nil, nil
	}
	// fmt.Println(report.Accepted)

	res.V = newV
	res.Updated, err = db.SelectCardsSyncRange(args[0], args[1], clientV+1, lastV)
	if err != nil {
		return nil, err
	}

	return &res, nil
	// return nil, nil
}

func Do(inputMsg []models.Msg) ([]*models.Msg, error) {
	outputMsg := make([]*models.Msg, 0, len(inputMsg))
	for _, m := range inputMsg {
		rs, err := sync2(m)
		if err != nil {
			return nil, err
		}
		if rs != nil {
			outputMsg = append(outputMsg, rs)
		}
	}

	return outputMsg, nil
}
