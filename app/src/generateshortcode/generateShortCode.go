package generateshortcode

import (
	"crypto/md5"

	"github.com/kellegous/base62"
)

//GenerateShortCode Tạo short code từ URL
func GenerateShortCode(url string, length int) string {
	hash := md5.Sum([]byte(url))
	code := base62.EncodeToString(hash[:])
	return code[0:length]
}
