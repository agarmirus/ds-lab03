package services

import (
	"log"

	"github.com/agarmirus/ds-lab02/internal/database"
	"github.com/agarmirus/ds-lab02/internal/models"
	"github.com/agarmirus/ds-lab02/internal/serverrors"
)

type LoyaltyService struct {
	loyaltyDAO database.IDAO[models.Loyalty]
}

func NewLoyaltyService(
	loyaltyDAO database.IDAO[models.Loyalty],
) ILoyaltyService {
	return &LoyaltyService{loyaltyDAO}
}

func (service *LoyaltyService) ReadLoyaltyByUsername(username string) (loyalty models.Loyalty, err error) {
	loyaltiesLst, err := service.loyaltyDAO.GetByAttribute(`username`, username)

	if err != nil {
		log.Println("[ERROR] LoyaltyService.ReadLoyaltyByUsername. loyaltyDAO.GetByAttribute returned error:", err)
		return loyalty, err
	}

	if loyaltiesLst.Len() == 0 {
		log.Println("[ERROR] LoyaltyService.ReadLoyaltyByUsername. Entity not found")
		return loyalty, serverrors.ErrEntityNotFound
	}

	return loyaltiesLst.Front().Value.(models.Loyalty), nil
}

func (service *LoyaltyService) UpdateLoyaltyById(loyalty *models.Loyalty) (updatedLoyalty models.Loyalty, err error) {
	updatedLoyalty, err = service.loyaltyDAO.Update(loyalty)

	if err != nil {
		log.Println("[ERROR] LoyaltyService.UpdateLoyaltyById. loyaltyDAO.Update returned error:", err)
	}

	return updatedLoyalty, err
}
