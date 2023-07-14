package models

import (
	"github.com/Yarik-xxx/CodeWarsRestApi/pkg/codewars"
	"time"
)

type Challenge struct {
	ID         string
	Info       codewars.Kata
	LastUpdate time.Time
}
