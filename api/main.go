package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-sql-driver/mysql"
)

// db -> notes_api
const (
	DBUSER     = "mysql"
	DBPASSWORD = "mysql"
	DBNET      = "tcp"
	DBHOST     = "127.0.0.1:3306"
	DBPORT     = "3306"
	DBNAME     = "notes_api"
)

type Note struct {
	ID      int    `json:"id"`
	Title   string `json:"title"`
	Content string `json:"content"`
}

var db *sql.DB

func main() {
	cfg := mysql.Config{
		User:   DBUSER,
		Passwd: DBPASSWORD,
		Net:    DBNET,
		Addr:   DBHOST,
		DBName: DBNAME,
	}

	var e error
	db, e = sql.Open("mysql", cfg.FormatDSN())
	if e != nil {
		log.Fatal(e)
	}

	var p error
	p = db.Ping()
	if p != nil {
		log.Fatal(p)
	}

	fmt.Println("Database Connection: Success!")
	defer db.Close()

	router := gin.Default()

	router.GET("/notes", getNotes)
	router.GET("/notes/:id", getNoteByID)

	router.Run(":8000")
}

func getNotes(c *gin.Context) {
	rows, e := db.Query("SELECT * FROM notes")
	if e != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch notes"})
		return
	}

	defer rows.Close()

	var notes []Note
	for rows.Next() {
		var note Note
		if e := rows.Scan(&note.ID, &note.Title, &note.Content); e != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to scan notes"})
			return
		}
		notes = append(notes, note)
	}

	c.JSON(http.StatusOK, notes)
}

func getNoteByID(c *gin.Context) {
	id := c.Param("id")
	var note Note

	// populate note item if query is successful
	e := db.QueryRow(
		"SELECT id, title, content FROM notes WHERE id = ?", id,
	).Scan(&note.ID, &note.Title, &note.Content)
	if e != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Failed to fetch note"})
		return
	}

	c.JSON(http.StatusOK, note)
}
