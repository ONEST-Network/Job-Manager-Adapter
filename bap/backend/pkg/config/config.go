package config

type Configuration struct {
	AllowedOrigins []string `split_words:"true" default:"(.)+.localhost:([0-9]+)?"`
	HttpProxy      string   `split_words:"true"`
	HttpsProxy     string   `split_words:"true"`
	NoProxy        string   `split_words:"true"`
	HTTPPort       string   `envconfig:"HTTP_PORT" default:"5000"`
	DbServer       string   `required:"true" split_words:"true"`
	DbUser         string   `required:"true" split_words:"true"`
	DbPassword     string   `required:"true" split_words:"true"`
	BppId          string   `required:"true" split_words:"true"`
	BppUri         string   `required:"true" split_words:"true"`
	BapId          string   `required:"true" split_words:"true"`
	BapUri         string   `required:"true" split_words:"true"`
	RecommendationServiceURL string `required:"true" split_words:"true"`
}

var Config Configuration
