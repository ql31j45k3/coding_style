package member

func newRepositoryMember() repositoryMember {
	return &memberMongo{}
}

type repositoryMember interface {
	GetMember() Members
}

type memberMongo struct {
	_ struct{}
}

func (mm *memberMongo) GetMember() Members {
	return Members{
		Account: "member001",
	}
}
