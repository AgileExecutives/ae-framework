package database

import (
	"strings"
	"testing"
)

func TestConnectionLogStringMasksPassword(t *testing.T) {
	config := Config{
		Host:     "localhost",
		Port:     "5432",
		User:     "postgres",
		Password: "super-secret",
		DBName:   "app_db",
		SSLMode:  "disable",
	}

	logLine := connectionLogString(config, config.DBName)
	if strings.Contains(logLine, config.Password) {
		t.Fatalf("connection log leaked password: %s", logLine)
	}
	if !strings.Contains(logLine, "password=********") {
		t.Fatalf("connection log did not contain masked password: %s", logLine)
	}
}

func TestConnectionLogStringShowsEmptyPasswordExplicitly(t *testing.T) {
	config := Config{Host: "localhost", Port: "5432", User: "postgres", DBName: "app_db", SSLMode: "disable"}

	logLine := connectionLogString(config, config.DBName)
	if !strings.Contains(logLine, "password=<empty>") {
		t.Fatalf("connection log did not show empty password marker: %s", logLine)
	}
}
