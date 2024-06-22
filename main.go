package main

import "log"

func init() {
	/*different approach, in the init section here we will test all connections before launching the app.
	Will use the same db-services to read and connect as a startup health check. init will not keep any data it access
	*/
}

func main() {
	log.Println("Hello world")
	/*
		will follow my standard approach of spawning go routines with the actual db-services
		from the main.go file. This allows me to handle unforseen panics or self inflicted panics while
		keeping the app running easier
	*/
}
