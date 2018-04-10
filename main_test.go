package main_test

import (
	"log"
	"os"
	"testing"

	"."
)

var a main.App

func TestMain(m *testing.M) {
	a = main.App{}
	a.Initialize(
		os.Getenv("APP_DB_USERNAME"),
		os.Getenv("APP_DB_PASSWORD"),
		os.Getenv("APP_DB_NAME"))

	ensureTableExists()

	code := m.Run()

	clearTable()

	os.Exit(code)
}

func ensureTableExists() {
	if _, err := a.DB.Exec(tableCreationQuery); err != nil {
		log.Fatal(err)
	}
}

func clearTable() {
	a.DB.Exec("DELETE FROM leads_test.leads")
	a.DB.Exec("ALTER SEQUENCE leads_test.leads_id_seq RESTART WITH 1")
}

const tableCreationQuery = `CREATE TABLE IF NOT EXISTS leads_test.leads
(
    id bigint NOT NULL,
    first_name text NOT NULL,
    last_name text NOT NULL,
    CONSTRAINT leads_pkey PRIMARY KEY (id)
)`
