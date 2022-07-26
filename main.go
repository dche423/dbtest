package main

import (
	"log"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"dbtest/pg"
)

func main() {
	db, err := gorm.Open(postgres.Open("host=myhost port=myport user=gorm dbname=gorm password=mypassword"), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	repo := pg.Repository{Db: db}
	if err := repo.Migrate(); err != nil {
		log.Fatal("migrate err", err)
	}
	// use the repo
	// repo.ListAll()
}
