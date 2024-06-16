package service

import (
	"context"
	"encoding/json"
	"strings"
	"time"

	"github.com/ljcbaby/form/config"
	"github.com/ljcbaby/form/database"
	"github.com/ljcbaby/form/model"
	openai "github.com/sashabaranov/go-openai"
)

type FormService struct{}

func (s *FormService) CheckFormExist(form model.Form) (bool, error) {
	db := database.DB

	query := db.Where("id = ?", form.ID)

	if form.OwnerID != 0 {
		query = query.Where("owner_id = ?", form.OwnerID).Where("status != ?", 3)
	} else {
		query = query.Where("status = ?", 2)
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
		Update("status", 3).Update("modifiedAt", time.Now())
	if result.Error != nil {
		return result.Error
	}

	return nil
}

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

func (s *FormService) OpenAIGenerate(question string) (Components json.RawMessage, err error) {
	conf := config.Conf.OpenAI

	llmConf := openai.DefaultConfig(conf.ApiKey)
	llmConf.BaseURL = conf.BaseURL
	c := openai.NewClientWithConfig(llmConf)

	resp, err := c.CreateChatCompletion(
		context.Background(),
		openai.ChatCompletionRequest{
			Model: conf.Model,
			Messages: []openai.ChatCompletionMessage{
				{
					Role: openai.ChatMessageRoleUser,
					Content: `[\n` +
						`{"type": "questionInfo", "fe_id": "c1", "props": { "desc": "问卷描述", "title": "问卷标题"}, "title": "问卷信息"},\n` +
						`{"type": "questionTitle", "fe_id": "c2", "props": {"text": "一行标题", "level": 1, "isCenter": false}, "title": "标题"},\n` +
						`{"type": "questionParagraph", "fe_id": "c3", "props": {"text": "一行段落", "isCenter": false}, "title": "段落"},\n` +
						`{"type": "questionInput", "fe_id": "c3", "props": {"title": "输入框标题", "placeholder": "请输入..."}, "title": "输入框"},\n` +
						`{"type": "questionTextarea", "fe_id": "c4", "props": {"title": "输入框标题", "placeholder": "请输入..."}, "title": "多行输入"},\n` +
						`{"type": "questionRadio", "fe_id": "c5", "props": {"title": "单选标题", "value": "", "options": [{"text": "选项1", "value": "item1"}, {"text": "选项2", "value": "item2"}, {"text": "选项3", "value": "item3"}], "isVertical": false}, "title": "单选"},\n` +
						`{"type": "questionCheckbox", "fe_id": "c7", "props": {"list": [{"text": "选项1", "value": "item1", "checked": false}, {"text": "选项2", "value": "item2", "checked": false},{"text": "选项3", "value": "item3", "checked": false}], "title": "多选标题", "isVertical": false}, "title": "多选"}]\n` +
						`以上为一个问卷中的所有组件格式,其中 fe_id 的值不可重复且不与type有关联性,type 的七种类型不可更改 \n` +
						`按照如上格式给出收集 ` + question + ` 的问卷，只给出问卷内容，不需要给出任何解释`,
				},
			},
		},
	)

	if err != nil {
		return nil, err
	}

	res := resp.Choices[0].Message.Content

	if strings.HasPrefix(res, "```json") {
		res = strings.TrimPrefix(res, "```json")
		res = strings.TrimSuffix(res, "```")
	}

	content := json.RawMessage(res)
	return content, nil
}
