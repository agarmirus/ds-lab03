package models

import (
	"time"

	"github.com/agarmirus/ds-lab02/internal/serverrors"
	"github.com/google/uuid"
)

type Reservation struct {
	Id         int    `json:"id"`
	Uid        string `json:"reservationUid"`
	Username   string `json:"username"`
	PaymentUid string `json:"paymentUid"`
	HotelId    int    `json:"hotelId"`
	Status     string `json:"status"`
	StartDate  string `json:"startDate"`
	EndDate    string `json:"endDate"`
}

type Payment struct {
	Id     int    `json:"id"`
	Uid    string `json:"paymentUid"`
	Status string `json:"status"`
	Price  int    `json:"price"`
}

type Loyalty struct {
	Id               int    `json:"id"`
	Username         string `json:"username"`
	ReservationCount int    `json:"reservationCount"`
	Status           string `json:"status"`
	Discount         int    `json:"discount"`
}

type Hotel struct {
	Id      int    `json:"id"`
	Uid     string `json:"hotelUid"`
	Name    string `json:"name"`
	Country string `json:"country"`
	City    string `json:"city"`
	Address string `json:"address"`
	Stars   int    `json:"stars"`
	Price   int    `json:"price"`
}

type PaymentInfo struct {
	Status string `json:"status"`
	Price  int    `json:"price"`
}

type LoyaltyInfoResponse struct {
	Status           string `json:"status"`
	Discount         int    `json:"discount"`
	ReservationCount int    `json:"reservationCount"`
}

type HotelResponse struct {
	HotelUid string `json:"hotelUid"`
	Name     string `json:"name"`
	Country  string `json:"country"`
	City     string `json:"city"`
	Address  string `json:"address"`
	Stars    int    `json:"stars"`
	Price    int    `json:"price"`
}

type HotelInfo struct {
	HotelUid    string `json:"hotelUid"`
	Name        string `json:"name"`
	FullAddress string `json:"fullAddress"`
	Stars       int    `json:"stars"`
}

type ReservationResponse struct {
	ReservationUid string      `json:"reservationUid"`
	Hotel          HotelInfo   `json:"hotel"`
	StartDate      string      `json:"startDate"`
	EndDate        string      `json:"endDate"`
	Status         string      `json:"status"`
	Payment        PaymentInfo `json:"payment"`
}

type UserInfoResponse struct {
	Reservations []ReservationResponse `json:"reservations"`
	Loyalty      LoyaltyInfoResponse   `json:"loyalty"`
}

type CreateReservationRequest struct {
	HotelUid  string `json:"hotelUid"`
	StartDate string `json:"startDate"`
	EndDate   string `json:"endDate"`
}

type CreateReservationResponse struct {
	ReservationUid string      `json:"reservationUid"`
	HotelUid       string      `json:"hotelUid"`
	StartDate      string      `json:"startDate"`
	EndDate        string      `json:"endDate"`
	Discount       int         `json:"discount"`
	Status         string      `json:"status"`
	Payment        PaymentInfo `json:"payment"`
}

type PagiationResponse struct {
	Page          int             `json:"page"`
	PageSize      int             `json:"pageSize"`
	TotalElements int             `json:"totalElements"`
	Items         []HotelResponse `json:"items"`
}

type ErrorDiscription struct {
	Field string `json:"field"`
	Error string `json:"error"`
}

type ValidationErrorResponse struct {
	Message string             `json:"message"`
	Errors  []ErrorDiscription `json:"errors"`
}

type ErrorResponse struct {
	Message string `json:"message"`
}

func HotelToHotelReponse(
	hotelRes *HotelResponse,
	hotel *Hotel,
) {
	hotelRes.HotelUid = hotel.Uid
	hotelRes.Name = hotel.Name
	hotelRes.Country = hotel.Country
	hotelRes.City = hotel.City
	hotelRes.Address = hotel.Address
	hotelRes.Stars = hotel.Stars
	hotelRes.Price = hotel.Price
}

func HotelsSliceToPagRes(
	pagRes *PagiationResponse,
	hotelsSlice []Hotel,
	page int,
	pageSize int,
) {
	pagRes.Page = page
	pagRes.PageSize = pageSize
	pagRes.TotalElements = len(hotelsSlice)

	for i := 0; i < pagRes.TotalElements; i++ {
		var hotelRes HotelResponse
		HotelToHotelReponse(&hotelRes, &hotelsSlice[i])

		pagRes.Items = append(pagRes.Items, hotelRes)
	}
}

