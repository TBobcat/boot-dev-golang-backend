package main

import (
	"encoding/json"
	"fmt"

	"io/ioutil"
	"os"
)

func saveJSONToFile(data interface{}) error {
	type Person struct {
		FirstName string `json:"first_name"`
		LastName  string `json:"last_name"`
		Age       int    `json:"age"`
	}

	// Sample data
	person := Person{
		FirstName: "John",
		LastName:  "Doe",
		Age:       30,
	}

	// create a file
	filePath := "file_path.json"
	outputFile, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer outputFile.Close()

	// Marshal the data into JSON format
	jsonData, err := json.MarshalIndent(person, "", "    ")
	if err != nil {
		fmt.Println("Error marshaling JSON:", err)
		return err
	}

	// write marshalled json data into created file
	_, err = outputFile.Write(jsonData)
	if err != nil {
		fmt.Println("Error writing to file:", err)
		return err
	}

	fmt.Println("JSON data has been written to", filePath)

	// Read the JSON data from the file
	dataRead, err := ioutil.ReadFile(filePath)
	if err != nil {
		return err
	}

	// Unmarshal the JSON data into the struct
	var read_person Person
	if err := json.Unmarshal(dataRead, &read_person); err != nil {
		return err
	}
	fmt.Println("Reading json file")
	fmt.Printf("%+v\n", read_person)

	return nil
}

func main() {

	saveJSONToFile("testingggggg")
}
