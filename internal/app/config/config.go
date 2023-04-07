package config

type Config interface {
	SetConfig(s, r string) ServerConfig
	GetConfig() ServerConfig
	GetResultAddress() string
}

type ServerConfig struct {
	serverAddress string
	resultAddress string
}

func NewServerConfig() *ServerConfig {
	return &ServerConfig{
		serverAddress: "",
		resultAddress: "",
	}
}

func (c *ServerConfig) SetConfig(s, r string) ServerConfig {
	c.serverAddress = s
	c.resultAddress = r
	return *c
}

func (c *ServerConfig) GetConfig() ServerConfig {
	return *c
}

func (c *ServerConfig) GetResultAddress() string {
	return c.resultAddress
}
