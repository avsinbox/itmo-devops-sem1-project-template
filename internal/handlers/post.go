package handlers

import (
	"archive/zip"
	"bytes"
	"database/sql"
	"encoding/csv"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"path/filepath"
	"strconv"
	"time"

	"itmo-devops-sem1-project-template/internal/models"
)

func PostPrices(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Check if the request method is POST, otherwise return an error
		if r.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		// Retrieve the uploaded file from the request
		file, _, err := r.FormFile("file")
		if err != nil {
			log.Printf("Error retrieving file: %v", err)
			http.Error(w, "Failed to retrieve file", http.StatusBadRequest)
			return
		}
		defer file.Close()

		// Read the ZIP file into a buffer
		zipBuffer := &bytes.Buffer{}
		if _, err := io.Copy(zipBuffer, file); err != nil {
			log.Printf("Error reading file: %v", err)
			http.Error(w, "Failed to read file", http.StatusInternalServerError)
			return
		}

		// Open the ZIP archive
		zipReader, err := zip.NewReader(bytes.NewReader(zipBuffer.Bytes()), int64(zipBuffer.Len()))
		if err != nil {
			log.Printf("Error opening zip: %v", err)
			http.Error(w, "Invalid zip file", http.StatusBadRequest)
			return
		}

		var validItems []models.Item

		for _, zipFile := range zipReader.File {
			if filepath.Ext(zipFile.Name) != ".csv" {
				continue
			}
			// Open the CSV file inside the ZIP
			csvFile, err := zipFile.Open()
			if err != nil {
				log.Printf("Error opening data.csv: %v", err)
				http.Error(w, "Error opening data.csv", http.StatusInternalServerError)
				return
			}
			defer csvFile.Close()

			csvReader := csv.NewReader(csvFile)

			// Skip the header row
			_, err = csvReader.Read()
			if err != nil {
				log.Printf("Error reading data.csv: %v", err)
				http.Error(w, "Error reading data.csv", http.StatusBadRequest)
				return
			}

			// Read and process CSV records
			for {
				record, err := csvReader.Read()
				if err == io.EOF {
					break
				}
				if err != nil {
					log.Printf("Error reading CSV record: %v", err)
					continue
				}
				if len(record) < 5 {
					log.Printf("Invalid record: %v", record)
					continue
				}

				idStr := record[0]
				name := record[1]
				category := record[2]
				priceStr := record[3]
				dateStr := record[4]

				// Convert id string to int64
				id, err := strconv.ParseInt(idStr, 10, 64)
				if err != nil {
					log.Printf("Invalid id: %v", priceStr)
					continue
				}

				// Convert price string to float64
				price, err := strconv.ParseFloat(priceStr, 64)
				if err != nil {
					log.Printf("Invalid price: %v", priceStr)
					continue
				}

				// Parse the date string
				createDate, err := time.Parse("2006-01-02", dateStr)
				if err != nil {
					log.Printf("Invalid date: %v", dateStr)
					continue
				}

				validItems = append(validItems, models.Item{
					ID:         id,
					Name:       name,
					Category:   category,
					Price:      price,
					CreateDate: createDate,
				})
			}
		}

		// Begin database transaction
		tx, err := db.Begin()
		if err != nil {
			log.Printf("Failed to begin transaction: %v", err)
			http.Error(w, "Failed to begin transaction", http.StatusInternalServerError)
			return
		}
		defer func() { _ = tx.Rollback() }()

		var validItemsCount int

		// Insert valid items into the database
		for _, item := range validItems {
			_, err := tx.Exec(`
                INSERT INTO prices (id, name, category, price, create_date)
                VALUES ($1, $2, $3, $4, $5)
                ON CONFLICT (id) DO NOTHING
            `, item.ID, item.Name, item.Category, item.Price, item.CreateDate)
			if err != nil {
				log.Printf("Error inserting item record: %v", err)
				continue
			}
			validItemsCount++
		}

		// Retrieve total categories and total price from database
		var dbCategories int
		var dbTotalPrice float64

		row := tx.QueryRow(`
            SELECT COUNT(DISTINCT category), COALESCE(SUM(price), 0)
            FROM prices
        `)
		if err := row.Scan(&dbCategories, &dbTotalPrice); err != nil {
			log.Printf("Failed to calculate totals: %v", err)
			http.Error(w, "Failed to calculate totals", http.StatusInternalServerError)
			return
		}

		// Commit the transaction
		if err := tx.Commit(); err != nil {
			log.Printf("Failed to commit transaction: %v", err)
			http.Error(w, "Failed to commit transaction", http.StatusInternalServerError)
			return
		}

		// Prepare JSON response with totals
		totalsResponse := models.Totals{
			TotalItems:      validItemsCount,
			TotalCategories: dbCategories,
			TotalPrice:      dbTotalPrice,
		}

		// Send JSON response
		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(totalsResponse); err != nil {
			log.Printf("Error encoding JSON: %v", err)
			http.Error(w, "Error encoding JSON", http.StatusInternalServerError)

		}
	}
}
