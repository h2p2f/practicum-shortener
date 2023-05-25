package config

type ServerConfig struct {
	serverAddress string
	resultAddress string
	useFile       bool
	useDB         bool
}

func NewServerConfig() *ServerConfig {
	return &ServerConfig{
		serverAddress: "",
		resultAddress: "",
		useFile:       false,
		useDB:         false,
	}
}

func (c *ServerConfig) SetConfig(s, r string, d, f bool) {
	c.serverAddress = s
	c.resultAddress = r
	c.useFile = f
	c.useDB = d
}

func (c *ServerConfig) GetConfig() (string, string, bool, bool) {
	return c.serverAddress, c.resultAddress, c.useFile, c.useDB
}

func (c *ServerConfig) GetResultAddress() string {
	return c.resultAddress
}

func (c *ServerConfig) UseDB() bool {
	return c.useDB
}

func (c *ServerConfig) UseFile() bool {
	return c.useFile
}
