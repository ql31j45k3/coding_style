package member

func newRepositoryMember() repositoryMember {
	return &memberMongo{}
}

type repositoryMember interface {
	GetMember() Members
}

type memberMongo struct {
}

func (mm *memberMongo) GetMember() Members {
	return Members{
		Account: "member001",
	}
}
