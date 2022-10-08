package student

import (
	"layout_2/internal/libs/response"
	"layout_2/internal/utils"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type Student struct {
	gorm.Model

	Name string `json:"name" binding:"required,min=1,max=10" gorm:"column:name; type:varchar(10); NOT NULL; DEFAULT:''; comment:姓名; index:idx_name;"`

	Gender int `json:"gender" binding:"studentGender" gorm:"column:gender; type:TINYINT(1); NOT NULL; DEFAULT:0; comment:性別 0:男, 1:女;"`
	Status int `json:"status" binding:"studentStatus" gorm:"column:status; type:TINYINT(1); NOT NULL; DEFAULT:1; comment:狀態 0:禁用, 1:啟用;"`
}

func (s *Student) TableName() string {
	return "student"
}

type StudentRepository interface {
	Create(student Student) (uint, error)
	UpdateID(cond StudentCond, student Student) error

	GetID(cond StudentCond) (Student, error)
	Get(cond StudentCond) ([]Student, int64, error)
}

type StudentUseCase interface {
	Create(student Student) (uint, error)
	UpdateID(cond StudentCond, student Student) error

	GetID(cond StudentCond) (Student, error)
	Get(cond StudentCond) ([]Student, int64, error)
}

type StudentCond struct {
	response.Pagination

	ID uint

	Name string

	Gender int
	Status int
}

func (cond *StudentCond) ParseID(c *gin.Context) error {
	var err error
	idStr := c.Param("id")

	cond.ID, err = cond.getID(idStr)
	if err != nil {
		return err
	}

	return nil
}

func (cond *StudentCond) getID(idStr string) (uint, error) {
	if utils.IsEmpty(idStr) {
		return 0, nil
	}

	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		return 0, err
	}

	return uint(id), nil
}

func (cond *StudentCond) ParseGet(c *gin.Context) error {
	var err error
	idStr := c.Query("id")

	cond.ID, err = cond.getID(idStr)
	if err != nil {
		return err
	}

	cond.PageIndex, err = utils.AtoiAndDefaultNotAssignInt(c.Query("page_index"))
	if err != nil {
		return err
	}

	cond.PageSize, err = utils.AtoiAndDefaultNotAssignInt(c.Query("page_size"))
	if err != nil {
		return err
	}

	cond.Name = strings.TrimSpace(c.Query("name"))

	cond.Gender, err = utils.AtoiAndDefaultNotAssignInt(c.Query("gender"))
	if err != nil {
		return err
	}

	cond.Status, err = utils.AtoiAndDefaultNotAssignInt(c.Query("status"))
	if err != nil {
		return err
	}

	return nil
}
