package dto

import "github.com/dariojcalo91/gym-backend-go-ver/internal/domain"

// MemberWithOwnerEmail pairs a member with their gym owner's email.
// It exists only to carry the result of a query that joins members -> gyms -> users.
// It is NOT a domain concept — internal/domain must never import this package.
type MemberWithOwnerEmail struct {
	Member     *domain.Member
	OwnerEmail string
}
