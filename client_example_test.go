package ileap_test

import (
	"context"
	"fmt"
	"os"

	"github.com/way-platform/ileap-go"
)

func ExampleClient() {
	client := ileap.NewClient(
		ileap.WithBaseURL(ileap.DemoBaseURL),
		ileap.WithOAuth2(os.Getenv("CLIENT_ID"), os.Getenv("CLIENT_SECRET")),
	)
	footprint, err := client.GetFootprint(context.Background(), &ileap.GetFootprintRequest{
		ID: "123",
	})
	if err != nil {
		panic(err)
	}
	fmt.Println(footprint)
}
