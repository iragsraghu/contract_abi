package mapping

import (
	"golang.org/x/exp/slices"
)

func GetLockDurationExist(protocol_name string, user_action string) bool {
	lock_duration_files := []string{"pancake", "apeswap"}
	exists := slices.Contains(lock_duration_files, protocol_name)
	if exists && user_action == "stake" {
		return true
	} else {
		return false
	}

}
