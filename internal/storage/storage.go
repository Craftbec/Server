package storage

import (
	"context"
	"errors"
	"fmt"
	"log"
	"math/rand"
	"os"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Storage interface {
	GetReport(ctx context.Context, report_id int) (string, error)
	PostReport(ctx context.Context, report_info string) error
	GetObservationTime(ctx context.Context, model_id int) (int, error)
	GracefulStopDB()
}

type DB struct {
	db *gorm.DB
}

type Report struct {
	Report_id     int    `gorm:"type:integer;primaryKey"`
	Creation_time string `gorm:"type:date;"`
	Report_info   string `gorm:"type:char;not null"`
	Model_id      int    `gorm:"type:integer;"`
}

func NewDB() (*DB, error) {
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable", os.Getenv("DB_HOST"),
		os.Getenv("DB_USER"), os.Getenv("DB_PASSWORD"), os.Getenv("DB_NAME"), os.Getenv("DB_PORT"))
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}
	err = db.AutoMigrate(&Report{})
	if err != nil {
		return nil, err
	}
	return &DB{db: db}, nil
}

func (s *DB) GetReport(ctx context.Context, report_id int) (string, error) {
	tmp := Report{Report_id: report_id}
	if result := s.db.WithContext(ctx).First(&tmp); result.RowsAffected == 0 {
		return "", errors.New("Not found")
	}
	return tmp.Report_info, nil

}

func (s *DB) PostReport(ctx context.Context, report_info string) error {
	date := time.Now()
	formatteDate := date.Format("2006-01-02")
	tmp := Report{Creation_time: formatteDate, Report_info: report_info, Model_id: rand.Intn(100)}
	result := s.db.WithContext(ctx).Create(&tmp)
	return result.Error
}

func (s *DB) GetObservationTime(ctx context.Context, model_id int) (int, error) {
	var maxDays int
	result := s.db.Raw(`WITH sorted_dates AS (SELECT Model_id, Creation_time,
	LEAD(Creation_time) OVER (PARTITION BY Model_id ORDER BY Creation_time) AS next_date
	FROM Reports
	WHERE Model_id = ?)
	SELECT COALESCE(MAX((next_date - Creation_time)::int),-1) AS max_days
	FROM sorted_dates
	WHERE next_date IS NOT NULL`, model_id).Scan(&maxDays)
	if result.Error != nil {
		return maxDays, result.Error
	}
	if maxDays == -1 {
		err := s.db.
			Model(&Report{}).
			Select("(current_date - MAX(creation_time))::int AS max_days").
			Where("model_id = ?", model_id).
			Group("model_id").
			Row().
			Scan(&maxDays)
		if err != nil {
			return maxDays, err
		}
	}
	return maxDays, nil
}

func (s *DB) GracefulStopDB() {
	links, _ := s.db.DB()
	err := links.Close()
	if err != nil {
		log.Fatal(err)
	}
}
