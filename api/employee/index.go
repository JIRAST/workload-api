package handler

import (
	"context"
	"encoding/json"
	"net/http"
	"os"

	"github.com/jackc/pgx/v5"
)

type Employee struct {
	ID           int    `json:"id"`
	FirstName    string `json:"first_name"`
	LastName     string `json:"last_name"`
	PositionName string `json:"position_name"`
	IsActive     bool   `json:"is_active"`
}

func Handler(w http.ResponseWriter, r *http.Request) {
	dbUrl := os.Getenv("DB_URL")
	ctx := context.Background()
	conn, err := pgx.Connect(ctx, dbUrl)
	if err != nil {
		http.Error(w, "Database connection failed", 500)
		return
	}
	defer conn.Close(ctx)

	switch r.Method {
	case "GET":
		// ดึงข้อมูลตามชื่อคอลัมน์จริงใน Schema
		rows, _ := conn.Query(ctx, "SELECT id, first_name, last_name, position_name, is_active FROM employees ORDER BY id DESC")
		var employees []Employee
		for rows.Next() {
			var e Employee
			rows.Scan(&e.ID, &e.FirstName, &e.LastName, &e.PositionName, &e.IsActive)
			employees = append(employees, e)
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(employees)

	case "POST":
		var e Employee
		json.NewDecoder(r.Body).Decode(&e)
		// เพิ่มข้อมูลใหม่ (สมมติว่ายังไม่ระบุ department_id)
		_, err = conn.Exec(ctx, "INSERT INTO employees (first_name, last_name, position_name, is_active) VALUES ($1, $2, $3, $4)",
			e.FirstName, e.LastName, e.PositionName, true)
		if err != nil {
			http.Error(w, err.Error(), 400)
			return
		}
		w.WriteHeader(http.StatusCreated)

	case "PUT":
		var e Employee
		json.NewDecoder(r.Body).Decode(&e)
		// อัปเดตข้อมูลตาม ID
		_, err = conn.Exec(ctx, "UPDATE employees SET first_name=$1, last_name=$2, position_name=$3, is_active=$4 WHERE id=$5",
			e.FirstName, e.LastName, e.PositionName, e.IsActive, e.ID)
		if err != nil {
			http.Error(w, err.Error(), 400)
			return
		}
		w.WriteHeader(http.StatusOK)
	}
}
