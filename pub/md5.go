package pub

import (
	"crypto/md5"
	"fmt"
)

func MD5String(data string) string {
	has := md5.Sum([]byte(data))
	return fmt.Sprintf("%x", has) //将[]byte转成16进制
}
