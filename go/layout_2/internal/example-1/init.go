package example_1

import (
	"fmt"
	studentRouter "layout_2/internal/example-1/student/delivery/http"
	studentRepo "layout_2/internal/example-1/student/repository"
	studentUseCase "layout_2/internal/example-1/student/usecase"
	"layout_2/internal/libs/container"
)

func Init() error {
	// example-1
	if err := container.Get().Provide(studentRepo.NewStudentRepository); err != nil {
		return fmt.Errorf("container.Provide(studentRepo.NewStudentRepository), err: %w", err)
	}

	if err := container.Get().Provide(studentUseCase.NewStudentUseCase); err != nil {
		return fmt.Errorf("container.Provide(studentUseCase.NewStudentUseCase), err: %w", err)
	}

	if err := container.Get().Invoke(func(cond studentRouter.StudentHandlerCond) {
		studentRouter.RegisterRouter(cond)
	}); err != nil {
		return err
	}

	return nil
}
