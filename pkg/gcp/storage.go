package gcp

import (
	"context"
	"fmt"
	"io/ioutil"
	"os"

	"cloud.google.com/go/storage"
	"github.com/hiromaily/go-bitcoin/pkg/enum"
	"github.com/pkg/errors"
	"google.golang.org/api/option"
	"strconv"
	"time"
)

// Google Cloud Platformラッパー

// Sotrage Storage操作オブジェクト
type Storage struct {
	bucketName  string
	keyFilePath string
	ctx         context.Context
	client      *storage.Client
	bkt         *storage.BucketHandle
}

//type Client struct {
//	ctx    context.Context
//	client *storage.Client
//	bkt    *storage.BucketHandle
//}

//TODO:どっかにまとめる。fileパッケージとかぶる
func CreateFilePath(actionType enum.ActionType, txType enum.TxType, txID int64) string {
	// receipt_8_unsigned_1534744535097796209
	return fmt.Sprintf("%s_%d_%s_", string(actionType), txID, txType)
}

// NewStorage バケット名を持つStorageオブジェクトを返す
func NewStorage(bucketName, keyFilePath string) *Storage {
	return &Storage{bucketName: bucketName, keyFilePath: keyFilePath}
}

// NewClient Storageを返す及び、認証処理(keyFilePathが空なら認証処理は行わない)
// clientのCloseを忘れないこと
func (s *Storage) NewClient(ctx context.Context) error {
	var err error

	if ctx != nil {
		s.ctx = ctx
	} else {
		s.ctx = context.Background()
	}

	// Authorization
	if s.keyFilePath != "" {
		//s.client, err = storage.NewClient(ctx, option.WithServiceAccountFile("../../api-keys/cayenne-dev-strage.json"))
		s.client, err = storage.NewClient(ctx, option.WithCredentialsFile(s.keyFilePath))
	} else {
		s.client, err = storage.NewClient(ctx)
	}
	if err != nil {
		return errors.Errorf("storage.NewClient() error: %v", err)
	}

	// バケット
	if s.bucketName != "" {
		//s.bkt = s.client.Bucket("cayenne-dev-exchanges-yasui-bucket")
		s.bkt = s.client.Bucket(s.bucketName)
	}

	return nil
}

// NewBucket 新しいBucketオブジェクトを作成する
func (s *Storage) NewBucket(bucketName string) {
	s.bkt = s.client.Bucket(bucketName)
}

// Close 切断処理
func (s *Storage) Close() error {
	return s.client.Close()
}

// Write バケットにファイルを書き込む
func (s *Storage) Write(path string, p []byte) (fileName string, err error) {
	ts := strconv.FormatInt(time.Now().UnixNano(), 10)
	fileName = path + ts

	w := s.bkt.Object(fileName).NewWriter(s.ctx)
	_, err = w.Write(p)
	if err != nil {
		err = errors.Errorf("w.Write() error: %v", err)
	}

	if err2 := w.Close(); err2 != nil {
		err = errors.Errorf("w.Close() error: %v", err2)
	}

	return
}

// ReadAndSave バケットからファイルを読み込み、保存する
func (s *Storage) ReadAndSave(readFileName, saveFileName string, perm os.FileMode) (err error) {
	// Read file from bucket.
	var r *storage.Reader
	r, err = s.bkt.Object(readFileName).NewReader(s.ctx)
	if err != nil {
		err = errors.Errorf("bkt.Object() ", err)
		return
	}
	defer func() {
		if errr := r.Close(); errr == nil {
			err = errr
		}
	}()

	var body []byte
	body, err = ioutil.ReadAll(r)
	if err != nil {
		err = errors.Errorf("ioutil.ReadAll() ", err)
		return
	}
	// Save
	err = ioutil.WriteFile(saveFileName, body, perm)
	if err != nil {
		err = errors.Errorf("ioutil.WriteFile() ", err)
	}

	return
}
