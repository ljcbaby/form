package service

import (
	"errors"

	"github.com/ljcbaby/form/database"
	"github.com/ljcbaby/form/model"
)

type ResultService struct{}

func (s *ResultService) SubmitForm(id int64, result model.Result) error {
	db := database.DB

	result.FormID = id
	if err := db.Create(&result).Error; err != nil {
		return err
	}

	return nil
}

func (s *ResultService) GetFormResultsCount(fid int64) (count int, err error) {
	db := database.DB

	var status int
	if err := db.Model(&model.Form{}).Where("id = ?", fid).
		Pluck("status", &status).Error; err != nil {
		return 0, err
	}

	if status != 2 {
		return 0, errors.New("FORM_STATUS_INVALID")
	}

	var t int64
	if err := db.Model(&model.Result{}).Where("form_id = ?", fid).
		Count(&t).Error; err != nil {
		return 0, err
	}

	count = int(t)
	return count, nil
}

func (s *ResultService) GetFormResultsList(fid int64, page int, size int, results *[]model.Result) error {
	db := database.DB

	if err := db.Where("form_id = ?", fid).Offset((page - 1) * size).
		Limit(size).Find(results).Error; err != nil {
		return err
	}

	return nil
}
