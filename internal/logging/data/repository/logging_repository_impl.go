package repository

import (
	LoggingSource "mini-accounting/internal/logging/data/source"
	LoggingEntity "mini-accounting/internal/logging/domain/entity"
	LoggingRepository "mini-accounting/internal/logging/domain/repository"
	Library "mini-accounting/library"
	CustomErrorPackage "mini-accounting/pkg/custom_error"
)

type LoggingRepositoryImpl struct {
	loggingSource LoggingSource.LoggingPersistent
	library       Library.Library
}

func NewLoggingRepository(
	loggingSource LoggingSource.LoggingPersistent,
	library Library.Library,
) LoggingRepository.LoggingRepository {
	return &LoggingRepositoryImpl{
		loggingSource: loggingSource,
		library:       library,
	}
}

func (r *LoggingRepositoryImpl) InsertInterfaceLog(param LoggingEntity.InterfaceLog) error {
	path := "LoggingRepository:InsertInterfaceLog"
	err := r.loggingSource.InsertInterfaceLog(*param.ToModel())
	if err != nil {
		return err.(*CustomErrorPackage.CustomError).UnshiftPath(path)
	}

	return nil
}

func (r *LoggingRepositoryImpl) InsertOutgoingLog(param LoggingEntity.OutgoingLog) error {
	path := "LoggingRepository:InsertOutgoingLog"
	err := r.loggingSource.InsertOutgoingLog(*param.ToModel())
	if err != nil {
		return err.(*CustomErrorPackage.CustomError).UnshiftPath(path)
	}

	return nil
}
