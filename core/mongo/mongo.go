package mongo

import (
	"context"
	"log"
	"time"

	mapset "github.com/deckarep/golang-set"
	"go.mongodb.org/mongo-driver/bson"

	"github.com/zeromicro/go-zero/core/logx"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Client struct {
	Client *mongo.Client
	Db     *mongo.Database
}

// MustNewClient returns a new mongo client
func MustNewClient(uri string, db string) *Client {
	client, err := NewClient(uri, db)
	if err != nil {
		log.Fatalf("create mongo client err: %v", err)
	}

	return client
}

// NewClient returns a Client.
func NewClient(uri string, db string) (*Client, error) {
	opts := options.Client().ApplyURI(uri)
	client, err := mongo.Connect(context.Background(), opts)
	if err != nil {
		logx.Errorf("[db] stores connect err: %v", err)
		return nil, err
	}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err = client.Ping(ctx, nil); err != nil {
		logx.Errorf("[db] stores ping err: %v", err)
		return nil, err
	}
	logx.Info("[db] stores ping success")
	// insert for ping
	driver := client.Database(db)
	return &Client{
		Client: client,
		Db:     driver,
	}, nil
}

// GetSession get tran session ctx
func GetSession(ctx context.Context, client *mongo.Client) (mongo.Session, mongo.SessionContext, error) {
	// 启动事务
	sess, err := client.StartSession()
	if err != nil {
		logx.Errorf("开启 session 失败: %v\n", err)
		return nil, nil, err
	}

	sessionCtx := mongo.NewSessionContext(ctx, sess)

	if err = sess.StartTransaction(); err != nil {
		logx.Errorf("开始事务失败: %v\n", err)
		return nil, nil, err
	}

	return sess, sessionCtx, nil
}

// InitCollection 初始化集合
func (c *Client) InitCollection(ctx context.Context, newCols []string) {
	cols, err := c.Db.ListCollectionNames(context.TODO(), bson.M{})
	if err != nil {
		logx.Errorf("列出 cols 失败: %v", err)
		return
	}

	logx.Infof("已有的集合: %v", cols)

	// 取出差集
	originSet := mapset.NewThreadUnsafeSet()
	for i := 0; i < len(cols); i++ {
		originSet.Add(cols[i])
	}
	newSet := mapset.NewThreadUnsafeSet()
	for i := 0; i < len(newCols); i++ {
		newSet.Add(newCols[i])
	}

	needCreateCols := newSet.Difference(originSet).ToSlice()
	logx.Infof("需要创建的集合: %v", needCreateCols)

	if len(needCreateCols) == 0 {
		logx.Infof("不需要重新创建集合")
		return
	}

	for i := 0; i < len(needCreateCols); i++ {
		name, ok := needCreateCols[i].(string)
		if !ok {
			continue
		}
		if err = c.Db.CreateCollection(ctx, name); err != nil {
			logx.Errorf("创建集合失败: %v", err)
			return
		}
		logx.Infof("创建集合 %s 成功", name)
	}

	return
}

func (c *Client) Close(ctx context.Context) error {
	return c.Client.Disconnect(ctx)
}
