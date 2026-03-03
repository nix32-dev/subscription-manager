package repository

import (
	"context"
	"fmt"
	"os"
	"strconv"
	"strings"

	"subscription/internal/model"
	"subscription/internal/service"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

var DataBaseConn pgx.Conn

func CreateDBConnection(ctx context.Context) error {
	// if err := godotenv.Load(); err != nil {
	// 	return err
	// }

	stringConn := fmt.Sprintf("postgres://%s:%s@%s:5432/%s?sslmode=disable",
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_HOST"),
		os.Getenv("DB_NAME"))

	conn, err := pgx.Connect(ctx, stringConn)
	if err != nil {
		return err
	}
	DataBaseConn = *conn
	return nil
}

func CreateSubscription(input model.Subscription) (error, model.Subscription) {
	if err := model.ValidateSevice(input, &DataBaseConn); err != nil {
		return err, model.Subscription{}
	}

	var existsUUID bool
	var newUUID string
	for {
		newUUID = uuid.NewString()
		if err := DataBaseConn.QueryRow(*service.CTX, `SELECT EXISTS(SELECT 1 FROM subscriptions WHERE user_id = $1)`, newUUID).Scan(&existsUUID); err != nil {
			return err, model.Subscription{}
		}
		if existsUUID {
			newUUID = uuid.NewString()
		} else {
			break
		}
	}

	var id int
	sqlRequest := `INSERT INTO subscriptions (user_id, service_name, price, start_date, end_date) VALUES ($1, $2, $3, $4, $5) RETURNING id`
	if err := DataBaseConn.QueryRow(*service.CTX, sqlRequest,
		newUUID, input.Service_name,
		input.Price, input.Start_date,
		input.End_date).Scan(&id); err != nil {
		return err, model.Subscription{}
	}
	return nil, model.Subscription{
		ID:           id,
		User_id:      newUUID,
		Service_name: input.Service_name,
		Price:        input.Price,
		Start_date:   input.Start_date,
		End_date:     input.End_date}
}

func ChangeSubscription(changes []string, getType string, id_uuid string) error {
	service.MTX.Lock()
	defer service.MTX.Unlock()
	changesOutput := strings.Join(changes, ", ")
	sqlRequest := fmt.Sprintf("UPDATE subscriptions SET %s WHERE %s = $1", changesOutput, getType)
	if _, err := DataBaseConn.Exec(*service.CTX, sqlRequest, id_uuid); err != nil {
		return err
	}
	return nil
}

func DeleteSubscription(getType string, id_uuid string) error {
	if err := model.ValidateIdExists(&DataBaseConn, id_uuid, getType); err != nil {
		return err
	}
	sqlRequest := fmt.Sprintf("DELETE FROM subscriptions WHERE %s = $1", getType)
	if _, err := DataBaseConn.Exec(*service.CTX, sqlRequest, id_uuid); err != nil {
		return err
	}
	return nil
}

func GetSubscription(id_uuid string, getType string) (model.Subscription, error) {
	service.RMTX.Lock()
	defer service.RMTX.Unlock()
	if err := model.ValidateIdExists(&DataBaseConn, id_uuid, getType); err != nil {
		return model.Subscription{}, err
	}
	var output model.Subscription
	sqlRequest := fmt.Sprintf("SELECT id, user_id, service_name, price, start_date, end_date FROM subscriptions WHERE %s = $1", getType)
	if err := DataBaseConn.QueryRow(*service.CTX, sqlRequest, id_uuid).Scan(&output.ID,
		&output.User_id,
		&output.Service_name,
		&output.Price,
		&output.Start_date,
		&output.End_date); err != nil {
		return model.Subscription{}, err
	}
	return output, nil
}

func GetSubscriptions(page string) ([]model.Subscription, error) {
	service.RMTX.Lock()
	defer service.RMTX.Unlock()
	pageInt, err := strconv.Atoi(page)
	if err != nil {
		return nil, err
	}
	if pageInt <= 0 {
		return nil, service.WrongPageNumber
	}
	page = strconv.Itoa(10 * (pageInt - 1))
	sqlRequest := fmt.Sprintf("SELECT * FROM subscriptions ORDER BY id LIMIT 10 OFFSET %s", page)

	rows, err := DataBaseConn.Query(*service.CTX, sqlRequest)
	if err != nil {
		return nil, err
	}

	defer rows.Close()
	var output []model.Subscription
	for rows.Next() {
		var input model.Subscription
		err := rows.Scan(
			&input.ID,
			&input.User_id,
			&input.Service_name,
			&input.Price,
			&input.Start_date,
			&input.End_date,
		)
		if err != nil {
			return nil, err
		}
		output = append(output, input)
	}
	return output, nil
}
