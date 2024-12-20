package main

import (
	"encoding/json"
	"io"
	"log"
	"os"

	"github.com/agarmirus/ds-lab02/internal/controllers"
	"github.com/agarmirus/ds-lab02/internal/serverrors"
	"github.com/agarmirus/ds-lab02/internal/services"
)

type gatewayConfigDataStruct struct {
	Host        string `json:"host"`
	LoayltyHost string `json:"loayltyHost"`
	PaymentHost string `json:"paymentHost"`
	ReservHost  string `json:"reservHost"`
	Port        int    `json:"port"`
	LoyaltyPort int    `json:"loyaltyPort"`
	PaymentPort int    `json:"paymentPort"`
	ReservPort  int    `json:"reservPort"`
}

func readConfig(path string, configData *gatewayConfigDataStruct) (err error) {
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

func buildService(configData *gatewayConfigDataStruct) (controller controllers.IController, err error) {
	service := services.NewGatewayService(
		configData.ReservHost,
		configData.ReservPort,
		configData.PaymentHost,
		configData.PaymentPort,
		configData.LoayltyHost,
		configData.LoyaltyPort,
	)

	controller = controllers.NewGatewayController(
		configData.Host,
		configData.Port,
		service,
	)

	return controller, nil
}

func main() {
	log.Println("[INFO] Starting server...")

	file, err := os.OpenFile("/logs/gateway.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)

	if err != nil {
		log.Fatalln("[FATAL] Main. Failed to open log file: ", err)
	}

	defer file.Close()

	log.SetOutput(file)

	var configData gatewayConfigDataStruct
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
