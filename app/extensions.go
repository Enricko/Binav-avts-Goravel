package extensions

import (
	"crypto/rand"
	"math/big"

	"github.com/goravel/framework/contracts/http"
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

// handleDBError handles database errors
func HandleBadRequestError(ctx http.Context, err error) http.Response {
	return ctx.Response().Json(http.StatusBadRequest, http.Json{
		"message": err.Error(),
		"status":http.StatusBadRequest,
	})
}
