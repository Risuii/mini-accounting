package entity

import (
	"time"

	LoggingModel "mini-accounting/internal/logging/data/model"
)

type InterfaceLog struct {
	TraceID         string    `json:"traceId"`
	ServiceName     string    `json:"serviceName"`
	ClientName      string    `json:"clientName"`
	RequestPayload  string    `json:"requestPayload"`
	ResponsePayload string    `json:"responsePayload"`
	RequestDate     time.Time `json:"requestDate"`
	ResponseDate    time.Time `json:"responseDate"`
}

func (e *InterfaceLog) ToModel() *LoggingModel.InterfaceLog {
	return &LoggingModel.InterfaceLog{
		TraceID:         e.TraceID,
		ServiceName:     e.ServiceName,
		ClientName:      e.ClientName,
		RequestPayload:  e.RequestPayload,
		ResponsePayload: e.ResponsePayload,
		RequestDate:     e.RequestDate,
		ResponseDate:    e.ResponseDate,
	}
}
