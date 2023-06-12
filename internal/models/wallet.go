package models

import (
	"database/sql"
	"time"

	"github.com/google/uuid"
	"google.golang.org/protobuf/types/known/timestamppb"

	walletsService "github.com/majidbl/wallet/proto/wallet"
)

// Wallet models
type Wallet struct {
	ID          uuid.UUID `json:"id"`
	Name        string    `json:"name"`
	Mobile      string    `json:"mobile,omitempty" validate:"required,min=3,max=250"`
	Balance     int64     `json:"balance,omitempty" validate:"required"`
	Avatar      *string   `json:"avatar,omitempty"`
	Description string    `json:"description,omitempty" validate:"required,min=3,max=500"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
}

type WalletBalanceResponse struct {
	Balance int64 `json:"balance"`
}

// WalletModel models
type WalletModel struct {
	ID          uuid.UUID
	Name        sql.NullString
	Mobile      sql.NullString
	Balance     sql.NullInt64
	Avatar      sql.NullString
	Description sql.NullString
	CreatedAt   sql.NullTime
	UpdatedAt   sql.NullTime
}

func (m WalletModel) Entity() *Wallet {
	return &Wallet{
		ID:          m.ID,
		Name:        m.Name.String,
		Mobile:      m.Mobile.String,
		Balance:     m.Balance.Int64,
		Avatar:      &m.Avatar.String,
		Description: m.Description.String,
		CreatedAt:   m.CreatedAt.Time,
		UpdatedAt:   m.UpdatedAt.Time,
	}
}

type ChargeWalletReq struct {
	Mobile    string    `json:"mobile"`
	Amount    int64     `json:"balance"`
	UpdatedAt time.Time `json:"updatedAt"`
}

type UpdateWalletBalanceReq struct {
	WalletID uuid.UUID `json:"wallet_id"`
	Amount   int64     `json:"amount"`
}

func (p *Wallet) GetImage() string {
	var avatar string
	if p.Avatar != nil {
		avatar = *p.Avatar
	}
	return avatar
}

// ToProto Convert wallet to proto
func (p *Wallet) ToProto() *walletsService.Wallet {
	return &walletsService.Wallet{
		ID:          p.ID.String(),
		Mobile:      p.Mobile,
		Name:        p.Name,
		Description: p.Description,
		Balance:     p.Balance,
		Avatar:      p.GetImage(),
		CreatedAt:   timestamppb.New(p.CreatedAt),
		UpdatedAt:   timestamppb.New(p.UpdatedAt),
	}
}

// WalletFromProto Get Wallet from proto
func WalletFromProto(wallet *walletsService.Wallet) (*Wallet, error) {
	walletID, err := uuid.Parse(wallet.ID)
	if err != nil {
		return nil, err
	}

	return &Wallet{
		ID:          walletID,
		Name:        wallet.GetName(),
		Description: wallet.GetDescription(),
		Balance:     wallet.Balance,
		Avatar:      &wallet.Avatar,
		CreatedAt:   wallet.GetCreatedAt().AsTime(),
		UpdatedAt:   wallet.GetUpdatedAt().AsTime(),
	}, nil
}

// WalletsList All Wallets response with pagination
type WalletsList struct {
	TotalCount int64     `json:"totalCount"`
	TotalPages int64     `json:"totalPages"`
	Page       int64     `json:"page"`
	Size       int64     `json:"size"`
	HasMore    bool      `json:"hasMore"`
	Wallets    []*Wallet `json:"wallets"`
}

// ToProtoList convert wallets list to proto
func (p *WalletsList) ToProtoList() []*walletsService.Wallet {
	walletsList := make([]*walletsService.Wallet, 0, len(p.Wallets))
	for _, wallet := range p.Wallets {
		walletsList = append(walletsList, wallet.ToProto())
	}
	return walletsList
}

type WalletErrorMsg struct {
	Subject   string    `json:"subject"`
	Sequence  uint64    `json:"sequence"`
	Data      []byte    `json:"data"`
	Timestamp int64     `json:"topic"`
	Error     string    `json:"error"`
	Time      time.Time `json:"time"`
}
