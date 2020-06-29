package do

import (
	"context"
	"fmt"
	"github.com/tencentyun/cos-go-sdk-v5"
	"log"
	"net/http"
	netUrl "net/url"
	"os"
	"path/filepath"
	"strings"
)

func NewCos(engine *Engine) *Cos {
	return &Cos{engine: engine}
}

type Cos struct {
	engine  *Engine
	process *Process
}

func (p *Cos) fork(process *Process) *Cos {
	n := new(Cos)
	n.process = process
	n.engine = p.engine
	return n
}

func NewCosClient(url, secretId, secretKey, sessionKey string) *cos.Client {

	u, _ := netUrl.Parse(url)
	b := &cos.BaseURL{BucketURL: u}

	client := cos.NewClient(b, &http.Client{
		Transport: &cos.AuthorizationTransport{
			SecretID:     secretId,
			SecretKey:    secretKey,
			SessionToken: sessionKey,
		},
	})

	if client == nil {
		log.Fatal("Init cos client error")
	}

	return client
}

func (p *Cos) UploadFile(client *cos.Client, file *File) *Process {

	if p.process.isAbort {
		return p.process.pass()
	}

	f, err := os.Open(file.Path)
	if err != nil {
		return p.process.Abort("【Cos】upload file error: ", err)
	}

	opt := &cos.ObjectPutOptions{
		ObjectPutHeaderOptions: &cos.ObjectPutHeaderOptions{
			ContentType: file.MimeType,
		},
	}
	res, err := client.Object.Put(context.Background(), file.Name, f, opt)
	if err != nil {
		return p.process.Abort("【Cos】upload file error: ", err)
	}

	if res == nil || res.Body == nil {
		return p.process.Abort("【Cos】response content is nil")
	}

	defer func() {
		_ = res.Body.Close()
	}()

	//data, err := ioutil.ReadAll(res.Body)
	//if err != nil {
	//	return p.process.Abort("【Cos】read response content error: ", err)
	//}

	file.Url = strings.TrimPrefix(client.BaseURL.BucketURL.String(), client.BaseURL.BucketURL.Scheme+"://")
	file.Url = filepath.Join(file.Url, file.Name)
	return p.process.next(fmt.Sprintf("upload %s success", file.Path), "", "")
}
