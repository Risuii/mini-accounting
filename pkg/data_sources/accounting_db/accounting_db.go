package uko_db

// import (
// 	"fmt"
// 	"log"
// 	"sync"
// 	"time"

// 	"gorm.io/gorm"
// 	"gorm.io/gorm/logger"

// 	Config "mini-accounting/config"
// 	Constants "mini-accounting/constants"
// 	Library "mini-accounting/library"
// 	CustomErrorPackage "mini-accounting/pkg/custom_error"

// 	GormPostgreServer "gorm.io/driver/postgres"
// )

// type AccountingDB interface {
// 	GetConnection() *gorm.DB
// }

// type AccountingDBImpl struct {
// 	connection *gorm.DB
// 	config     Config.Config
// 	setup      string
// 	library    Library.Library
// }

// var (
// 	AccountingDBInstance AccountingDBImpl
// 	AccountingDBOnce     sync.Once
// )

// // New creates a new instance of InfinitiumImpl and initializes the database connection.
// func New(
// 	config Config.Config,
// 	library Library.Library,
// ) AccountingDB {
// 	AccountingDBOnce.Do(func() {
// 		path := "AccountingDB:New"
// 		// Setup configuration string
// 		setup := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=Asia/Jakarta",
// 			config.GetConfig().DBAccounting.HOST,
// 			config.GetConfig().DBAccounting.USER,
// 			config.GetConfig().DBAccounting.PASSWORD,
// 			config.GetConfig().DBAccounting.DATABASE,
// 			config.GetConfig().DBAccounting.PORT,
// 		)

// 		// Open connection
// 		connection, err := gorm.Open(GormPostgreServer.Open(setup), &gorm.Config{
// 			Logger: logger.Default.LogMode(logger.Silent),
// 		})

// 		// Handle connection error
// 		if err != nil {
// 			err = CustomErrorPackage.New(Constants.ErrConnectionPostgres, err, path, library)
// 			log.Println(err)
// 			return
// 		}

// 		// Get underlying sql.DB object
// 		db, err := connection.DB()
// 		if err != nil {
// 			err = CustomErrorPackage.New(Constants.ErrConnectionPostgres, err, path, library)
// 			log.Println(err)
// 			return
// 		}

// 		// Configure database connection pool
// 		db.SetMaxIdleConns(config.GetConfig().DB.MaxIdleConns)
// 		db.SetMaxOpenConns(config.GetConfig().DB.MaxOpenConns)
// 		db.SetConnMaxLifetime(time.Duration(config.GetConfig().DB.ConnMaxLifetime) * time.Second)

// 		// Initialize singleton instance
// 		AccountingDBInstance = AccountingDBImpl{
// 			connection: connection,
// 			config:     config,
// 			setup:      setup,
// 			library:    library,
// 		}
// 	})
// 	return &AccountingDBInstance
// }

// // GetConnection returns the database connection instance.
// func (c *AccountingDBImpl) GetConnection() *gorm.DB {
// 	return c.connection
// }
