package pg

import (
	"time"

	"github.com/lib/pq"
)

type Blog struct {
	ID        uint
	Title     string
	Content   string
	Tags      pq.StringArray // string array for tags
	CreatedAt time.Time
}
