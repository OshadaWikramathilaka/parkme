package repositories

import (
	"context"

	"github.com/dfanso/parkme-backend/internal/models"
	"github.com/qiniu/qmgo"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type WalletRepository struct {
	db *qmgo.Database
}

func NewWalletRepository(db *qmgo.Database) *WalletRepository {
	return &WalletRepository{db: db}
}

func (r *WalletRepository) walletCollection() *qmgo.Collection {
	return r.db.Collection("wallets")
}

func (r *WalletRepository) transactionCollection() *qmgo.Collection {
	return r.db.Collection("transactions")
}

func (r *WalletRepository) FindByUserID(ctx context.Context, userID primitive.ObjectID) (*models.Wallet, error) {
	var wallet models.Wallet
	err := r.walletCollection().Find(ctx, bson.M{"user_id": userID}).One(&wallet)
	if err != nil {
		if err == qmgo.ErrNoSuchDocuments {
			return nil, ErrNotFound
		}
		return nil, err
	}
	return &wallet, nil
}

func (r *WalletRepository) Create(ctx context.Context, wallet *models.Wallet) error {
	_, err := r.walletCollection().InsertOne(ctx, wallet)
	return err
}

func (r *WalletRepository) Update(ctx context.Context, wallet *models.Wallet) error {
	err := r.walletCollection().UpdateOne(ctx, bson.M{"_id": wallet.ID}, bson.M{"$set": wallet})
	if err != nil {
		if err == qmgo.ErrNoSuchDocuments {
			return ErrNotFound
		}
		return err
	}
	return nil
}

func (r *WalletRepository) CreateTransaction(ctx context.Context, transaction *models.Transaction) error {
	_, err := r.transactionCollection().InsertOne(ctx, transaction)
	return err
}

func (r *WalletRepository) GetTransactions(ctx context.Context, walletID primitive.ObjectID) ([]models.Transaction, error) {
	var transactions []models.Transaction
	err := r.transactionCollection().Find(ctx, bson.M{"wallet_id": walletID}).Sort("-created_at").All(&transactions)
	if err != nil {
		return nil, err
	}
	return transactions, nil
}
