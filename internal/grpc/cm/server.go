package cm

import (
	"context"
	cmv1 "github.com/tendze/gRPC_ContactManager_Protos/gen/go/cm"
	"google.golang.org/grpc"
)

type serverAPI struct {
	cmv1.UnimplementedContactManagerServer
}

func Register(gRPC *grpc.Server) {
	cmv1.RegisterContactManagerServer(gRPC, &serverAPI{})
}

func (s *serverAPI) CreateContact(
	ctx context.Context,
	req *cmv1.CreateContactRequest,
) (*cmv1.CreateContactResponse, error) {
	panic("implement me")
}

func (s *serverAPI) GetContactByName(
	ctx context.Context,
	req *cmv1.GetContactByNameRequest,
) (*cmv1.GetContactResponse, error) {
	panic("implement me")
}

func (s *serverAPI) GetContactByEmail(
	ctx context.Context,
	req *cmv1.GetContactByEmailRequest,
) (*cmv1.GetContactResponse, error) {
	panic("implement me")
}

func (s *serverAPI) GetContactByPhone(
	ctx context.Context,
	req *cmv1.GetContactByPhoneRequest,
) (*cmv1.GetContactResponse, error) {
	panic("implement me")
}

func (s *serverAPI) DeleteContact(
	ctx context.Context,
	req *cmv1.DeleteContactRequest,
) (*cmv1.DeleteContactResponse, error) {
	panic("implement me")
}
