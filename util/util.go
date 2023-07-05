package util

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"math/rand"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt"
)

func MD5(str string) string {
	h := md5.New()
	h.Write([]byte(str))
	return hex.EncodeToString(h.Sum(nil))
}

func GenerateSalt() string {
	rand.NewSource(time.Now().UnixNano())
	randNum := rand.Intn(100000000)

	salt := fmt.Sprintf("%d%d", time.Now().UnixNano(), randNum)
	salt = MD5(salt)

	return salt[:16]
}

func GenerateToken(id int64) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"userId": id,
		"exp":    time.Now().Add(time.Hour * 2).Unix(),
	})

	byteSlice, _ := hex.DecodeString("0x75001ea42539b2dd087f53eec22a714b4fc7cfdfd2c408914315b2ba20c05108a3b67bac62d5fc2ddf4db7f2094a6be50375e8d82abab650746ad4ddd1e1963c")
	tokenString, err := token.SignedString(byteSlice)
	if err != nil {
		return "", err
	}
	return tokenString, nil
}

func IsAscii(s string) bool {
	for _, c := range s {
		if !strconv.IsPrint(c) {
			return false
		}
	}
	return true
}

func IsHex(s string) bool {
	for _, c := range s {
		if !((c >= '0' && c <= '9') || (c >= 'a' && c <= 'f')) {
			return false
		}
	}
	return true
}
