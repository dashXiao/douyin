package db

import (
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
	"tiktok/pkg/constants"
	"tiktok/pkg/utils"
)

var DB *gorm.DB
var SF *utils.Snowflake

func Init() {
	var err error

	DB, err = gorm.Open(mysql.Open(utils.GetMysqlDSN()),
		&gorm.Config{
			PrepareStmt:            true,
			SkipDefaultTransaction: true,                                // Disable default transactions.
			Logger:                 logger.Default.LogMode(logger.Info), // Log database queries.
			NamingStrategy: schema.NamingStrategy{
				SingularTable: true, // Use singular table names.
			},
		})

	if err != nil {
		panic(err)
	}

	sqlDB, err := DB.DB()

	if err != nil {
		panic(err)
	}

	sqlDB.SetMaxIdleConns(constants.MaxIdleConns)       // Maximum idle connections
	sqlDB.SetMaxOpenConns(constants.MaxConnections)     // Maximum open connections
	sqlDB.SetConnMaxLifetime(constants.ConnMaxLifetime) // Maximum connection lifetime

	SF, err = utils.NewSnowflake(constants.SnowflakeDatacenterID, constants.SnowflakeWorkerID)
	if err != nil {
		panic(err)
	}
}
