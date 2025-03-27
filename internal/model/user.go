package model

import "gorm.io/gorm"

type User struct {
	gorm.Model
	Username string `form:"username" json:"username" binding:"required" gorm:"unique;not null"` //not null表示数据库中必须含有此字段
	Password string `form:"password" json:"password" binding:"required" gorm:"not null"`        //：json:"-" 表示该字段不会被包含在 JSON 编码或解码中
	Nickname string `form:"nickname" json:"nickname" gorm:"not null"`
	Avatar   string `json:"avatar"` //头像url
}
