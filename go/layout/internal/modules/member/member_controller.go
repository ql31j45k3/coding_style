package member

import (
	"fmt"

	"go.uber.org/dig"
)

func RegisterContainer(container *dig.Container) error {
	if err := container.Provide(newRepositoryMember); err != nil {
		return fmt.Errorf("container.Provide(newRepositoryMember) - %w", err)
	}

	if err := container.Provide(newUseCaseMember); err != nil {
		return fmt.Errorf("container.Provide(newUseCaseMember) - %w", err)
	}

	return nil
}
