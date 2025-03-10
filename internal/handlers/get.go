package handlers

import (
	"archive/zip"
	"bytes"
	"database/sql"
	"encoding/csv"
	"log"
	"net/http"
	"strconv"

	"itmo-devops-sem1-project-template/internal/models"
)

func GetPrices(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Check if the request method is GET, otherwise return an error
		if r.Method != http.MethodGet {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		// Query the database to retrieve data
		rows, err := db.Query(`
            SELECT id, name, category, price, create_date 
            FROM prices
        `)
		if err != nil {
			log.Printf("DB query error: %v", err)
			http.Error(w, "Failed to retrieve data", http.StatusInternalServerError)
			return
		}
		defer rows.Close()

		var items []models.Item
		// Iterate over rows and scan into struct
		for rows.Next() {
			var item models.Item
			if err := rows.Scan(&item.ID, &item.Name, &item.Category, &item.Price, &item.CreateDate); err != nil {
				log.Printf("Row scan error: %v", err)
				continue
			}

			items = append(items, item)
		}
		if err := rows.Err(); err != nil {
			log.Printf("Rows iteration error: %v", err)
			http.Error(w, "Failed to read rows", http.StatusInternalServerError)
			return
		}

		// Create CSV buffer and writer
		csvBuffer := &bytes.Buffer{}
		writer := csv.NewWriter(csvBuffer)

		// Write CSV header row
		writer.Write([]string{"id", "name", "category", "price", "create_date"})

		// Write data to CSV
		for _, item := range items {
			writer.Write([]string{
				strconv.FormatInt(item.ID, 10),
				item.Name,
				item.Category,
				strconv.FormatFloat(item.Price, 'f', 2, 64),
				item.CreateDate.Format("2006-01-02"),
			})
		}
		writer.Flush()

		// Check for any errors in writing CSV
		if err := writer.Error(); err != nil {
			log.Printf("CSV error: %v", err)
			http.Error(w, "Failed to write CSV", http.StatusInternalServerError)
			return
		}

		// Create ZIP buffer and writer
		zipBuffer := &bytes.Buffer{}
		zipWriter := zip.NewWriter(zipBuffer)

		// Create a CSV file inside the ZIP archive
		csvFile, err := zipWriter.Create("data.csv")
		if err != nil {
			log.Printf("Error creating file in ZIP: %v", err)
			http.Error(w, "Failed to create ZIP", http.StatusInternalServerError)
			return
		}

		// Write CSV data into the ZIP file
		if _, err := csvFile.Write(csvBuffer.Bytes()); err != nil {
			log.Printf("Error writing CSV to ZIP: %v", err)
			http.Error(w, "Failed to write ZIP", http.StatusInternalServerError)
			return
		}

		// Close the ZIP writer
		if err := zipWriter.Close(); err != nil {
			log.Printf("Error processing ZIP writer: %v", err)
			http.Error(w, "Failed to process ZIP", http.StatusInternalServerError)
			return
		}

		// Set response headers for ZIP file download
		w.Header().Set("Content-Type", "application/zip")
		w.Header().Set("Content-Disposition", "attachment; filename=data.zip")
		w.WriteHeader(http.StatusOK)

		// Send the ZIP file to the client
		if _, err := w.Write(zipBuffer.Bytes()); err != nil {
			log.Printf("Error sending ZIP file: %v", err)
		}
	}
}
