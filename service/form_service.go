package service

import (
	"errors"

	"github.com/ljcbaby/form/database"
	"github.com/ljcbaby/form/model"
)

type FormService struct{}

func (s *FormService) CheckFormNameExist(id int64) (bool, error) {
	db := database.DB

	var form model.Form
	if err := db.Where("id = ?", id).Where("status != ?", 3).First(&form).Error; err != nil {
		if err.Error() == "record not found" {
			return false, errors.New("FORM_NOT_FOUND")
		}
		return false, err
	}

	return true, nil
}

func (s *FormService) CreateForm(uid int64) (id int64, err error) {
	db := database.DB

	form := model.Form{
		OwnerID:    uid,
		Title:      "未命名表单",
		Status:     1,
		Components: []byte(`[]`),
	}
	if err := db.Create(&form).Error; err != nil {
		return 0, err
	}

	return form.ID, nil
}

func (s *FormService) GetFormList() {}

func (s *FormService) GetFormDetail(form *model.Form) error {
	db := database.DB

	if err := db.Where("id = ?", form.ID).First(form).Error; err != nil {
		return err
	}

	return nil
}

func (s *FormService) UpdateForm(form model.Form) {

}

func (s *FormService) DeleteForm(id int64) error {
	db := database.DB

	result := db.Model(&model.Form{}).Where("id = ?", id).Update("status", 3)
	if result.Error != nil {
		return result.Error
	}

	return nil
}

func (s *FormService) SubmitForm() {}

func (s *FormService) GetFormResults() {}

func (s *FormService) DuplicateForm() {}
