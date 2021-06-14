package config

type Config interface {
	GetAWSRegion() string
	Profiles() (Profiles, error)
	SetProfile(string, Profile)
	Write(string) error
}
