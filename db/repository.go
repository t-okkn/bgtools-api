package db

import (
	"bgtools-api/models"

	"github.com/go-gorp/gorp"
)

type BgRepository struct {
	*gorp.DbMap
}

func NewRepository(dm *gorp.DbMap) *BgRepository {
	return &BgRepository{dm}
}

func (r *BgRepository) GetScoreSupported() ([]models.BgScoreSupport, error) {
	var result []models.BgScoreSupport
	query := GetSQL("get-score-supported-games", "")

	if _, err := r.Select(&result, query); err != nil {
		return []models.BgScoreSupport{}, err
	}

	return result, nil
}
