package main

import (
	"github.com/Luthor91/Tenshi/config"
	"github.com/Luthor91/Tenshi/core"
	"github.com/Luthor91/Tenshi/database"
)

func main() {
	config.LoadConfig()
	database.InitDatabase()
	config.CheckConfig()
	core.Run()

}
