package services

import "github.com/agarmirus/ds-lab02/internal/models"

type IPaymentService interface {
	ReadPaymentByUid(string) (models.Payment, error)
	UpdatePaymentByUid(*models.Payment) (models.Payment, error)
	CreatePayment(*models.Payment) (models.Payment, error)
}
