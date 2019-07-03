package utils

import (
	"fmt"
	"os"
	"github.com/joho/godotenv"
)

func LoadEnv(){
	godotenv.Load("github.com/andresvillavicenciowizeline/proxy-app/.env")
	fmt.Println(os.Getenv("Port"))
}
