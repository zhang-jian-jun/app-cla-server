package mongodb

import (
	"context"
	"fmt"
	"time"

	"github.com/huaweicloud/golangsdk"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"

	"github.com/zengchen1024/cla-server/models"
)

const orgRepoCollection = "org_repos"

func (c *client) CreateOrgRepo(orgRepo models.OrgRepo) (string, error) {
	body, err := golangsdk.BuildRequestBody(orgRepo, "")
	if err != nil {
		return "", fmt.Errorf("build body failed, err:%v", err)
	}

	col := c.db.Collection(orgRepoCollection)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	r, err := col.InsertOne(ctx, bson.M(body))
	if err != nil {
		return "", fmt.Errorf("write db failed, err:%v", err)
	}

	v, ok := r.InsertedID.(primitive.ObjectID)
	if !ok {
		return "", fmt.Errorf("retrieve id failed")
	}

	return v.String(), nil
}