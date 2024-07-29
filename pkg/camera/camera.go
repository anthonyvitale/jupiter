// Package camera provides a way to take photos with a connected camera.
package camera

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"time"
)

// Camera is the interface that provides a way to take pictures via a connected camera.
type Camera interface {
	Take() (io.Reader, error)
}

// ShellCamera is an implementation of Camera that uses the command line.
type ShellCamera struct {
	dstDir string
}

func New(dstDir string) (*ShellCamera, error) {
	err := setupDst(dstDir)
	if err != nil {
		return nil, err
	}

	return &ShellCamera{
		dstDir: dstDir,
	}, nil
}

func (c *ShellCamera) Take() (string, error) {
	now := time.Now().UTC()
	filePath, err := c.nameFile(now)
	if err != nil {
		return "", err
	}

	// cmd := exec.Command("libcamera-still", "--shutter", "5000000", "--gain", "1", "--awbgains", "2.2,2.3", "--immediate", "-o", filePath)
	// cmd := exec.Command("fswebcam", "--no-banner", "-r", "1920x1080", "--jpeg", "100", "-D", "3", "-S", "13", filePath)
	cmd := exec.Command("libcamera-still", "--immediate", "-o", filePath)
	stderr, err := cmd.StderrPipe()
	if err != nil {
		return "", err
	}

	if err := cmd.Start(); err != nil {
		return "", err
	}

	slurp, _ := io.ReadAll(stderr)
	fmt.Printf("%s\n", slurp)

	if err := cmd.Wait(); err != nil {
		return "", err
	}

	return filePath, nil
}

func setupDst(dir string) error {
	err := os.MkdirAll(dir, os.ModePerm)
	return err
}

func (c *ShellCamera) nameFile(now time.Time) (string, error) {
	imgPath := now.Format("2006/01/02")
	dirPath := filepath.Join(c.dstDir, imgPath)

	err := os.MkdirAll(dirPath, os.ModePerm)
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("%s.jpeg", filepath.Join(dirPath, now.Format("150405Z"))), nil
}
