package libraries

import (
	"math/rand"
	"time"
)

func WpGeneratePassword(length int, specialChars, extraSpecialChars bool) string {
	chars := []byte("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")
	if specialChars {
		chars = append(chars, []byte("!@#$%^&*()")...)
	}
	if extraSpecialChars {
		chars = append(chars, []byte("-_ []{}<>~`+=,.;:/?|")...)
	}
	result := []byte{}
	random := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < length; i++ {
		result = append(result, chars[random.Intn(len(chars))])
	}
	return string(result)
}
