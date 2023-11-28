package users

import (
	"crypto/md5"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"io"
)

func md5Encoder(password, salt string) string {
	MD5 := md5.New()
	io.WriteString(MD5, password+"$"+salt)
	return hex.EncodeToString(MD5.Sum(nil))
}

func getSalt() string {
	b := make([]byte, 16)
	_, err := io.ReadFull(rand.Reader, b)
	if err != nil {
		return ""
	}
	return fmt.Sprintf("%x", b)
}

func encryptPwd(password, salt string) string {
	return fmt.Sprintf(md5Encoder(password, salt))
}

func checkPwd(password, salt, savedPwd string) bool {
	return md5Encoder(password, salt) == savedPwd
}
