package client

type Config struct {
	Id 			string
	Name 		string
	Description string
	Version 	string
}

type ConfigFunc func(*Config)

var DefaultConfig = Config{
	Name: 			"default.client",
}
