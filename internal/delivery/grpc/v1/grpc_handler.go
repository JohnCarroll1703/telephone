package v1

import (
	"context"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"strconv"
	pb "telephone/internal/proto"
	"telephone/internal/schema"
)

func (srv *Server) CreateUser(ctx context.Context, request *pb.CreateUserRequest,
) (_ *pb.CreateUserResponse, err error) {

	modelUser := schema.NewFromProtoToModelUserRequest(request)

	err = modelUser.Validate()
	if err != nil {
		return nil, status.Error(
			codes.InvalidArgument,
			err.Error(),
		)
	}
	res, err := srv.services.TelephoneService.CreateUser(ctx, modelUser)
	if err != nil {
		return nil, err
	}
	return &pb.CreateUserResponse{
		UserId: uint64(res.ID),
	}, nil
}

func (srv *Server) CreateContact(ctx context.Context, req *pb.CreateContactRequest,
) (_ *pb.CreateContactResponse, err error) {

	modelContact := schema.NewFromProtoToModelCreateContactRequest(req)

	err = modelContact.Validate()
	if err != nil {
		return nil, status.Error(
			codes.InvalidArgument,
			err.Error(),
		)
	}

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
	resp, err := srv.services.TelephoneService.GetUserByID(ctx, uint(req.Id))
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

func (srv *Server) AddUserContact(ctx context.Context, request *pb.AddContactRequest,
) (_ *pb.AddContactResponse, err error) {
	modelContact := schema.NewFromProtoToModelAddContactRequest(request)

	err = modelContact.Validate()
	if err != nil {
		return nil, status.Error(
			codes.InvalidArgument,
			err.Error(),
		)
	}

	user, err := srv.services.TelephoneService.GetUserByID(ctx, uint(request.Id))
	if err != nil {
		return nil, status.Error(codes.NotFound, "пользователя "+
			strconv.Itoa(int(user.ID))+" не существует")
	}

	resp, err := srv.services.UserContactService.AddContacts(ctx, request.Id, modelContact)
	if err != nil {
		return nil, err
	}

	return &pb.AddContactResponse{
		Phone:  resp.PhoneNumber,
		Status: "the contact is added successfully!",
	}, nil
}

func (srv *Server) GetUserContact(ctx context.Context, request *pb.GetUserRequest,
) (resp *pb.GetUserContactRelationResponse, err error) {
	relations, err := srv.services.UserContactService.ListFav(ctx, int(request.Id))
	if err != nil {
		return nil, err
	}

	contactsList := make([]*pb.Contact, len(relations))

	for i, contact := range relations {
		contactsList[i] = &pb.Contact{
			ContactId: uint64(contact.ContactID),
			Phone:     contact.PhoneNumber,
		}
	}

	return &pb.GetUserContactRelationResponse{
		UserId:   uint64(request.Id),
		Contacts: contactsList,
	}, nil
}

func (srv *Server) GetAllRelations(ctx context.Context, empty *pb.GetAllRelationsRequest) (
	resp *pb.GetAllRelationsResponse, err error) {
	relations, err := srv.services.UserContactService.GetAllRelations()
	if err != nil {
		return nil, err
	}

	resp = &pb.GetAllRelationsResponse{
		Relations: []*pb.UserContacts{},
	}

	for _, c := range relations {
		res := &pb.UserContacts{
			ContactId: uint64(c.ContactID),
			IsFav:     true,
			UserId:    uint64(c.UserID),
		}

		resp.Relations = append(resp.Relations, res)
	}
	return resp, nil
}

//func (srv *Server) GetUsersWithPaginationAndFiltering(ctx context.Context,
//	input *pb.GetUsersWithPaginationAndFilteringRequest) (
//	resp *pb.GetUsersWithPaginationAndFilteringResponse, err error) {
//	if input.Page < 1 {
//		input.Page = 1
//	}
//
//	if input.Limit < 1 {
//		input.Limit = 10
//	}
//
//	filter := map[string]interface{}{
//		"name": input.Name,
//	}
//
//	users, pagination, err := srv.services.TelephoneService.GetAllUsersWithPaginationAndFiltering()
//
//	resp = &pb.GetUsersWithPaginationAndFilteringResponse{
//		Users: []*pb.User{},
//	}
//	for _, val := range users {
//		user := &pb.User{}
//	}
//}
