package controllers

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/agarmirus/ds-lab02/internal/models"
	"github.com/agarmirus/ds-lab02/internal/services"
	"github.com/google/uuid"
)

type GatewayController struct {
	host string
	port int

	service services.IGatewayService
}

func NewGatewayController(
	host string,
	port int,
	service services.IGatewayService,
) IController {
	return &GatewayController{host, port, service}
}

func (controller *GatewayController) handleAllHotelsGet(res http.ResponseWriter, req *http.Request) {
	log.Println("[INFO] GatewayController.handleAllHotelsGet. Handling hotels GET request")

	page, pageParseErr := strconv.Atoi(req.FormValue(`page`))
	pageSize, pageSizeParseErr := strconv.Atoi(req.FormValue(`size`))

	if pageParseErr != nil || pageSizeParseErr != nil || page <= 0 || pageSize <= 0 {
		log.Println("[ERROR] GatewayController.handleAllHotelsGet. Invalid URL parameters")
		res.WriteHeader(http.StatusBadRequest)
		return
	}

	pagRes, err := controller.service.ReadAllHotels(page, pageSize)

	if err != nil {
		log.Println("[ERROR] GatewayController.handleAllHotelsGet. service.ReadAllHotels returned error: ", err)
		res.WriteHeader(http.StatusInternalServerError)
		return
	}

	if len(pagRes.Items) == 0 {
		res.WriteHeader(http.StatusNoContent)
		return
	}

	var pageResJSON []byte
	pageResJSON, err = json.Marshal(pagRes)

	if err == nil {
		res.Header().Add(`Content-Type`, `application/json`)
		res.WriteHeader(http.StatusOK)
		res.Write(pageResJSON)
	} else {
		log.Println("[ERROR] GatewayController.handleAllHotelsGet. Cannot convert result into JSON format: ", err)
		res.WriteHeader(http.StatusInternalServerError)
	}
}

func (controller *GatewayController) handleUserInfoGet(res http.ResponseWriter, req *http.Request) {
	log.Println("[INFO] GatewayController.handleUserInfoGet. Handling user info GET request")

	username := req.Header.Get(`X-User-Name`)

	if strings.Trim(username, ` `) == `` {
		log.Println("[ERROR] GatewayController.handleUserInfoGet. Invalid username: " + username)
		res.WriteHeader(http.StatusBadRequest)
		return
	}

	userInfoRes, err := controller.service.ReadUserInfo(username)

	if err != nil {
		log.Println("[ERROR] GatewayController.handleUserInfoGet. service.ReadUserInfo returned error: ", err)
		res.WriteHeader(http.StatusInternalServerError)
		return
	}

	var userInfoResJSON []byte
	userInfoResJSON, err = json.Marshal(userInfoRes)

	if err == nil {
		res.Header().Add(`Content-Type`, `application/json`)
		res.WriteHeader(http.StatusOK)
		res.Write(userInfoResJSON)
	} else {
		log.Println("[ERROR] GatewayController.handleUserInfoGet. Cannot convert result into JSON format: ", err)
		res.WriteHeader(http.StatusInternalServerError)
	}
}

func (controller *GatewayController) handleUserReservationsGet(res http.ResponseWriter, req *http.Request) {
	log.Println("[INFO] GatewayController.handleUserReservationsGet. Handling user reservations GET request")

	username := req.Header.Get(`X-User-Name`)

	if strings.Trim(username, ` `) == `` {
		log.Println("[ERROR] GatewayController.handleUserReservationsGet. Invalid username: " + username)
		res.WriteHeader(http.StatusBadRequest)
		return
	}

	reservsResSlice, err := controller.service.ReadUserReservations(username)

	if err != nil {
		log.Println("[ERROR] GatewayController.handleUserReservationsGet. service.ReadUserReservations returned error: ", err)
		res.WriteHeader(http.StatusInternalServerError)
		return
	}

	var reservsResSliceJSON []byte
	reservsResSliceJSON, err = json.Marshal(reservsResSlice)

	if err == nil {
		log.Println("[TRACE] GatewayController.handleUserReservationsGet. Answer:", string(reservsResSliceJSON))
		res.Header().Add(`Content-Type`, `application/json`)
		res.WriteHeader(http.StatusOK)
		res.Write(reservsResSliceJSON)
	} else {
		log.Println("[ERROR] GatewayController.handleUserReservationsGet. Cannot convert result into JSON format: ", err)
		res.WriteHeader(http.StatusInternalServerError)
	}
}

