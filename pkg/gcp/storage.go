package gcp

import (
	"context"
	"io/ioutil"
	"os"

	"cloud.google.com/go/storage"
	"github.com/pkg/errors"
	"google.golang.org/api/option"
)

// Google Cloud Platformラッパー

// Strage Strage操作オブジェクト
type Strage struct {
	ctx    context.Context
	client *storage.Client
	bkt    *storage.BucketHandle
}

// NewStrage Strageを返す及び、認証処理(keyFilePathが空なら認証処理は行わない)
// clientのCloseを忘れないこと
func NewStrage(ctx context.Context, keyFilePath, buketName string) (*Strage, error) {
	var err error
	s := new(Strage)

	if ctx != nil {
		s.ctx = ctx
	} else {
		s.ctx = context.Background()
	}

	// Authorization
	if keyFilePath != "" {
		//s.client, err = storage.NewClient(ctx, option.WithServiceAccountFile("../../api-keys/cayenne-dev-strage.json"))
		//s.client, err = storage.NewClient(ctx, option.WithServiceAccountFile(keyFilePath))
		s.client, err = storage.NewClient(ctx, option.WithCredentialsFile(keyFilePath))
	} else {
		s.client, err = storage.NewClient(ctx)
	}
	if err != nil {
		return nil, errors.Errorf("storage.NewClient() ", err)
	}

	// バケット
	if buketName != "" {
		//s.bkt = s.client.Bucket("cayenne-dev-exchanges-yasui-bucket")
		s.bkt = s.client.Bucket(buketName)
	}
	return s, nil
}

// NewBucket 新しいBucketオブジェクトを作成する
func (s *Strage) NewBucket(buketName string) {
	s.bkt = s.client.Bucket(buketName)
}

// Close 切断処理
func (s *Strage) Close() error {
	return s.client.Close()
}

// Write バケットにファイルを書き込む
func (s *Strage) Write(fileName string, p []byte) (err error) {
	w := s.bkt.Object(fileName).NewWriter(s.ctx)
	defer func() {
		if errr := w.Close(); errr == nil {
			err = errr
		}
	}()

	_, err = w.Write(p)
	return
}

// ReadAndSave バケットからファイルを読み込み、保存する
func (s *Strage) ReadAndSave(readFileName, saveFileName string, perm os.FileMode) (err error) {
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
