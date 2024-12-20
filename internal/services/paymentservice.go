package services

import (
	"log"

	"github.com/agarmirus/ds-lab02/internal/database"
	"github.com/agarmirus/ds-lab02/internal/models"
	"github.com/agarmirus/ds-lab02/internal/serverrors"
)

type PaymentService struct {
	paymentDAO database.IDAO[models.Payment]
}

func NewPaymentService(
	paymentDAO database.IDAO[models.Payment],
) IPaymentService {
	return &PaymentService{paymentDAO}
}

func (service *PaymentService) ReadPaymentByUid(paymentUid string) (payment models.Payment, err error) {
	paymentsLst, err := service.paymentDAO.GetByAttribute(`payment_uid`, paymentUid)

	if err != nil {
		log.Println("[ERROR] PaymentService.ReadPaymentByUid. paymentDAO.GetByAttribute returned error:", err)
		return payment, err
	}

	if paymentsLst.Len() == 0 {
		log.Println("[ERROR] PaymentService.ReadPaymentByUid. Entity not found")
		return payment, serverrors.ErrEntityNotFound
	}

	return paymentsLst.Front().Value.(models.Payment), nil
}

func (service *PaymentService) UpdatePaymentByUid(payment *models.Payment) (newPayment models.Payment, err error) {
	newPayment, err = service.paymentDAO.Update(payment)

	if err != nil {
		log.Println("[ERROR] PaymentService.UpdatePaymentByUid. paymentDAO.Update returned error:", err)
	}

	return newPayment, err
}

func (service *PaymentService) CreatePayment(payment *models.Payment) (newPayment models.Payment, err error) {
	newPayment, err = service.paymentDAO.Create(payment)

	if err != nil {
		log.Println("[ERROR] PaymentService.CreatePayment. paymentDAO.Create returned error:", err)
	}

	return newPayment, err
}
