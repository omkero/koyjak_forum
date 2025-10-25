package main

import (
	"fmt"
	"log"
	"os/exec"
	"strings"
)

func main() {
	database_name := "koyjak"
	cmd := exec.Command("pg_restore", "--clean", "-U", "postgres", "-d", database_name, "./migrations/database.sql")
	msg_db_exists := []string{"database", database_name, "not", "exist"}

	output, err := cmd.CombinedOutput()
	if err != nil {
		msg := string(output)
		log.Printf("Output:\n %s", msg)

		for _, ms := range msg_db_exists {
			if strings.Contains(msg, ms) {
				var s string
				fmt.Printf("Would you like to create the database ? (Y,n) : ")
				fmt.Scanf("%s", &s)

				if s == "Y" || s == "y" {
					create_db(database_name)
				}
				break
			}
		}
	}

	fmt.Println("pg_restore has loaded .sql successfully")
}

func create_db(db_name string) {
	// createdb -U postgres mydatabase
	cmd := exec.Command("createdb", "-U", "postgres", db_name)
	output, err := cmd.CombinedOutput()
	if err != nil {
		log.Printf("Output: %s", string(output))
	}

	fmt.Printf("database created %s", db_name)
}
