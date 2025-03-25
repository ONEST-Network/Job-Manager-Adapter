package config

type Configuration struct {
	AllowedOrigins []string `split_words:"true" default:"(.)+.localhost:([0-9]+)?"`
	HttpProxy      string   `split_words:"true"`
	HttpsProxy     string   `split_words:"true"`
	NoProxy        string   `split_words:"true"`
	HTTPPort       string   `envconfig:"HTTP_PORT" default:"8080"`
	DbServer                 string   `required:"true" split_words:"true" default:"mongodb://localhost:27017"`
	DbUser                   string   `required:"true" split_words:"true" default:"admin"`
	DbPassword               string   `required:"true" split_words:"true" default:"secretpassword"`
	BppId                    string   `required:"true" split_words:"true" default:"bpp"`
	BppUri                   string   `required:"true" split_words:"true" default:"http://localhost:8080"`
}

var Config Configuration
