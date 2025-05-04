package repositories

import (
	"context"
	"time"

	"github.com/dfanso/parkme-backend/internal/models"

	"github.com/qiniu/qmgo"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type UserRepository struct {
	collection *qmgo.Collection
}

func NewUserRepository(db *qmgo.Database) *UserRepository {
	return &UserRepository{
		collection: db.Collection("users"),
	}
}

func (r *UserRepository) FindOne(ctx context.Context, filter interface{}) (*models.User, error) {
	var user models.User

	err := r.collection.Find(ctx, filter).One(&user)
	if err != nil {
		if err == qmgo.ErrNoSuchDocuments {
			return nil, ErrNotFound
		}
		return nil, err
	}

	return &user, nil
}

func (r *UserRepository) FindAll(ctx context.Context) ([]models.User, error) {
	var users []models.User
	err := r.collection.Find(ctx, bson.M{}).All(&users)
	return users, err
}

func (r *UserRepository) FindByID(ctx context.Context, id primitive.ObjectID) (*models.User, error) {
	var user models.User
	err := r.collection.Find(ctx, bson.M{"_id": id}).One(&user)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *UserRepository) Create(ctx context.Context, user *models.User) error {
	user.CreatedAt = time.Now()
	user.UpdatedAt = time.Now()

	_, err := r.collection.InsertOne(ctx, user)
	return err
}

func (r *UserRepository) Update(ctx context.Context, user *models.User) error {
	user.UpdatedAt = time.Now()

	return r.collection.UpdateOne(
		ctx,
		bson.M{"_id": user.ID},
		bson.M{"$set": user},
	)
}

func (r *UserRepository) Delete(ctx context.Context, id primitive.ObjectID) error {
	return r.collection.RemoveId(ctx, id)
}
