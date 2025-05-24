package functions

import (
	"fmt"
	"log"
)

// this function must be executed before any db query to check if pool established
func Failed_db_connection() {
	log.Fatal("Error cannot establish database connection try again")
}

func Something_wnt_wrong() error {
	return fmt.Errorf("Ops Something Went Wrong !!")
}
