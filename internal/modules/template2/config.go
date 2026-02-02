package template2

type Config struct {
	Prefix  string `yaml:"prefix" validate:"required"`
	Enabled bool   `yaml:"enabled"`
	MaxLogs int    `yaml:"max_logs" validate:"gte=0,lte=1000"`
}
