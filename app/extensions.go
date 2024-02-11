package extensions

import (
	"crypto/rand"
	"math/big"
)

func GenerateRandomString(length int) string {
    const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
    var result string
    for i := 0; i < length; i++ {
        index, err := rand.Int(rand.Reader, big.NewInt(int64(len(charset))))
        if err != nil {
            panic(err)
        }
        result += string(charset[index.Int64()])
    }
    return result
}