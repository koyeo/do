package do

import (
	"fmt"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common"
	"github.com/tencentyun/vod-go-sdk"
	"strings"
)

func NewVod(engine *Engine) *Vod {
	return &Vod{engine: engine}
}

type Vod struct {
	process *Process
	engine  *Engine
}

func (p *Vod) fork(process *Process) *Vod {

	n := new(Vod)
	n.process = process
	n.engine = p.engine
	return n
}

func (p *Vod) UploadFile(client *VodClient, file *File) *Process {

	if p.process.isAbort {
		return p.process.pass()
	}

	t := strings.ToUpper(strings.TrimPrefix(file.MimeType, "video/"))
	req := vod.NewVodUploadRequest()
	req.MediaFilePath = common.StringPtr(file.Path)
	req.MediaName = common.StringPtr(file.Name)
	req.MediaType = common.StringPtr(t)
	req.StorageRegion = common.StringPtr(client.region)

	rsp, err := client.client.Upload(client.region, req)
	if err != nil {
		return p.process.Abort("【Vod】upload file error: ", err)
	}

	fmt.Printf("%+v\n", rsp.Response)
	fmt.Printf("%+v\n", rsp.CommitUploadResponse)
	fmt.Printf("%+v\n", rsp.BaseResponse)
	file.Uid = *rsp.Response.FileId
	file.Url = *rsp.Response.MediaUrl
	return p.process.next(fmt.Sprintf("upload %s success", file.Path), "", "")
}
