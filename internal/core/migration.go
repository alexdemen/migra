package core

type Migration struct {
	Name            string
	Query           string
	ReversQuery     string
	TransactionMode bool
}

func NewMigration(name string, query string, reversQuery string) *Migration {
	return &Migration{
		Name:            name,
		Query:           query,
		ReversQuery:     reversQuery,
		TransactionMode: true,
	}
}
