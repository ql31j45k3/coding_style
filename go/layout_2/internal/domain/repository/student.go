package repository

import (
	"layout_2/internal/domain/entity"
	"layout_2/internal/domain/vo"
)

type StudentRepository interface {
	Create(student entity.Student) (uint, error)
	UpdateID(cond vo.StudentCond, student entity.Student) error

	GetID(cond vo.StudentCond) (entity.Student, error)
	Get(cond vo.StudentCond) ([]entity.Student, int64, error)
}
