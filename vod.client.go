package do

import "github.com/tencentyun/vod-go-sdk"

type VodClient struct {
	name   string
	region string
	client *vod.VodUploadClient
}

func NewVodClient(region, name, secretId, secretKey, sessionToken string) *VodClient {

	client := &VodClient{
		name:   name,
		region: region,
	}
	client.client = &vod.VodUploadClient{}
	client.client.SecretId = secretId
	client.client.SecretKey = secretKey
	client.client.Token = sessionToken
	return client
}
