/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"context"
	"fmt"
	"github.com/minio/minio-go/v7"
	"io/ioutil"
	"log"
	"strconv"
	"strings"

	"github.com/minio/minio-go/v7/pkg/credentials"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// upfileCmd represents the upfile command
type MinIoConfig struct {
	Host string `mapstructure:"host"`
	Port int    `mapstructure:"port"`
}
type Config struct {
	MinIoInfo MinIoConfig `mapstructure:"minio_url"`
	AccessKey string      `mapstructure:"access_key"`
	SecretKey string      `mapstructure:"secret_key"`
	DataDir   string      `mapstructure:"data_dir"`
}
type client struct {
	client         *minio.Client
	bucketName     string
	targetFilePath string
	useSSL         bool
	location       string
}

func readConfig() *Config {
	v := viper.New()
	v.SetConfigFile("./config.yaml")
	if err := v.ReadInConfig(); err != nil {
		panic(err)
	}
	Config := Config{}
	if err := v.Unmarshal(&Config); err != nil {
		panic(err)
	}
	return &Config
}

func NewClient(bucketName, location, targetFilePath string) *client {
	Config := readConfig()
	endpoint := fmt.Sprintf("%s:%s", Config.MinIoInfo.Host, strconv.Itoa(Config.MinIoInfo.Port))
	accessKeyID := Config.AccessKey
	secretAccessKey := Config.SecretKey
	useSSL := false
	target, err := minio.New(endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(accessKeyID, secretAccessKey, ""),
		Secure: useSSL,
	})
	if err != nil {
		log.Fatalln(err)
		return nil
	}
	return &client{
		client:         target,
		bucketName:     bucketName,
		location:       location,
		targetFilePath: targetFilePath,
	}
}

func (c *client) checkBucket() {
	ctx := context.Background()
	isExists, err := c.client.BucketExists(ctx, c.bucketName)
	if err != nil {
		log.Println("check bucket exist error ")
		return
	}
	if !isExists {
		err2 := c.client.MakeBucket(ctx, c.bucketName, minio.MakeBucketOptions{Region: c.location})
		if err2 != nil {
			log.Println("MakeBucket error ")
			fmt.Println(err2)
			return
		}
		// 权限设置
		policy := fmt.Sprintf(`{"Version":"2012-10-17","Statement":[{"Effect":"Allow","Principal":{"AWS":["*"]},"Action":["s3:GetBucketLocation","s3:ListBucketMultipartUploads"],"Resource":["arn:aws:s3:::%s"]},{"Effect":"Allow","Principal":{"AWS":["*"]},"Action":["s3:ListBucket"],"Resource":["arn:aws:s3:::%s"],"Condition":{"StringEquals":{"s3:prefix":["*"]}}},{"Effect":"Allow","Principal":{"AWS":["*"]},"Action":["s3:AbortMultipartUpload","s3:DeleteObject","s3:GetObject","s3:ListMultipartUploadParts","s3:PutObject"],"Resource":["arn:aws:s3:::%s/**"]}]}`, c.bucketName, c.bucketName, c.bucketName)
		err = c.client.SetBucketPolicy(ctx, c.bucketName, policy)
		if err != nil {
			log.Println(err)
			return
		}
		log.Printf("Successfully created %s\n", c.bucketName)
	}
}
func (c *client) UpLoadFile(path string) {
	ctx := context.Background()
	rd, _ := ioutil.ReadDir(path)
	for _, fi := range rd {
		if fi.IsDir() {
			c.UpLoadFile(path + "/" + fi.Name())
		} else {
			fullPath := path + "/" + fi.Name()
			rawPathLength := len(c.targetFilePath)
			objectName := fullPath[rawPathLength:]
			objectName = strings.TrimLeft(objectName, "/")
			//log.Printf("fullPath=%s  ,objectName=%s\n", fullPath, objectName)
			n, err := c.client.FPutObject(ctx, c.bucketName, objectName, fullPath, minio.PutObjectOptions{
				ContentType: "",
			})
			if err != nil {
				fmt.Println(err)
				return
			}
			fmt.Println("Successfully uploaded bytes: ", n)
		}
	}
}

var upfileCmd = &cobra.Command{
	Use:   "upfile",
	Short: "Create Minio bucket name you DataDir...",
	Long: `Create Minio bucket name you DataDir and Upload DataDir file in bucket
            config is Peer directory, name is config.yaml`,
	Run: func(cmd *cobra.Command, args []string) {
		Config := readConfig()
		rd, _ := ioutil.ReadDir(Config.DataDir)
		for _, fi := range rd {
			targetFilePath := Config.DataDir + "/" + fi.Name()
			c := NewClient(fi.Name(), "us-east-1", targetFilePath)
			c.checkBucket()
			c.UpLoadFile(c.targetFilePath)
		}
	},
}

func init() {
	rootCmd.AddCommand(upfileCmd)
}
