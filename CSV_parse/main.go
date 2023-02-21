package main

import (
	"encoding/csv"
	"fmt"
	"os"
)

type Record struct {
	Name    string
	Age     int
	Country string
}

func main() {
	// Open the CSV file
	file, err := os.Open("example.csv")
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	defer file.Close()

	// Create a new CSV reader
	reader := csv.NewReader(file)

	// Read all the records
	records, err := reader.ReadAll()
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	// Parse the records into a slice of structs
	var data []Record
	for _, record := range records {
		age := 0
		fmt.Sscanf(record[1], "%d", &age)
		data = append(data, Record{
			Name:    record[0],
			Age:     age,
			Country: record[2],
		})
	}

	// Output the data as a table
	fmt.Printf("%-20s %-10s %-10s\n", "Name", "Age", "Country")
	fmt.Println("-----------------------------------------")
	for _, record := range data {
		fmt.Printf("%-20s %-10d %-10s\n", record.Name, record.Age, record.Country)
	}
}
