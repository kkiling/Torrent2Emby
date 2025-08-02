package delivery

import (
	"fmt"
	"os/user"
	"strconv"
	"syscall"
)

func setGroup(path string, groupName string) error {
	group, err := user.LookupGroup(groupName)
	if err != nil {
		return fmt.Errorf("user.LookupGroup: %w", err)
	}
	gid, _ := strconv.Atoi(group.Gid)
	err = syscall.Chown(path, -1, gid)
	if err != nil {
		return fmt.Errorf("syscall.Chown: %w", err)
	}
	return nil
}
