package main

import (
	"encoding/json"
	"io"
	"log"
	"os"

	"github.com/agarmirus/ds-lab02/internal/controllers"
	"github.com/agarmirus/ds-lab02/internal/database"
	"github.com/agarmirus/ds-lab02/internal/serverrors"
	"github.com/agarmirus/ds-lab02/internal/services"
)

type reservConfigDataStruct struct {
	Host    string `json:"host"`
	Port    int    `json:"port"`
	ConnStr string `json:"connDb"`
}

func readConfig(path string, configData *reservConfigDataStruct) (err error) {
	configFile, err := os.Open(path)

	if err != nil {
		return err
	}

	defer configFile.Close()

	configJSON, err := io.ReadAll(configFile)

	if err != nil {
		return err
	}

	return json.Unmarshal(configJSON, configData)
}

func buildService(configData *reservConfigDataStruct) (controller controllers.IController, err error) {
	hotelDAO := database.NewPostgresHotelDAO(configData.ConnStr)
	reservDAO := database.NewPostgresReservationDAO(configData.ConnStr)
	service := services.NewReservationService(reservDAO, hotelDAO)
	controller = controllers.NewReservationController(configData.Host, configData.Port, service)

	return controller, nil
}

func main() {
	log.Println("[INFO] Starting server...")

	file, err := os.OpenFile("/logs/reservation.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)

	if err != nil {
		log.Fatalln("[FATAL] Main. Failed to open log file: ", err)
	}

	defer file.Close()

	log.SetOutput(file)

	var configData reservConfigDataStruct
	err = readConfig(`/configs/config.json`, &configData)

	if err != nil {
		log.Fatalln("[FATAL] Main. Failed to read config file: ", err)
		panic(serverrors.ErrConfigRead)
	}

	controller, err := buildService(&configData)

	if err != nil {
		log.Fatalln("[FATAL] Main. Failed to build service: ", err)
		panic(serverrors.ErrServiceBuild)
	}

	err = controller.Prepare()

	if err != nil {
		log.Fatalln("[FATAL] Main. Failed to prepare API: ", err)
		panic(serverrors.ErrControllerPrepare)
	}

	log.Println("[INFO] Running server...")
	controller.Run()
}
