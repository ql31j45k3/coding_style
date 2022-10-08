package usecase

import (
	"layout_2/internal/domain/student"
)

type studentUseCase struct {
	studentRepo student.StudentRepository
}

func NewStudentUseCase(studentRepo student.StudentRepository) student.StudentUseCase {
	return &studentUseCase{
		studentRepo: studentRepo,
	}
}

func (suc *studentUseCase) Create(student student.Student) (uint, error) {
	return suc.studentRepo.Create(student)
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
