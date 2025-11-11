package cfg

import "github.com/molpeDE/spark/pkg/framework"

type Config struct {
	Server struct {
		Port uint16 `default:"3999"`
		Host string `default:"0.0.0.0"`
	}
	Production bool    `default:"false"`
	SomeFloat  float64 `default:"1.0"`
}

var inst = &Config{}

func Parse(path string) error {
	return framework.ParseConfig(path, inst)
}

func Get() Config {
	return *inst // copy
}
