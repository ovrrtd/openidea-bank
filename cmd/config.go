package main

import "os"

var (

	// APP ENV VARS
	APP_PORT = "8080"
	// DB ENV VARS
	DB_HOST     = os.Getenv("DB_HOST")
	DB_USERNAME = os.Getenv("DB_USERNAME")
	DB_PASSWORD = os.Getenv("DB_PASSWORD")
	DB_NAME     = os.Getenv("DB_NAME")
	DB_PORT     = os.Getenv("DB_PORT")
	DB_PARAMS   = os.Getenv("DB_PARAMS")
)
