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

func handleStandard(msg *models.Msg) (*models.Msg, error) {
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

// func Do(inputMsg []models.Msg) ([]*models.Msg, error) {
// func Do(inputMsg models.Message) ([]*models.Msg, error) {
func Do(inputMsg models.Message) (*models.Message, error) {
	// standard := make([]*models.Msg, 0, len(inputMsg.Standard))
	var standard []*models.Msg
	for _, m := range inputMsg.Standard {
		rs, err := handleStandard(m)
		if err != nil {
			return nil, err
		}
		if rs != nil {
			standard = append(standard, rs)
		}
	}

	outputMsg := models.Message{Standard: standard}
	if (inputMsg.DeletedWords != nil) {
		outputMsg.DeletedWords = inputMsg.DeletedWords
	}

	// return standard, nil
	return &outputMsg, nil
}
