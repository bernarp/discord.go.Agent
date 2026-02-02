package template

type Config struct {
	Enabled    *bool      `yaml:"enabled" validate:"required"`
	LogDetails LogDetails `yaml:"logDetails" validate:"required"`
}

type LogDetails struct {
	Guild   bool `yaml:"guild"`
	Channel bool `yaml:"channel"`
	Author  bool `yaml:"author"`
	Content bool `yaml:"content"`
}
