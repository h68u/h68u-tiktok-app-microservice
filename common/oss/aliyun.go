package oss

import (
	"fmt"
	"github.com/aliyun/aliyun-oss-go-sdk/oss"
	"mime/multipart"
)

//TODO:重构成接口类型，每种存储介质只需实现对应的方法

type AliyunCfg struct {
	Endpoint        string
	AccessKeyID     string
	AccessKeySecret string
	BucketName      string
}

func AliyunInit(AliyunCfg AliyunCfg) *oss.Client {
	client, err := oss.New(AliyunCfg.Endpoint, AliyunCfg.AccessKeyID, AliyunCfg.AccessKeySecret)
	if err != nil {
		panic(err)
	}
	return client
}

func UploadVideoToOss(AliyunClient *oss.Client, bucketName string, objectName string, reader multipart.File) (bool, error) {
	bucket, err := AliyunClient.Bucket(bucketName)
	if err != nil {
		return false, err
	}
	err = bucket.PutObject(objectName, reader)
	if err != nil {
		fmt.Println(err)
		return false, err
	}
	return true, nil
}

func GetOssVideoUrlAndImgUrl(AliyunCfg AliyunCfg, objectName string) (string, string) {
	url := "https://" + AliyunCfg.BucketName + "." + AliyunCfg.Endpoint + "/" + objectName
	return url, url + "?x-oss-process=video/snapshot,t_0,f_jpg,w_0,h_0,m_fast,ar_auto"
}