func (controller *GatewayController) handleNewReservationPost(res http.ResponseWriter, req *http.Request) {
	log.Println("[INFO] GatewayController.handleNewReservationPost. Handling reservation POST request")

	username := req.Header.Get(`X-User-Name`)

	if strings.Trim(username, ` `) == `` {
		log.Println("[ERROR] GatewayController.handleNewReservationPost. Invalid username: " + username)
		res.WriteHeader(http.StatusBadRequest)
		return
	}

	reqBody, err := io.ReadAll(req.Body)

	if err != nil {
		log.Println("[ERROR] GatewayController.handleNewReservationPost. Error while reading request body: ", err)
		res.WriteHeader(http.StatusBadRequest)
		return
	}

	var crReservReq models.CreateReservationRequest
	err = json.Unmarshal(reqBody, &crReservReq)

	if err != nil {
		log.Println("[ERROR] GatewayController.handleNewReservationPost. Error while parsing JSON request body: ", err)
		res.WriteHeader(http.StatusInternalServerError)
		return
	}

	var validErrRes models.ValidationErrorResponse
	validErrRes, err = models.ValidateCrReservReq(&crReservReq)

	if err == nil {
		var crReservRes models.CreateReservationResponse
		crReservRes, err = controller.service.CreateReservation(username, &crReservReq)

		if err != nil {
			log.Println("[ERROR] GatewayController.handleNewReservationPost. service.CreateReservation returned error: ", err)
			res.WriteHeader(http.StatusInternalServerError)
			return
		}

		var crReservResJSON []byte
		crReservResJSON, err = json.Marshal(crReservRes)

		if err == nil {
			res.Header().Add(`Content-Type`, `application/json`)
			res.WriteHeader(http.StatusOK)
			res.Write(crReservResJSON)
		} else {
			log.Println("[ERROR] GatewayController.handleNewReservationPost. Cannot convert result into JSON format: ", err)
			res.WriteHeader(http.StatusInternalServerError)
		}
	} else {
		log.Println("[ERROR] GatewayController.handleNewReservationPost. Invalid create reservation request:", err)
		var validErrResJSON []byte
		validErrResJSON, err = json.Marshal(validErrRes)

		if err == nil {
			res.Header().Add(`Content-Type`, `application/json`)
			res.WriteHeader(http.StatusBadRequest)
			res.Write(validErrResJSON)
		} else {
			log.Println("[ERROR] GatewayController.handleNewReservationPost. Cannot convert error result into JSON format: ", err)
			res.WriteHeader(http.StatusInternalServerError)
		}
	}
}

func (controller *GatewayController) handleSingleReservationGet(res http.ResponseWriter, req *http.Request) {
	log.Println("[INFO] GatewayController.handleSingleReservationGet. Handling single reservation GET request")

	reservationUid := req.PathValue("reservationUid")
	username := req.Header.Get(`X-User-Name`)

	if strings.Trim(username, ` `) == `` || uuid.Validate(reservationUid) != nil {
		log.Println("[ERROR] GatewayController.handleSingleReservationGet. Invalid parameters or headers")
		res.WriteHeader(http.StatusBadRequest)
		return
	}

	reservRes, err := controller.service.ReadReservation(reservationUid, username)

	if err != nil {
		log.Println("[ERROR] GatewayController.handleSingleReservationGet. service.ReadReservation returned error: ", err)
		res.WriteHeader(http.StatusInternalServerError)
		return
	}

	reservResJSON, err := json.Marshal(reservRes)

	if err != nil {
		log.Println("[ERROR] GatewayController.handleSingleReservationGet. Cannot convert result into JSON format: ", err)
		res.WriteHeader(http.StatusInternalServerError)
		return
	}

	res.Header().Add(`Content-Type`, `application/json`)
	res.WriteHeader(http.StatusOK)
	res.Write(reservResJSON)
}

func (controller *GatewayController) handleSingleReservationDelete(res http.ResponseWriter, req *http.Request) {
	log.Println("[INFO] GatewayController.handleSingleReservationDelete. Handling reservation DELETE request")

	reservationUid := req.PathValue("reservationUid")
	username := req.Header.Get(`X-User-Name`)

	if strings.Trim(username, ` `) == `` || uuid.Validate(reservationUid) != nil {
		log.Println("[ERROR] GatewayController.handleSingleReservationDelete. Invalid parameters or headers")
		res.WriteHeader(http.StatusBadRequest)
		return
	}

	err := controller.service.DeleteReservation(reservationUid, username)

	if err != nil {
		log.Println("[ERROR] GatewayController.handleSingleReservationDelete. service.DeleteReservation returned error: ", err)
		res.WriteHeader(http.StatusInternalServerError)
		return
	}

	res.WriteHeader(http.StatusNoContent)
}

