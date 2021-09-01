// consignment-service/main.go
package main

import (

	// Import the generated protobuf code
	"errors"
	"fmt"
	"log"
	"os"

	"github.com/asim/go-micro/v3"
	"github.com/asim/go-micro/v3/client"
	"github.com/asim/go-micro/v3/metadata"
	"github.com/asim/go-micro/v3/server"
	pb "github.com/davebehr1/micro/consignment-service/proto/consignment"
	userService "github.com/davebehr1/micro/user-service/proto/user"
	vesselProto "github.com/davebehr1/micro/vessel-service/proto/vessel"
	"golang.org/x/net/context"
)

const (
	port        = ":50051"
	defaultHost = "localhost:27017"
)

// type IRepository interface {
// 	Create(*pb.Consignment) (*pb.Consignment, error)
// 	GetAll() []*pb.Consignment
// }

// Repository - Dummy repository, this simulates the use of a datastore
// of some kind. We'll replace this with a real implementation later on.
// type Repository struct {
// 	consignments []*pb.Consignment
// }

// func (repo *Repository) Create(consignment *pb.Consignment) (*pb.Consignment, error) {
// 	updated := append(repo.consignments, consignment)
// 	repo.consignments = updated
// 	return consignment, nil
// }

// func (repo *Repository) GetAll() []*pb.Consignment {
// 	return repo.consignments
// }

// Service should implement all of the methods to satisfy the service
// we defined in our protobuf definition. You can check the interface
// in the generated code itself for the exact method signatures etc
// to give you a better idea.
// type service struct {
// 	repo         IRepository
// 	vesselClient vesselProto.VesselService
// }

// CreateConsignment - we created just one method on our service,
// which is a create method, which takes a context and a request as an
// argument, these are handled by the gRPC server.
// func (s *service) CreateConsignment(ctx context.Context, req *pb.Consignment, res *pb.Response) error {

// 	vesselResponse, err := s.vesselClient.FindAvailable(context.Background(), &vesselProto.Specification{
// 		MaxWeight: req.Weight,
// 		Capacity:  int32(len(req.Containers)),
// 	})

// 	log.Printf("Found vessel: %s \n", vesselResponse.Vessel.Name)
// 	if err != nil {
// 		return err
// 	}

// 	req.VesselId = vesselResponse.Vessel.Id

// 	// Save our consignment
// 	consignment, err := s.repo.Create(req)
// 	if err != nil {
// 		return err
// 	}

// 	// Return matching the `Response` message we created in our
// 	// protobuf definition.
// 	res.Created = true
// 	res.Consignment = consignment
// 	return nil
// }

// func (s *service) GetConsignments(ctx context.Context, req *pb.GetRequest, res *pb.Response) error {
// 	consignments := s.repo.GetAll()
// 	res.Consignments = consignments
// 	return nil
// }

func AuthWrapper(fn server.HandlerFunc) server.HandlerFunc {
	return func(ctx context.Context, req server.Request, resp interface{}) error {
		meta, ok := metadata.FromContext(ctx)
		if !ok {
			return errors.New("no auth meta-data found in request")
		}

		// Note this is now uppercase (not entirely sure why this is...)
		token := meta["Token"]
		log.Println("Authenticating with token: ", token)

		// Auth here
		authClient := userService.NewUserService("go.micro.srv.user", client.DefaultClient)
		_, err := authClient.ValidateToken(context.Background(), &userService.Token{
			Token: token,
		})
		if err != nil {
			return err
		}
		err = fn(ctx, req, resp)
		return err
	}
}

func main() {

	host := os.Getenv("DB_HOST")

	if host == "" {
		host = defaultHost
	}

	session, err := CreateSession(host)

	defer session.Close()

	if err != nil {

		// We're wrapping the error returned from our CreateSession
		// here to add some context to the error.
		log.Panicf("Could not connect to datastore with host %s - %v", host, err)
	}

	// repo := &Repository{}

	srv := micro.NewService(

		// This name must match the package name given in your protobuf definition
		micro.Name("go.micro.srv.consignment"),
		micro.Version("latest"),
		micro.WrapHandler(AuthWrapper),
	)

	vesselClient := vesselProto.NewVesselService("go.micro.srv.vessel", srv.Client())

	srv.Init()

	pb.RegisterShippingServiceHandler(srv.Server(), &service{session: session, vesselClient: vesselClient})

	// Run the server
	if err := srv.Run(); err != nil {
		fmt.Println(err)
	}
}
