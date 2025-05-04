package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type TransactionType string

const (
	TransactionTypeTopUp  TransactionType = "top_up"
	TransactionTypeDeduct TransactionType = "deduct"
	TransactionTypeRefund TransactionType = "refund"
)

type Transaction struct {
	ID          primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	WalletID    primitive.ObjectID `bson:"wallet_id" json:"walletId"`
	Type        TransactionType    `bson:"type" json:"type"`
	Amount      float64            `bson:"amount" json:"amount"`
	Description string             `bson:"description" json:"description"`
	CreatedAt   time.Time          `bson:"created_at" json:"createdAt"`
}

type Wallet struct {
	ID        primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	UserID    primitive.ObjectID `bson:"user_id" json:"userId"`
	Balance   float64            `bson:"balance" json:"balance"`
	CreatedAt time.Time          `bson:"created_at" json:"createdAt"`
	UpdatedAt time.Time          `bson:"updated_at" json:"updatedAt"`
}
