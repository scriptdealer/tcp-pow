package config

type Common struct {
	Host string `json:"host"`
	Port int    `json:"port"`
}

func New() *Common {
	return &Common{
		// Host: "localhost",
		Port: 8080,
	}
}
