package cm

import (
	"context"
	"errors"
	"gRPC_ContactManagement_Service/internal/domain/models"
	"gRPC_ContactManagement_Service/internal/service/cm"
	cmv1 "github.com/tendze/gRPC_ContactManager_Protos/gen/go/cm"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"regexp"
)

const emailContextKey = "email"

type ContactManager interface {
	CreateContact(
		ctx context.Context,
		creatorEmail, name, email, phone string,
	) (uid int64, err error)

	GetContactByName(
		ctx context.Context,
		creatorEmail string,
		name string,
	) (models.Contact, error)

	GetContactByEmail(
		ctx context.Context,
		creatorEmail string,
		email string,
	) (models.Contact, error)

	GetContactByPhone(
		ctx context.Context,
		creatorEmail string,
		phone string,
	) (models.Contact, error)

	DeleteContact(
		ctx context.Context,
		creatorEmail string,
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
	creatorEmail, err := getEmailFromContext(ctx)
	if err != nil {
		return nil, err
	}
	uid, err := s.cm.CreateContact(ctx, creatorEmail, req.GetName(), req.GetEmail(), req.GetPhone())
	if err != nil {
		if errors.Is(err, cm.ErrContactExists) {
			return nil, status.Error(codes.AlreadyExists, "contact already exists")
		}
		return nil, status.Error(codes.Internal, "cannot add new contact")
	}
	return &cmv1.CreateContactResponse{Id: uid, Success: true}, nil
}

func (s *serverAPI) GetContactByName(
	ctx context.Context,
	req *cmv1.GetContactByNameRequest,
) (*cmv1.GetContactResponse, error) {
	if req.GetName() == "" {
		return nil, status.Error(codes.InvalidArgument, "name required")
	}
	creatorEmail, err := getEmailFromContext(ctx)
	if err != nil {
		return nil, err
	}

	contact, err := s.cm.GetContactByName(ctx, creatorEmail, req.GetName())
	if err != nil {
		if errors.Is(err, cm.ErrContactNotFound) {
			return nil, status.Error(codes.NotFound, "contact not found")
		}
		return nil, status.Error(codes.Internal, "cannot find contact")
	}
	return &cmv1.GetContactResponse{
		Id:    contact.ID,
		Name:  contact.Name,
		Email: contact.Email,
		Phone: contact.Phone,
	}, nil
}

func (s *serverAPI) GetContactByEmail(
	ctx context.Context,
	req *cmv1.GetContactByEmailRequest,
) (*cmv1.GetContactResponse, error) {
	if req.GetEmail() == "" {
		return nil, status.Error(codes.InvalidArgument, "email required")
	}

	creatorEmail, err := getEmailFromContext(ctx)
	if err != nil {
		return nil, err
	}

	contact, err := s.cm.GetContactByName(ctx, creatorEmail, req.GetEmail())
	if err != nil {
		if errors.Is(err, cm.ErrContactNotFound) {
			return nil, status.Error(codes.NotFound, "contact not found")
		}
		return nil, status.Error(codes.Internal, "cannot find contact")
	}
	return &cmv1.GetContactResponse{
		Id:    contact.ID,
		Name:  contact.Name,
		Email: contact.Email,
		Phone: contact.Phone,
	}, nil
}

func (s *serverAPI) GetContactByPhone(
	ctx context.Context,
	req *cmv1.GetContactByPhoneRequest,
) (*cmv1.GetContactResponse, error) {
	if req.GetPhone() == "" {
		return nil, status.Error(codes.InvalidArgument, "phone required")
	}

	creatorEmail, err := getEmailFromContext(ctx)
	if err != nil {
		return nil, err
	}

	contact, err := s.cm.GetContactByName(ctx, creatorEmail, req.GetPhone())
	if err != nil {
		if errors.Is(err, cm.ErrContactNotFound) {
			return nil, status.Error(codes.NotFound, "contact not found")
		}
		return nil, status.Error(codes.Internal, "cannot find contact")
	}
	return &cmv1.GetContactResponse{
		Id:    contact.ID,
		Name:  contact.Name,
		Email: contact.Email,
		Phone: contact.Phone,
	}, nil
}

func (s *serverAPI) DeleteContact(
	ctx context.Context,
	req *cmv1.DeleteContactRequest,
) (*cmv1.DeleteContactResponse, error) {
	if req.GetId() <= 0 {
		return nil, status.Error(codes.InvalidArgument, "id required")
	}

	creatorEmail, err := getEmailFromContext(ctx)
	if err != nil {
		return nil, err
	}

	err = s.cm.DeleteContact(ctx, creatorEmail, req.GetId())

	if err != nil {
		if errors.Is(err, cm.ErrContactNotFound) {
			return nil, status.Error(codes.InvalidArgument, "contact not found")
		}
		return nil, status.Error(codes.Internal, "cannot find contact")
	}

	return &cmv1.DeleteContactResponse{Success: true}, nil
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

// Extracts email from context
func getEmailFromContext(ctx context.Context) (string, error) {
	creatorEmail, ok := ctx.Value(emailContextKey).(string)
	if !ok {
		return "", status.Error(codes.Internal, "cannot get user email")
	}
	return creatorEmail, nil
}
