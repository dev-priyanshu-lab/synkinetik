package main

import (
	"dnslookup/service"
)

func main() {
	svc, err := service.New("lookup.conf")
	if err != nil {
		panic(err)
	}

	svc.Start()

	select {}
}
