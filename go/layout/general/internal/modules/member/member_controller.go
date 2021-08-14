package member

import "sync"

var (
	once sync.Once
	m    UseCaseMember
)

func NewUseCaseMember() UseCaseMember {
	once.Do(func() {
		m = newUseCaseMember(newRepositoryMember())
	})

	return m
}
