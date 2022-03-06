package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
	Id             primitive.ObjectID `bson:"_id,omitempty" json:"_id"`
	ChatId         int64              `bson:"chat_id,omitempty" json:"chatId"`
	Name           string             `bson:"name,omitempty" json:"name"`
	LearnWordsMode bool               `bson:"learn_words_mode," json:"learnMode"`
}

func NewUserIngot() *User {
	return &User{}
}

func NewUser(chatId int64, name string) *User {
	return &User{
		Id:             primitive.NewObjectID(),
		ChatId:         chatId,
		Name:           name,
		LearnWordsMode: false,
	}
}

func GetUserCollectionName() string {
	return "users"
}

func (u *User) Empty() bool {
	return u.Id.IsZero()
}
