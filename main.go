package main

import (
	"fmt"
	"log"
	"os"

	"github.com/HarrisonLeach1/xero-tui/internal/api/auth"
	"github.com/HarrisonLeach1/xero-tui/internal/ui"
	"github.com/VladimirMarkelov/clui"
	"github.com/joho/godotenv"
)

func main() {
	// loads values from .env into the system
	if err := godotenv.Load(); err != nil {
		log.Print("No .env file found")
	}

	clientId, exists := os.LookupEnv("XERO_CLIENT_ID")
	if !exists {
		fmt.Println("XERO_CLIENT_ID environment variable is not present")
		return
	}

	auth.AuthorizeUser(clientId, "http://localhost:5003")

	mainLoop()
}

func createView() {

	ui.RenderHomePage()
}

func mainLoop() {
	// Every application must create a single Composer and
	// call its intialize method
	clui.InitLibrary()
	defer clui.DeinitLibrary()

	clui.SetThemePath("themes")

	createView()

	// start event processing loop - the main core of the library
	clui.MainLoop()
}
