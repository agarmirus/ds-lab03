package controllers

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/agarmirus/ds-lab02/internal/models"
	"github.com/agarmirus/ds-lab02/internal/serverrors"
	"github.com/agarmirus/ds-lab02/internal/services"
)

type ReservationController struct {
	host string
	port int

	service services.IReservationService
}

func NewReservationController(
	host string,
	port int,
	service services.IReservationService,
) IController {
	return &ReservationController{host, port, service}
}

func (controller *ReservationController) handleAllHotelsGet(res http.ResponseWriter, req *http.Request) {
	log.Println("[INFO] ReservationController.handleAllHotelsGet. Handling hotels GET request")

	page, pageParseErr := strconv.Atoi(req.FormValue(`page`))
	pageSize, pageSizeParseErr := strconv.Atoi(req.FormValue(`size`))

	if pageParseErr != nil || pageSizeParseErr != nil || page <= 0 || pageSize <= 0 {
		log.Println("[ERROR] ReservationController.handleAllHotelsGet. Invalid URL parameters")
		res.WriteHeader(http.StatusBadRequest)
		return
	}

	hotelsLst, err := controller.service.ReadPaginatedHotels(page, pageSize)

	if err != nil {
		log.Println("[ERROR] ReservationController.handleAllHotelsGet. service.ReadPaginatedHotels returned error: ", err)
		res.WriteHeader(http.StatusInternalServerError)
		return
	}

	if hotelsLst.Len() == 0 {
		log.Println("[INFO] ReservationController.handleAllHotelsGet. Entities not found. Sending No Content code")
		res.WriteHeader(http.StatusNotFound)
		return
	}

	hotelsSlice := make([]models.Hotel, 0)
	hotelsLstEl := hotelsLst.Front()

	for hotelsLstEl != nil {
		hotelsSlice = append(hotelsSlice, hotelsLstEl.Value.(models.Hotel))
		hotelsLstEl = hotelsLstEl.Next()
	}

	hotelsSliceJSON, err := json.Marshal(hotelsSlice)

	if err != nil {
		log.Println("[ERROR] ReservationController.handleAllHotelsGet. Cannot convert result into JSON format: ", err)
		res.WriteHeader(http.StatusInternalServerError)
		return
	}

	log.Println("[INFO] hotelsSliceJSON", string(hotelsSliceJSON), page, pageSize)

	res.Header().Add(`Content-Type`, `application/json`)
	res.WriteHeader(http.StatusOK)
	res.Write(hotelsSliceJSON)
}

func (controller *ReservationController) handleHotelByIdGet(res http.ResponseWriter, req *http.Request) {
	log.Println("[INFO] ReservationController.handleHotelByIdGet. Handling hotel by id GET request")

	hotelId, err := strconv.Atoi(req.Header.Get(`Hotel-Id`))

	if err != nil {
		log.Println("[ERROR] ReservationController.handleHotelByIdGet. Invalid hotel id")
		res.WriteHeader(http.StatusBadRequest)
		return
	}

	hotel, err := controller.service.ReadHotelById(hotelId)

	if err != nil {
		log.Println("[ERROR] ReservationController.handleHotelByIdGet. service.ReadHotelById returned error: ", err)
		if errors.Is(err, serverrors.ErrEntityNotFound) {
			res.WriteHeader(http.StatusNotFound)
		} else {
			res.WriteHeader(http.StatusInternalServerError)
		}

		return
	}

	hotelJSON, err := json.Marshal(hotel)

	if err != nil {
		log.Println("[ERROR] ReservationController.handleHotelByIdGet. Cannot convert result into JSON format: ", err)
		res.WriteHeader(http.StatusInternalServerError)
		return
	}

	res.Header().Add(`Content-Type`, `application/json`)
	res.WriteHeader(http.StatusOK)
	res.Write(hotelJSON)
}

func (controller *ReservationController) handleHotelByUidGet(res http.ResponseWriter, req *http.Request) {
	log.Println("[INFO] ReservationController.handleHotelByUidGet. Handling hotel by uid GET request")

	hotelUid := req.PathValue("hotelUid")

	if strings.Trim(hotelUid, ` `) == `` {
		log.Println("[ERROR] ReservationController.handleHotelByUidGet. Invalid username")
		res.WriteHeader(http.StatusBadRequest)
		return
	}

	hotel, err := controller.service.ReadHotelByUid(hotelUid)

	if err != nil {
		log.Println("[ERROR] ReservationController.handleHotelByUidGet. service.ReadHotelByUid returned error: ", err)
		if errors.Is(err, serverrors.ErrEntityNotFound) {
			res.WriteHeader(http.StatusNotFound)
		} else {
			res.WriteHeader(http.StatusInternalServerError)
		}

		return
	}

	hotelJSON, err := json.Marshal(hotel)

	if err != nil {
		log.Println("[ERROR] ReservationController.handleHotelByUidGet. Cannot convert result into JSON format: ", err)
		res.WriteHeader(http.StatusInternalServerError)
		return
	}

	res.Header().Add(`Content-Type`, `application/json`)
	res.WriteHeader(http.StatusOK)
	res.Write(hotelJSON)
}

