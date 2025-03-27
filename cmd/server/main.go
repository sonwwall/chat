package main

import "chat/initialize"

func main() {
	initialize.SetUpViper()
	initialize.SetupLogger()
	initialize.SetupDatabase()
}
