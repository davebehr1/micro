package main

import (
	"context"
	"log"
	"os"

	microclient "github.com/asim/go-micro/v3/client"
	"github.com/asim/go-micro/v3/cmd"
	pb "github.com/davebehr1/micro/user-service/proto/user"
	"github.com/urfave/cli/v2"
)

func car(*cli.Context) error {
	return nil
}
func main() {

	cmd.Init()

	client := pb.NewUserService("go.micro.srv.user", microclient.DefaultClient)
	app := &cli.App{
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:  "name",
				Usage: "You full name",
			},
			&cli.StringFlag{
				Name:  "email",
				Usage: "Your email",
			},
			&cli.StringFlag{
				Name:  "password",
				Usage: "Your password",
			},
			&cli.StringFlag{
				Name:  "company",
				Usage: "Your company",
			},
		},
		Action: func(c *cli.Context) error {
			name := c.String("name")
			email := c.String("email")
			password := c.String("password")
			company := c.String("company")

			r, err := client.Create(context.TODO(), &pb.User{
				Name:     name,
				Email:    email,
				Password: password,
				Company:  company,
			})
			if err != nil {
				log.Fatalf("Could not create: %v", err)
			}
			log.Printf("Created: %t", r.User.Id)

			getAll, err := client.GetAll(context.Background(), &pb.Request{})
			if err != nil {
				log.Fatalf("Could not list users: %v", err)
			}
			for _, v := range getAll.Users {
				log.Println(v)
			}

			return nil
		},
	}

	// Run the server
	if err := app.Run(os.Args); err != nil {
		log.Println(err)
	}
}
