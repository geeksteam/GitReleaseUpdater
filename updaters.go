package updater

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"syscall"

	"github.com/kardianos/osext"
)

func UMUntarGz(destDir string) func(string) error {
	return func(srcPath string) error {
		out, err := exec.Command("tar", "-C"+destDir, "-x", "-z", "-f"+srcPath).CombinedOutput()
		if err != nil {
			return fmt.Errorf("Failed to unpack %s to %s:\n\t%s", srcPath, destDir, out)
		}

		return nil
	}
}

func UMReplaceBin() func(string) error {
	return func(srcPath string) error {
		src, err := os.Open(srcPath)
		if err != nil {
			return fmt.Errorf("Failed to open source file '%s': %s", srcPath, err)
		}
		defer src.Close()

		execPath, err := osext.Executable()
		if err != nil {
			return fmt.Errorf("Failed to get executable path: %s", err)
		}

		mod, err := getFileMode(execPath)
		if err != nil {
			return fmt.Errorf("Failed to get file mode of executable: %s", err)
		}

		if err := syscall.Unlink(execPath); err != nil {
			return fmt.Errorf("Can't Unlink file: %s", err)
		}

		f, err := os.Create(execPath)
		if err != nil {
			return fmt.Errorf("Failed to open executable: %s", err)
		}
		defer f.Close()

		if _, err := io.Copy(f, src); err != nil && err != io.EOF {
			f.Close()
			return fmt.Errorf("Failed to rewtite ececutable: %s", err)
		}

		os.Chmod(execPath, mod)

		return nil
	}
}

func getFileMode(path string) (os.FileMode, error) {
	info, err := os.Stat(path)
	if err != nil {
		return 0, err
	}

	return info.Mode(), nil
}
