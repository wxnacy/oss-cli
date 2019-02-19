package goss

import (

    "testing"
    "fmt"
)


func TestParsePWD(t *testing.T) {

    var o *Terminal
    var pwd string

    pwd = fmt.Sprintf("/%s", BRANDS_NAME)
    o = NewTerminal(pwd)
    o.parsePWD()
    if o.brandName != "" {
        t.Errorf("brandName error %s", o.brandName)
    }

    pwd = fmt.Sprintf("%s/%s", pwd, "oss")
    o = NewTerminal(pwd)
    o.parsePWD()
    if o.brandName != "oss" {
        t.Errorf("brandName error %s", o.brandName)
    }
    if o.credentailName != "" {
        t.Errorf("credentailName error %s", o.credentailName)
    }

    pwd = fmt.Sprintf("%s/%s/%s", pwd, CREDENTIALS_NAME, "wxnacy")
    o = NewTerminal(pwd)
    o.parsePWD()
    if o.brandName != "oss" {
        t.Errorf("brandName error %s", o.brandName)
    }
    if o.credentailName != "wxnacy" {
        t.Errorf("credentailName error %s", o.credentailName)
    }
    if o.bucketName != "" {
        t.Errorf("bucketName error %s", o.bucketName)
    }

    pwd = fmt.Sprintf("%s/%s/%s", pwd, BUCKETS_NAME, "wxnacy-img")
    o = NewTerminal(pwd)
    o.parsePWD()
    if o.brandName != "oss" {
        t.Errorf("brandName error %s", o.brandName)
    }
    if o.credentailName != "wxnacy" {
        t.Errorf("credentailName error %s", o.credentailName)
    }
    if o.bucketName != "wxnacy-img" {
        t.Errorf("bucketName error %s", o.bucketName)
    }
    if o.keyPrefix != "" {
        t.Errorf("keyPrefix error %s", o.keyPrefix)
    }

    pwd = fmt.Sprintf("%s/%s/%s", pwd, KEYS_NAME, "api/test")
    o = NewTerminal(pwd)
    o.parsePWD()
    if o.brandName != "oss" {
        t.Errorf("brandName error %s", o.brandName)
    }
    if o.credentailName != "wxnacy" {
        t.Errorf("credentailName error %s", o.credentailName)
    }
    if o.bucketName != "wxnacy-img" {
        t.Errorf("bucketName error %s", o.bucketName)
    }
    if o.keyPrefix != "api/test" {
        t.Errorf("keyPrefix error %s", o.keyPrefix)
    }
}

func TestCd(t *testing.T) {


    var o *Terminal
    var pwd string

    pwd = "/"
    o = NewTerminal(pwd)
    o.Cd("Brands")
    if o.PWD() != "/Brands" {
        t.Error("cd error", o.PWD())
    }

    pwd = "/Brands"
    o = NewTerminal(pwd)
    o.Cd("oss")
    if o.PWD() != "/Brands/oss" {
        t.Error("cd error", o.PWD())
    }

    pwd = "/Brands/oss"
    o = NewTerminal(pwd)
    o.Cd("/Brands")
    if o.PWD() != "/Brands" {
        t.Error("cd error", o.PWD())
    }
}

