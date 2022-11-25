package http

import (
	"fmt"
	httpStudent "layout_2/internal/example-2/delivery/http/student"
	repoStudent "layout_2/internal/example-2/repository/student"
	userCaseStudent "layout_2/internal/example-2/usecase/student"
	"layout_2/internal/libs/container"

	"go.uber.org/dig"
)

func Init() error {
	// example-2
	if err := container.Get().Provide(repoStudent.NewStudentRepository, dig.Name("NewStudentRepository2")); err != nil {
		return fmt.Errorf("container.Provide(repository.NewStudentRepository), err: %w", err)
	}

	if err := container.Get().Provide(userCaseStudent.NewStudentUseCase, dig.Name("NewStudentUseCase2")); err != nil {
		return fmt.Errorf("container.Provide(usecase.NewStudentUseCase), err: %w", err)
	}

	if err := container.Get().Invoke(func(cond httpStudent.StudentHandlerCond) {
		httpStudent.RegisterRouter(cond)
	}); err != nil {
		return err
	}

	return nil
}
