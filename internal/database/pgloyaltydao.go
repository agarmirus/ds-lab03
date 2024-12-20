package database

import (
	"container/list"
	"context"
	"errors"
	"fmt"
	"log"

	"github.com/jackc/pgx/v5"

	"github.com/agarmirus/ds-lab02/internal/models"
	"github.com/agarmirus/ds-lab02/internal/serverrors"
)

type PostgresLoyaltyDAO struct {
	connStr string
}

func NewPostgresLoyaltyDAO(connStr string) IDAO[models.Loyalty] {
	return &PostgresLoyaltyDAO{connStr}
}

func (dao *PostgresLoyaltyDAO) SetConnectionString(connStr string) {
	dao.connStr = connStr
}

func (dao *PostgresLoyaltyDAO) Create(loyalty *models.Loyalty) (models.Loyalty, error) {
	log.Println("[ERROR] PostgresLoyaltyDAO.Create. Method is not implemented")
	return models.Loyalty{}, serverrors.ErrMethodIsNotImplemented
}

func (dao *PostgresLoyaltyDAO) Get() (list.List, error) {
	log.Println("[ERROR] PostgresLoyaltyDAO.Get. Method is not implemented")
	return list.List{}, serverrors.ErrMethodIsNotImplemented
}

func (dao *PostgresLoyaltyDAO) GetPaginated(
	page int,
	pageSize int,
) (resLst list.List, err error) {
	log.Println("[ERROR] PostgresLoyaltyDAO.GetPaginated. Method is not implemented")
	return list.List{}, serverrors.ErrMethodIsNotImplemented
}

func (dao *PostgresLoyaltyDAO) GetById(loyalty *models.Loyalty) (models.Loyalty, error) {
	log.Println("[ERROR] PostgresLoyaltyDAO.GetById. Method is not implemented")
	return models.Loyalty{}, serverrors.ErrMethodIsNotImplemented
}

func (dao *PostgresLoyaltyDAO) GetByAttribute(attrName string, attrValue string) (resLst list.List, err error) {
	conn, err := pgx.Connect(context.Background(), dao.connStr)

	if err != nil {
		log.Println("[ERROR] PostgresLoyaltyDAO.GetByAttribute. Cannot connect to database:", err)
		return resLst, serverrors.ErrDatabaseConnection
	}

	defer conn.Close(context.Background())

	queryStr := fmt.Sprintf(
		`select * from loyalty where %s = $1;`,
		attrName,
	)

	rows, err := conn.Query(context.Background(), queryStr, attrValue)

	if err != nil {
		if !errors.Is(err, pgx.ErrNoRows) {
			log.Println("[ERROR] PostgresLoyaltyDAO.GetByAttribute. Error while executing query:", err)
			return resLst, serverrors.ErrQueryResRead
		}

		return resLst, nil
	}

	defer rows.Close()

	for rows.Next() {
		var loyalty models.Loyalty
		err = rows.Scan(
			&loyalty.Id, &loyalty.Username,
			&loyalty.ReservationCount, &loyalty.Status,
			&loyalty.Discount,
		)

		if err != nil {
			log.Println("[ERROR] PostgresLoyaltyDAO.GetByAttribute. Error while reading query result:", err)
			return list.List{}, serverrors.ErrQueryResRead
		}

		resLst.PushBack(loyalty)
	}

	return resLst, nil
}

func (dao *PostgresLoyaltyDAO) Update(loyalty *models.Loyalty) (updatedLoyalty models.Loyalty, err error) {
	conn, err := pgx.Connect(context.Background(), dao.connStr)

	if err != nil {
		log.Println("[ERROR] PostgresLoyaltyDAO.Update. Cannot connect to database:", err)
		return updatedLoyalty, serverrors.ErrDatabaseConnection
	}

	defer conn.Close(context.Background())

	row := conn.QueryRow(
		context.Background(),
		`update loyalty
		set username = $1, reservation_count = $2, status = $3, discount = $4
		where id = $5
		returning *`,
		loyalty.Username, loyalty.ReservationCount, loyalty.Status, loyalty.Discount,
		loyalty.Id,
	)

	err = row.Scan(
		&updatedLoyalty.Id, &updatedLoyalty.Username,
		&updatedLoyalty.ReservationCount, &updatedLoyalty.Status,
		&updatedLoyalty.Discount,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			log.Println("[ERROR] PostgresLoyaltyDAO.Update. Entity not found")
			err = serverrors.ErrEntityNotFound
		} else {
			log.Println("[ERROR] PostgresLoyaltyDAO.Update. Error while reading query result:", err)
			err = serverrors.ErrQueryResRead
		}
	}

	return updatedLoyalty, err
}

func (dao *PostgresLoyaltyDAO) Delete(loyalty *models.Loyalty) error {
	log.Println("[ERROR] PostgresLoyaltyDAO.Delete. Method is not implemented")
	return serverrors.ErrMethodIsNotImplemented
}

func (dao *PostgresLoyaltyDAO) DeleteByAttr(attrName string, attrValue string) error {
	log.Println("[ERROR] PostgresLoyaltyDAO.DeleteByAttr. Method is not implemented")
	return serverrors.ErrMethodIsNotImplemented
}
