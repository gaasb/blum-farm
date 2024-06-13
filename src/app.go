package app

import (
	"fmt"
	"log"
)

var (
	AppIsFinished = fmt.Errorf("Все игры успешно выполнены!")
)

type App struct {
	endpoints EndpointImpl
	blum      BlumDataImpl
}

func NewApp(blum BlumDataImpl) *App {
	return &App{
		blum: blum,
	}
}

func (a *App) RunApp() (err error) {

	log.Println("App is started...\n")

	if err = a.blum.UpdateUserData(); err != nil {
		log.Fatal(err)
		return
	}

	for a.blum.IsAviableGames() {
		if err = a.blum.GameAndClaim(); err != nil {
			log.Fatal(err)
		}
	}
	log.Println(AppIsFinished)
	return
}
