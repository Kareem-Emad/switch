package main

import (
	"github.com/Kareem-Emad/switch/producer"
	"github.com/Kareem-Emad/switch/server"
)

func main() {
	var pm producer.ProductionManager

	pm.InitalizeFaktoryConnection()
	pm.SeedTopicSubcriptionMap()

	server.Start(pm)
}
