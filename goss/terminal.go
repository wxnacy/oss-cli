package goss

import (

	"github.com/c-bata/go-prompt"
    "strings"
    "fmt"
)

type DirLevel uint8


const (
    BRANDS_NAME = "Brands"
    CREDENTIALS_NAME = "Credentials"
    BUCKETS_NAME = "Buckets"
    KEYS_NAME = "Keys"

    DirLevelRoot DirLevel = iota
    DirLevelBrands
    DirLevelCredentials
    DirLevelBuckets
    DirLevelKeys
)

type File struct {
    Name string
    IsFile bool
    Size int64
}


type Terminal struct {
    pwd string
    brand IBrand
    brandName string
    credentailName string
    bucketName string
    keyPrefix string
    dirLevel DirLevel
}

func NewTerminal(pwd string) (this *Terminal) {
    this = &Terminal{
        pwd: pwd,
        dirLevel: DirLevelRoot,
    }
    Log.Info(this.dirLevel)
    this.parsePWD()
    return
}

func (this *Terminal) PWD() string {
    return this.pwd
}
func (this *Terminal) BrandName() string {
    return this.brandName
}
func (this *Terminal) CredentialName() string {
    return this.credentailName
}
func (this *Terminal) BucketName() string {
    return this.bucketName
}

func (this *Terminal) Get(name string, localfile string) {
    key := name
    if localfile == "" {
        localfile = name
    }
    if strings.HasSuffix(localfile, "/") {
        localfile = localfile + name
    }

    this.Brand().Get(key, localfile)
}

func (this *Terminal) Post(name string, localfile string) {
    key := name
    if key == "" {
        key = localfile
    }
    Log.Info(key, localfile)

    this.Brand().Post(key, localfile)
}

func (this *Terminal) Cd(dir string) {
    this.cd(dir)
    this.parsePWD()
}
func (this *Terminal) cd(dir string) {
    if dir == ".." {
        pwds := strings.Split(this.pwd, "/")
        this.pwd = strings.Join(pwds[0:len(pwds) - 1], "/")
        return
    }
    if dir == "..." {
        pwds := strings.Split(this.pwd, "/")
        this.pwd = strings.Join(pwds[0:len(pwds) - 2], "/")
        return
    }
    dir = strings.TrimRight(dir, "/")
    ls := this.LS()
    if !strings.HasPrefix(dir, "/") {
        if index := InArray(dir, ls); index == -1 {
            fmt.Printf("cd: no such file or directory: %s\n", dir)
            return
        }
    }
    if this.pwd == "/" {
        this.pwd = "/" + dir
        return
    }
    if strings.HasPrefix(dir, "/") {
        this.pwd = dir
        return
    }

    this.pwd = fmt.Sprintf("%s/%s", this.pwd, dir)

}

func (this *Terminal) ChangePWD(pwd string) {
    this.pwd = pwd
    this.parsePWD()
}

func (this *Terminal) clear() {
    this.brandName = ""
    this.credentailName = ""
    this.bucketName = ""
    this.keyPrefix = ""
    this.dirLevel = DirLevelRoot
}

func (this *Terminal) parsePWD() {
    this.clear()
    pwd := this.pwd
    Log.Info(strings.Split(pwd, "/"), this.dirLevel)
    pwds := strings.Split(pwd, "/")
    for i, d := range pwds {
        Log.Info(i, d)
    }
    Log.Info("pwds length", len(pwds))
    if len(pwds) >=2 && pwds[1] == BRANDS_NAME {
        this.dirLevel = DirLevelBrands
    }
    if len(pwds) >= 3 && pwds[1] == BRANDS_NAME {
        this.brandName = pwds[2]
    }

    if len(pwds) >=4 && pwds[3] == CREDENTIALS_NAME {
        this.dirLevel = DirLevelCredentials
    }
    if len(pwds) >= 5 && pwds[3] == CREDENTIALS_NAME {
        this.credentailName = pwds[4]
    }

    if len(pwds) >=6 && pwds[5] == BUCKETS_NAME {
        this.dirLevel = DirLevelBuckets
    }
    if len(pwds) >= 7 && pwds[5] == BUCKETS_NAME {
        this.bucketName = pwds[6]
    }

    if len(pwds) >=8 && pwds[7] == KEYS_NAME {
        this.dirLevel = DirLevelKeys
    }
    if len(pwds) >= 9 && pwds[7] == KEYS_NAME {
        keys := pwds[8:len(pwds)]
        this.keyPrefix = strings.Join(keys, "/")
    }

    Log.Info(this.dirLevel)

    this.resetBrand()

}

func (this *Terminal) resetBrand() {
    var b IBrand
    switch this.brandName {
        case "oss": {
            b = NewOSS()
        }
    }
    if this.credentailName != "" {
        b.UseCredential(this.credentailName)
    }
    if this.bucketName != "" {
        b.UseBucket(this.bucketName)
    }
    this.brand = b
    return
}

func (this *Terminal) Brand() (b IBrand) {
    b = this.brand
    return
}

func (this *Terminal) PrintLL() {
    files := this.LL()
    for i, d := range files {
        dir := "dir"
        if d.IsFile {
            dir = "file"
        }
        fmt.Printf("%d\t%s\t%s\n", i, dir, d.Name)
    }
    fmt.Println("Totals: ", len(files))
}

func (this *Terminal) PrintLS() {
    for _, d := range this.LS() {
        fmt.Print(d + " ")
    }
    fmt.Println("")
}

func (this *Terminal) LS() (names []string){
    names = make([]string, 0)
    for _, d := range this.LL() {
        names = append(names, d.Name)
    }
    return

}

func (this *Terminal) LL() (files []File){
    files = make([]File, 0)
    switch this.dirLevel {
        case DirLevelRoot: {
            files = []File{
                {Name: "Brands"},
            }
        }
        case DirLevelBrands: {
            if this.brandName == "" {
                files = []File{
                    {Name: "oss"},
                    {Name: "s3"},
                    {Name: "qiniu"},
                }
            } else {
                files = []File{
                    {Name: CREDENTIALS_NAME},
                }
            }
        }
        case DirLevelCredentials: {
            if this.credentailName == "" {
                files = this.Brand().Credentials()
            } else {
                files = []File{
                    {Name: BUCKETS_NAME},
                }
            }
        }
        case DirLevelBuckets: {
            if this.bucketName == "" {
                files = this.Brand().Buckets()
            } else {
                files = []File{
                    {Name: KEYS_NAME},
                }
            }
        }
        case DirLevelKeys: {
            if this.keyPrefix == "" {
                files = this.Brand().Keys("")
            } else {
                files = this.Brand().Keys(this.keyPrefix + "/")
            }
        }
    }
    return
}


func (this *Terminal) Prompt(cmd string) (prompts []prompt.Suggest) {
    prompts = make([]prompt.Suggest, 0)
    if cmd == "" {
        return
    }

    ll := this.LL()
    dirs := make([]string, 0)
    for _, d := range ll {
        if !d.IsFile {
            dirs = append(dirs, d.Name + "/")
        }
    }
    files := make([]string, 0)
    for _, d := range ll {
        if d.IsFile {
            files = append(files, d.Name)
        }
    }

    if cmd == "cd " {
        for _, d := range dirs {
            prompts = append(prompts, prompt.Suggest{Text: d})
        }
    }

    if cmd == "get " {

        for _, d := range files {
            prompts = append(prompts, prompt.Suggest{Text: d})
        }
    }

    return
}
