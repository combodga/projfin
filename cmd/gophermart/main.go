package main

import (
	"flag"
	"log"
	"os"

	"github.com/combodga/projfin/internal/app"
)

func main() {
	var run string
	var db string
	var accr string

	flag.StringVar(&run, "a", os.Getenv("RUN_ADDRESS"), "run address")
	flag.StringVar(&db, "d", os.Getenv("DATABASE_URI"), "db connection")
	flag.StringVar(&accr, "r", os.Getenv("ACCRUAL_SYSTEM_ADDRESS"), "accrual system address")
	flag.Parse()

	if run == "" {
		run = ":8080"
	}

	err := app.Go(run, db, accr)
	if err != nil {
		log.Fatal(err)
	}
}
