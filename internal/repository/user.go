package repository

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"linnoxlewis/tg-bot-eng-dictionary/internal/db"
	"linnoxlewis/tg-bot-eng-dictionary/internal/models"
)

type UserRepo struct {
	db             *db.Mongodb
	collectionName string
}

func NewUserRepo(db *db.Mongodb) *UserRepo {
	return &UserRepo{db: db, collectionName: models.GetUserCollectionName()}
}

func (u *UserRepo) GetByTgId(ctx context.Context, id int64) (*models.User, error) {
	filter := bson.D{{"chat_id", id}}
	user := models.NewUserIngot()
	collections := u.db.Database.Collection(u.collectionName)
	if err := collections.FindOne(ctx, filter).Decode(&user); err != nil {
		if err == mongo.ErrNoDocuments {
			return user, nil
		}
		return nil, err
	}

	return user, nil
}

func (u *UserRepo) CreateUser(ctx context.Context, user *models.User) error {
	collections := u.db.Database.Collection(u.collectionName)
	_, err := collections.InsertOne(ctx, user)

	return err
}

func (u *UserRepo) UpdateLearnStatusUser(ctx context.Context, user *models.User) error {
	collection := u.db.Database.Collection(u.collectionName)
	_, err := collection.UpdateOne(ctx,
		bson.M{"_id": user.Id},
		bson.D{
			{"$set", bson.D{{"learn_words_mode", true}}},
		},
	)

	return err
}
