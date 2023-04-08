package entity

import (
	"gorm.io/gorm"
)

type Student struct {
	gorm.Model

	Name     string `json:"name" binding:"required,min=1,max=10" gorm:"column:name; type:varchar(10); NOT NULL; DEFAULT:''; comment:姓名; index:idx_name;"`
	NickName string `json:"nick_name" gorm:"column:nick_name; type:varchar(10); NOT NULL; DEFAULT:''; comment:暱稱;"`

	Gender int `json:"gender" binding:"studentGender" gorm:"column:gender; type:TINYINT(1); NOT NULL; DEFAULT:0; comment:性別 0:男, 1:女;"`
	Status int `json:"status" binding:"studentStatus" gorm:"column:status; type:TINYINT(1); NOT NULL; DEFAULT:1; comment:狀態 0:禁用, 1:啟用;"`
}

func (s *Student) TableName() string {
	return "student"
}
