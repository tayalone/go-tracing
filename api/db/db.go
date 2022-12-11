package db

import (
	"fmt"
	"log"
	"os"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/plugin/opentelemetry/tracing"
)

// Todo is db schema `barcode_condition`
type Todo struct {
	ID        uint   `gorm:"primaryKey" json:"id"`
	UserID    uint   `gorm:"not null;index;" json:"userId"`
	Title     string `gorm:"not null;" json:"title"`
	Completed bool   `gorm:"not null;index;" json:"completed"`

	CreatedAt time.Time `gorm:"index;autoCreateTime" json:"createdAt"`
	UpdatedAt time.Time `gorm:"index;autoUpdateTime" json:"updatedAt"`
}

func migrate(db *gorm.DB) {
	db.Set("gorm:table_options", "ENGINE=InnoDB")

	// // /  about 'barcode_condition'
	if (db.Migrator().HasTable(&Todo{})) {
		log.Println("Table Existing, Drop IT")

		db.Migrator().DropTable(&Todo{})
	}
	db.AutoMigrate(&Todo{})
	log.Println("Create 'todo'")

	// / Add Initail Data
	initTodos := []Todo{
		{
			UserID:    1,
			Title:     "delectus aut autem",
			Completed: false,
		},
	}
	db.Create(initTodos)

	log.Println("Create Initial 'Todo' Data")
}

/*New do Create Rdb Connection*/
func New() (*gorm.DB, error) {
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=%s",
		os.Getenv("RDM_HOST"),
		os.Getenv("RDM_USER"),
		os.Getenv("RDM_PASSWORD"),
		os.Getenv("RDM_DB"),
		os.Getenv("RDM_PORT"),
		os.Getenv("TIME_ZONE"))

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {

		log.Println("FAIL: Connect RDB Error", err.Error())

		return nil, err

	}
	log.Println("Connect RDB Success!!!")

	// --------- setup trancing ----------------------------
	if err := db.Use(tracing.NewPlugin(tracing.WithoutMetrics())); err != nil {
		panic(err)
	}
	// ----------------------------------------------------
	log.Println("db is added tracing plugin!!!")

	migrate(db)

	return db, nil
}
