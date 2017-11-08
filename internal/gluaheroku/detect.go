package gluaheroku

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"syscall"
)

func isNotFound(err error) bool {
	if ee, ok := err.(*exec.ExitError); ok {
		if ws, ok := ee.ProcessState.Sys().(syscall.WaitStatus); ok {
			return ws.ExitStatus() == 1
		}
	}
	return false
}

func gitRemotes() (map[string]string, error) {
	b, err := exec.Command("git", "remote", "-v").Output()
	if err != nil {
		return nil, err
	}

	return parseGitRemoteOutput(b)
}

func parseGitRemoteOutput(b []byte) (results map[string]string, err error) {
	s := bufio.NewScanner(bytes.NewBuffer(b))
	s.Split(bufio.ScanLines)

	results = make(map[string]string)

	for s.Scan() {
		by := s.Bytes()
		f := bytes.Fields(by)
		if len(f) != 3 || string(f[2]) != "(push)" {
			// this should have 3 tuples + be a push remote, skip it if not
			continue
		}

		if appName := appNameFromGitURL(string(f[1])); appName != "" {
			results[string(f[0])] = appName
		}
	}
	if err = s.Err(); err != nil {
		return nil, err
	}
	return
}

// AppFromGitRemote returns the heroku app name as detected from a git remote
// by name.
func AppFromGitRemote(remote string) (string, error) {
	if remote != "" {
		b, err := exec.Command("git", "config", "remote."+remote+".url").Output()
		if err != nil {
			if isNotFound(err) {
				wdir, _ := os.Getwd()
				return "", fmt.Errorf("could not find git remote "+remote+" in %s", wdir)
			}
			return "", err
		}

		out := strings.TrimSpace(string(b))

		appName := appNameFromGitURL(out)
		if appName == "" {
			return "", fmt.Errorf("could not find app name in " + remote + " git remote")
		}
		return appName, nil
	}

	// no remote specified, see if there is a single Heroku app remote
	remotes, err := gitRemotes()
	if err != nil {
		return "", nil // hide this error
	}
	if len(remotes) > 1 {
		return "", errors.New("multiple heroku remotes")
	}
	for _, v := range remotes {
		return v, nil
	}
	return "", fmt.Errorf("no apps in git remotes")
}

func gitHost() string {
	if herokuGitHost := os.Getenv("HEROKU_GIT_HOST"); herokuGitHost != "" {
		return herokuGitHost
	}
	if herokuHost := os.Getenv("HEROKU_HOST"); herokuHost != "" {
		return herokuHost
	}
	return "heroku.com"
}

func httpGitHost() string {
	if herokuHTTPGitHost := os.Getenv("HEROKU_HTTP_GIT_HOST"); herokuHTTPGitHost != "" {
		return herokuHTTPGitHost
	}
	return "git." + gitHost()
}

func sshGitURLPre() string {
	return "git@" + gitHost() + ":"
}

func httpGitURLPre() string {
	return "https://" + httpGitHost() + "/"
}

func appNameFromGitURL(remote string) string {
	if !strings.HasSuffix(remote, gitURLSuf) {
		return ""
	}

	if strings.HasPrefix(remote, sshGitURLPre()) {
		return remote[len(sshGitURLPre()) : len(remote)-len(gitURLSuf)]
	}

	if strings.HasPrefix(remote, httpGitURLPre()) {
		return remote[len(httpGitURLPre()) : len(remote)-len(gitURLSuf)]
	}

	return ""
}
