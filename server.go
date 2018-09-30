package main

import (
	"net/http"
	"fmt"
	"strings"
	"os"
	"crypto/md5"
	"bytes"
	"io/ioutil"
	"encoding/json"
	"log"
	"reflect"
	"unsafe"
	//"os/exec"
)

type Document struct {
   Id       int       
   Hash  	string    
   Status string
}
type DocList struct{
	List []Document
}

func BytesToString(b [16]byte) string{
	bh:= (*reflect.SliceHeader)(unsafe.Pointer(&b))
	sh:= reflect.StringHeader{bh.Data, bh.Len}
	return *(*string) (unsafe.Pointer(&sh))
}
func getMd5(file string) [16]byte{
	 f,err := ioutil.ReadFile(file)
	if(err!=nil){
		log.Fatal(err)
	}
	 return md5.Sum(f)
}

func readjson() DocList{
	jsonFile, err := os.Open("C:\\Users\\nimisha\\Documents\\Angular projects\\eClosing\\src\\assets\\DocList.json")
   if err != nil {
      fmt.Println("Error opening JSON file:", err)
   }
   defer jsonFile.Close()
   jsonData, err := ioutil.ReadAll(jsonFile)
   if err != nil {
      fmt.Println("Error reading JSON data:", err)
   	}
	bytes:=[]byte(jsonData)
   var document DocList
   json.Unmarshal(bytes, &document)
   return document
}

func main() {

	saveDocument := func(w http.ResponseWriter, req *http.Request) {
		body, err := ioutil.ReadAll(req.Body)
		req.Body=ioutil.NopCloser(bytes.NewBuffer(body))
		if(err == nil){
			var fileName = strings.Replace(string(body),"C:\\fakepath\\","G:\\GO\\server\\Documents\\",1) // Replace the value to location from where you are keeping the document
			fmt.Printf(fileName)
			//cmd:=exec.Command("sleep","1") // Path to batch file  that valdates the PDF and executes Smart contract needs to be added here, also per Kiran we need the Hash value should be the hash generated post executing the smart contract
			doclist:=readjson()
			length:=len(doclist.List)+1
			 value:=getMd5(fileName)
			document := Document{
				Id:      length,
				Hash: string(value[:]) ,
				Status: "Signature validated, funds transferred", 
				}
			target:= make([]Document, length-1)
			copy(target,doclist.List[:(length -1)] )
			target=append(target,document)
			output, err := json.Marshal(target)
			if err != nil {
				fmt.Println("Error marshalling to JSON:", err)
				return
			}			
			f, err := os.OpenFile("G:\\GO\\server\\docs\\DocList.json", os.O_WRONLY, 0777) 
			n, err := f.WriteString("{\"List\" : "+string(output)+"}")
			fmt.Printf("%v",n)
			if err != nil {
				fmt.Println("Error writing JSON to file:", err)
				return
			}			
		}
			
	}
	
	getDocList:=func(w http.ResponseWriter, req *http.Request) {
		w.Header().Set("Content-Type","application/json")
		doclist:=readjson()
		output, err := json.Marshal(doclist)
		fmt.Println(string(output))
		w.Write(output)
		if err!=nil {
			fmt.Println(err) }
	}
	bringNetworkup:= func(w http.ResponseWriter, req *http.Request) {
		body, err := ioutil.ReadAll(req.Body)
		req.Body=ioutil.NopCloser(bytes.NewBuffer(body))
		if(err == nil){
			fmt.Println(body);
			//cmd:=exec.Command("sleep","1") // Path to batch file needs to be added here
		}
	}

	http.HandleFunc("/save", saveDocument)
	http.HandleFunc("/bringNetworkup", bringNetworkup)
	http.HandleFunc("/getDocList", getDocList)
	http.Handle("/", http.FileServer(http.Dir("G:/Angular projects/eClosing/dist")))
	http.ListenAndServe(":3000", nil)

}

