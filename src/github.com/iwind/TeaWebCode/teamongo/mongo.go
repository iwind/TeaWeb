package teamongo

import (
	"github.com/mongodb/mongo-go-driver/mongo"
	"github.com/iwind/TeaGo/files"
	"github.com/iwind/TeaGo/Tea"
	"github.com/iwind/TeaGo/logs"
	"github.com/syndtr/goleveldb/leveldb/errors"
	"context"
)

var sharedClient *mongo.Client

func SharedClient() *mongo.Client {
	if sharedClient == nil {
		configFile := files.NewFile(Tea.ConfigFile("mongo.conf"))
		if !configFile.Exists() {
			logs.Fatal(errors.New("'mongo.conf' not found"))
			return nil
		}
		reader, err := configFile.Reader()
		if err != nil {
			logs.Fatal(err)
			return nil
		}


		config := &Config{}
		err = reader.ReadYAML(config)
		if err != nil {
			logs.Fatal(err)
			return nil
		}

		sharedClient, err = mongo.NewClient(config.URI)
		if err != nil {
			logs.Fatal(err)
			return nil
		}

		err = sharedClient.Connect(context.Background())
		if err != nil {
			logs.Fatal(err)
			return nil
		}
	}

	return sharedClient
}
