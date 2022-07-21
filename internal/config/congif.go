package config

// Config takes config parameters from environment, or uses default.
type Config struct {
	Server  Server
	MongoDB MongoDB
}

type Server struct {
	Port string `default:":9090"`
}

type MongoDB struct {
	URI      string `default:"mongodb://localhost:27017"`
	User     string `default:"admin"`
	Password string `default:"admin"`
	DB       string `default:"reports"`
}
