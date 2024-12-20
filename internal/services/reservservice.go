package services

import (
	"container/list"
	"log"

	"github.com/agarmirus/ds-lab02/internal/database"
	"github.com/agarmirus/ds-lab02/internal/models"
	"github.com/agarmirus/ds-lab02/internal/serverrors"
)

type ReservationService struct {
	reservsDAO database.IDAO[models.Reservation]
	hotelsDAO  database.IDAO[models.Hotel]
}

func NewReservationService(
	reservsDAO database.IDAO[models.Reservation],
	hotelsDAO database.IDAO[models.Hotel],
) IReservationService {
	return &ReservationService{reservsDAO, hotelsDAO}
}

func (service *ReservationService) ReadPaginatedHotels(
	page int,
	pageSize int,
) (hotelsLst list.List, err error) {
	if page <= 0 || pageSize <= 0 {
		log.Println("[ERROR] ReservationService.ReadPaginatedHotels. Invalid page and pageSize values")
		return hotelsLst, serverrors.ErrInvalidPagesData
	}

	hotelsLst, err = service.hotelsDAO.GetPaginated(page, pageSize)

	if err != nil {
		log.Println("[ERROR] ReservationService.ReadPaginatedHotels. hotelsDAO.GetPaginated returned error:", err)
	}

	return hotelsLst, err
}

func (service *ReservationService) ReadHotelById(hotelId int) (hotel models.Hotel, err error) {
	desiredHotel := models.Hotel{Id: hotelId}
	hotel, err = service.hotelsDAO.GetById(&desiredHotel)

	if err != nil {
		log.Println("[ERROR] ReservationService.ReadHotelById. hotelsDAO.GetById returned error:", err)
	}

	return hotel, err
}

func (service *ReservationService) ReadHotelByUid(hotelUid string) (hotel models.Hotel, err error) {
	hotelsLst, err := service.hotelsDAO.GetByAttribute(`hotel_uid`, hotelUid)

	if err != nil {
		log.Println("[ERROR] ReservationService.ReadHotelByUid. hotelsDAO.GetByAttribute returned error:", err)
		return hotel, err
	}

	if hotelsLst.Len() == 0 {
		log.Println("[ERROR] ReservationService.ReadHotelByUid. Entity not found")
		return hotel, serverrors.ErrEntityNotFound
	}

	return hotelsLst.Front().Value.(models.Hotel), nil
}

func (service *ReservationService) ReadReservsByUsername(username string) (reservsLst list.List, err error) {
	reservsLst, err = service.reservsDAO.GetByAttribute(`username`, username)

	if err != nil {
		log.Println("[ERROR] ReservationService.ReadReservsByUsername. hotelsDAO.GetByAttribute returned error:", err)
		return reservsLst, err
	}

	if reservsLst.Len() == 0 {
		log.Println("[ERROR] ReservationService.ReadReservsByUsername. Entity not found")
		return reservsLst, serverrors.ErrEntityNotFound
	}

	return reservsLst, nil
}

func (service *ReservationService) ReadReservByUid(reservUid string) (reservation models.Reservation, err error) {
	reservsLst, err := service.reservsDAO.GetByAttribute(`reservation_uid`, reservUid)

	if err != nil {
		log.Println("[ERROR] ReservationService.ReadReservByUid. hotelsDAO.GetByAttribute returned error:", err)
		return reservation, err
	}

	if reservsLst.Len() == 0 {
		log.Println("[ERROR] ReservationService.ReadReservByUid. Entity not found")
		return reservation, serverrors.ErrEntityNotFound
	}

	return reservsLst.Front().Value.(models.Reservation), nil
}

func (service *ReservationService) UpdateReservByUid(reservation *models.Reservation) (updatedReservation models.Reservation, err error) {
	updatedReservation, err = service.reservsDAO.Update(reservation)

	if err != nil {
		log.Println("[ERROR] ReservationService.UpdateReservByUid. reservsDAO.Update returned error:", err)
	}

	return updatedReservation, err
}

func (service *ReservationService) CreateReserv(reservation *models.Reservation) (newReservation models.Reservation, err error) {
	newReservation, err = service.reservsDAO.Create(reservation)

	if err != nil {
		log.Println("[ERROR] ReservationService.CreateReserv. reservsDAO.Create returned error:", err)
	}

	return newReservation, err
}
