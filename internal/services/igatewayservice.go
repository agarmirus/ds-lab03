package services

import (
	"github.com/agarmirus/ds-lab02/internal/models"
)

type IGatewayService interface {
	ReadAllHotels(int, int) (models.PagiationResponse, error)
	ReadUserInfo(string) (models.UserInfoResponse, error)
	ReadUserReservations(string) ([]models.ReservationResponse, error)
	CreateReservation(string, *models.CreateReservationRequest) (models.CreateReservationResponse, error)
	ReadReservation(string, string) (models.ReservationResponse, error)
	DeleteReservation(string, string) error
	ReadUserLoyalty(string) (models.LoyaltyInfoResponse, error)
}
