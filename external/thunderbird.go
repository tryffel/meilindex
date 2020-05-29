package external

import (
	"github.com/sirupsen/logrus"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
)

// OpenById opens given mail id with thunderbird email client.
func OpenById(id string) error {
	cmd := exec.Command("thunderbird", "-thunderlink",
		"thunderlink://messageid="+id)
	return cmd.Run()
}

// files that are not mbox files
var skipFile = regexp.MustCompile(`.+\.(dat)|(html)|(msf)|(sbd)`)

// MboxFiles lists all mbox files found in given directory
func MboxFiles(baseDir string, recursive bool) ([]MboxFile, error) {
	if strings.HasSuffix(baseDir, "/") {
		baseDir = strings.TrimSuffix(baseDir, "/")
	}
	return mboxFiles(baseDir, "")
}

func mboxFiles(baseDir string, parentName string) ([]MboxFile, error) {
	// Directory / Folder ends with .sbd.
	// MBox file has no ending. Skip '.msf'
	var files []MboxFile

	dirFiles, err := filepath.Glob(baseDir + "/*")
	if err != nil {
		return nil, err
	}
	for _, v := range dirFiles {
		if strings.HasSuffix(v, ".sbd") {
			base := filepath.Base(v)
			base = strings.TrimSuffix(base, ".sbd")
			parent := base
			if parentName != "" {
				parent = parentName + "/" + base
			}
			found, err := mboxFiles(v, parent)
			if err != nil {
				logrus.Error(err)
			} else {
				files = append(files, found...)
			}

		} else {
			base := filepath.Base(v)
			if !skipFile.MatchString(base) {
				file := MboxFile{
					File: v,
				}
				if parentName == "" {
					file.Name = base
				} else {
					file.Name = parentName + "/" + base
				}
				files = append(files, file)
			}
		}
	}
	return files, nil
}
