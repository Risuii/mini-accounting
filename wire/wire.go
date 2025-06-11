//go:build wireinject
// +build wireinject

package wire

import (
	"github.com/gin-gonic/gin"
	"github.com/google/wire"

	Config "mini-accounting/config"
	Library "mini-accounting/library"

	LoggingRepository "mini-accounting/internal/logging/data/repository"
	LoggingSource "mini-accounting/internal/logging/data/source"

	CountingHandler "mini-accounting/internal/counting/delivery/presenter/http"
	CountingUsecase "mini-accounting/internal/counting/domain/usecase"

	Middleware "mini-accounting/middlewares"

	CustomValidationPackage "mini-accounting/pkg/custom_validation"
	AccountingDBPackage "mini-accounting/pkg/data_sources/accounting_db"

	Routes "mini-accounting/routes"
)

var ProviderSet = wire.NewSet(
	// FRAMEWORK
	gin.New,
	// PACKAGE
	CustomValidationPackage.NewCustomValidation,
	// DATABASE
	AccountingDBPackage.New,

	// DATASOURCE
	LoggingSource.NewLoggingPersistent,

	// REPOSITORY
	LoggingRepository.NewLoggingRepository,

	// USECASE
	CountingUsecase.NewCountingUsecase,

	// HANDLER
	CountingHandler.NewCountingHandler,

	// MIDDLEWARE
	Middleware.NewMiddleware,
	// ROUTE
	Routes.New,
)

func InjectRoute(config Config.Config, library Library.Library) Routes.Routes {
	wire.Build(
		ProviderSet,
	)
	return nil
}
