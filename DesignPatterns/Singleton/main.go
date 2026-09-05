package main

import (
	"fmt"
	"sync"
)

type Config struct {
	AppName string
}

var configInstance *Config
var once sync.Once

func GetConfig() *Config {
	once.Do(func() {
		configInstance = &Config{
			AppName: "myapp",
		}
	})
	return configInstance
}

func main() {
	c1 := GetConfig()
	c2 := GetConfig()

	fmt.Println(c1 == c2)
}