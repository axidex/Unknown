package api

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
)

type ErrorCode int

// Basic errors
const (
	InternalServerErrorText = "Internal Server Error" // 500
	TaskIdErrorText         = "Provide valid TaskID"  // 400

	InternalStorageErrorText = "Cannot upload file to internal storage" // 500
	KafkaErrorText           = "Error task creating"                    // 500

	UnknownStatusErrorText = "Unknown status" // 500
)

// Scan Errors
const (
	FileProvideErrorText = "No file provided" // 400

	ArchiveTypeErrorText       = "Uploaded file must be a ZIP/TAR file"   // 400
	ArchiveSizeErrorText       = "Uploaded Archive file too large"        // 400
	ArchiveValidationErrorText = "Bad file"                               // 400
	ArchiveSaving              = "Failed saving archive file to the path" // 500

	FailedToReadReport = "Failed to read report file"

	ScanFailed = "Scan failed"
)

// Report Errors
const (
	ReportNotFoundErrorText = "Report not found" // 400
)

// Statistic Errors
const (
	CalculatingStatusesErrorText   = "Can't calculate task statuses" // 500
	CalculatingStatisticsErrorText = "Can't calculate statistics"    // 500
)

const (
	ErrTaskId ErrorCode = iota
	ErrFileNotProvided

	ErrUploadedFileTooLarge

	ErrArchiveType
	ErrArchiveValidation
	ErrArchiveSaving

	ErrFailedToReadReport

	ErrScanFailed

	ErrDatabase
	ErrInternalStorage
	ErrKafka

	ErrReportNotFound
	ErrUnknownStatus

	ErrCalculatingStatuses
	ErrCalculatingStatistics

	ErrInternalServer
)

var ErrorResponses = map[ErrorCode]ResponseError{
	ErrTaskId:          NewError(TaskIdErrorText, http.StatusBadRequest),
	ErrFileNotProvided: NewError(FileProvideErrorText, http.StatusBadRequest),

	ErrUploadedFileTooLarge: NewError(ArchiveSizeErrorText, http.StatusBadRequest),
	ErrArchiveType:          NewError(ArchiveTypeErrorText, http.StatusBadRequest),
	ErrArchiveValidation:    NewError(ArchiveValidationErrorText, http.StatusBadRequest),
	ErrArchiveSaving:        NewError(ArchiveSaving, http.StatusInternalServerError),
	ErrFailedToReadReport:   NewError(FailedToReadReport, http.StatusInternalServerError),
	ErrScanFailed:           NewError(ScanFailed, http.StatusInternalServerError),

	ErrDatabase:        NewError(InternalServerErrorText, http.StatusInternalServerError),
	ErrInternalStorage: NewError(InternalStorageErrorText, http.StatusInternalServerError),
	ErrKafka:           NewError(KafkaErrorText, http.StatusInternalServerError),

	ErrReportNotFound: NewError(ReportNotFoundErrorText, http.StatusNotFound),
	ErrUnknownStatus:  NewError(UnknownStatusErrorText, http.StatusInternalServerError),

	ErrCalculatingStatuses:   NewError(CalculatingStatusesErrorText, http.StatusInternalServerError),
	ErrCalculatingStatistics: NewError(CalculatingStatisticsErrorText, http.StatusInternalServerError),

	ErrInternalServer: NewError(InternalServerErrorText, http.StatusInternalServerError),
}

type ResponseError struct {
	ErrorResponse string `json:"error"`
	StatusCode    int    `json:"-"`
}

var UnknownErrorResponse = ResponseError{
	ErrorResponse: InternalServerErrorText,
	StatusCode:    http.StatusInternalServerError,
}

func ResponseJSON(c *gin.Context, errorCode ErrorCode) {
	response, valid := ErrorResponses[errorCode]
	if !valid {
		response = UnknownErrorResponse
	}
	c.JSON(response.StatusCode, response)
}

func (e *ResponseError) Error() string {
	return fmt.Sprintf("Error: %s (Code: %d)", e.ErrorResponse, e.StatusCode)
}

func NewError(message string, code int) ResponseError {
	return ResponseError{
		ErrorResponse: message,
		StatusCode:    code,
	}
}
