package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"

	"github.com/asim/go-micro/v3/broker"
	pb "github.com/davebehr1/micro/user-service/proto/user"
	"golang.org/x/crypto/bcrypt"
	"golang.org/x/net/context"
)

const topic = "user.created"

type service struct {
	repo         Repository
	tokenService Authable
	pubsub       broker.Broker
}

func (srv *service) Get(ctx context.Context, req *pb.User, res *pb.Response) error {
	user, err := srv.repo.Get(req.Id)
	if err != nil {
		return err
	}
	res.User = user
	return nil
}

func (srv *service) GetAll(ctx context.Context, req *pb.Request, res *pb.Response) error {
	users, err := srv.repo.GetAll()
	if err != nil {
		return err
	}
	res.Users = users
	return nil
}

func (srv *service) Auth(ctx context.Context, req *pb.User, res *pb.Token) error {
	log.Println("Logging in with:", req.Email, req.Password)
	user, err := srv.repo.GetByEmail(req.Email)
	if err != nil {
		return err
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		return err
	}
	token, err := srv.tokenService.Encode(user)
	if err != nil {
		return err
	}
	res.Token = token
	return nil
}

func (srv *service) Create(ctx context.Context, req *pb.User, res *pb.Response) error {
	log.Println("Creating user: ", req)
	hashedPass, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf(fmt.Sprintf("error hashing password: %v", err))
	}
	req.Password = string(hashedPass)

	if err := srv.repo.Create(req); err != nil {
		return fmt.Errorf(fmt.Sprintf("error creating user: %v", err))
	}

	res.User = req
	if err := srv.publishEvent(req); err != nil {
		return fmt.Errorf(fmt.Sprintf("error publishing event: %v", err))
	}
	return nil
}

func (srv *service) publishEvent(user *pb.User) error {
	// Marshal to JSON string
	body, err := json.Marshal(user)
	if err != nil {
		return err
	}

	msg := &broker.Message{
		Header: map[string]string{
			"id": user.Id,
		},
		Body: body,
	}
	if err := srv.pubsub.Publish(topic, msg); err != nil {
		log.Printf("[pub] failed: %v", err)
	}

	return nil
}

func (srv *service) ValidateToken(ctx context.Context, req *pb.Token, res *pb.Token) error {

	claims, err := srv.tokenService.Decode(req.Token)

	if err != nil {
		return err
	}

	if claims.User.Id == "" {
		return errors.New("invalid user")
	}

	res.Valid = true

	return nil
}
