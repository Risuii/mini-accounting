package model

import "time"

type InterfaceLog struct {
	TraceID         string    `DB:"trace_id"`
	ServiceName     string    `DB:"service_name"`
	ClientName      string    `DB:"client_name"`
	RequestPayload  string    `DB:"request_payload"`
	ResponsePayload string    `DB:"response_payload"`
	RequestDate     time.Time `DB:"request_date"`
	ResponseDate    time.Time `DB:"response_date"`
}
