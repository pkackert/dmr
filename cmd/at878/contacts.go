// Creates DMR contacts CSV file for AnyTone 878uv using data from
// https://www.radioid.net/static/users.json
//
// * Download json file manually
// CSV format compatible with AnyTone CSP 1.21

package main

import (
	"bufio"
	"encoding/json"
	"flag"
	"fmt"
	"os"
)

type User struct {
	Fname    string      `json:"fname"`
	Name     string      `json:"name"`
	Country  string      `json:"country"`
	Callsign string      `json:"callsign"`
	City     string      `json:"city"`
	Surname  string      `json:"surname"`
	RadioID  int         `json:"radio_id"`
	ID       int         `json:"id"`
	Remarks  interface{} `json:"remarks"`
	State    string      `json:"state"`
}

type DMRUsers struct {
	Users []User `json:"users"`
}

func main() {
	var src, dst string

	flag.StringVar(&src, "i", "users.json", "Name of json source file.  Default is \"users.json\"")
	flag.StringVar(&dst, "o", "users.csv", "Name of CSV output file.  Default is \"users.csv\"")

	flag.Parse()

	// Read json file
	jsonBlob, err := os.ReadFile(src)
	if err != nil {
		fmt.Printf("Error: %s\n", err.Error())
		os.Exit(1)
	}

	var dmrUsers DMRUsers

	err = json.Unmarshal(jsonBlob, &dmrUsers)
	if err != nil {
		fmt.Println("error:", err)
		os.Exit(1)
	}

	f, err := os.Create(dst)
	if err != nil {
		fmt.Printf("Error: Failed to open %s for writing: %s\n", dst, err.Error())
		os.Exit(1)
	}
	defer f.Close()
	w := bufio.NewWriter(f)

	// Output CSV header
	fmt.Fprintf(w, "%s\r\n", `"No.","Radio ID","Callsign","Name","City","State","Country","Remarks","Call Type","Call Alert"`)

	// Process each record in dataset
	for n, u := range dmrUsers.Users {
		remarks := ""
		if u.Remarks != nil {
			remarks = u.Remarks.(string)
		}

		fmt.Fprintf(w, "\"%d\",\"%d\",\"%s\",\"%s %s\",\"%s\",\"%s\",\"%s\",\"%s\",\"Private Call\",\"None\"\r\n",
			n+1, u.RadioID, u.Callsign, u.Fname, u.Surname, u.City, u.State, u.Country, remarks)
	}
	w.Flush()
}
