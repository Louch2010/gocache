package main

import (
	"log"

	"github.com/louch2010/gocache/server"
)

func main() {
	log.Println("启动...")
	server.Start(1334, 30)
}
