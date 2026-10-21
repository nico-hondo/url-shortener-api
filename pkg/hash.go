package pkg

import (
	"crypto/rand"
	"crypto/subtle"
	"encoding/base64"
	"errors"
	"fmt"
	"strings"

	"golang.org/x/crypto/argon2"
)

type HashConfig struct {
	memory  uint32
	time    uint32
	threads uint8
	keyLen  uint32
	saltLen uint32
}

func NewHashConfig(memory, time, keyLen, saltLen uint32, threads uint8) *HashConfig {
	return &HashConfig{
		memory:  memory,
		time:    time,
		keyLen:  keyLen,
		saltLen: saltLen,
		threads: threads,
	}
}

// OWASP =, mei 2023
func NewRecommendedHashConfig() *HashConfig {
	return &HashConfig{
		memory:  64 * 1024,
		time:    2,
		threads: 2,
		keyLen:  32,
		saltLen: 16,
	}
}

// Kita buat dulu salt nya dengan mengguankan crypto rand
func (hc *HashConfig) genSalt() []byte {
	salt := make([]byte, hc.saltLen)

	rand.Read(salt)

	return salt
}

// Generate Hash
func (hc *HashConfig) GenHash(password string) string {
	salt := hc.genSalt()                                                                    //Bikin salt random agar 2 user dengan password yang sama bisa berbeda hasil hashingnya
	hash := argon2.IDKey([]byte(password), salt, hc.time, hc.memory, hc.threads, hc.keyLen) //Inti dari teknik hashing yang dimana nantinya disini digabungkan mulai dari password + salt + parameter menjadi hash binner

	base64Salt := base64.RawStdEncoding.EncodeToString(salt)
	base64Hash := base64.RawStdEncoding.EncodeToString(hash)

	// $alg$version$parameters$salt$hash
	completeHash := fmt.Sprintf("$argon2id$v=%d$m=%d,t=%d,p=%d$%s$%s", argon2.Version, hc.memory, hc.time, hc.threads, base64Salt, base64Hash)

	return completeHash
}

// Kita bandingkan lalu ambil error jika error dan kembalikan nil jika tidak error
func Compare(password string, hashedPwd string) error {
	// $alg$version$parameters$salt$hash
	result := strings.Split(hashedPwd, "$")

	if len(result) != 6 {
		return ErrInvalidHash
	}

	//cek apakah sama algoritmanya
	if result[1] != "argon2id" {
		return ErrIncompatibleVariant
	}

	//cek apakah sama versionnya
	var version int
	if _, err := fmt.Sscanf(result[2], "v=%d", &version); err != nil {
		return err
	}

	if version != argon2.Version {
		return ErrIncompatibleVersion
	}

	var memory, time uint32
	var threads uint8

	if _, err := fmt.Sscanf(result[3], "m=%d,t=%d,p=%d", &memory, &time, &threads); err != nil {
		return err
	}

	//kita ambil saltnya
	salt, err := base64.RawStdEncoding.DecodeString(result[4])
	if err != nil {
		return err
	}

	hash, err := base64.RawStdEncoding.DecodeString(result[5])
	if err != nil {
		return err
	}

	//kita generate hash dari password
	newHash := argon2.IDKey([]byte(password), salt, time, memory, threads, uint32(len(hash)))

	if subtle.ConstantTimeCompare(hash, newHash) == 0 {
		return ErrMismatchedHashAndPassword
	}

	return nil
}

var (
	ErrInvalidHash               = errors.New("argon2id: hash is not in the correct format")
	ErrIncompatibleVariant       = errors.New("argon2id: incompatible variant of argon2")
	ErrIncompatibleVersion       = errors.New("argon2id: incompatible version of argon2")
	ErrMismatchedHashAndPassword = errors.New("argon2id: the provided password does not match the hash")
)
