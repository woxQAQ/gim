package users

import (
	"crypto/md5"
	"crypto/rand"
	"encoding/hex"
)

func md5Encoder(password, salt string) string {
	data := []byte(password + "$" + salt)
	hash := md5.Sum(data)
	return hex.EncodeToString(hash[:])
}

func getSalt() (string, error) {
	b := make([]byte, 16)
	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}
	return hex.EncodeToString(b), nil
}

func encryptPwd(password, salt string) string {
	return md5Encoder(password, salt)
}

func checkPwd(password, salt, savedPwd string) bool {
	return md5Encoder(password, salt) == savedPwd
}
