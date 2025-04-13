package urlhasher

import (
	"strconv"

	"github.com/spaolacci/murmur3"
)

// HashLength hash length
const HashLength = 20

func GetHash(URL string) string {
	m := murmur3.Sum64([]byte(URL))
	return strconv.FormatUint(m, 10)
}
