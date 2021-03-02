package igdl

import (
	"bufio"
	"os"
	"os/exec"
)

// Create directory if it does not exist
func CreateDirIfNotExist(dir string) (err error) {
	if _, err = os.Stat(dir); os.IsNotExist(err) {
		err = os.MkdirAll(dir, 0755)
	}
	return
}

// Call shell command wget to download. The reason to use wget is that wget
// supports automatically resume download. So this package only runs on Linux
// systems.
func Wget(url, filepath string) error {
	// run shell command `wget URL -O filepath`
	cmd := exec.Command("wget", url, "-O", filepath)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func mergePostliveVideoAndAudio(vpath, apath string) error {
	// run shell command `ffmpeg -i video.mp4 -i audio.mp4 -c:v copy -c:a copy merged.mp4`
	path := getPostLiveMergedFilePath(vpath, apath)
	cmd := exec.Command("ffmpeg", "-i", vpath, "-i", apath, "-c:v", "copy", "-c:a", "copy", path)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func IsCommandAvailable(name string) bool {
	cmd := exec.Command("/bin/sh", "-c", "command -v "+name)
	if err := cmd.Run(); err != nil {
		return false
	}
	return true
}

func FileToLines(filePath string) (lines []string, err error) {
	f, err := os.Open(filePath)
	if err != nil {
		return
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	err = scanner.Err()
	return
}
