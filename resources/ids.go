package resources

import (
	"time"

	"go.rls.moe/nyx/resources/snowflakes"
)

var fountain = snowflakes.Generator{
	StartTime: time.Date(
		2017, 03, 11,
		11, 12, 29,
		0, time.UTC).Unix(),
}

func getID() (int, error) {
	id, err := fountain.NewID()
	return int(id), err
}

func DateFromId(id int) time.Time {
	return time.Unix(int64(fountain.IDToUnix(id)), 0).Local()
}
