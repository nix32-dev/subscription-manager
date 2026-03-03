package service

import "errors"

var WrongMethod error = errors.New("Wrong request method!")
var IncorrectLenth error = errors.New("Incorrect service name length!")
var IncorrectPrice error = errors.New("Incorrect service price!")
var IncorrectDateFormat error = errors.New("Incorrect date format!")
var IncorrectEndData error = errors.New("Incorrect date: end_date must be after start_date")
var AlreadyExistsName error = errors.New("Service with this name already exists!")
var NotExistsID_UUID error = errors.New("Subscription ID or UUID not specified!")
var NotExistsByID error = errors.New("Subsription with this ID/UUID not exists!")
var WrongPageNumber error = errors.New("Wrong page number!")
