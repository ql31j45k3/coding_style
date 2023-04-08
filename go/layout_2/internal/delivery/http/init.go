package http

import (
	"fmt"
	httpStudent "layout_2/internal/delivery/http/student"
	"layout_2/internal/libs/container"
	repoStudent "layout_2/internal/repository/student"
	userCaseStudent "layout_2/internal/usecase/student"
)

func Init() error {
	if err := container.Get().Provide(repoStudent.NewStudentRepository); err != nil {
		return fmt.Errorf("container.Provide(repository.NewStudentRepository), err: %w", err)
	}

	if err := container.Get().Provide(userCaseStudent.NewStudentUseCase); err != nil {
		return fmt.Errorf("container.Provide(usecase.NewStudentUseCase), err: %w", err)
	}

	if err := container.Get().Invoke(func(cond httpStudent.StudentHandlerCond) {
		httpStudent.RegisterRouter(cond)
	}); err != nil {
		return err
	}

	return nil
}
