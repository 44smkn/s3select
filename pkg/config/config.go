package config

type Config interface {
	GetAWSRegion() string
	Profiles() (*Profiles, error)
}
