// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.20.0

package timescale

import (
	"database/sql"
)

type File struct {
	CreatedAt sql.NullTime
	ID        string
	Name      string
	Mime      string
	Size      int64
	Buffer    []byte
	Hash      string
}
