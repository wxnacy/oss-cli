package goss

import (
    "github.com/aliyun/aliyun-oss-go-sdk/oss"
    "gopkg.in/ini.v1"
    "fmt"
    "strings"
)

var OSS_HOME = HOME + "/brands/qiniu"
var CREDENTIALS = OSS_HOME + "/credentials"
var err error

type QNClient struct {
    *Brand
    client *oss.Client
    bucket *oss.Bucket
}

func NewOSS() (this *OSSClient) {
    this = &OSSClient{
        Brand: NewBrand("qiniu"),
    }
    if !Exists(CREDENTIALS) {
        SaveFile(CREDENTIALS, "")
    }

    this.conf, err = ini.Load(CREDENTIALS)
    this.Brand.conf = this.conf

    HandleError(err)
    return
}


func (this *OSSClient) Buckets() (files []File){
    // 列举存储空间。
    files = make([]File, 0)
    marker := ""
    for {
        lsRes, err := this.client.ListBuckets(oss.Marker(marker))
        if err != nil {
            fmt.Println("Error:", err)
            Log.Error("Error: ", err)
            // os.Exit(-1)
        }
        // 默认情况下一次返回100条记录。
        for _, bucket := range lsRes.Buckets {
            // fmt.Println("Bucket: ", bucket.Name)
            files = append(files, File{Name: bucket.Name})
        }
        if lsRes.IsTruncated {
            marker = lsRes.NextMarker
        } else {
            break
        }
    }
    return
}

func (this *OSSClient) UseCredential(name string) {
    this.Brand.UseCredential(name)
    this.client, err = oss.New(
        fmt.Sprintf("https://%s", this.credentail.Endpoint),
        this.credentail.AccessKeyId,
        this.credentail.SecretAccessKey,
    )
    HandleError(err)
}

func (this *OSSClient) UseBucket(name string) {
	this.bucket, err = this.client.Bucket(name)
    HandleError(err)
}

func (this *OSSClient) Keys(prefix string) (objs []File){
    objs = make([]File, 0)
	marker := ""
    // if !strings.HasSuffix(prefix, "/") && prefix != "" {
        // prefix = prefix + "/"
    // }
    Log.Info("prefix", prefix)

    pref := oss.Prefix(prefix)
	for {
		lsRes, err := this.bucket.ListObjects(pref, oss.Marker(marker))
        HandleError(err)
		// 打印列举文件，默认情况下一次返回100条记录。
		for _, o := range lsRes.Objects {
			// fmt.Println("Bucket:", o.Key, o.Size)
            name := o.Key
            if name == prefix {
                continue
            }
            IsFile := true
            if strings.HasSuffix(name, "/") {
                IsFile = false
            }
            name = strings.TrimRight(name, "/")
            objs = append(objs, File{Name: name, IsFile: IsFile, Size: o.Size})
		}
		if lsRes.IsTruncated {
			marker = lsRes.NextMarker
		} else {
			break
		}
	}
    return
}
// 定义进度条监听器。
type OssProgressListener struct {
}

// 定义进度变更事件处理函数。
func (listener *OssProgressListener) ProgressChanged(event *oss.ProgressEvent) {
	switch event.EventType {
	case oss.TransferStartedEvent:
		fmt.Printf("Transfer Started, ConsumedBytes: %d, TotalBytes %d.\n",
			event.ConsumedBytes, event.TotalBytes)
	case oss.TransferDataEvent:
        cons := event.ConsumedBytes*100/event.TotalBytes
        if cons == 100 {

            fmt.Printf(
                "\r下载中, 完成: %d, 总数据: %d, 进度: %d%%.\n",
                event.ConsumedBytes, event.TotalBytes, cons,
            )
            fmt.Println("下载完成。")
        } else {

            fmt.Printf(
                "\rTransfer Data, ConsumedBytes: %d, TotalBytes %d, %d%%.",
                event.ConsumedBytes, event.TotalBytes, cons,
            )
        }
	case oss.TransferCompletedEvent:
		fmt.Printf("\nTransfer Completed, ConsumedBytes: %d, TotalBytes %d.\n",
			event.ConsumedBytes, event.TotalBytes)
	case oss.TransferFailedEvent:
		fmt.Printf("\nTransfer Failed, ConsumedBytes: %d, TotalBytes %d.\n",
			event.ConsumedBytes, event.TotalBytes)
	default:
	}
}

func (this *OSSClient) Get(key string, localfile string) {
    err = this.bucket.GetObjectToFile(key, localfile, oss.Progress(&OssProgressListener{}))
	if err != nil {
		fmt.Println("Error:", err)
        Log.Error("Error:", err)
	}
}

func (this *OSSClient) Post(key string, localfile string) {
    err = this.bucket.PutObjectFromFile(key, localfile, oss.Progress(&OssProgressListener{}))
	if err != nil {
		fmt.Println("Error:", err)
        Log.Error("Error:", err)
	}
}
