package main

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"testing"

	aw "github.com/deanishe/awgo"
)

const (
	full_repo_path      = "/hoge/fuga/github.com/user/repo"
	full_repo_bucket    = "/hoge/fuga/bitbucket.org/user/repo"
	full_repo_other     = "/hoge/fuga/other.git/user/repo"
	user_repo           = "user/repo"
	github_user_repo    = "github.com/user/repo"
	bitbucket_user_repo = "bitbucket.org/user/repo"
)

func TestExcludeDomain(t *testing.T) {
	testcases := []struct {
		path     string
		expected string
		domain   bool
	}{
		{
			full_repo_path,
			user_repo,
			true,
		}, {
			full_repo_path,
			github_user_repo,
			false,
		}, {
			full_repo_bucket,
			user_repo,
			true,
		}, {
			full_repo_bucket,
			bitbucket_user_repo,
			false,
		},
	}
	for _, tc := range testcases {
		repoPath := strings.Split(tc.path, "/")
		actual := excludeDomain(repoPath, tc.domain)
		if actual != tc.expected {
			t.Errorf("%s is expected, but actual %s\n", tc.expected, actual)
		}
	}
}

func TestGetDomainName(t *testing.T) {
	testcases := []struct {
		path   string
		domain string
	}{
		{
			full_repo_path,
			"github.com",
		}, {
			full_repo_bucket,
			"bitbucket.org",
		}, {
			full_repo_other,
			"other.git",
		},
	}
	for _, tc := range testcases {
		repoPath := strings.Split(tc.path, "/")
		if actual := getDomainName(repoPath); actual != tc.domain {
			t.Errorf("%s is expected, but actual %s\n", tc.domain, actual)
		}
	}
}

func TestIsWorktree(t *testing.T) {
	tmp, err := ioutil.TempDir("", "go-ghq-alfred-worktree-test")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmp)

	primary := filepath.Join(tmp, "primary")
	if err := os.MkdirAll(filepath.Join(primary, ".git"), 0o755); err != nil {
		t.Fatal(err)
	}

	worktree := filepath.Join(tmp, "worktree")
	if err := os.MkdirAll(worktree, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := ioutil.WriteFile(
		filepath.Join(worktree, ".git"),
		[]byte("gitdir: "+filepath.Join(primary, ".git", "worktrees", "wt")+"\n"),
		0o644,
	); err != nil {
		t.Fatal(err)
	}

	plain := filepath.Join(tmp, "plain")
	if err := os.MkdirAll(plain, 0o755); err != nil {
		t.Fatal(err)
	}

	cases := []struct {
		name string
		path string
		want bool
	}{
		{"primary clone (.git dir)", primary, false},
		{"linked worktree (.git file)", worktree, true},
		{"path without .git", plain, false},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			if got := isWorktree(tc.path); got != tc.want {
				t.Errorf("isWorktree(%q) = %v, want %v", tc.path, got, tc.want)
			}
		})
	}
}

func TestGetIconPath(t *testing.T) {
	testcases := []struct {
		path string
		icon *aw.Icon
	}{
		{
			full_repo_path,
			gitHubIcon,
		}, {
			full_repo_bucket,
			bitBucketIcon,
		}, {
			full_repo_other,
			gitIcon,
		},
	}
	for _, tc := range testcases {
		repoPath := strings.Split(tc.path, "/")
		if icon := getIcon(repoPath); icon != tc.icon {
			t.Errorf("Expect: %s\nResult: %s\n", tc.icon, icon)
		}
	}
}
