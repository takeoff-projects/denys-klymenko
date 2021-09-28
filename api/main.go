package main

import (
	_ "github.com/swaggo/files"       // swagger embed files
	_ "github.com/swaggo/gin-swagger" // gin-swagger middleware
	"takeoff-projects/denys-klymenko/api/pets"
)

func main() {
	pets.Run()
}
