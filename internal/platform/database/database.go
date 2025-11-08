package database

import (
	"os"

	"github.com/Dzox13524/PA_Funcpro_Kel12/internal/domain"
	"github.com/Dzox13524/PA_Funcpro_Kel12/internal/middleware"
	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
)

func NewConnection() *gorm.DB {
	dbPath := os.Getenv("DB_PATH")
	if dbPath == "" {
		dbPath = "database.db"
	}
	db, err := gorm.Open(sqlite.Open(dbPath), &gorm.Config{})
	
	if err != nil {
		// log.Println("k",err)
		// middleware.HandleError("A" + err)
	}
	middleware.HandleLog("Sukses terhubung ke database SQLite di: " + dbPath)

	err = db.AutoMigrate(
        &domain.User{},
    )
	if err != nil {
        // middleware.HandleError(err)
    }
	return db
}