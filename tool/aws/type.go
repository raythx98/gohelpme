package aws

// ConfigProvider is the interface that wraps the GetAwsAccessKey and GetAwsSecretKey methods.
type ConfigProvider interface {
	GetAwsRegion() string
	GetAwsAccessKey() string
	GetAwsSecretKey() string
}
