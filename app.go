package main

import (
	"bytes"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/kms"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3crypto"
	"io"
	"io/ioutil"
	"os"
	"strings"
)

type App struct {
	Bucket string
	Prefix string
	KeyID  string

	s3  *s3.S3
	enc *s3crypto.EncryptionClient
	dec *s3crypto.DecryptionClient
}

func NewApp(bucket string, prefix string, keyID string) *App {
	if prefix != "" && !strings.HasSuffix(prefix, "/") {
		prefix += "/"
	}

	sess := session.Must(session.NewSession())
	svc := s3.New(sess)

	handler := s3crypto.NewKMSKeyGenerator(kms.New(sess), keyID)
	enc := s3crypto.NewEncryptionClient(sess, s3crypto.AESGCMContentCipherBuilder(handler))
	dec := s3crypto.NewDecryptionClient(sess)

	return &App{
		Bucket: bucket,
		Prefix: prefix,
		KeyID:  keyID,

		s3:  svc,
		enc: enc,
		dec: dec,
	}
}

func (a *App) ListIter(fn func(string) bool) error {
	return a.s3.ListObjectsPages(&s3.ListObjectsInput{
		Bucket: aws.String(a.Bucket),
		Prefix: aws.String(a.Prefix),
	}, func(out *s3.ListObjectsOutput, b bool) bool {
		for _, o := range out.Contents {
			key := aws.StringValue(o.Key)
			name := key[len(a.Prefix):]
			if ret := fn(name); !ret {
				return false
			}
		}
		return true
	})
}

func (a *App) Put(name string, val string) error {
	key := a.Prefix + name

	if a.KeyID == "" {
		return fmt.Errorf("key ID is not specified")
	}

	var body io.ReadSeeker
	if val == "-" {
		b, err := ioutil.ReadAll(os.Stdin)
		if err != nil {
			return err
		}

		body = bytes.NewReader(b)
	} else if val[0] == '@' {
		f, err := os.Open(val[1:])
		if err != nil {
			return err
		}
		defer f.Close()

		body = f
	} else {
		body = bytes.NewReader([]byte(val))
	}

	_, err := a.enc.PutObject(&s3.PutObjectInput{
		Bucket: aws.String(a.Bucket),
		Key:    aws.String(key),
		Body:   body,
	})

	return err
}

func (a *App) Get(name string) error {
	key := a.Prefix + name

	out, err := a.dec.GetObject(&s3.GetObjectInput{
		Bucket: aws.String(a.Bucket),
		Key:    aws.String(key),
	})
	if err != nil {
		return err
	}

	body, err := ioutil.ReadAll(out.Body)
	if err != nil {
		return err
	}

	fmt.Print(string(body))

	return nil
}

func (a *App) Delete(name string) error {
	_, err := a.s3.DeleteObject(&s3.DeleteObjectInput{
		Bucket: aws.String(a.Bucket),
		Key:    aws.String(a.Prefix + name),
	})

	return err
}
