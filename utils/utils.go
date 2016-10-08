package utils

import (
	"crypto/md5"
	"encoding/base64"
	"encoding/hex"
	"log"
	"net/url"
)

func Urlencode(originalUrl string) string {
	return url.QueryEscape(originalUrl)
}

func MD5(original string) string {
	cipherbyte := md5.Sum([]byte(original))
	return hex.EncodeToString(cipherbyte[:])
}

func Base64Encode(original string) string {
	return base64.StdEncoding.EncodeToString([]byte(original))
}

func Base64Decode(original string) string {
	cipherbyte, err := base64.StdEncoding.DecodeString(original)
	if err != nil {
		log.Println("base64 decode failed...")
		return ""
	}
	return string(cipherbyte)
}
