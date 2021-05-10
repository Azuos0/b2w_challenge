package controller

import "go.mongodb.org/mongo-driver/mongo"

type Controller interface {
	SetService(*mongo.Database)
}
