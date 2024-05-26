package dockerCLI

import "strconv"

func GetCommandForDocker(fileName string, image string, memoryLimitKb int, timeLimitMs int, args ...string) []string {
	res := []string{
		"docker", "run",
		"--rm",
		"--mount", "type=bind,source=./tmp,target=/home/tmp",
		"-i",
		"-e", "FILE_NAME=" + fileName,
		image,
		"./unprivrun", strconv.Itoa(timeLimitMs), strconv.Itoa(memoryLimitKb),
	}
	for _, arg := range args {
		res = append(res, arg)
	}
	return res
}
