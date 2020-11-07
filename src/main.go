package main

import (
	"fmt"

	"github.com/spf13/viper"
	config "github.com/urbinaalex17/go-vipper-app/src/config"
)

func main() {
	// Set the file name of the configurations file
	viper.SetConfigName("auth")

	// Set the path to look for the configurations file
	viper.AddConfigPath(".")
	viper.AddConfigPath("/vault/secrets")

	// Enable VIPER to read Environment Variables
	viper.AutomaticEnv()

	viper.SetConfigType("yml")
	var auth config.AppAuth

	if err := viper.ReadInConfig(); err != nil {
		fmt.Printf("Error reading config file, %s", err)
	}

	err := viper.Unmarshal(&auth)
	if err != nil {
		fmt.Printf("Unable to decode into struct, %v", err)
	}

	// Reading variables using the model
	fmt.Println("Reading variables using the model..")
	fmt.Println("API_KEY is\t", auth.API_KEY)
	fmt.Println("API_SECRET is\t", auth.API_SECRET)

	// Reading variables without using the model
	fmt.Println("\nReading variables without using the model..")
	apiKey := viper.GetString("API_KEY")
	secretKey := viper.GetString("API_SECRET")
	fmt.Println("API_KEY is\t", apiKey)
	fmt.Println("API_SECRET is\t", secretKey)

}
