package service

import (
	"time"

	"github.com/ljcbaby/form/database"
	"github.com/ljcbaby/form/model"
)

type FormService struct{}

func (s *FormService) CheckFormExist(form model.Form) (bool, error) {
	db := database.DB

	query := db.Where("id = ?", form.ID).Where("status != ?", 3)

	if form.OwnerID != 0 {
		query = query.Where("owner_id = ?", form.OwnerID)
	}

	if err := query.First(&form).Error; err != nil {
		if err.Error() == "record not found" {
			return false, nil
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

func (s *FormService) GetFormListCount(uid int64) (count int, err error) {
	db := database.DB

	var t int64
	if err := db.Model(&model.Form{}).Where("owner_id = ?", uid).Where("status != ?", 3).
		Count(&t).Error; err != nil {
		return 0, err
	}

	count = int(t)
	return count, nil
}

func (s *FormService) GetFormList(uid int64, page int, size int, forms *[]model.FormBase) error {
	db := database.DB

	if err := db.Where("owner_id = ?", uid).Offset((page - 1) * size).
		Limit(size).Find(forms).Error; err != nil {
		return err
	}

	return nil
}

func (s *FormService) GetFormDetail(form *model.Form) error {
	db := database.DB

	if err := db.Where("id = ?", form.ID).First(form).Error; err != nil {
		return err
	}

	form.IsPublish = int64(form.Status) - 1

	return nil
}

func (s *FormService) UpdateForm(form *model.Form) error {
	db := database.DB

	r := db.Model(&form).Updates(form)
	if r.Error != nil {
		return r.Error
	}

	form.IsPublish = int64(form.Status) - 1

	return nil
}

func (s *FormService) DeleteForm(id int64) error {
	db := database.DB

	result := db.Model(&model.Form{}).Where("id = ?", id).
		Update("status", 3).Update("modified_at", time.Now())
	if result.Error != nil {
		return result.Error
	}

	return nil
}

func (s *FormService) SubmitForm(id int64, result model.Result) error {
	db := database.DB

	result.FormID = id
	if err := db.Create(&result).Error; err != nil {
		return err
	}

	return nil
}

func (s *FormService) GetFormResults() {}

func (s *FormService) DuplicateForm(form model.Form) (id int64, err error) {
	db := database.DB

	s.GetFormDetail(&form)
	form.ID = 0
	form.Status = 1
	form.Title = form.Title + " - 副本"
	if err := db.Create(&form).Error; err != nil {
		return 0, err
	}

	return form.ID, nil
}
