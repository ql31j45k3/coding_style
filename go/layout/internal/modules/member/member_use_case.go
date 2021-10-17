package member

func newUseCaseMember(repositoryMember repositoryMember) UseCaseMember {
	return &member{
		repositoryMember: repositoryMember,
	}
}

type UseCaseMember interface {
	GetMember() Members
}

type member struct {
	_ struct{}

	repositoryMember
}
