package sync

import (
	"fmt"
	"japv6/db"
	"japv6/models"
)

var typeToArgs = map[string][]string{
	"wordCards": {"word", "card"},
	"wordProgs": {"word", "prog"},
	"kanjiCards": {"kanji", "card"},
	"kanjiProgs": {"kanji", "prog"},
}

func handleStandard(block *models.SyncBlock) (*models.SyncBlock, error) {
	var res models.SyncBlock
	res.Type = block.Type

	args := typeToArgs[block.Type]

	clientV := block.V
	// fmt.Println(args, clientV)
	lastV, err := db.GetVersion(fmt.Sprintf("%s_%ss", args[0], args[1]))
	if err != nil {
		return nil, err
	}

	newV := lastV
	if block.Updated != nil {
		res.Accepted, newV, err = db.UpsertCards(block.Updated, lastV, args[0], args[1])
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
	if inputMsg.DeletedWords != nil {
		err := db.DeleteWords(inputMsg.DeletedWords)
		if err != nil {
			return nil, err
		}
		// outputMsg.DeletedWords = inputMsg.DeletedWords
	}

	var standard []*models.SyncBlock
	for _, block := range inputMsg.Standard {
		rs, err := handleStandard(block)
		if err != nil {
			return nil, err
		}
		if rs != nil {
			standard = append(standard, rs)
		}
	}

	outputMsg := models.Message{Standard: standard}

	// return standard, nil
	return &outputMsg, nil
}
