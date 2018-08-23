package gcp

import (
	"context"
	"io/ioutil"
	"os"

	"cloud.google.com/go/storage"
	"github.com/pkg/errors"
	"google.golang.org/api/option"
	"strconv"
	"time"
)

// Storage Storage操作オブジェクト
type Storage struct {
	bucketName  string
	keyFilePath string
	ext         *ExtClient
}

// ExtClient client拡張オブジェクト
type ExtClient struct {
	ctx    context.Context
	client *storage.Client
	bkt    *storage.BucketHandle
}

// NewStorage バケット名を持つStorageオブジェクトを返す
func NewStorage(bucketName, keyFilePath string) *Storage {
	return &Storage{bucketName: bucketName, keyFilePath: keyFilePath}
}

// NewClient Storageを返す及び、認証処理(keyFilePathが空なら認証処理は行わない)
// clientのCloseを忘れないこと
func (s *Storage) NewClient(ctx context.Context) (*ExtClient, error) {
	var err error

	ext := new(ExtClient)

	if ctx != nil {
		ext.ctx = ctx
	} else {
		ext.ctx = context.Background()
	}

	// Authorization
	if s.keyFilePath != "" {
		ext.client, err = storage.NewClient(ctx, option.WithCredentialsFile(s.keyFilePath))
	} else {
		ext.client, err = storage.NewClient(ctx)
	}
	if err != nil {
		return nil, errors.Errorf("storage.NewClient() error: %v", err)
	}

	// バケット
	if s.bucketName != "" {
		ext.bkt = ext.client.Bucket(s.bucketName)
	}

	s.ext = ext

	return ext, nil
}

// WriteOnce セッションを保持しない、一度きりの書き込み処理
func (s *Storage) WriteOnce(path, hex string) (string, error) {
	cli, err := s.NewClient(context.Background())
	if err != nil {
		return "", errors.Errorf("storage.NewClient(): error: %v", err)
	}

	generatedFileName, err := cli.Write(path, []byte(hex))
	if err != nil {
		return "", errors.Errorf("storage.Write(): error: %v", err)
	}

	err = cli.Close()
	if err != nil {
		return "", errors.Errorf("storage.Close(): error: %v", err)
	}

	return generatedFileName, nil
}

// ReadOnce セッションを保持しない、一度きりの読み込み処理
func (s *Storage) ReadOnce(path, outputPath string) error {
	cli, err := s.NewClient(context.Background())
	if err != nil {
		return errors.Errorf("storage.NewClient(): error: %v", err)
	}

	err = cli.ReadAndSave(path, outputPath, 0666)
	if err != nil {
		return errors.Errorf("storage.Write(): error: %v", err)
	}

	err = cli.Close()
	if err != nil {
		return errors.Errorf("storage.Close(): error: %v", err)
	}

	return nil
}

// NewBucket 新しいBucketオブジェクトを作成する
func (e *ExtClient) NewBucket(bucketName string) {
	e.bkt = e.client.Bucket(bucketName)
}

// Close 切断処理
func (e *ExtClient) Close() error {
	return e.client.Close()
}

// Write バケットにファイルを書き込む
func (e *ExtClient) Write(path string, p []byte) (fileName string, err error) {
	ts := strconv.FormatInt(time.Now().UnixNano(), 10)
	fileName = path + ts

	w := e.bkt.Object(fileName).NewWriter(e.ctx)
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
func (e *ExtClient) ReadAndSave(readFileName, saveFileName string, perm os.FileMode) (err error) {
	// Read file from bucket.
	var r *storage.Reader
	r, err = e.bkt.Object(readFileName).NewReader(e.ctx)
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
