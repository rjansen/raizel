package cassandra

type repository struct{}

func NewRepository() *repository {
	return new(repository)
}
