package tbcrypto

import (
	"crypto/md5"
	"crypto/sha1"
	"encoding/hex"
	"hash"
)

func GetMD5Hash(text string) string {
	return getHash(md5.New(), text)
}

func GetSHA1Hash(text string) string {
	return getHash(sha1.New(), text)
}

func getHash(h hash.Hash, text string) string {
	_, err := h.Write([]byte(text))
	if err != nil {
		return ""
	}
	return hex.EncodeToString(h.Sum(nil))
}
