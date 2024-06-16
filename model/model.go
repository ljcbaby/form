package model

import (
	"encoding/json"
	"time"

	"gorm.io/gorm"
)

type User struct {
	ID       int64  `gorm:"primaryKey;autoIncrement" json:"id"`
	Username string `gorm:"column:username" json:"username"`
	Password string `gorm:"column:password" json:"password"`
	Salt     string `gorm:"column:salt" json:"-"`
	Nickname string `gorm:"column:nickname" json:"nickname"`
}

type FormBase struct {
	OwnerID     int64     `gorm:"column:owner_id" json:"-"`
	ID          int64     `gorm:"primaryKey;autoIncrement" json:"id"`
	Title       string    `gorm:"column:title" json:"title"`
	IsPublish   int64     `gorm:"column:isPublish" json:"isPublish"`
	AnswerCount int64     `gorm:"column:answerCount" json:"answerCount"`
	ModifiedAt  time.Time `gorm:"column:modifiedAt" json:"modifiedAt"`
}

type Form struct {
	OwnerID     int64           `gorm:"column:owner_id" json:"-"`
	ID          int64           `gorm:"primaryKey;autoIncrement" json:"id"`
	Title       string          `gorm:"column:title" json:"title"`
	Status      int             `gorm:"column:status" json:"-"` // 1: 未发布 2: 已发布 3: 已删除
	IsPublish   int64           `gorm:"-" json:"isPublish"`
	AnswerCount int64           `gorm:"-" json:"answerCount"`
	ModifiedAt  time.Time       `gorm:"column:modifiedAt" json:"modifiedAt"`
	Components  json.RawMessage `gorm:"column:components" json:"components"`
}

func (form *Form) BeforeSave(tx *gorm.DB) (err error) {
	form.ModifiedAt = time.Now()
	return nil
}

type ResultRequest struct {
	ID     int64           `gorm:"primaryKey;autoIncrement" json:"id"`
	FormID int64           `gorm:"column:form_id" json:"-"`
	Res    json.RawMessage `gorm:"column:result;not null" json:"answerList" binding:"required"`
}

func (ResultRequest) TableName() string {
	return "results"
}

type ResultResponse struct {
	Title       string      `gorm:"column:title" json:"title"`
	AnswerCount int64       `gorm:"column:answerCount" json:"answerCount"`
	Components  []Component `json:"components"`
}

type Component struct {
	FeID  string          `gorm:"column:fe_id" json:"fe_id,omitempty"`
	Title string          `gorm:"column:title" json:"title,omitempty"`
	Type  string          `gorm:"column:type" json:"type,omitempty"`
	Props json.RawMessage `gorm:"column:props" json:"props,omitempty"`
	Value string          `gorm:"column:value" json:"value,omitempty"`
}

type ResultList struct {
	ID     int64  `gorm:"primaryKey;autoIncrement" json:"id"`
	ToView string `gorm:"-" json:"toView"`
	Res    string `gorm:"column:value" json:"-"`
}
