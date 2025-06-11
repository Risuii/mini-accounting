package entity

import (
	"time"

	LoggingModel "mini-accounting/internal/logging/data/model"
)

type OutgoingLog struct {
	TraceID         string    `json:"trace_id"`
	BackendSystem   string    `json:"backend_system"`
	ServiceName     string    `json:"service_name"`
	HttpStatus      int       `json:"http_status"`
	RequestPayload  string    `json:"request_payload"`
	ResponsePayload string    `json:"response_payload"`
	RequestDate     time.Time `json:"request_date"`
	ResponseDate    time.Time `json:"response_date"`
}

func (e *OutgoingLog) ToModel() *LoggingModel.OutgoingLog {
	return &LoggingModel.OutgoingLog{
		TraceID:         e.TraceID,
		BackendSystem:   e.BackendSystem,
		ServiceName:     e.ServiceName,
		HttpStatus:      e.HttpStatus,
		RequestPayload:  e.RequestPayload,
		ResponsePayload: e.ResponsePayload,
		RequestDate:     e.RequestDate,
		ResponseDate:    e.ResponseDate,
	}
}