func (controller *GatewayController) handleLoyaltyGet(res http.ResponseWriter, req *http.Request) {
	log.Println("[INFO] GatewayController.handleLoyaltyGet. Handling loyalty GET request")

	username := req.Header.Get(`X-User-Name`)

	if strings.Trim(username, ` `) == `` {
		log.Println("[ERROR] GatewayController.handleLoyaltyGet. Invalid username: " + username)
		res.WriteHeader(http.StatusBadRequest)
		return
	}

	loyaltyInfoRes, err := controller.service.ReadUserLoyalty(username)

	if err != nil {
		log.Println("[ERROR] GatewayController.handleLoyaltyGet. service.ReadUserLoyalty returned error: ", err)
		res.WriteHeader(http.StatusInternalServerError)
		return
	}

	loyaltyInfoResJSON, err := json.Marshal(loyaltyInfoRes)

	if err != nil {
		log.Println("[ERROR] GatewayController.handleLoyaltyGet. Cannot convert result into JSON format: ", err)
		res.WriteHeader(http.StatusInternalServerError)
		return
	}

	res.Header().Add(`Content-Type`, `application/json`)
	res.WriteHeader(http.StatusOK)
	res.Write(loyaltyInfoResJSON)
}

func (controller *GatewayController) handleHotelsRequest(res http.ResponseWriter, req *http.Request) {
	if req.Method == `GET` {
		log.Println("[INFO] GatewayController.handleHotelsRequest. Got hotels GET request")
		controller.handleAllHotelsGet(res, req)
	} else {
		log.Println("[ERROR] GatewayController.handleHotelsRequest. Method not allowed")
		res.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func (controller *GatewayController) handleUserRequest(res http.ResponseWriter, req *http.Request) {
	if req.Method == `GET` {
		log.Println("[INFO] GatewayController.handleUserRequest. Got user info GET request")
		controller.handleUserInfoGet(res, req)
	} else {
		log.Println("[ERROR] GatewayController.handleUserRequest. Method not allowed")
		res.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func (controller *GatewayController) handleReservationsRequest(res http.ResponseWriter, req *http.Request) {
	if req.Method == `GET` {
		log.Println("[INFO] GatewayController.handleReservationsRequest. Got user reservations GET request")
		controller.handleUserReservationsGet(res, req)
	} else if req.Method == `POST` {
		log.Println("[INFO] GatewayController.handleReservationsRequest. Got user reservation POST request")
		controller.handleNewReservationPost(res, req)
	} else {
		log.Println("[ERROR] GatewayController.handleReservationsRequest. Method not allowed")
		res.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func (controller *GatewayController) handleSingleReservationRequest(res http.ResponseWriter, req *http.Request) {
	if req.Method == `GET` {
		log.Println("[INFO] GatewayController.handleSingleReservationRequest. Got single reservation GET request")
		controller.handleSingleReservationGet(res, req)
	} else if req.Method == `DELETE` {
		log.Println("[INFO] GatewayController.handleSingleReservationRequest. Got reservation DELETE request")
		controller.handleSingleReservationDelete(res, req)
	} else {
		log.Println("[ERROR] GatewayController.handleSingleReservationRequest. Method not allowed")
		res.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func (controller *GatewayController) handleLoyaltyRequest(res http.ResponseWriter, req *http.Request) {
	if req.Method == `GET` {
		log.Println("[INFO] GatewayController.handleLoyaltyRequest. Got loyalty GET request")
		controller.handleLoyaltyGet(res, req)
	} else {
		log.Println("[ERROR] GatewayController.handleLoyaltyRequest. Method not allowed")
		res.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func (controller *GatewayController) handleHealthRequest(res http.ResponseWriter, req *http.Request) {
	if req.Method == `GET` {
		log.Println("[INFO] GatewayController.handleHealthRequest. Got health GET request")
		res.WriteHeader(http.StatusOK)
	} else {
		log.Println("[ERROR] GatewayController.handleHealthRequest. Method not allowed")
		res.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func (controller *GatewayController) Prepare() error {
	http.HandleFunc(`/api/v1/hotels`, controller.handleHotelsRequest)
	http.HandleFunc(`/api/v1/me`, controller.handleUserRequest)
	http.HandleFunc(`/api/v1/reservations`, controller.handleReservationsRequest)
	http.HandleFunc(`/api/v1/reservations/{reservationUid}`, controller.handleSingleReservationRequest)
	http.HandleFunc(`/api/v1/loyalty`, controller.handleLoyaltyRequest)

	http.HandleFunc(`/manage/health`, controller.handleHealthRequest)

	return nil
}

func (controller *GatewayController) Run() error {
	return http.ListenAndServe(
		fmt.Sprintf(
			`%s:%d`,
			controller.host,
			controller.port,
		),
		nil,
	)
}
