package serverrors

import "errors"

// Startup errors
var ErrConfigRead error = errors.New(`error while reading config file`)

var ErrServiceBuild error = errors.New(`error while building service`)
var ErrControllerPrepare error = errors.New(`error while preparing controller`)

// Implementation errors
var ErrMethodIsNotImplemented error = errors.New(`method is not implemented`)

// Database errors
var ErrDatabaseConnection error = errors.New(`cannot connect to database`)
var ErrEntityInsert error = errors.New(`error while inserting new entity`)
var ErrQueryResRead error = errors.New(`error while reading SQL-query result`)
var ErrQueryExec error = errors.New(`error while executing SQL-query`)

// Data errors
var ErrInvalidPagesData error = errors.New(`invalid pages data`)
var ErrInvalidUsername error = errors.New(`invalid username`)
var ErrInvalidCrReservReq error = errors.New(`invalid create reservation request data`)

var ErrInvalidReservUid error = errors.New(`invalid reservation UID`)

var ErrInvalidReservUsername error = errors.New(`invalid reservation username field`)
var ErrInvalidReservPayUID error = errors.New(`invalid reservation payment UID field`)
var ErrInvalidReservHotelId error = errors.New(`invalid reservation hotel ID field`)
var ErrInvalidReservStatus error = errors.New(`invalid reservation status field`)
var ErrInvalidReservDates error = errors.New(`invalid reservation dates`)

var ErrInvalidHotelId error = errors.New(`invalid hotel ID field`)
var ErrInvalidHoteUid error = errors.New(`invalid hotel UID`)

var ErrInvalidPaymentUid error = errors.New(`invalid payment UID`)
var ErrInvalidPaymentStatus error = errors.New(`invalid payment status`)
var ErrInvalidPaymentPrice error = errors.New(`invalid payment price`)

// Result errors
var ErrEntityNotFound error = errors.New(`entity not found in database`)
var ErrReservNotFound error = errors.New(`reservation not found`)
var ErrHotelNotFound error = errors.New(`hotel not found`)
var ErrPaymentNotFound error = errors.New(`payment not found`)
var ErrLoyaltyNotFound error = errors.New(`loyalty not found`)

// HTTP errors
var ErrNewRequestForming error = errors.New(`error while creating new request`)
var ErrRequestSend error = errors.New(`error while sending request to service`)
var ErrResponseRead error = errors.New(`error while reading service response`)
var ErrResponseParse error = errors.New(`error while parsing service response`)

// Internal errors
var ErrJSONParse error = errors.New(`error while writting entity into json`)

// Unknown :P
var ErrUnknown error = errors.New(`unknown error`)
