package model

import "time"

type OutgoingLog struct {
	TraceID         string    `DB:"trace_id"`
	BackendSystem   string    `DB:"backend_system"`
	ServiceName     string    `DB:"service_name"`
	HttpStatus      int       `DB:"http_status"`
	RequestPayload  string    `DB:"request_payload"`
	ResponsePayload string    `DB:"response_payload"`
	RequestDate     time.Time `DB:"request_date"`
	ResponseDate    time.Time `DB:"response_date"`
}
