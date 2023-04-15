package config

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

func (c *ServerConfig) SetConfig(s, r string) {
	c.serverAddress = s
	c.resultAddress = r
}

func (c *ServerConfig) GetConfig() (string, string) {
	return c.serverAddress, c.resultAddress
}

func (c *ServerConfig) GetResultAddress() string {
	return c.resultAddress
}
