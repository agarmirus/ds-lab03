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

type loyaltyConfigDataStruct struct {
	Host    string `json:"host"`
	Port    int    `json:"port"`
	ConnStr string `json:"connDb"`
}

func readConfig(path string, configData *loyaltyConfigDataStruct) (err error) {
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

func buildService(configData *loyaltyConfigDataStruct) (controller controllers.IController, err error) {
	loyaltyDAO := database.NewPostgresLoyaltyDAO(configData.ConnStr)
	service := services.NewLoyaltyService(loyaltyDAO)
	controller = controllers.NewLoyaltyController(configData.Host, configData.Port, service)

	return controller, nil
}

func main() {
	log.Println("[INFO] Starting server...")

	file, err := os.OpenFile("/logs/loyalty.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)

	if err != nil {
		log.Fatalln("[FATAL] Main. Failed to open log file: ", err)
	}

	defer file.Close()

	log.SetOutput(file)

	var configData loyaltyConfigDataStruct
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