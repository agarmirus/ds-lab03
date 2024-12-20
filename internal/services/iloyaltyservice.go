package services

import "github.com/agarmirus/ds-lab02/internal/models"

type ILoyaltyService interface {
	ReadLoyaltyByUsername(string) (models.Loyalty, error)
	UpdateLoyaltyById(*models.Loyalty) (models.Loyalty, error)
}
