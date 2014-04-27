package config

type Config struct {
	SiteName      string
	SessionSecret string
	DB
}

type DB struct {
	Address  string
	Database string
	Tables   []string
}
