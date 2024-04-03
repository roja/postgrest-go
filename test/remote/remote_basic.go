// This is basic example for postgrest-go library usage.
// For now this example is represent wanted syntax and bindings for library.
// After core development this test files will be used for CI tests.

package main

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/supabase-community/gotrue-go"
	"github.com/supabase-community/gotrue-go/types"
	postgrest "github.com/supabase-community/postgrest-go"
)

var (
	RestUrl = `http://localhost:3000`
	headers = map[string]string{}
	schema  = "public"
)

func loginAndGetToken(projectReference string, anonKey string) (string, gotrue.Client) {

	// Initialise client
	client := gotrue.New(
		projectReference,
		anonKey,
	)

	// Log in a user (get access and refresh tokens)
	resp, err := client.Token(types.TokenRequest{
		GrantType: "password",
		Email:     os.Getenv("TESTUSER"),
		Password:  os.Getenv("TESTUSERPASSWORD"),
	})
	fmt.Println(err)
	return resp.AccessToken, client.WithToken(resp.AccessToken)
}

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	projectURL := os.Getenv("SUPABASE_PROJECT")
	anonKey := os.Getenv("SUPABASE_ANON_KEY")

	accessToken, authClient := loginAndGetToken(projectURL, anonKey)
	headers := make(map[string]string)
	headers["Authorization"] = "Bearer " + accessToken
	headers["apikey"] = anonKey

	fmt.Println(authClient)
	client := postgrest.NewClient(os.Getenv("SUPABASE_URL"), "public", headers)

	//
	rooms, _, err := client.From("rooms").Select("*", "", false).ExecuteString()
	if err != nil {
		panic(err)
	}
	fmt.Println(rooms)

}
