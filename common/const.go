package common

import "strconv"

const (
	KeyUserTokenPrefix = "ws:user:token:"
)

func GetUserTokenKey(userId int64, remoteIp string) string {
	return KeyUserTokenPrefix +
		strconv.FormatInt(userId, 10) + ":" + remoteIp
}
