package urlhasher

import (
	"strconv"

	"github.com/spaolacci/murmur3"
)

// HashLength hash length
const HashLength = 20

func GetHash(shortUrl string) string {
	m := murmur3.Sum64([]byte(shortUrl))
	return strconv.FormatUint(m, 10)
}
