package db

import (
	"context"
	"fmt"
	"github.com/qiniu/qmgo"
)

func NewMgoClient(host, port, user, password, database string) (*qmgo.Client, error) {
	ctx := context.Background()
	uri := fmt.Sprintf("mongodb://%s:%s/%s", host, port, database)
	client, err := qmgo.NewClient(ctx, &qmgo.Config{Uri: uri})

	if err != nil {
		return nil, err
	}

	return client, err
}
