package analytics

import (
	"database/sql"
	"fmt"
	"github.com/go-sql-driver/mysql"
	_ "github.com/go-sql-driver/mysql"
	"log"
	"os"
	"time"
)

var db *sql.DB

func UpdateBookings(reportinDate time.Time, numDays int) {
	log.Printf("Store analytics data for date: %s ...", reportinDate.String())

	connect()

	query := "INSERT INTO `bookings_history` (`reporting_date`, `days_booked`) VALUES (?, ?);"
	insert, err := db.Prepare(query)

	if err != nil {
		log.Fatalf("Impossible insert bookings_history: %s", err)
		return
	}
	resp, err := insert.Exec(reportinDate, numDays)
	insert.Close()
  if err != nil {
      if mysqlErr, ok := err.(*mysql.MySQLError); ok && mysqlErr.Number == 1062 {
          log.Printf("Duplicate entry for reporting_date: %s, skipping insert.", reportinDate)
          return
      }
      log.Fatalf("Error while storing data: %s", err)
      return
  }

	lastInsertId, err := resp.LastInsertId()
	if err != nil {
		log.Fatalf("Error getting LastInsertId: %s", err)
	}
	fmt.Printf("LastInsertId: %d, Error: %v\n", lastInsertId, err)
}
func connect() {
	// ssh -N -L 3306:kbatchdb.k-cloud.io:3306 acharton@shell001.ek-prod.dus1.cloud

	// Capture connection properties.
	cfg := mysql.Config{
		User:                 os.Getenv("DB_USER"),
		Passwd:               os.Getenv("DB_PASS"),
		Net:                  "tcp",
		Addr:                 os.Getenv("DB_URL"),
		DBName:               "defaultdb",
		AllowNativePasswords: true,
	}
	// Get a database handle.
	var err error
	db, err = sql.Open("mysql", cfg.FormatDSN())
	if err != nil {
		log.Fatal(err)
	}

	pingErr := db.Ping()
	if pingErr != nil {
		log.Fatal(pingErr)
	}

	// For test we read user mail here.
	if err := db.QueryRow("SELECT reproting_date FROM bookings_history limit 1"); err != nil {
		fmt.Errorf("error while reading: %s", err)
	}

	fmt.Println("Connected!")
}
