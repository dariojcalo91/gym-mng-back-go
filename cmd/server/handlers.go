package main

import (
	"context"

	"github.com/dariojcalo91/gym-backend-go-ver/internal/domain"
	"github.com/dariojcalo91/gym-backend-go-ver/internal/middleware"
	sv "github.com/dariojcalo91/gym-backend-go-ver/internal/service"
	pb "github.com/dariojcalo91/gym-backend-go-ver/proto"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type grpcHandler struct {
	pb.UnimplementedGymServiceServer
	pb.UnimplementedUserServiceServer
	pb.UnimplementedAuthServiceServer
	members         *sv.MemberService
	identityService *sv.IdentityService
}

// --- GymService ---

func (h *grpcHandler) CreateMember(ctx context.Context, req *pb.MemberRequest) (*pb.MemberResponse, error) {
	auth, ok := middleware.FromContext(ctx)
	if !ok {
		return nil, status.Error(codes.Unauthenticated, "missing auth context")
	}

	m := &domain.Member{
		GymID: auth.GymID,
		Name:  req.Name,
		Phone: req.Phone,
		Type:  toDomainMemberType(req.Type),
	}
	if req.MembershipStart != nil {
		t := req.MembershipStart.AsTime()
		m.MembershipStart = &t
	}
	if req.MembershipEnd != nil {
		t := req.MembershipEnd.AsTime()
		m.MembershipEnd = &t
	}

	if err := h.members.CreateMember(ctx, m); err != nil {
		return nil, err
	}

	resp := toMemberResponse(m)
	resp.Message = "Member registered successfully"
	return resp, nil
}

func (h *grpcHandler) ListMembers(ctx context.Context, _ *pb.ListMembersRequest) (*pb.ListMembersResponse, error) {
	auth, ok := middleware.FromContext(ctx)
	if !ok {
		return nil, status.Error(codes.Unauthenticated, "missing auth context")
	}

	members, err := h.members.ListMembers(ctx, auth.GymID)
	if err != nil {
		return nil, err
	}

	resp := &pb.ListMembersResponse{Members: make([]*pb.MemberResponse, 0, len(members))}
	for _, m := range members {
		resp.Members = append(resp.Members, toMemberResponse(m))
	}
	return resp, nil
}

func (h *grpcHandler) UpdateMember(ctx context.Context, req *pb.UpdateMemberRequest) (*pb.MemberResponse, error) {
	auth, ok := middleware.FromContext(ctx)
	if !ok {
		return nil, status.Error(codes.Unauthenticated, "missing auth context")
	}

	m := &domain.Member{
		ID:    req.Id,
		GymID: auth.GymID,
		Name:  req.Name,
		Phone: req.Phone,
	}
	if req.MembershipStart != nil {
		t := req.MembershipStart.AsTime()
		m.MembershipStart = &t
	}
	if req.MembershipEnd != nil {
		t := req.MembershipEnd.AsTime()
		m.MembershipEnd = &t
	}

	if err := h.members.UpdateMember(ctx, m); err != nil {
		return nil, err
	}

	resp := toMemberResponse(m)
	resp.Message = "Member updated successfully"
	return resp, nil
}

func (h *grpcHandler) CheckIn(ctx context.Context, req *pb.CheckInRequest) (*pb.CheckInResponse, error) {
	auth, ok := middleware.FromContext(ctx)
	if !ok {
		return nil, status.Error(codes.Unauthenticated, "missing auth context")
	}

	checkIn, err := h.members.CheckIn(ctx, auth.GymID, req.MemberId, toDomainPaymentStatus(req.PaymentStatus), req.PaymentNote)
	if err != nil {
		return nil, err
	}

	return toCheckInResponse(checkIn), nil
}

func (h *grpcHandler) GetMemberHistory(ctx context.Context, req *pb.IdRequest) (*pb.MemberHistoryResponse, error) {
	auth, ok := middleware.FromContext(ctx)
	if !ok {
		return nil, status.Error(codes.Unauthenticated, "missing auth context")
	}

	history, err := h.members.GetMemberHistory(ctx, auth.GymID, req.Id)
	if err != nil {
		return nil, err
	}

	resp := &pb.MemberHistoryResponse{CheckIns: make([]*pb.CheckInResponse, 0, len(history))}
	for _, c := range history {
		resp.CheckIns = append(resp.CheckIns, toCheckInResponse(c))
	}
	return resp, nil
}

// --- UserService / AuthService ---

func (h *grpcHandler) RegisterUser(ctx context.Context, req *pb.RegisterUserRequest) (*pb.RegisterUserResponse, error) {
	u := &domain.User{
		Username: req.Username,
		Password: req.Password,
		Email:    req.Email,
		Role:     "OWNER", // sin roles diferenciados en el MVP
	}

	result, err := h.identityService.RegisterUser(ctx, u, req.GymName)
	if err != nil {
		return nil, err
	}

	return &pb.RegisterUserResponse{
		UserId:  result.UserID,
		GymId:   result.GymID,
		Message: "User and gym registered successfully",
	}, nil
}

func (h *grpcHandler) GetUserProfile(ctx context.Context, req *pb.UserRequest) (*pb.User, error) {
	user, err := h.identityService.GetUserByUsername(ctx, req.UserId)
	if err != nil {
		return nil, err
	}

	return &pb.User{
		Id:       user.ID,
		Username: user.Username,
		Role:     user.Role,
		GymId:    user.GymID,
	}, nil
}

func (h *grpcHandler) Login(ctx context.Context, req *pb.LoginRequest) (*pb.LoginResponse, error) {
	result, err := h.identityService.LoginUser(ctx, req.Username, req.Password)
	if err != nil {
		return nil, err
	}

	return &pb.LoginResponse{
		Token: result.Token,
		User: &pb.User{
			Id:       result.User.ID,
			Username: result.User.Username,
			Role:     result.User.Role,
			GymId:    result.User.GymID,
		},
	}, nil
}

func (h *grpcHandler) Logout(_ context.Context, _ *pb.LogoutRequest) (*pb.LogoutResponse, error) {
	return &pb.LogoutResponse{Message: "Logged out successfully"}, nil
}