func (controller *ReservationController) handleReservsByUsernameGet(res http.ResponseWriter, req *http.Request) {
	log.Println("[INFO] ReservationController.handleReservsByUsernameGet. Handling reservations by username GET request")

	username := req.Header.Get(`X-User-Name`)

	if strings.Trim(username, ` `) == `` {
		log.Println("[ERROR] ReservationController.handleReservsByUsernameGet. Invalid username: " + username)
		res.WriteHeader(http.StatusBadRequest)
		return
	}

	reservsLst, err := controller.service.ReadReservsByUsername(username)

	if err != nil {
		log.Println("[ERROR] ReservationController.handleReservsByUsernameGet. service.ReadReservsByUsername returned error: ", err)
		if errors.Is(err, serverrors.ErrEntityNotFound) {
			res.WriteHeader(http.StatusNotFound)
		} else {
			res.WriteHeader(http.StatusInternalServerError)
		}

		return
	}

	reservsSlice := make([]models.Reservation, 0)
	reservsLstEl := reservsLst.Front()

	for reservsLstEl != nil {
		reservsSlice = append(reservsSlice, reservsLstEl.Value.(models.Reservation))
		reservsLstEl = reservsLstEl.Next()
	}

	reservsSliceJSON, err := json.Marshal(reservsSlice)

	if err != nil {
		log.Println("[ERROR] ReservationController.handleReservsByUsernameGet. Cannot convert result into JSON format: ", err)
		res.WriteHeader(http.StatusInternalServerError)
		return
	}

	res.Header().Add(`Content-Type`, `application/json`)
	res.WriteHeader(http.StatusOK)
	res.Write(reservsSliceJSON)
}

func (controller *ReservationController) handleReservByUidGet(res http.ResponseWriter, req *http.Request) {
	log.Println("[INFO] ReservationController.handleReservByUidGet. Handling reservation by uid GET request")

	reservUid := req.PathValue(`reservUid`)

	if strings.Trim(reservUid, ` `) == `` {
		log.Println("[ERROR] ReservationController.handleReservByUidGet. Invalid reservation uid")
		res.WriteHeader(http.StatusBadRequest)
		return
	}

	reservation, err := controller.service.ReadReservByUid(reservUid)

	if err != nil {
		log.Println("[ERROR] ReservationController.handleReservByUidGet. service.ReadReservByUid returned error: ", err)
		if errors.Is(err, serverrors.ErrEntityNotFound) {
			res.WriteHeader(http.StatusNotFound)
		} else {
			res.WriteHeader(http.StatusInternalServerError)
		}

		return
	}

	reservJSON, err := json.Marshal(reservation)

	if err != nil {
		log.Println("[ERROR] ReservationController.handleReservByUidGet. Cannot convert result into JSON format: ", err)
		res.WriteHeader(http.StatusInternalServerError)
		return
	}

	res.Header().Add(`Content-Type`, `application/json`)
	res.WriteHeader(http.StatusOK)
	res.Write(reservJSON)
}

func (controller *ReservationController) handleReservByUidPut(res http.ResponseWriter, req *http.Request) {
	log.Println("[INFO] ReservationController.handleReservByUidPut. Handling reservation by uid PUT request")

	reservUid := req.PathValue("reservUid")

	if strings.Trim(reservUid, ` `) == `` {
		log.Println("[ERROR] ReservationController.handleReservByUidPut. Invalid reservation uid")
		res.WriteHeader(http.StatusBadRequest)
		return
	}

	defer req.Body.Close()

	reqBody, err := io.ReadAll(req.Body)

	if err != nil {
		log.Println("[ERROR] ReservationController.handleReservByUidPut. Error while reading request body: ", err)
		res.WriteHeader(http.StatusBadRequest)
		return
	}

	var reservation models.Reservation
	err = json.Unmarshal(reqBody, &reservation)

	if err != nil {
		log.Println("[ERROR] ReservationController.handleReservByUidPut. Error while parsing JSON request body: ", err)
		res.WriteHeader(http.StatusBadRequest)
		return
	}

	reservation.Uid = reservUid
	_, err = controller.service.UpdateReservByUid(&reservation)

	if err != nil {
		log.Println("[ERROR] ReservationController.handleReservByUidPut. service.UpdateReservByUid returned error: ", err)
		if errors.Is(err, serverrors.ErrEntityNotFound) {
			res.WriteHeader(http.StatusNotFound)
		} else {
			res.WriteHeader(http.StatusInternalServerError)
		}

		return
	}

	res.WriteHeader(http.StatusNoContent)
}

