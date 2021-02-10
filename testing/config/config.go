package config

func NewConfig() Config {
	return varConfig
}

// Config
type Config interface {
	Server() Server
	Retry() int
}
type objConfig struct {
	varServer Server `yaml:"server"`
	varRetry  int    `yaml:"retry"`
}

var varConfig = &objConfig{
	varServer: varServer,
	varRetry:  1000,
}

func (m *objConfig) Server() Server {
	return m.varServer
}
func (m *objConfig) Retry() int {
	return m.varRetry
}

// Server
type Server interface {
	Address() string
	Port() int
}
type objServer struct {
	varAddress string `yaml:"address"`
	varPort    int    `yaml:"port"`
}

var varServer = &objServer{
	varAddress: "google.com",
	varPort:    8888,
}

func (m *objServer) Address() string {
	return m.varAddress
}
func (m *objServer) Port() int {
	return m.varPort
}
