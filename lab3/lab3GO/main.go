package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"
	"task/tasks"
	"time"

	"github.com/go-redis/redis/v8"
	_ "github.com/lib/pq"
)

func main() {
	host := os.Getenv("POSTGRES_HOST")
	port := os.Getenv("POSTGRES_PORT")
	login := os.Getenv("POSTGRES_LOGIN")
	password := os.Getenv("POSTGRES_PASSWORD")
	database := os.Getenv("POSTGRES_DB")

	redisHost := os.Getenv("REDIS_HOST")
	redisPort := os.Getenv("REDIS_PORT")
	redisPassword := os.Getenv("REDIS_PASSWORD")

	redisAddr := fmt.Sprintf("%s:%s", redisHost, redisPort)

	if redisAddr == "" {
		log.Fatal("REDIS_ADDR не задан")
	}

	rdb := redis.NewClient(&redis.Options{
		Addr:     redisAddr,
		Password: redisPassword,
		DB:       0,
	})

	ctx := context.Background()
	err1 := rdb.Ping(ctx).Err()
	if err1 != nil {
		log.Fatalf("Не удалось подключиться к Redis: %v", err1)
	}

	fmt.Println("Task 1 - MultiplicationTable")
	tasks.MultiplicationTable(7)
	key1 := fmt.Sprintf("task1_%v", time.Now().Unix())
	
	// 🔥 ИСПРАВЛЕНИЕ: Добавь .Err()
	err5 := rdb.Set(ctx, key1, "completed for 7", 10*time.Minute).Err()
	if err5 != nil {
		log.Fatalf("Cannot write the answer for task1: %v", err5)
	} else {
		fmt.Printf("✅ Task1 записан в Redis: %s\n", key1)
	}

	// Task 2
	fmt.Println("Task 2 - Increment")
	test2 := []int{1, 2, 3, 4, 7, 2}
	res2 := tasks.Increment(test2)
	key2 := fmt.Sprintf("task2_%v", time.Now().Unix())
	
	// 🔥 ИСПРАВЛЕНИЕ: Добавь .Err()
	err6 := rdb.Set(ctx, key2, fmt.Sprintf("%v", res2), 10*time.Minute).Err()
	if err6 != nil {
		log.Fatalf("Cannot write the answer for task2: %v", err6)
	} else {
		fmt.Printf("✅ Task2 записан в Redis: %s = %v\n", key2, res2)
	}

	// Task 3
	fmt.Println("Task 3 - IsVow")
	arr2 := []interface{}{115, 101, 122, 105, 122}
	res3 := tasks.IsVow(arr2)
	key3 := fmt.Sprintf("task3_%v", time.Now().Unix())
	
	// 🔥 ИСПРАВЛЕНИЕ: Добавь .Err()
	err7 := rdb.Set(ctx, key3, fmt.Sprintf("%v", res3), 10*time.Minute).Err()
	if err7 != nil {
		log.Fatalf("Cannot write the answer for task3: %v", err7)
	} else {
		fmt.Printf("✅ Task3 записан в Redis: %s = %v\n", key3, res3)
	}
	
	keys, err22 := rdb.Keys(ctx, "task*").Result()
	if err22 != nil {
		log.Printf("Cannot to get the keys")
	} else {
		fmt.Printf("Founded keys: %v", keys)
		for _, key := range keys {
			val, err := rdb.Get(ctx, key).Result()
			if err != nil {
				fmt.Printf("Error in reading the key %s: %v", key, err)
			} else {
				fmt.Printf("%s = %v\n", key, val)
			}
		}
	}

	defer rdb.Close()

	fmt.Println("Подключение к Redis установлено")

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
	_, err = db.Exec("INSERT INTO results(task_name, task_result) VALUES($1, $2)",
		"Increment", fmt.Sprintf("%v", answ2))
	if err != nil {
		log.Printf("Failed to insert task2 result: %v", err)
	}

	fmt.Println("Task 3")
	arr1 := []interface{}{115, 101, 122, 105, 122}
	result1 := tasks.IsVow(arr1)
	_, err = db.Exec("INSERT INTO results(task_name, task_result) VALUES($1, $2)",
		"IsVow", fmt.Sprintf("%v", result1))
	if err != nil {
		log.Printf("Failed to insert task3 result: %v", err)
	}

	fmt.Println("All tasks completed and results saved to database!")

	fmt.Println("Application finished. Press Ctrl+C to exit.")
	select {}
}
