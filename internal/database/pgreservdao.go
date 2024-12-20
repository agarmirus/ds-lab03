package database

import (
	"container/list"
	"context"
	"errors"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/google/uuid"
	_ "github.com/jackc/pgx"
	"github.com/jackc/pgx/v5"

	"github.com/agarmirus/ds-lab02/internal/models"
	"github.com/agarmirus/ds-lab02/internal/serverrors"
)

type PostgresReservationDAO struct {
	connStr string
}

func NewPostgresReservationDAO(connStr string) IDAO[models.Reservation] {
	return &PostgresReservationDAO{connStr}
}

func (dao *PostgresReservationDAO) SetConnectionString(connStr string) {
	dao.connStr = connStr
}

func validateReservation(reservation *models.Reservation) (err error) {
	if uuid.Validate(reservation.Uid) != nil {
		err = serverrors.ErrInvalidReservUid
	} else if strings.Trim(reservation.Username, ` `) == `` {
		err = serverrors.ErrInvalidReservUsername
	} else if uuid.Validate(reservation.Uid) != nil {
		err = serverrors.ErrInvalidReservPayUID
	} else if reservation.HotelId <= 0 {
		err = serverrors.ErrInvalidReservHotelId
	} else if reservation.Status != `PAID` && reservation.Status != `CANCELED` {
		err = serverrors.ErrInvalidReservStatus
	} else if reservation.StartDate > reservation.EndDate {
		err = serverrors.ErrInvalidReservDates
	}

	return err
}

func (dao *PostgresReservationDAO) Create(reservation *models.Reservation) (newReservation models.Reservation, err error) {
	err = validateReservation(reservation)

	if err != nil {
		log.Println("[ERROR] PostgresReservationDAO.Create. Invalid reservation data:", err)
		return newReservation, err
	}

	conn, err := pgx.Connect(context.Background(), dao.connStr)

	if err != nil {
		log.Println("[ERROR] PostgresReservationDAO.Create. Cannot connect to database:", err)
		return newReservation, serverrors.ErrDatabaseConnection
	}

	defer conn.Close(context.Background())

	row := conn.QueryRow(
		context.Background(),
		`insert into reservation (reservation_uid, username, payment_uid, hotel_id, status, start_date, end_date)
		values ($1, $2, $3, $4, $5, $6, $7)
		returning id;`,
		reservation.Uid, reservation.Username,
		reservation.PaymentUid, reservation.HotelId,
		reservation.Status, reservation.StartDate,
		reservation.EndDate,
	)

	newReservation = *reservation
	err = row.Scan(&newReservation.Id)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			log.Println("[ERROR] PostgresReservationDAO.Create. Entity not found")
			err = serverrors.ErrEntityNotFound
		} else {
			log.Println("[ERROR] PostgresReservationDAO.Create. Error while reading query result:", err)
			err = serverrors.ErrQueryResRead
		}
	}

	return newReservation, err
}

func (dao *PostgresReservationDAO) Get() (list.List, error) {
	log.Println("[ERROR] PostgresReservationDAO.Get. Method is not implemented")
	return list.List{}, serverrors.ErrMethodIsNotImplemented
}

func (dao *PostgresReservationDAO) GetPaginated(
	page int,
	pageSize int,
) (resLst list.List, err error) {
	log.Println("[ERROR] PostgresLoyaltyDAO.GetPaginated. Method is not implemented")
	return list.List{}, serverrors.ErrMethodIsNotImplemented
}

func (dao *PostgresReservationDAO) GetById(reservation *models.Reservation) (models.Reservation, error) {
	log.Println("[ERROR] PostgresReservationDAO.GetById. Method is not implemented")
	return models.Reservation{}, serverrors.ErrMethodIsNotImplemented
}

