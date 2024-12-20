package services

import (
	"container/list"

	"github.com/agarmirus/ds-lab02/internal/models"
)

type IReservationService interface {
	ReadPaginatedHotels(int, int) (list.List, error)
	ReadHotelById(int) (models.Hotel, error)
	ReadHotelByUid(string) (models.Hotel, error)
	ReadReservsByUsername(string) (list.List, error)
	ReadReservByUid(string) (models.Reservation, error)
	UpdateReservByUid(*models.Reservation) (models.Reservation, error)
	CreateReserv(*models.Reservation) (models.Reservation, error)
}
