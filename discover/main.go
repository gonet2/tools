package main

import (
	"log"
	"os"
	"text/template"
)

func main() {
	if len(os.Args) <= 1 {
		return
	}

	tmpl, err := template.New("discover").Parse(t)
	if err != nil {
		log.Fatal(err)
	}

	f, err := os.Create("services.go")
	if err != nil {
		log.Fatal(err)
	}

	err = tmpl.Execute(f, os.Args[1:])
	if err != nil {
		log.Fatal(err)
	}
	f.Close()
}
