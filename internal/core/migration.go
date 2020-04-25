package core

type Migration struct {
	name            string
	query           string
	reversQuery     string
	transactionMode bool
}

func NewMigration(name string, query string, reversQuery string) *Migration {
	return &Migration{
		name:            name,
		query:           query,
		reversQuery:     reversQuery,
		transactionMode: true}
}

func (m Migration) TransactionMode() bool {
	return m.transactionMode
}

func (m Migration) ReversQuery() string {
	return m.reversQuery
}

func (m Migration) Query() string {
	return m.query
}

func (m Migration) Name() string {
	return m.name
}
