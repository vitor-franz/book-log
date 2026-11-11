package setup

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/vitor-franz/book-log/infrastructure/settings"
)

func OpenDatabase(provider settings.Provider) (*sql.DB, error) {
	databaseSettings := provider.GetDatabaseSettings()

	connStr := fmt.Sprintf(
		"%s:%s@tcp(%s:%d)/%s?parseTime=true",
		databaseSettings.Username,
		databaseSettings.Password,
		databaseSettings.Host,
		databaseSettings.Port,
		databaseSettings.Name,
	)

	db, err := sql.Open("mysql", connStr)
	if err != nil {
		return nil, errors.Join(errors.New("failed to open database"), err)
	}

	err = db.PingContext(context.Background())
	if err != nil {
		db.Close()
		return nil, errors.Join(errors.New("failed to connect to database"), err)
	}

	db.SetMaxOpenConns(50)
	db.SetMaxIdleConns(10)
	db.SetConnMaxLifetime(time.Minute * 10)

	return db, nil
}
