package password

import "github.com/alexedwards/argon2id"

func Hash(pwd string) (string, error) {
	return argon2id.CreateHash(pwd, argon2id.DefaultParams)
}

func Verify(pwd, hash string) (bool, error) {
	return argon2id.ComparePasswordAndHash(pwd, hash)
}
