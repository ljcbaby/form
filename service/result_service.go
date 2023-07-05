package service

import (
	"encoding/json"
	"strings"

	"github.com/ljcbaby/form/database"
	"github.com/ljcbaby/form/model"
)

type ResultService struct{}

func (s *ResultService) CheckFormResultExist(form model.Form) (bool, error) {
	db := database.DB

	if err := db.Where("id = ?", form.ID).Where("owner_id = ?", form.OwnerID).
		Where("status = ?", 2).First(&form).Error; err != nil {
		if err.Error() == "record not found" {
			return false, nil
		}
		return false, err
	}

	return true, nil
}

func (s *ResultService) SubmitForm(id int64, result model.ResultRequest) error {
	db := database.DB

	result.FormID = id
	if err := db.Create(&result).Error; err != nil {
		return err
	}

	return nil
}

func (s *ResultService) GetFormResult(fid int64, result *model.ResultResponse) error {
	db := database.DB

	err := db.Raw("CALL GetComponents(?)", fid).Scan(&result.Components).Error
	if err != nil {
		return err
	}

	var fb model.FormBase
	if err := db.Where("id = ?", fid).First(&fb).Error; err != nil {
		return err
	}

	result.Title = fb.Title
	result.AnswerCount = fb.AnswerCount

	return nil
}

func (s *ResultService) GetFormResultsCount(fid int64) (count int, err error) {
	db := database.DB

	var t int64
	if err := db.Model(&model.ResultRequest{}).Where("form_id = ?", fid).
		Count(&t).Error; err != nil {
		return 0, err
	}

	count = int(t)
	return count, nil
}

func (s *ResultService) GetFormResultsList(fid int64, fe_id string, page int, size int, results *[]model.ResultList) error {
	db := database.DB

	if err := db.Raw("CALL GetResults(?, ?, ?, ?)", fid, fe_id, size, size*(page-1)).Find(results).Error; err != nil {
		return err
	}

	var f struct {
		T string `gorm:"column:value"`
		P string `gorm:"column:props"`
	}
	if err := db.Raw("CALL GetComponentType(? ,?)", fid, fe_id).Scan(&f).Error; err != nil {
		return err
	}

	switch f.T {
	case "questionInput", "questionTextarea":
		for i := range *results {
			(*results)[i].ToView = (*results)[i].Res
		}
	case "questionRadio":
		var t struct {
			Options []struct {
				K string `json:"value"`
				V string `json:"text"`
			} `json:"options"`
		}
		err := json.Unmarshal([]byte(f.P), &t)
		if err != nil {
			return err
		}
		for i := range *results {
			for j := range t.Options {
				if (*results)[i].Res == t.Options[j].K {
					(*results)[i].ToView = t.Options[j].V
					break
				}
			}
		}
	case "questionCheckbox":
		var t struct {
			Options []struct {
				K string `json:"value"`
				V string `json:"text"`
			} `json:"list"`
		}
		err := json.Unmarshal([]byte(f.P), &t)
		if err != nil {
			return err
		}
		for i := range *results {
			s := strings.Split((*results)[i].Res, ",")
			for k := range t.Options {
				for j := range s {
					if s[j] == t.Options[k].K {
						(*results)[i].ToView += t.Options[k].V + ", "
						break
					}
				}
			}
			(*results)[i].ToView = (*results)[i].ToView[:len((*results)[i].ToView)-2]
		}
	}

	return nil
}

func (s *ResultService) GetFormResultsDetail(fid int64, rid int64, res *[]model.Component) error {
	db := database.DB

	var str string
	err := db.Table("forms").Select("components").Where("id = ?", fid).Find(&str).Error
	if err != nil {
		return err
	}

	if err := json.Unmarshal([]byte(str), res); err != nil {
		return err
	}

	err = db.Table("results").Select("result").Where("id = ?", rid).Find(&str).Error
	if err != nil {
		return err
	}

	var ans []model.Component
	if err := json.Unmarshal([]byte(str), &ans); err != nil {
		return err
	}

	// attach value to res[x] from ans[x] where res[x].fe_id == ans[x].fe_id
	for i := range *res {
		for j := range ans {
			if (*res)[i].FeID == ans[j].FeID {
				(*res)[i].Value = ans[j].Value
				break
			}
		}
	}

	return nil
}
