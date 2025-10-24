package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"time"
	"task/tasks"

	_ "github.com/lib/pq"
)

func main() {
	host := os.Getenv("POSTGRES_HOST")
	port := os.Getenv("POSTGRES_PORT")
	login := os.Getenv("POSTGRES_LOGIN")
	password := os.Getenv("POSTGRES_PASSWORD")
	database := os.Getenv("POSTGRES_DB")

	connStr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		host, port, login, password, database)

	var db *sql.DB
	var err error
	
	for i := 0; i < 10; i++ {
		db, err = sql.Open("postgres", connStr)
		if err != nil {
			log.Printf("Attempt %d: Failed to open connection: %v", i+1, err)
			time.Sleep(2 * time.Second)
			continue
		}
		
		err = db.Ping()
		if err != nil {
			log.Printf("Attempt %d: Failed to ping DB: %v", i+1, err)
			db.Close()
			time.Sleep(2 * time.Second)
			continue
		}
		
		break
	}
	
	if err != nil {
		log.Fatalf("Failed to connect to DB after retries: %v", err)
	}
	defer db.Close()

	fmt.Println("Connected to PostgreSQL!")

	createTableSQL := `
	CREATE TABLE IF NOT EXISTS results (
		id SERIAL PRIMARY KEY,
		task_name VARCHAR(255) NOT NULL,
		task_result TEXT NOT NULL,
		created_at TIMESTAMP DEFAULT NOW()
	);`
	
	_, err = db.Exec(createTableSQL)
	if err != nil {
		log.Fatalf("Failed to create table: %v", err)
	}
	fmt.Println("Table 'results' created or already exists!")

	fmt.Println("Task 1")
	tasks.MultiplicationTable(7)
	_, err = db.Exec("INSERT INTO results(task_name, task_result) VALUES($1, $2)", 
		"MultiplicationTable", "Completed for number 7")
	if err != nil {
		log.Printf("Failed to insert task1 result: %v", err)
	}

	fmt.Println("Task 2")
	test := []int{1, 2, 3, 4, 7, 2}
	answ2 := tasks.Increment(test)
	fmt.Printf("Increment result: %v\n", answ2)
	_, err = db.Exec("INSERT INTO results(task_name, task_result) VALUES($1, $2)", 
		"Increment", fmt.Sprintf("%v", answ2))
	if err != nil {
		log.Printf("Failed to insert task2 result: %v", err)
	}

	fmt.Println("Task 3")
	arr1 := []interface{}{115, 101, 122, 105, 122}
	result1 := tasks.IsVow(arr1)
	fmt.Printf("IsVow result: %v\n", result1)
	_, err = db.Exec("INSERT INTO results(task_name, task_result) VALUES($1, $2)", 
		"IsVow", fmt.Sprintf("%v", result1))
	if err != nil {
		log.Printf("Failed to insert task3 result: %v", err)
	}

	fmt.Println("All tasks completed and results saved to database!")
	
	fmt.Println("Application finished. Press Ctrl+C to exit.")
	select {} 
}