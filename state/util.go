package state

import (
	"hash/fnv"
	"os"
)

// HashPassword - properly gets salt and hashes password
func HashPassword(p string) uint32 {
	salt := os.Getenv("PASSWORD_SALT")
	if salt == "" {
		panic("env variable PASSWORD_SALT not set")
	}
	p = p + salt
	h := fnv.New32a()
	h.Write([]byte(p))
	return h.Sum32()
}
