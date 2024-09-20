package utils

import (
	"crypto/sha256"
	"math/big"
	"strconv"
	"strings"
)

const base62Chars = "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

func Base62Encode(num *big.Int) string {
	if num.Sign() == 0 {
		return string(base62Chars[0])
	}
	var result strings.Builder
	base := big.NewInt(62)
	zero := big.NewInt(0)
	mod := new(big.Int)
	for num.Cmp(zero) > 0 {
		num.DivMod(num, base, mod)
		result.WriteByte(base62Chars[mod.Int64()])
	}
	return result.String()
}

func HashURL(url string, increment int) string {
	urlWithSalt := url + strconv.Itoa(increment)

	hash := sha256.Sum256([]byte(urlWithSalt))

	num := new(big.Int).SetBytes(hash[:])

	shortURL := Base62Encode(num)

	return shortURL[:8]
}
