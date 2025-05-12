package config

type MilvusConfig struct {
	Host string
	Port string
}

func NewMilvusConfig(host, port string) *MilvusConfig {
	return &MilvusConfig{
		Host: host,
		Port: port,
	}
}
func (c *MilvusConfig) GetHost() string {
	return c.Host
}

func (c *MilvusConfig) GetPort() string {
	return c.Port
}

func (c *MilvusConfig) GetURL() string {
	return c.Host + ":" + c.Port
}
