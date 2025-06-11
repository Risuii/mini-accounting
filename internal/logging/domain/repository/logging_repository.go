package respository

import LoggingEntity "mini-accounting/internal/logging/domain/entity"

type LoggingRepository interface {
	InsertInterfaceLog(param LoggingEntity.InterfaceLog) error
	InsertOutgoingLog(param LoggingEntity.OutgoingLog) error
}
