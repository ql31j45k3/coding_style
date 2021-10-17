package dependency

import (
	"fmt"

	"go.uber.org/dig"
)

func RegisterContainerTransaction(container *dig.Container) error {
	if err := container.Provide(newRepositoryTransaction); err != nil {
		return fmt.Errorf("container.Provide(newRepositoryTransaction) - %w", err)
	}

	if err := container.Provide(newUseCaseTransaction); err != nil {
		return fmt.Errorf("container.Provide(newUseCaseTransaction) - %w", err)
	}

	return nil
}
