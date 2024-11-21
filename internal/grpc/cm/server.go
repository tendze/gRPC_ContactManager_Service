package cm

import (
	"context"
	"gRPC_ContactManagement_Service/internal/domain/models"
	cmv1 "github.com/tendze/gRPC_ContactManager_Protos/gen/go/cm"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"regexp"
)

type ContactManager interface {
	CreateContact(
		ctx context.Context,
		creatorEmail, name, email, phone string,
	) (uid int64, err error)

	GetContactByName(
		ctx context.Context,
		name string,
	) (models.Contact, error)

	GetContactByEmail(
		ctx context.Context,
		email string,
	) (models.Contact, error)

	GetContactByPhone(
		ctx context.Context,
		phone string,
	) (models.Contact, error)

	DeleteContact(
		ctx context.Context,
		id int64,
	) error
}

type serverAPI struct {
	cmv1.UnimplementedContactManagerServer
	cm ContactManager
}

func Register(gRPC *grpc.Server, cm ContactManager) {
	cmv1.RegisterContactManagerServer(gRPC, &serverAPI{cm: cm})
}

func (s *serverAPI) CreateContact(
	ctx context.Context,
	req *cmv1.CreateContactRequest,
) (*cmv1.CreateContactResponse, error) {
	if err := validateCreateContactRequest(req); err != nil {
		return nil, err
	}
	if err := validateEmail(req.GetEmail()); err != nil {
		return nil, err
	}
	if err := validatePhone(req.GetPhone()); err != nil {
		return nil, err
	}

	// TODO: implement
	return &cmv1.CreateContactResponse{}, nil
}

func (s *serverAPI) GetContactByName(
	ctx context.Context,
	req *cmv1.GetContactByNameRequest,
) (*cmv1.GetContactResponse, error) {
	if req.GetName() == "" {
		return nil, status.Error(codes.InvalidArgument, "name required")
	}

	return &cmv1.GetContactResponse{}, nil
}

func (s *serverAPI) GetContactByEmail(
	ctx context.Context,
	req *cmv1.GetContactByEmailRequest,
) (*cmv1.GetContactResponse, error) {
	if req.GetEmail() == "" {
		return nil, status.Error(codes.InvalidArgument, "email required")
	}

	return &cmv1.GetContactResponse{}, nil
}

func (s *serverAPI) GetContactByPhone(
	ctx context.Context,
	req *cmv1.GetContactByPhoneRequest,
) (*cmv1.GetContactResponse, error) {
	if req.GetPhone() == "" {
		return nil, status.Error(codes.InvalidArgument, "phone required")
	}

	return &cmv1.GetContactResponse{}, nil
}

func (s *serverAPI) DeleteContact(
	ctx context.Context,
	req *cmv1.DeleteContactRequest,
) (*cmv1.DeleteContactResponse, error) {
	if req.GetId() <= 0 {
		return nil, status.Error(codes.InvalidArgument, "id required")
	}

	return &cmv1.DeleteContactResponse{}, nil
}

func validateCreateContactRequest(req *cmv1.CreateContactRequest) error {
	if req.GetName() == "" {
		return status.Error(codes.InvalidArgument, "name required")
	}
	if req.GetEmail() == "" {
		return status.Error(codes.InvalidArgument, "email required")
	}
	if req.GetPhone() == "" {
		return status.Error(codes.InvalidArgument, "phone required")
	}
	return nil
}

func validateEmail(email string) error {
	// Регулярное выражение для проверки email
	re := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	if !re.MatchString(email) {
		return status.Error(codes.InvalidArgument, "invalid email format")
	}
	return nil
}

func validatePhone(phone string) error {
	re := regexp.MustCompile(`^\+7\d{10}$`)
	if !re.MatchString(phone) {
		return status.Error(codes.InvalidArgument, "invalid phone format")
	}
	return nil
}
