package main

import (
	"log"
	"os"

	micro "github.com/asim/go-micro/v3"
	microclient "github.com/asim/go-micro/v3/client"
	pb "github.com/davebehr1/micro/user-service/proto/user"
	"golang.org/x/net/context"
)

func main() {

	srv := micro.NewService(

		micro.Name("go.micro.srv.user-cli"),
		micro.Version("latest"),
	)

	// Init will parse the command line flags.
	srv.Init()

	client := pb.NewUserService("go.micro.srv.user", microclient.DefaultClient)

	name := "Ewan Valentine"
	email := "ewan.valentine89@gmail.com"
	password := "test123"
	company := "BBC"

	r, err := client.Create(context.TODO(), &pb.User{
		Name:     name,
		Email:    email,
		Password: password,
		Company:  company,
	})
	if err != nil {
		log.Fatalf("Could not create: %v", err)
	}
	log.Printf("Created: %s", r.User.Id)

	getAll, err := client.GetAll(context.Background(), &pb.Request{})
	if err != nil {
		log.Fatalf("Could not list users: %v", err)
	}
	for _, v := range getAll.Users {
		log.Println(v)
	}

	authResponse, err := client.Auth(context.TODO(), &pb.User{
		Email:    email,
		Password: password,
	})

	if err != nil {
		log.Fatalf("Could not authenticate user: %s error: %v\n", email, err)
	}

	log.Printf("Your access token is: %s \n", authResponse.Token)

	// let's just exit because
	os.Exit(0)
}
