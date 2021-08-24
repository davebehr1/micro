package main

import (
	"fmt"
	"log"
	"os"

	"github.com/asim/go-micro/v3"
	pb "github.com/davebehr1/micro/vessel-service/proto/vessel"
)

const (
	defaultHost = "localhost:27017"
)

func createDummyData(repo Repository) {
	defer repo.Close()
	vessels := []*pb.Vessel{
		{Id: "vessel001", Name: "Kane's Salty Secret", MaxWeight: 200000, Capacity: 500},
	}
	for _, v := range vessels {
		repo.Create(v)
	}
}

// type Repository interface {
// 	FindAvailable(*pb.Specification) (*pb.Vessel, error)
// }

// type VesselRepository struct {
// 	vessels []*pb.Vessel
// }

// FindAvailable - checks a specification against a map of vessels,
// if capacity and max weight are below a vessels capacity and max weight,
// then return that vessel.
// func (repo *VesselRepository) FindAvailable(spec *pb.Specification) (*pb.Vessel, error) {
// 	for _, vessel := range repo.vessels {
// 		if spec.Capacity <= vessel.Capacity && spec.MaxWeight <= vessel.MaxWeight {
// 			return vessel, nil
// 		}
// 	}
// 	return nil, errors.New("No vessel found by that spec")
// }

// Our grpc service handler
// type service struct {
// 	repo Repository
// }

// func (s *service) FindAvailable(ctx context.Context, req *pb.Specification, res *pb.Response) error {

// 	// Find the next available vessel
// 	vessel, err := s.repo.FindAvailable(req)
// 	if err != nil {
// 		return err
// 	}

// 	// Set the vessel as part of the response message type
// 	res.Vessel = vessel
// 	return nil
// }

func main() {

	host := os.Getenv("DB_HOST")

	if host == "" {
		host = defaultHost
	}

	session, err := CreateSession(host)
	defer session.Close()

	if err != nil {
		log.Fatalf("Error connecting to datastore: %v", err)
	}

	repo := &VesselRepository{session.Copy()}

	createDummyData(repo)

	srv := micro.NewService(
		micro.Name("go.micro.srv.vessel"),
		micro.Version("latest"),
	)

	srv.Init()

	// Register our implementation with
	pb.RegisterVesselServiceHandler(srv.Server(), &service{session})

	if err := srv.Run(); err != nil {
		fmt.Println(err)
	}
}
