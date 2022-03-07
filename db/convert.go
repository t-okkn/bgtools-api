package db

import (
	"errors"

	"bgtools-api/models"
)

// <summary>: ボードゲームのデータを読み込みます
func LoadBgDataForScore(r *BgRepository) error {
	if r == nil {
		e := errors.New("DBの接続に失敗しました")
		return e
	}

	list, err := r.GetScoreSupported()
	if err != nil {
		return err
	}

	for _, data := range list {
		_, ok := models.BgScore[data.GameId]
		bgd := models.BgPartialData{}

		if !ok {
			bgd.Title = data.Title
			bgd.MinPlayers = data.MinPlayers
			bgd.MaxPlayers = data.MaxPlayers

			bgd.Colors = make([]string, 0, data.MaxPlayers)
			bgd.Colors = append(bgd.Colors, data.Color)

			models.BgScore[data.GameId] = bgd

		} else {
			d := models.BgScore[data.GameId]
			d.Colors = append(d.Colors, data.Color)

			models.BgScore[data.GameId] = d
		}
	}

	return nil
}
