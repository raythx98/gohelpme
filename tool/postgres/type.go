package postgres

// ConfigProvider is the interface for providing the db configurations.
type ConfigProvider interface {
	GetDbUsername() string
	GetDbPassword() string
	GetDbHost() string
	GetDbPort() int
	GetDbDefaultName() string
}
