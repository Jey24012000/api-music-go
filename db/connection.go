package db

import (
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

//Conecta la base de datos.
func Connect() {
	var db *gorm.DB
	var err error

	dsn := "host=postgres user=songs password=songs dbname=songs port=5432"
	db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})

	if err != nil {
		panic(err)
	}
	DB = db
}


