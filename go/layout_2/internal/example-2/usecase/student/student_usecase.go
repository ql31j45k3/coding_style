package student

import (
	"layout_2/internal/domain/student"
	"layout_2/internal/utils"
)

type studentUseCase struct {
	studentRepo student.StudentRepository `name:"NewStudentRepository2"`
}

func NewStudentUseCase(studentRepo student.StudentRepository) student.StudentUseCase {
	return &studentUseCase{
		studentRepo: studentRepo,
	}
}

func (suc *studentUseCase) Create(s student.Student) (uint, error) {
	if utils.IsEmpty(s.NickName) {
		return 0, student.ErrNickNameIsEmpty
	}

	return suc.studentRepo.Create(s)
}

func (suc *studentUseCase) UpdateID(cond student.StudentCond, student student.Student) error {
	return suc.studentRepo.UpdateID(cond, student)
}

func (suc *studentUseCase) GetID(cond student.StudentCond) (student.Student, error) {
	return suc.studentRepo.GetID(cond)
}

func (suc *studentUseCase) Get(cond student.StudentCond) ([]student.Student, int64, error) {
	return suc.studentRepo.Get(cond)
}
