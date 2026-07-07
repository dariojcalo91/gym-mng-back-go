package main

import (
	"time"

	"github.com/dariojcalo91/gym-backend-go-ver/internal/domain"
	pb "github.com/dariojcalo91/gym-backend-go-ver/proto"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// Mapping helpers: proto <-> domain. Nothing in internal/domain or
// internal/service knows these types exist; the translation lives only here,
// at the adapter boundary.

func toDomainMemberType(t pb.MemberType) domain.MemberType {
	if t == pb.MemberType_MEMBER_TYPE_MONTHLY {
		return domain.MemberTypeMonthly
	}
	return domain.MemberTypeOccasional
}

func toProtoMemberType(t domain.MemberType) pb.MemberType {
	if t == domain.MemberTypeMonthly {
		return pb.MemberType_MEMBER_TYPE_MONTHLY
	}
	return pb.MemberType_MEMBER_TYPE_OCCASIONAL
}

func toDomainPaymentStatus(s pb.PaymentStatus) domain.PaymentStatus {
	switch s {
	case pb.PaymentStatus_PAYMENT_STATUS_PAID:
		return domain.PaymentStatusPaid
	case pb.PaymentStatus_PAYMENT_STATUS_ON_CREDIT:
		return domain.PaymentStatusPending
	default:
		return domain.PaymentStatusNotApplicable
	}
}

func toProtoPaymentStatus(s domain.PaymentStatus) pb.PaymentStatus {
	switch s {
	case domain.PaymentStatusPaid:
		return pb.PaymentStatus_PAYMENT_STATUS_PAID
	case domain.PaymentStatusPending:
		return pb.PaymentStatus_PAYMENT_STATUS_ON_CREDIT
	default:
		return pb.PaymentStatus_PAYMENT_STATUS_NOT_APPLICABLE
	}
}

func toMemberResponse(m *domain.Member) *pb.MemberResponse {
	resp := &pb.MemberResponse{
		Id:               m.ID,
		Name:             m.Name,
		Phone:            m.Phone,
		Type:             toProtoMemberType(m.Type),
		MembershipActive: m.IsMembershipActive(time.Now()),
	}
	if m.MembershipEnd != nil {
		resp.MembershipEnd = timestamppb.New(*m.MembershipEnd)
	}
	return resp
}

func toCheckInResponse(c *domain.CheckIn) *pb.CheckInResponse {
	return &pb.CheckInResponse{
		Id:            c.ID,
		CheckedInAt:   timestamppb.New(c.CheckedInAt),
		PaymentStatus: toProtoPaymentStatus(c.PaymentStatus),
	}
}
