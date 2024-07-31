package config

import (
	"fmt"
	"log"
	"time"

	"github.com/kelseyhightower/envconfig"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
)

// mysql config
type database struct {
	// SSLMode to enable/disable SSL connection
	MysqlSSLMode bool `envconfig:"MYSQL_SSL_MODE" default:"true"`
	// MaxIdleConnection to set max idle connection pooling
	MysqlMaxIdleConnection int `envconfig:"MYSQL_MAX_IDLE_CONNECTION" default:"10"`
	// MaxOpenConnection to set max open connection pooling
	MysqlMaxOpenConnection int `envconfig:"MYSQL_MAX_OPEN_CONNECTION" default:"50"`
	// MaxLifetimeConnectionn to set max lifetime of pooling | minutes unit
	MysqlMaxLifetimeConnection int `envconfig:"MYSQL_MAX_LIFETIME_CONNECTION" default:"10"`
	// Host is host of mysql service
	MysqlHost string `envconfig:"MYSQL_HOST" default:"localhost"`
	// Port is port of mysql service
	MysqlPort string `envconfig:"MYSQL_PORT" default:"3306"`
	// Username is name of registered user in mysql service
	MysqlUsername string `envconfig:"MYSQL_USERNAME" default:""`
	// DBName is name of registered database in mysql service
	MysqlDBName string `envconfig:"MYSQL_DB_NAME" default:""`
	// Password is password of used Username in mysql service
	MysqlPassword string `envconfig:"MYSQL_PASSWORD" default:""`
	// LogMode is toggle to enable/disable log query in your service by default false
	MysqlLogMode int `envconfig:"MYSQL_LOG_MODE" default:"0"`
	// ParseTime to parse to local time
	MysqlParseTime bool `envconfig:"MYSQL_PARSE_TIME" default:"true"`
	// Charset to define charset of database
	MysqlCharset string `envconfig:"MYSQL_CHARSET" default:"utf8mb4"`
	// Charset to define charset of database
	MysqlLoc string `envconfig:"MYSQL_LOC" default:"Local"`
}

func Connect() *gorm.DB {
	var dbConfig database
	envconfig.MustProcess("", &dbConfig)

	// construct connection string
	dsn := fmt.Sprintf(
		"%s:%s@tcp(%s:%s)/%s?charset=%s&parseTime=%+v&loc=%s",
		dbConfig.MysqlUsername,
		dbConfig.MysqlPassword,
		dbConfig.MysqlHost,
		dbConfig.MysqlPort,
		dbConfig.MysqlDBName,
		dbConfig.MysqlCharset,
		dbConfig.MysqlParseTime,
		dbConfig.MysqlLoc)
	log.Println(dsn)

	// open mysql connection
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.LogLevel(dbConfig.MysqlLogMode)),
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true,
		},
	})
	if err != nil {
		panic(err)
	}

	// set configuration pooling connection
	mysqlDb, _ := db.DB()
	mysqlDb.SetMaxOpenConns(dbConfig.MysqlMaxOpenConnection)
	mysqlDb.SetConnMaxLifetime(time.Duration(dbConfig.MysqlMaxLifetimeConnection) * time.Minute)
	mysqlDb.SetMaxIdleConns(dbConfig.MysqlMaxIdleConnection)

	return db
}
