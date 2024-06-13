package main

import (
	app "blum-points/src"
	"flag"
	"log"
)

var token = flag.String("token", "eyJhbGciOi.....", "токен для авторизации")

func main() {
	blum := app.NewBlumData(*token)
	blum.WithEndpoints(app.NewEndpoints())
	app.NewApp(blum).RunApp()

}

func init() {
	flag.Parse()
	if *token == "eyJhbGciOi....." || *token == "" {
		log.Fatal("неверно задан флаг -token")
	}
}
