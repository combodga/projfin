package user

import (
	"crypto/hmac"
	"crypto/md5"
	"crypto/sha256"
	"encoding/hex"
	"fmt"

	"github.com/combodga/projfin/internal/cookie"
	"github.com/labstack/echo/v4"
)

var key = "gophermart secret user key"

func SetAuthCookie(c echo.Context, user string) {
	sign, err := cookie.Read(c, "sign")
	if err == nil && sign == getSign(user) {
		return
	}

	cookie.Write(c, "user", user)
	cookie.Write(c, "sign", getSign(user))
}

func GetAuthCookie(c echo.Context) (string, error) {
	user, err := cookie.Read(c, "user")
	sign, err1 := cookie.Read(c, "sign")
	if err == nil && err1 == nil && sign == getSign(user) {
		return user, nil
	}

	return user, fmt.Errorf("failed authentification error")
}

func PasswordHasher(p string) string {
	hash := md5.Sum([]byte(p + "salt and pepper"))
	str := hex.EncodeToString(hash[:])
	return str
}

func ValidateOrderNumber(orderNumber int) bool {
	l := luhn(orderNumber / 10)
	return (orderNumber%10+l)%10 == 0
}

func luhn(n int) int {
	var l int

	for i := 0; n > 0; i++ {
		cur := n % 10
		if i%2 == 0 {
			cur = cur * 2
			if cur > 9 {
				cur = cur%10 + cur/10
			}
		}
		l += cur
		n = n / 10
	}

	return l % 10
}

func getSign(user string) string {
	h := hmac.New(sha256.New, []byte(key))
	h.Write([]byte(user))
	s := h.Sum(nil)
	return hex.EncodeToString(s)[:32]
}
