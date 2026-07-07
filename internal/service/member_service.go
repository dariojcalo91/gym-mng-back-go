package service

import (
	"context"
	"time"

	"github.com/dariojcalo91/gym-backend-go-ver/internal/domain"
)

type MemberRepository interface {
	SaveMember(ctx context.Context, member *domain.Member) error
	GetMemberByID(ctx context.Context, gymID, id string) (*domain.Member, error)
	ListMembers(ctx context.Context, gymID string) ([]*domain.Member, error)
	UpdateMember(ctx context.Context, member *domain.Member) error
}

type CheckInRepository interface {
	SaveCheckIn(ctx context.Context, checkIn *domain.CheckIn) error
	ListCheckInsByMember(ctx context.Context, gymID, memberID string) ([]*domain.CheckIn, error)
}

type MemberService struct {
	members  MemberRepository
	checkIns CheckInRepository
}

func NewMemberService(members MemberRepository, checkIns CheckInRepository) *MemberService {
	return &MemberService{members: members, checkIns: checkIns}
}

func (s *MemberService) CreateMember(ctx context.Context, m *domain.Member) error {
	if err := m.Validate(); err != nil {
		return err
	}
	return s.members.SaveMember(ctx, m)
}

func (s *MemberService) ListMembers(ctx context.Context, gymID string) ([]*domain.Member, error) {
	return s.members.ListMembers(ctx, gymID)
}

func (s *MemberService) UpdateMember(ctx context.Context, m *domain.Member) error {
	if err := m.Validate(); err != nil {
		return err
	}
	return s.members.UpdateMember(ctx, m)
}

// CheckIn records a visit. If the member is an active monthly member, payment
// is never asked — the server decides this itself, it does not trust whatever
// payment status the client sent.
func (s *MemberService) CheckIn(ctx context.Context, gymID, memberID string, requested domain.PaymentStatus, note string) (*domain.CheckIn, error) {
	member, err := s.members.GetMemberByID(ctx, gymID, memberID)
	if err != nil {
		return nil, err
	}

	status := requested
	if member.IsMembershipActive(time.Now()) {
		status = domain.PaymentStatusNotApplicable
	}

	checkIn := &domain.CheckIn{
		GymID:         gymID,
		MemberID:      memberID,
		CheckedInAt:   time.Now(),
		PaymentStatus: status,
		PaymentNote:   note,
	}
	if err := checkIn.Validate(); err != nil {
		return nil, err
	}
	if err := s.checkIns.SaveCheckIn(ctx, checkIn); err != nil {
		return nil, err
	}
	return checkIn, nil
}

func (s *MemberService) GetMemberHistory(ctx context.Context, gymID, memberID string) ([]*domain.CheckIn, error) {
	return s.checkIns.ListCheckInsByMember(ctx, gymID, memberID)
}
