package v1

import (
	"context"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	pb "telephone/internal/proto"
	"telephone/internal/schema"
)

func (srv *Server) CreateUser(ctx context.Context, request *pb.CreateUserRequest,
) (*pb.CreateUserResponse, error) {

	modelUser := schema.NewFromProtoToModelUserRequest(request)

	res, err := srv.services.TelephoneService.CreateUser(ctx, modelUser)
	if err != nil {
		return nil, err
	}
	return &pb.CreateUserResponse{
		UserId: uint64(res.ID),
	}, nil
}

func (srv *Server) CreateContact(ctx context.Context, req *pb.CreateContactRequest,
) (*pb.CreateContactResponse, error) {
	if req.Contact.Phone == "" {
		return nil, status.Error(codes.InvalidArgument,
			"phone is required")
	}

	modelContact := schema.NewFromProtoToModelCreateContactRequest(req)
	res, err := srv.services.ContactService.CreateContact(ctx, modelContact)
	if err != nil {
		return nil, err
	}
	return &pb.CreateContactResponse{
		ContactId: uint64(res.ContactID),
	}, nil
}

func (srv *Server) GetAllUsers(ctx context.Context, empty *pb.GetUsersRequest) (resp *pb.GetUsersResponse, err error) {
	users, err := srv.services.TelephoneService.GetAllUsers(ctx)
	if err != nil {
		return nil, err
	}
	resp = &pb.GetUsersResponse{
		Users: []*pb.User{},
	}

	for _, user := range users {
		res := &pb.User{
			UserId: uint64(user.ID),
			Name:   user.Name,
			Email:  user.Email,
		}

		resp.Users = append(resp.Users, res)
	}
	return resp, nil
}

func (srv *Server) GetUser(ctx context.Context, req *pb.GetUserRequest) (*pb.GetUserResponse, error) {
	resp, err := srv.services.TelephoneService.GetUserByID(ctx, int(req.Id))
	if err != nil {
		return nil, err
	}

	return &pb.GetUserResponse{
		User: &pb.User{
			UserId: uint64(resp.ID),
			Name:   resp.Name,
			Email:  resp.Email,
		},
	}, nil
}

func (srv *Server) GetContact(ctx context.Context, req *pb.GetContactRequest,
) (*pb.GetContactResponse, error) {
	resp, err := srv.services.ContactService.GetContactByID(ctx, uint64(req.ContactId))
	if err != nil {
		return nil, err
	}

	return &pb.GetContactResponse{Contact: &pb.Contact{
		ContactId: resp.ContactID,
		Phone:     resp.PhoneNumber,
	}}, nil
}

func (srv *Server) GetContacts(ctx context.Context, empty *pb.GetContactsRequest) (resp *pb.GetContactsResponse, err error) {
	contacts, err := srv.services.ContactService.GetContacts()
	if err != nil {
		return nil, err
	}
	resp = &pb.GetContactsResponse{
		Contacts: []*pb.Contact{},
	}

	for _, c := range contacts {
		res := &pb.Contact{
			ContactId: uint64(c.ContactID),
			Phone:     c.PhoneNumber,
		}

		resp.Contacts = append(resp.Contacts, res)
	}
	return resp, nil
}

func (srv *Server) AddUserContact(ctx context.Context, request *pb.AddUserContactRequest) (*pb.AddUserContactResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (srv *Server) GetUserContact(ctx context.Context, request *pb.GetUserRequest,
) (*pb.GetUserContactRelationResponse, error) {
	relations, err := srv.services.UserContactService.ListFav(ctx, int(request.Id))
	if err != nil {
		return nil, err
	}

	resp := &pb.GetUserContactRelationResponse{
		ContactId: uint64(relations.ContactID),
	}

	return resp, nil
}
