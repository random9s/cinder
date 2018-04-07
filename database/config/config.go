package config

//Config ...
type Config interface {
	Register([]byte) error
}

//Register ...
func Register(data []byte, c Config) error {
	return c.Register(data)
}
