package utils

import (
	"errors"
	"gopkg.in/ini.v1"
	"os"
	"path/filepath"
	"strings"
)

func FindGitRoot(start string) (string, error) {
	dir, err := filepath.Abs(start)
	if err != nil {
		return "", err
	}

	for {
		gitPath := filepath.Join(dir, ".git")

		if _, err := os.Stat(gitPath); err == nil {
			return dir, nil
		}

		parent := filepath.Dir(dir)
		if parent == dir {
			return "", errors.New("git repository not found")
		}

		dir = parent
	}
}

func GetGitRemote(folder string) (string, error) {
	if folder == "" {
		folder = "."
	}
	cfg, err := ini.Load(filepath.Join(folder, ".git", "config"))
	if err != nil {
		return "", err
	}

	for _, section := range cfg.Sections() {
		if len(section.Name()) > 8 && section.Name()[:8] == `remote "` {
			name := section.Name()[8 : len(section.Name())-1]
			if name == "origin" {
				return section.Key("url").String(), nil
			}
		}
	}
	return "", errors.New("Repo address not found in " + folder)
}

func SSHtoHTTPS(ssh string) (https string) {
	// git@github.com:user/repo.git -> https://github.com/user/repo.git
	if strings.HasPrefix(ssh, "git@") {
		parts := strings.SplitN(strings.TrimPrefix(ssh, "git@"), ":", 2)
		if len(parts) == 2 {
			host := parts[0]
			path := parts[1]
			https = "https://" + host + "/" + path
		}
	} else if strings.HasPrefix(ssh, "ssh://git@") {
		// ssh://git@github.com/user/repo.git -> https://github.com/user/repo.git
		https = "https://" + strings.TrimPrefix(ssh, "ssh://git@")
	} else if strings.HasPrefix(ssh, "https://") {
		// In case it is already https
		https = ssh
	}
	https = strings.TrimSuffix(https, ".git")
	return https
}

func GetGoPackage(protocol string) (string, error) {
	gitRoot, err := FindGitRoot(protocol)
	if err != nil {
		return "", err
	}
	remoteAddr, err := GetGitRemote(gitRoot)
	if err != nil {
		return "", err
	}
	remoteAddr = SSHtoHTTPS(remoteAddr)
	goPackage := remoteAddr + "/" + strings.TrimPrefix(protocol, gitRoot)
	return goPackage, nil
}
