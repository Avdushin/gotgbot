package config

import (
	"os"
	"strconv"
)

func GetGroupID() int64 {
	groupIDStr := os.Getenv("GROUP_ID")
	groupID, err := strconv.ParseInt(groupIDStr, 10, 64)
	if err != nil {
		panic(err)
	}
	return groupID
}
