package teamongo

import (
	"github.com/mongodb/mongo-go-driver/mongo"
	"context"
	"github.com/mongodb/mongo-go-driver/bson"
)

type Collection struct {
	*mongo.Collection
}

func FindCollection(collName string) *Collection {
	return &Collection{
		SharedClient().Database("teaweb").Collection(collName),
	}
}

// 创建索引
func (this *Collection) CreateIndex(indexes map[string]bool) {
	manager := this.Indexes()

	doc := bson.NewDocument()
	for index, b := range indexes {
		if b {
			doc.Append(bson.EC.Int32(index, 1))
		} else {
			doc.Append(bson.EC.Int32(index, -1))
		}
	}

	manager.CreateOne(context.Background(), mongo.IndexModel{
		Keys:    doc,
		Options: bson.NewDocument(bson.EC.Boolean("background", true)),
	})
}
