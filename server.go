package main

import (
	"net/http"
	"fmt"
	"strings"
)

func main() {

	sayHelloAgain := func(w http.ResponseWriter, req *http.Request) {
		message := req.URL.Path
		 message = strings.TrimPrefix(message, "/")
		message = "Hello " + message
		fmt.Printf("succesfully redirected")
	}

	sayHello := func(w http.ResponseWriter, req *http.Request) {		
		http.Redirect(w, req, "/helloAgain", 302)
	}
	
	// we can use of the following handles to sign and validate & other to save
	http.HandleFunc("/hello", sayHello)
	http.HandleFunc("/helloAgain", sayHelloAgain)
	
	http.Handle("/", http.FileServer(http.Dir("C:/Documents/angular2-registration-login-example-master/angular2-registration-login-example-master")))
	http.ListenAndServe(":3000", nil)

}
