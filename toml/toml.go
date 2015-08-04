package main

/*
http://qiita.com/futoase/items/fd697a708fcbcee104de
*/

import (
	"fmt"
	"github.com/BurntSushi/toml"
)

type Config struct {
	Server ServerConfig
	Db     DbConfig
}

type ServerConfig struct {
	Host  string        `toml:"host"`
	Port  string        `toml:"port"`
	Slave []SlaveServer `toml:"slave"`
}

type DbConfig struct {
	User string `toml:"user"`
	Pass string `toml:"pass"`
}

type SlaveServer struct {
	Weight int    `toml:"wight"`
	Ip     string `toml:"ip"`
}

func main() {
	var config Config
	_, err := toml.DecodeFile("config.tml", &config)
	if err != nil {
		panic(err)
	}
	fmt.Printf("Host is :%s\n", config.Server.Host)
	fmt.Printf("Port is :%s\n", config.Server.Port)
	for k, v := range config.Server.Slave {
		fmt.Printf("Slave %d\n", k)
		fmt.Printf("  weight is %d\n", v.Weight)
		fmt.Printf("  ip is %s\n", v.Ip)
	}
	fmt.Printf("DB Username is :%s\n", config.Db.User)
	fmt.Printf("DB Password is :%s\n", config.Db.Pass)
}
