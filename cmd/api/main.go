package main

import "github.com/afandimsr/go-gin-api/internal/bootstrap"

// @title           Go Gin API
// @version         1.0
// @description     Clean Architecture Go Gin API
// @termsOfService  http://swagger.io/terms/

// @contact.name   Afandi
// @contact.email  mohamadafandi71@gmail.com

// @host      localhost:8080
// @BasePath  /
func main() {
	bootstrap.Run()
}
