package goss

import (
    "gopkg.in/ini.v1"

    "fmt"
)

type IBrand interface {
    Name() string
    CredentialName() string
    Credentials() []File
    Buckets() []File
    UseCredential(string)
    UseBucket(string)
    Keys(string) []File
    Get(string, string)
    Post(string, string)
}

type Brand struct {
    name string
    conf *ini.File
    credentails []Credential
    credentail Credential
    credentailName string
}

func NewBrand(name string) (this *Brand) {

    this = &Brand{
        name: name,
        credentails: make([]Credential, 0),
        credentail: Credential{},
        credentailName: "",
    }
    return
}

func (this *Brand) Name() string {
    return this.name
}

func (this *Brand) Credentials() (files []File){
    files = make([]File, 0)
    Log.Info("Brand.conf", this.conf)
    ses := this.conf.Sections()
    for _, d := range ses {
        files = append(files, File{Name: d.Name()})
    }
    return
}
func (this *Brand) CredentialName() (name string){
    name = this.credentailName
    return
}

func (this *Brand) UseCredential(name string) {
    ses := this.conf.Sections()
    for _, d := range ses {
        if d.Name() == name {
            d.MapTo(&this.credentail)
            this.credentailName = name
        }
    }
    return
}

func HandleError(err error) {
    if err != nil {
        Log.Error(err)
        fmt.Println(err)
        // os.Exit(-1)
    }
}
type Credential struct {
    AccessKeyId string `ini:"access_key_id"`
    SecretAccessKey string `ini:"secret_access_key"`
    Endpoint string `ini:"endpoint"`
}
