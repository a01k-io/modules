package nanoid

import gonanoid "github.com/matoous/go-nanoid"

func NewID(size int) string {
	id, err := gonanoid.Generate("0123456789qwertyuiopasdfghjklzxcvbnm", size)
	if err != nil {
		return ""
	}
	return id
}
