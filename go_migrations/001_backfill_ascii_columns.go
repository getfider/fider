package main

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/getfider/fider/app/pkg/dbx"
	"github.com/joho/godotenv"
	"github.com/mozillazg/go-unidecode"
)

func main() {
	godotenv.Load()
	db := dbx.Connection()
	rows, err := db.Query("SELECT id, title, description FROM posts")
	if err != nil {
		log.Fatalf("Failed to query posts: %v", err)
	}
	defer rows.Close()

	for rows.Next() {
		var id int
		var title, description sql.NullString
		if err := rows.Scan(&id, &title, &description); err != nil {
			log.Printf("Failed to scan post %d: %v", id, err)
			continue
		}

		titleAscii := ""
		descriptionAscii := ""
		if title.Valid {
			titleAscii = unidecode.Unidecode(title.String)
		}
		if description.Valid {
			descriptionAscii = unidecode.Unidecode(description.String)
		}

		_, err := db.Exec("UPDATE posts SET title_ascii = $1, description_ascii = $2 WHERE id = $3", titleAscii, descriptionAscii, id)
		if err != nil {
			log.Printf("Failed to update post %d: %v", id, err)
		} else {
			fmt.Printf("Updated post %d\n", id)
		}
	}

	if err := rows.Err(); err != nil {
		log.Fatalf("Row error: %v", err)
	}

	fmt.Println("Backfill complete.")
}
