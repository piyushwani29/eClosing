package main

import (
   "encoding/json"
   "fmt"
   "io/ioutil"
)

type Document struct {
   Id       int       `json:"id"`
   hash  string    `json:"Hash"`
   path string `json:"path"`
}


func main() {
   document := Document{
      Id:      1,
      hash: "xyzesf" ,
	  path: "C:\\Doclist" ,
	  }

   output, err := json.MarshalIndent(&document, "", "\t\t")
   if err != nil {
      fmt.Println("Error marshalling to JSON:", err)
      return
   }
   err = ioutil.WriteFile("DocList.json", output, 0644)
   if err != nil {
      fmt.Println("Error writing JSON to file:", err)
      return
   }
}
func readjson(){
	jsonFile, err := os.Open("DocList.json")
   if err != nil {
      fmt.Println("Error opening JSON file:", err)
      return
   }
   defer jsonFile.Close()
   jsonData, err := ioutil.ReadAll(jsonFile)
   if err != nil {
      fmt.Println("Error reading JSON data:", err)
      return
   	}

   fmt.Println(string(jsonData))
   var document Document
   json.Unmarshal(jsonData, &document)
   fmt.Println(document)
   
}