func hotelToHotelInfo(
	hotelInfo *HotelInfo,
	hotel *Hotel,
) {
	hotelInfo.HotelUid = hotel.Uid
	hotelInfo.Name = hotel.Name
	hotelInfo.FullAddress = hotel.Country + `, ` + hotel.City + `, ` + hotel.Address
	hotelInfo.Stars = hotel.Stars
}

func paymentToPaymentInfo(
	paymentInfo *PaymentInfo,
	payment *Payment,
) {
	paymentInfo.Status = payment.Status
	paymentInfo.Price = payment.Price
}

func ReservToReservRes(
	reservRes *ReservationResponse,
	reservation *Reservation,
	hotel *Hotel,
	payment *Payment,
) {
	reservRes.ReservationUid = reservation.Uid
	reservRes.StartDate = reservation.StartDate
	reservRes.EndDate = reservation.EndDate
	reservRes.Status = reservation.Status

	hotelToHotelInfo(&reservRes.Hotel, hotel)
	paymentToPaymentInfo(&reservRes.Payment, payment)
}

func ReservsSliceToReservRes(
	reservsResSlice *[]ReservationResponse,
	userReservsSlice []Reservation,
	hotelsMap map[int]Hotel,
	paymentsMap map[string]Payment,
) {
	totalElements := len(userReservsSlice)

	for i := 0; i < totalElements; i++ {
		var reservRes ReservationResponse
		hotel := hotelsMap[userReservsSlice[i].HotelId]
		payment := paymentsMap[userReservsSlice[i].PaymentUid]

		ReservToReservRes(&reservRes, &userReservsSlice[i], &hotel, &payment)

		*reservsResSlice = append(*reservsResSlice, reservRes)
	}
}

func LoyaltyToLoyaltyInfoRes(
	loyatlyInfoRes *LoyaltyInfoResponse,
	loyalty *Loyalty,
) {
	loyatlyInfoRes.Status = loyalty.Status
	loyatlyInfoRes.Discount = loyalty.Discount
	loyatlyInfoRes.ReservationCount = loyalty.ReservationCount
}

func ReservsSliceToUserInfoRes(
	userInfoRes *UserInfoResponse,
	userReservsSlice []Reservation,
	hotelsMap map[int]Hotel,
	paymentsMap map[string]Payment,
	loyalty *Loyalty,
) {
	userInfoRes.Reservations = make([]ReservationResponse, 0)
	ReservsSliceToReservRes(&userInfoRes.Reservations, userReservsSlice, hotelsMap, paymentsMap)
	LoyaltyToLoyaltyInfoRes(&userInfoRes.Loyalty, loyalty)
}

func ReservToCrReservRes(
	crReservRes *CreateReservationResponse,
	reservation *Reservation,
	payment *Payment,
	loyalty *Loyalty,
	hotelUid string,
) {
	crReservRes.ReservationUid = reservation.Uid
	crReservRes.HotelUid = hotelUid
	crReservRes.StartDate = reservation.StartDate
	crReservRes.EndDate = reservation.EndDate
	crReservRes.Discount = loyalty.Discount
	crReservRes.Status = reservation.Status

	paymentToPaymentInfo(&crReservRes.Payment, payment)
}

func ValidateCrReservReq(
	createReservReq *CreateReservationRequest,
) (validErrRes ValidationErrorResponse, err error) {
	if uuid.Validate(createReservReq.HotelUid) != nil {
		validErrRes.Errors = append(validErrRes.Errors, ErrorDiscription{Field: `hotelUid`, Error: `invalid uid`})
	}

	startDate, err := time.Parse(time.DateOnly, createReservReq.StartDate)

	if err != nil {
		validErrRes.Errors = append(validErrRes.Errors, ErrorDiscription{Field: `startDate`, Error: `invalid date format`})
	}

	endDate, err := time.Parse(time.DateOnly, createReservReq.EndDate)

	if err != nil {
		validErrRes.Errors = append(validErrRes.Errors, ErrorDiscription{Field: `endDate`, Error: `invalid date format`})
	} else if startDate.Unix() > endDate.Unix() {
		validErrRes.Errors = append(validErrRes.Errors, ErrorDiscription{Field: `startDate`, Error: `invalid date period`})
		err = serverrors.ErrInvalidReservDates
	}

	if err != nil {
		validErrRes.Message = `invalid reservation request data`
	}

	return validErrRes, err
}

func UpdateLoyaltyStatus(
	loyalty *Loyalty,
) {
	if loyalty.ReservationCount > 20 {
		loyalty.Status = `GOLD`
		loyalty.Discount = 10
	} else if loyalty.ReservationCount > 10 {
		loyalty.Status = `SILVER`
		loyalty.Discount = 7
	} else {
		loyalty.Status = `BRONZE`
		loyalty.Discount = 5
	}
}