func (controller *ReservationController) handleReservPost(res http.ResponseWriter, req *http.Request) {
	log.Println("[INFO] ReservationController.handleReservPost. Handling reserfvation POST request")

	reqBody, err := io.ReadAll(req.Body)

	if err != nil {
		log.Println("[ERROR] ReservationController.handleReservPost. Error while reading request body: ", err)
		res.WriteHeader(http.StatusBadRequest)
		return
	}

	var reservation models.Reservation
	err = json.Unmarshal(reqBody, &reservation)

	if err != nil {
		log.Println("[ERROR] ReservationController.handleReservPost. Error while parsing JSON request body: ", err)
		res.WriteHeader(http.StatusBadRequest)
		return
	}

	newReservation, err := controller.service.CreateReserv(&reservation)

	if err != nil {
		log.Println("[ERROR] ReservationController.handleReservPost. service.CreateReserv returned error: ", err)
		if errors.Is(err, serverrors.ErrEntityNotFound) {
			res.WriteHeader(http.StatusNotFound)
		} else {
			res.WriteHeader(http.StatusInternalServerError)
		}

		return
	}

	newReservationJSON, err := json.Marshal(&newReservation)

	if err != nil {
		log.Println("[ERROR] ReservationController.handleReservPost. Cannot convert result into JSON format: ", err)
		res.WriteHeader(http.StatusInternalServerError)
		return
	}

	res.Header().Add(`Content-Type`, `application/json`)
	res.WriteHeader(http.StatusOK)
	res.Write(newReservationJSON)
}

func (controller *ReservationController) handleHotelsRequest(res http.ResponseWriter, req *http.Request) {
	if req.Method == `GET` {
		if strings.Trim(req.Header.Get(`Hotel-Id`), ` `) == `` {
			log.Println("[INFO] ReservationController.handleHotelsRequest. Got hotels GET request")
			controller.handleAllHotelsGet(res, req)
		} else {
			log.Println("[INFO] ReservationController.handleHotelsRequest. Got hotel by id GET request")
			controller.handleHotelByIdGet(res, req)
		}
	} else {
		log.Println("[ERROR] ReservationController.handleHotelsRequest. Method not allowed")
		res.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func (controller *ReservationController) handleHotelWithUidRequest(res http.ResponseWriter, req *http.Request) {
	if req.Method == `GET` {
		log.Println("[INFO] ReservationController.handleHotelWithUidRequest. Got hotel with uid GET request")
		controller.handleHotelByUidGet(res, req)
	} else {
		log.Println("[ERROR] ReservationController.handleHotelWithUidRequest. Method not allowed")
		res.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func (controller *ReservationController) handleReservsRequest(res http.ResponseWriter, req *http.Request) {
	if req.Method == `GET` {
		if strings.Trim(req.Header.Get(`X-User-Name`), ` `) != `` {
			log.Println("[INFO] ReservationController.handleReservsRequest. Got reservation GET request")
			controller.handleReservsByUsernameGet(res, req)
		} else {
			log.Println("[ERROR] PaymentController.handleReservsRequest. Invalid request")
			res.WriteHeader(http.StatusMethodNotAllowed)
		}
	} else if req.Method == `POST` {
		log.Println("[INFO] ReservationController.handleReservsRequest. Got reservation POST request")
		controller.handleReservPost(res, req)
	} else {
		log.Println("[ERROR] ReservationController.handleReservsRequest. Method not allowed")
		res.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func (controller *ReservationController) handleReservWithUidRequest(res http.ResponseWriter, req *http.Request) {
	if req.Method == `GET` {
		log.Println("[INFO] ReservationController.handleReservWithUidRequest. Got reservation with uid GET request")
		controller.handleReservByUidGet(res, req)
	} else if req.Method == `PUT` {
		log.Println("[INFO] ReservationController.handleReservWithUidRequest. Got reservation with uid PUT request")
		controller.handleReservByUidPut(res, req)
	} else {
		log.Println("[ERROR] ReservationController.handleReservWithUidRequest. Method not allowed")
		res.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func (controller *ReservationController) handleHealthRequest(res http.ResponseWriter, req *http.Request) {
	if req.Method == `GET` {
		log.Println("[INFO] ReservationController.handleHealthRequest. Got health GET request")
		res.WriteHeader(http.StatusOK)
	} else {
		log.Println("[ERROR] ReservationController.handleHealthRequest. Method not allowed")
		res.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func (controller *ReservationController) Prepare() error {
	http.HandleFunc(`/api/v1/hotels`, controller.handleHotelsRequest)
	http.HandleFunc(`/api/v1/hotels/{hotelUid}`, controller.handleHotelWithUidRequest)
	http.HandleFunc(`/api/v1/reservations`, controller.handleReservsRequest)
	http.HandleFunc(`/api/v1/reservations/{reservUid}`, controller.handleReservWithUidRequest)

	http.HandleFunc(`/manage/health`, controller.handleHealthRequest)

	return nil
}

func (controller *ReservationController) Run() error {
	return http.ListenAndServe(
		fmt.Sprintf(
			`%s:%d`,
			controller.host,
			controller.port,
		),
		nil,
	)
}
