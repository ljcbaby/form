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
	Status      int       `gorm:"column:status" json:"-"` // 1: 未发布 2: 已发布 3: 已删除
	IsPublish   int64     `gorm:"-" json:"isPublish"`
	AnswerCount int64     `gorm:"-" json:"answerCount"`
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
