package database

import (
	"container/list"
	"context"
	"errors"
	"fmt"
	"log"

	"github.com/google/uuid"
	_ "github.com/jackc/pgx"
	"github.com/jackc/pgx/v5"

	"github.com/agarmirus/ds-lab02/internal/models"
	"github.com/agarmirus/ds-lab02/internal/serverrors"
)

type PostgresPaymentDAO struct {
	connStr string
}

func NewPostgresPaymentDAO(connStr string) IDAO[models.Payment] {
	return &PostgresPaymentDAO{connStr}
}

func (dao *PostgresPaymentDAO) SetConnectionString(connStr string) {
	dao.connStr = connStr
}

func validatePayment(payment *models.Payment) (err error) {
	if uuid.Validate(payment.Uid) != nil {
		err = serverrors.ErrInvalidPaymentUid
	} else if payment.Status != `PAID` && payment.Status != `CANCELED` {
		err = serverrors.ErrInvalidPaymentStatus
	} else if payment.Price < 0 {
		err = serverrors.ErrInvalidPaymentPrice
	}

	return err
}

func (dao *PostgresPaymentDAO) Create(payment *models.Payment) (newPayment models.Payment, err error) {
	err = validatePayment(payment)

	if err != nil {
		log.Println("[ERROR] PostgresPaymentDAO.Create. Invalid payments data:", err)
		return newPayment, err
	}

	conn, err := pgx.Connect(context.Background(), dao.connStr)

	if err != nil {
		log.Println("[ERROR] PostgresPaymentDAO.Create. Cannot connect to database:", err)
		return newPayment, serverrors.ErrDatabaseConnection
	}

	defer conn.Close(context.Background())

	row := conn.QueryRow(
		context.Background(),
		`insert into payment (payment_uid, status, price)
		values ($1, $2, $3)
		returning *;`,
		payment.Uid, payment.Status, payment.Price,
	)

	err = row.Scan(
		&newPayment.Id, &newPayment.Uid,
		&newPayment.Status, &newPayment.Price,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			log.Println("[ERROR] PostgresPaymentDAO.Create. Entity not found")
			err = serverrors.ErrEntityNotFound
		} else {
			log.Println("[ERROR] PostgresPaymentDAO.Create. Error while reading query result:", err)
			err = serverrors.ErrQueryResRead
		}
	}

	return newPayment, err
}

func (dao *PostgresPaymentDAO) Get() (list.List, error) {
	log.Println("[ERROR] PostgresPaymentDAO.Get. Method is not implemented")
	return list.List{}, serverrors.ErrMethodIsNotImplemented
}

func (dao *PostgresPaymentDAO) GetPaginated(
	page int,
	pageSize int,
) (resLst list.List, err error) {
	log.Println("[ERROR] PostgresLoyaltyDAO.GetPaginated. Method is not implemented")
	return list.List{}, serverrors.ErrMethodIsNotImplemented
}

func (dao *PostgresPaymentDAO) GetById(payment *models.Payment) (models.Payment, error) {
	log.Println("[ERROR] PostgresPaymentDAO.GetById. Method is not implemented")
	return models.Payment{}, serverrors.ErrMethodIsNotImplemented
}

func (dao *PostgresPaymentDAO) GetByAttribute(attrName string, attrValue string) (resLst list.List, err error) {
	conn, err := pgx.Connect(context.Background(), dao.connStr)

	if err != nil {
		log.Println("[ERROR] PostgresPaymentDAO.GetByAttribute. Cannot connect to database:", err)
		return resLst, serverrors.ErrDatabaseConnection
	}

	defer conn.Close(context.Background())

	queryStr := fmt.Sprintf(
		`select * from payment where %s = $1;`,
		attrName,
	)

	rows, err := conn.Query(context.Background(), queryStr, attrValue)

	if err != nil {
		if !errors.Is(err, pgx.ErrNoRows) {
			log.Println("[ERROR] PostgresPaymentDAO.GetByAttribute. Error while executing query:", err)
			return resLst, serverrors.ErrQueryResRead
		}

		return resLst, nil
	}

	defer rows.Close()

	for rows.Next() {
		var payment models.Payment
		err = rows.Scan(
			&payment.Id, &payment.Uid,
			&payment.Status, &payment.Price,
		)

		if err != nil {
			log.Println("[ERROR] PostgresPaymentDAO.GetByAttribute. Error while reading query result:", err)
			return list.List{}, serverrors.ErrQueryResRead
		}

		resLst.PushBack(payment)
	}

	return resLst, nil
}

func (dao *PostgresPaymentDAO) Update(payment *models.Payment) (updatedPayment models.Payment, err error) {
	conn, err := pgx.Connect(context.Background(), dao.connStr)

	if err != nil {
		log.Println("[ERROR] PostgresPaymentDAO.Update. Cannot connect to database:", err)
		return updatedPayment, serverrors.ErrDatabaseConnection
	}

	defer conn.Close(context.Background())

	row := conn.QueryRow(
		context.Background(),
		`update payment
		set status = $1, price = $2
		where payment_uid = $3
		returning *`,
		payment.Status, payment.Price,
		payment.Uid,
	)

	err = row.Scan(
		&updatedPayment.Id, &updatedPayment.Uid,
		&updatedPayment.Status, &updatedPayment.Price,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			log.Println("[ERROR] PostgresPaymentDAO.Update. Entity not found")
			err = serverrors.ErrEntityNotFound
		} else {
			log.Println("[ERROR] PostgresPaymentDAO.Update. Error while reading query result:", err)
			err = serverrors.ErrQueryResRead
		}
	}

	return updatedPayment, err
}

func (dao *PostgresPaymentDAO) Delete(payment *models.Payment) error {
	log.Println("[ERROR] PostgresPaymentDAO.Delete. Method is not implemented")
	return serverrors.ErrMethodIsNotImplemented
}

func (dao *PostgresPaymentDAO) DeleteByAttr(attrName string, attrValue string) error {
	log.Println("[ERROR] PostgresPaymentDAO.DeleteByAttr. Method is not implemented")
	return serverrors.ErrMethodIsNotImplemented
}
