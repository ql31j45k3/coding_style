package student

import (
	"layout_2/internal/domain"
	"layout_2/internal/domain/entity"
	"layout_2/internal/domain/repository"
	"layout_2/internal/domain/usecase"
	"layout_2/internal/domain/vo"
	"layout_2/internal/utils"
)

type studentUseCase struct {
	studentRepo repository.StudentRepository
}

func NewStudentUseCase(studentRepo repository.StudentRepository) usecase.StudentUseCase {
	return &studentUseCase{
		studentRepo: studentRepo,
	}
}

func (suc *studentUseCase) Create(s entity.Student) (uint, error) {
	if utils.IsEmpty(s.NickName) {
		return 0, domain.ErrNickNameIsEmpty
	}

	return suc.studentRepo.Create(s)
}

func (suc *studentUseCase) UpdateID(cond vo.StudentCond, student entity.Student) error {
	return suc.studentRepo.UpdateID(cond, student)
}

func (suc *studentUseCase) GetID(cond vo.StudentCond) (entity.Student, error) {
	return suc.studentRepo.GetID(cond)
}

func (suc *studentUseCase) Get(cond vo.StudentCond) ([]entity.Student, int64, error) {
	return suc.studentRepo.Get(cond)
}
