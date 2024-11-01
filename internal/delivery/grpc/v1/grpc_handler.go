package v1

import (
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	pb "telephone/internal/proto"
)

func (srv *Server) CreateUser(request *pb.CreateUserRequest,
) (*pb.CreateUserResponse, error) {

	if request.User.Name == "" {
		return nil, status.Error(codes.InvalidArgument,
			"user name is required")
	}

	return &pb.CreateUserResponse{}, nil
}

func (srv *Server) CreateContact(request *pb.CreateContactRequest) (*pb.CreateContactResponse, error) {
	if request.Contact.Phone == "" {
		return nil, status.Error(codes.InvalidArgument,
			"phone is required")
	}
	//res, err := srv.services.CreateContact()
	return &pb.CreateContactResponse{
		Message: "contact is created successfully"}, nil
}

func (s *Server) GetAllUsers(request *pb.GetUsersRequest) *pb.GetUsersResponse {
	res, err := s.services.GetAllUsers()
	if err != nil {
		return &pb.GetUsersResponse{}
	}
	return &pb.GetUsersResponse{
		Users: res,
	}
}
