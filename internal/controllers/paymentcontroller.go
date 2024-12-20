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
	"github.com/google/uuid"
)

type PaymentController struct {
	host string
	port int

	service services.IPaymentService
}

func NewPaymentController(
	host string,
	port int,
	service services.IPaymentService,
) IController {
	return &PaymentController{host, port, service}
}

func (controller *PaymentController) handlePaymentByPricePost(res http.ResponseWriter, req *http.Request) {
	log.Println("[INFO] PaymentController.handlePaymentByPricePost. Handling payment by price POST request")

	price, err := strconv.Atoi(req.Header.Get(`Price`))

	if err != nil || price <= 0 {
		log.Println("[ERROR] PaymentController.handlePaymentByPricePost. Ivalid price value:", err)
		res.WriteHeader(http.StatusBadRequest)
		return
	}

	payment := models.Payment{Uid: uuid.New().String(), Status: `PAID`, Price: price}

	newPayment, err := controller.service.CreatePayment(&payment)

	if err != nil {
		log.Println("[ERROR] PaymentController.handlePaymentByPricePost. service.CreatePayment returned error: ", err)
		res.WriteHeader(http.StatusInternalServerError)
		return
	}

	newPaymentJSON, err := json.Marshal(newPayment)

	if err != nil {
		log.Println("[ERROR] PaymentController.handlePaymentByPricePost. Cannot convert result into JSON format: ", err)
		res.WriteHeader(http.StatusInternalServerError)
		return
	}

	res.Header().Add(`Content-Type`, `application/json`)
	res.WriteHeader(http.StatusOK)
	res.Write(newPaymentJSON)
}

func (controller *PaymentController) handlePaymentByUidGet(res http.ResponseWriter, req *http.Request) {
	log.Println("[INFO] PaymentController.handlePaymentByUidGet. Handling payment by uid GET request")

	paymentUid := req.PathValue("paymentUid")

	if paymentUid == `` {
		log.Println("[ERROR] PaymentController.handlePaymentByUidGet. Invalid payment uid")
		res.WriteHeader(http.StatusBadRequest)
		return
	}

	payment, err := controller.service.ReadPaymentByUid(paymentUid)

	if err != nil {
		log.Println("[ERROR] PaymentController.handlePaymentByUidGet. service.ReadPaymentByUid returned error: ", err)
		if errors.Is(err, serverrors.ErrEntityNotFound) {
			res.WriteHeader(http.StatusNotFound)
		} else {
			res.WriteHeader(http.StatusInternalServerError)
		}

		return
	}

	paymentJSON, err := json.Marshal(payment)

	if err != nil {
		log.Println("[ERROR] PaymentController.handlePaymentByUidGet. Cannot convert result into JSON format: ", err)
		res.WriteHeader(http.StatusInternalServerError)
		return
	}

	res.Header().Add(`Content-Type`, `application/json`)
	res.WriteHeader(http.StatusOK)
	res.Write(paymentJSON)
}

func (controller *PaymentController) handlePaymentByUidPut(res http.ResponseWriter, req *http.Request) {
	log.Println("[INFO] PaymentController.handlePaymentByUidPut. Handling payment by uid PUT request")

	paymentUid := req.PathValue("paymentUid")

	if strings.Trim(paymentUid, ` `) == `` {
		log.Println("[ERROR] PaymentController.handlePaymentByUidPut. Invalid payment uid")
		res.WriteHeader(http.StatusBadRequest)
		return
	}

	defer req.Body.Close()

	reqBody, err := io.ReadAll(req.Body)

	if err != nil {
		log.Println("[ERROR] PaymentController.handlePaymentByUidPut. Error while reading request body: ", err)
		res.WriteHeader(http.StatusBadRequest)
		return
	}

	var payment models.Payment
	err = json.Unmarshal(reqBody, &payment)

	if err != nil {
		log.Println("[ERROR] PaymentController.handlePaymentByUidPut. Error while parsing JSON request body: ", err)
		res.WriteHeader(http.StatusBadRequest)
		return
	}

	payment.Uid = paymentUid
	_, err = controller.service.UpdatePaymentByUid(&payment)

	if err != nil {
		log.Println("[ERROR] PaymentController.handlePaymentByUidPut. service.UpdatePaymentByUid returned error: ", err)
		if errors.Is(err, serverrors.ErrEntityNotFound) {
			res.WriteHeader(http.StatusNotFound)
		} else {
			res.WriteHeader(http.StatusInternalServerError)
		}

		return
	}

	res.WriteHeader(http.StatusNoContent)
}

func (controller *PaymentController) handlePaymentRequest(res http.ResponseWriter, req *http.Request) {
	if req.Method == `POST` {
		if strings.Trim(req.Header.Get(`Price`), ` `) != `` {
			log.Println("[INFO] PaymentController.handlePaymentByUidRequest. Got payment by price POST request")
			controller.handlePaymentByPricePost(res, req)
		} else {
			log.Println("[ERROR] PaymentController.handlePaymentRequest. Invalid request")
			res.WriteHeader(http.StatusBadRequest)
		}
	} else {
		log.Println("[ERROR] PaymentController.handlePaymentRequest. Method not allowed")
		res.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func (controller *PaymentController) handlePaymentByUidRequest(res http.ResponseWriter, req *http.Request) {
	if req.Method == `GET` {
		log.Println("[INFO] PaymentController.handlePaymentByUidRequest. Got payment by uid GET request")
		controller.handlePaymentByUidGet(res, req)
	} else if req.Method == `PUT` {
		log.Println("[INFO] PaymentController.handlePaymentByUidRequest. Got payment by uid PUT request")
		controller.handlePaymentByUidPut(res, req)
	} else {
		log.Println("[ERROR] PaymentController.handlePaymentByUidRequest. Method not allowed")
		res.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func (controller *PaymentController) handleHealthRequest(res http.ResponseWriter, req *http.Request) {
	if req.Method == `GET` {
		log.Println("[INFO] PaymentController.handleHealthRequest. Got health GET request")
		res.WriteHeader(http.StatusOK)
	} else {
		log.Println("[ERROR] PaymentController.handleHealthRequest. Method not allowed")
		res.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func (controller *PaymentController) Prepare() error {
	http.HandleFunc(`/api/v1/payment`, controller.handlePaymentRequest)
	http.HandleFunc(`/api/v1/payment/{paymentUid}`, controller.handlePaymentByUidRequest)

	http.HandleFunc(`/manage/health`, controller.handleHealthRequest)

	return nil
}

func (controller *PaymentController) Run() error {
	return http.ListenAndServe(
		fmt.Sprintf(
			`%s:%d`,
			controller.host,
			controller.port,
		),
		nil,
	)
}
