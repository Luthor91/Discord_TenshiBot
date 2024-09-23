package main

import (
	"github.com/Luthor91/Tenshi/config"
	"github.com/Luthor91/Tenshi/core"
	"github.com/Luthor91/Tenshi/database"
)

func main() {
	//go utils.PrintMemoryUsage(10)

	config.LoadConfig()
	database.InitDatabase()
	config.CheckConfig()
	core.Run()

}
