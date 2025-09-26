package utils

import (
	"os/exec"
	"strconv"
	"strings"
)

func parseQuotaOutput(output string) (int64, int64, error) {
	lines := strings.Split(strings.TrimSpace(output), "\n")
	fields := strings.Fields(lines[len(lines)-1])

	usedKb, err := strconv.ParseInt(fields[1], 10, 64)
	if err != nil {
		return -1, -1, err
	}

	totalKb, err := strconv.ParseInt(fields[2], 10, 64)
	if err != nil {
		return -1, -1, err
	}

	return totalKb * 1024, usedKb * 1024, nil
}

func GetDiskQuota() (int64, int64, error) {
	cmd := exec.Command("quota")
	output, err := cmd.Output()
	if err != nil {
		return -1, -1, err
	}

	return parseQuotaOutput(string(output))
}
