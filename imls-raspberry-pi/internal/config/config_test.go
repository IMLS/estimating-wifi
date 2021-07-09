package config

import (
	"log"
	"testing"
)

func Test_Config(t *testing.T) {
	cfg := NewConfig()
	//cfg.Validate()
	log.Println(cfg)
}
