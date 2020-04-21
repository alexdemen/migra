package driver

type Storage interface {
	Ping() error
	InitScheme() error
}
