package main

import (
	"log"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"

	"dbtest/pg"
)

func main() {
	db, err := gorm.Open("postgres",
		"host=myhost port=myport user=gorm dbname=gorm password=mypassword")
	if err != nil {
		panic("failed to connect database")
	}
	defer db.Close()

	repo := pg.NewRepository(db)
	if err := repo.Migrate(); err != nil {
		log.Fatal("migrate err", err)
	}
	// use the repo
	// repo.ListAll()
}
