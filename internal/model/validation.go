package model

import (
	"fmt"
	"strconv"
	"subscription/internal/service"
	"time"

	"github.com/jackc/pgx/v5"
)

func ValidateSevice(model Subscription, db *pgx.Conn) error {
	timeValidation := model.ValidateSeviceDate()
	switch {
	case model.ValidateSeviceName(db) != nil:
		return service.IncorrectLenth
	case model.ValidateSevicePrice() != nil:
		return service.IncorrectPrice
	case timeValidation != nil:
		return timeValidation
	default:
		return nil
	}
}

func ValidateIdExists(db *pgx.Conn, id_uuid string, getType string) error {
	var existsID bool
	request := fmt.Sprintf("SELECT EXISTS(SELECT 1 FROM subscriptions WHERE %s = $1)", getType)
	if err := db.QueryRow(*service.CTX, request, id_uuid).Scan(&existsID); err != nil {
		return err
	}
	if !existsID {
		return service.NotExistsByID
	}
	return nil
}

func (input *Subscription) ValidateSeviceName(db *pgx.Conn) error {
	var existsNAME bool
	if err := db.QueryRow(*service.CTX, "SELECT EXISTS(SELECT 1 FROM subscriptions WHERE service_name = $1)", *input.Service_name).Scan(&existsNAME); err != nil {
		return err
	}
	if existsNAME {
		return service.AlreadyExistsName
	}

	if len([]rune(*input.Service_name)) < 3 || len([]rune(*input.Service_name)) > 20 {
		return service.IncorrectLenth
	}
	return nil
}

func (input *Subscription) ValidateSevicePrice() error {
	if *input.Price <= 0 {
		return service.IncorrectPrice
	}
	return nil
}

func (input *Subscription) ValidateSeviceDate() error {
	timeFormat := "01-2006"
	if *input.Start_date == "" {
		return service.IncorrectDateFormat
	}
	inputDateStart, err := time.Parse(timeFormat, *input.Start_date)
	if err != nil {
		return service.IncorrectDateFormat
	}
	if *input.End_date != "" {
		if inputDateEnd, err := time.Parse(timeFormat, *input.End_date); err != nil {
			return service.IncorrectDateFormat
		} else if !inputDateEnd.After(inputDateStart) {
			return service.IncorrectEndData
		}
	}
	return nil
}

func (input *Subscription) ValidateInputExists(getType string, db *pgx.Conn, id_uuid string) ([]string, error) {
	var db_changes []string

	if err := ValidateIdExists(db, id_uuid, getType); err != nil {
		return nil, err
	}

	if input.Service_name != nil {
		if err := input.ValidateSeviceName(db); err != nil {
			return nil, err
		}
		db_changes = append(db_changes, "service_name = "+*input.Service_name)
	}

	if input.Price != nil {
		if err := input.ValidateSevicePrice(); err != nil {
			return nil, err
		}
		db_changes = append(db_changes, "price = "+strconv.Itoa(*input.Price))
	}

	if input.Start_date != nil || input.End_date != nil {
		dateStart, dateEnd, err := getDates(getType, db, id_uuid)
		if err != nil {
			return nil, err
		}
		switch {
		case input.End_date != nil && input.Start_date != nil:
			if err := input.ValidateSeviceDate(); err != nil {
				return nil, err
			}
			db_changes = append(db_changes, "start_date = '"+*input.Start_date+"'")
			db_changes = append(db_changes, "end_date = '"+*input.End_date+"'")
		case input.Start_date != nil:
			input.End_date = &dateEnd
			if err := input.ValidateSeviceDate(); err != nil {
				return nil, err
			}
			db_changes = append(db_changes, "start_date = '"+*input.Start_date+"'")
		case input.End_date != nil:
			input.Start_date = &dateStart
			if err := input.ValidateSeviceDate(); err != nil {
				return nil, err
			}
			db_changes = append(db_changes, "end_date = '"+*input.End_date+"'")
		}
	}
	return db_changes, nil
}

func getDates(getType string, db *pgx.Conn, id_uuid string) (string, string, error) {
	var startDate, endDate string
	request := fmt.Sprintf("SELECT start_date, end_date FROM subscriptions WHERE %s = $1", getType)
	if err := db.QueryRow(*service.CTX, request, id_uuid).Scan(&startDate, &endDate); err != nil {
		return "", "", err
	}
	return startDate, endDate, nil
}
