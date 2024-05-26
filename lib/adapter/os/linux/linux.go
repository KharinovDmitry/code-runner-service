package linux

import (
	"fmt"
	"os"
	"os/exec"
	"os/user"
	"strconv"
	"syscall"
	"time"
)

type LinuxAdapter struct {
}

func NewLinuxAdapter() LinuxAdapter { return LinuxAdapter{} }

func (l LinuxAdapter) CreateTempFileWithText(text string, extension string) (fileName string, err error) {
	fileName = strconv.FormatInt(time.Now().Unix(), 10) + extension
	tmpFile, err := os.Create(fileName)
	if err != nil {
		return "", fmt.Errorf("In LinuxAdapter(CreateTempFileWithText): %s", err.Error())
	}
	defer func() {
		tmpFile.Close()
		os.Remove(tmpFile.Name())
	}()

	if _, err = tmpFile.WriteString(text); err != nil {
		return "", fmt.Errorf("In utils(CreateFileWithText): %w", err)
	}
	tmpFile.Close()

	cmd := exec.Command("sudo", "-S", "mv", fileName, "/home/test_user/"+fileName)
	stdin, err := cmd.StdinPipe()
	if err != nil {
		return "", fmt.Errorf("In utils(CreateTempFileWithText): %w", err)
	}
	defer stdin.Close()
	fmt.Fprintln(stdin, "2004dxDX")

	if err = cmd.Run(); err != nil {
		return "", fmt.Errorf("In LinuxAdapter(CreateTempFileWithText): %s", err.Error())
	}

	return "/home/test_user/" + fileName, nil
}

func (l LinuxAdapter) AddFileExecutablePermission(fileName string) error {

	if err := exec.Command("echo", "qwerty", "|", "sudo", "-u", "test_user", "chmod", "+x", fileName).Run(); err != nil {
		return fmt.Errorf("In LinuxAdapter(AddFileExecutablePermission): %s", err.Error())
	}

	if err := exec.Command("echo", "qwerty", "|", "sudo", "-u", "test_user", "chown", "test_user", fileName).Run(); err != nil {
		return fmt.Errorf("In utils(AddFileExecutablePermission): %s", err.Error())
	}
	return nil
}

func (l LinuxAdapter) GetUnprivilegedProcAttr() (*syscall.SysProcAttr, error) {
	// Получаем информацию о непривилегированном пользователе
	u, err := user.Lookup("test_user")
	if err != nil {
		return nil, fmt.Errorf("In LinuxAdapter(GetUnprivilegedProcAttr): %w", err)
	}

	groupsStr, err := u.GroupIds()
	if err != nil {
		return nil, fmt.Errorf("In LinuxAdapter(GetUnprivilegedProcAttr): %w", err)
	}

	groups := make([]uint32, len(groupsStr))
	for i, group := range groupsStr {
		groupsId, _ := strconv.Atoi(group)
		groups[i] = uint32(groupsId)
	}

	// Преобразуем строковые идентификаторы в числа
	uid, _ := strconv.Atoi(u.Uid)
	gid, _ := strconv.Atoi(u.Gid)

	// Создаем атрибуты процесса с ограниченными привилегиями
	return &syscall.SysProcAttr{
		Chroot: "./testFiles",
		Credential: &syscall.Credential{
			Uid:    uint32(uid),
			Gid:    uint32(gid),
			Groups: groups,
		},

		//AmbientCaps: []uintptr{syscall.PR_SET_DUMPABLE},
		Setctty:     true,
		UseCgroupFD: false,
	}, nil

}
