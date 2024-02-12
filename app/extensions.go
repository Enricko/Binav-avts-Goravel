package extensions

import (
	"crypto/rand"
	"math/big"

	"github.com/goravel/framework/contracts/http"
	"github.com/goravel/framework/facades"
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

func PasswordHash(stringPassword string, ctx http.Context) (string, bool, http.Response) {
	if stringPassword == "" {
		return "", true, ctx.Response().Json(http.StatusInternalServerError, http.Json{
			"message": "Password is required.",
			"status":  http.StatusInternalServerError,
		})
	}
	password, err := facades.Hash().Make(stringPassword)
	if err != nil {
		return "", true, HandleInternalServerError(ctx, err)
	}
	return password, false, nil
}

// handleDBError handles database errors
func HandleBadRequestError(ctx http.Context, err error) http.Response {
	return ctx.Response().Json(http.StatusBadRequest, http.Json{
		"message": err.Error(),
		"status":  http.StatusBadRequest,
	})
}

// handleDBError handles database errors
func HandleInternalServerError(ctx http.Context, err error) http.Response {
	return ctx.Response().Json(http.StatusInternalServerError, http.Json{
		"message": err.Error(),
		"status":  http.StatusInternalServerError,
	})
}

// GenerateOTP generates a random one-time password of a specified length.
func GenerateOTP(length int) (string, error) {
	const charset = "0123456789"
	otp := make([]byte, length)
	max := big.NewInt(int64(len(charset)))

	for i := range otp {
		num, err := rand.Int(rand.Reader, max)
		if err != nil {
			return "", err
		}
		otp[i] = charset[num.Int64()]
	}

	return string(otp), nil
}

// VerifyOTP verifies if a given OTP matches the expected OTP.
func VerifyOTP(expectedOTP, inputOTP string) bool {
	return expectedOTP == inputOTP
}