func (dao *PostgresReservationDAO) GetByAttribute(attrName string, attrValue string) (resLst list.List, err error) {
	conn, err := pgx.Connect(context.Background(), dao.connStr)

	if err != nil {
		log.Println("[ERROR] PostgresReservationDAO.GetByAttribute. Cannot connect to database:", err)
		return resLst, serverrors.ErrDatabaseConnection
	}

	defer conn.Close(context.Background())

	queryStr := fmt.Sprintf(
		`select id, reservation_uid, username, payment_uid, hotel_id, status, start_date, end_date from reservation where %s = $1;`,
		attrName,
	)

	rows, err := conn.Query(context.Background(), queryStr, attrValue)

	if err != nil {
		if !errors.Is(err, pgx.ErrNoRows) {
			log.Println("[ERROR] PostgresReservationDAO.GetByAttribute. Error while executing query:", err)
			return resLst, serverrors.ErrQueryResRead
		}

		return resLst, nil
	}

	defer rows.Close()

	for rows.Next() {
		var reservation models.Reservation
		var startDate, endDate time.Time
		err = rows.Scan(
			&reservation.Id, &reservation.Uid,
			&reservation.Username, &reservation.PaymentUid,
			&reservation.HotelId, &reservation.Status,
			&startDate, &endDate,
		)

		if err != nil {
			log.Println("[ERROR] PostgresReservationDAO.GetByAttribute. Error while reading query result:", err)
			return list.List{}, serverrors.ErrQueryResRead
		}

		reservation.StartDate = startDate.Format(time.DateOnly)
		reservation.EndDate = endDate.Format(time.DateOnly)

		resLst.PushBack(reservation)
	}

	return resLst, nil
}

func (dao *PostgresReservationDAO) Update(reservation *models.Reservation) (updatedReservation models.Reservation, err error) {
	conn, err := pgx.Connect(context.Background(), dao.connStr)

	if err != nil {
		log.Println("[ERROR] PostgresReservationDAO.Update. Cannot connect to database:", err)
		return updatedReservation, serverrors.ErrDatabaseConnection
	}

	defer conn.Close(context.Background())

	log.Println("[TRACE] PostgresReservationDAO.Update. Status =", reservation.Status)

	row := conn.QueryRow(
		context.Background(),
		`update reservation
		set username = $1, payment_uid = $2, hotel_id = $3, status = $4, start_date = $5, end_date = $6
		where reservation_uid = $7
		returning *`,
		reservation.Username, reservation.PaymentUid, reservation.HotelId,
		reservation.Status, reservation.StartDate, reservation.EndDate,
		reservation.Uid,
	)

	var startDate, endDate time.Time

	err = row.Scan(
		&updatedReservation.Id, &updatedReservation.Uid,
		&updatedReservation.Username, &updatedReservation.PaymentUid,
		&updatedReservation.HotelId, &updatedReservation.Status,
		&startDate, &endDate,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			log.Println("[ERROR] PostgresReservationDAO.Update. Entity not found")
			err = serverrors.ErrEntityNotFound
		} else {
			log.Println("[ERROR] PostgresReservationDAO.Update. Error while reading query result:", err)
			err = serverrors.ErrQueryResRead
		}
	}

	updatedReservation.StartDate = startDate.Format(time.DateOnly)
	updatedReservation.EndDate = endDate.Format(time.DateOnly)

	return updatedReservation, err
}

func (dao *PostgresReservationDAO) Delete(reservation *models.Reservation) error {
	log.Println("[ERROR] PostgresReservationDAO.Delete. Method is not implemented")
	return serverrors.ErrMethodIsNotImplemented
}

func (dao *PostgresReservationDAO) DeleteByAttr(attrName string, attrValue string) (err error) {
	conn, err := pgx.Connect(context.Background(), dao.connStr)

	if err != nil {
		log.Println("[ERROR] PostgresReservationDAO.DeleteByAttr. Cannot connect to database:", err)
		return serverrors.ErrDatabaseConnection
	}

	defer conn.Close(context.Background())

	_, err = conn.Exec(
		context.Background(),
		`delete from reservation where $1 = $2;`,
		attrName, attrValue,
	)

	if err != nil {
		log.Println("[ERROR] PostgresReservationDAO.DeleteByAttr. Error while executing query:", err)
		return serverrors.ErrQueryExec
	}

	return nil
}
