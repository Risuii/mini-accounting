package source

import (
	Constants "mini-accounting/constants"

	Library "mini-accounting/library"

	LoggingModel "mini-accounting/internal/logging/data/model"

	CustomErrorPackage "mini-accounting/pkg/custom_error"
	AccountingDBPackage "mini-accounting/pkg/data_sources/accounting_db"
)

type LoggingPersistent interface {
	InsertInterfaceLog(param LoggingModel.InterfaceLog) error
	InsertOutgoingLog(param LoggingModel.OutgoingLog) error
}

type LoggingPersistentImpl struct {
	AccountingDBPackage AccountingDBPackage.AccountingDB
	library             Library.Library
}

func NewLoggingPersistent(
	AccountingDBPackage AccountingDBPackage.AccountingDB,
	library Library.Library,
) LoggingPersistent {
	return &LoggingPersistentImpl{
		AccountingDBPackage: AccountingDBPackage,
		library:             library,
	}
}

func (p *LoggingPersistentImpl) InsertInterfaceLog(param LoggingModel.InterfaceLog) error {
	path := "LoggingAppPersistent:InsertInterfaceLog"

	err := p.AccountingDBPackage.GetConnection().Exec(
		`INSERT INTO interface_log (
			trace_id, 
			service_name, 
			client_name, 
			request_payload, 
			response_payload, 
			request_date, 
			response_date
			) VALUES (?, ?, ?, ?, ?, ?, ?)`,
		param.TraceID, param.ServiceName, param.ClientName, param.RequestPayload, param.ResponsePayload, param.RequestDate, param.ResponseDate).Error

	if err == nil {
		return nil
	}

	return CustomErrorPackage.New(Constants.ErrInternalServerError, err, path, p.library)
}

func (p *LoggingPersistentImpl) InsertOutgoingLog(param LoggingModel.OutgoingLog) error {
	path := "LoggingAppPersistent:InsertOutgoingLog"

	err := p.AccountingDBPackage.GetConnection().Exec(
		`INSERT INTO outgoing_log (
			trace_id, 
			backend_system, 
			service_name, 
			http_status, 
			request_payload, 
			response_payload, 
			request_date, 
			response_date
			) VALUES (?, ?, ?, ?, ?, ?, ?, ?)`,
		param.TraceID, param.BackendSystem, param.ServiceName, param.HttpStatus,
		param.RequestPayload, param.ResponsePayload, param.RequestDate, param.ResponseDate).Error

	if err == nil {
		return nil
	}

	return CustomErrorPackage.New(Constants.ErrInternalServerError, err, path, p.library)
}
