package main

import (
	"flag"
)

func main() {
	var configFilePath string
	flag.StringVar(&configFilePath, "c", "", "config file path")
	flag.Parse()

	config, err := ReadYml(configFilePath)
	if err != nil {
		panic(err)
	}

	host := func() []string {
		var hosts []string
		for _, host := range config["hosts"].([]interface{}) {
			hosts = append(hosts, host.(string))
		}
		return hosts
	}()

	sshServer := SshServer{
		Hosts: host,
		User: &User{
			Username: config["username"].(string),
			Password: config["password"].(string),
		},
	}

	api := New(":8080", &sshServer)
	err = api.Start()

}
