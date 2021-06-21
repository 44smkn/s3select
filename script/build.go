package main

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"
	"time"
)

var tasks = map[string]func(string) error{
	"bin/s3select": func(exe string) error {
		ldflags := os.Getenv("GO_LDFLAGS")
		ldflags = fmt.Sprintf("-X github.com/44smkn/s3select/pkg/build.Version=%s %s", version(), ldflags)
		ldflags = fmt.Sprintf("-X github.com/44smkn/s3select/pkg/build.Date=%s %s", date(), ldflags)
		return run("go", "build", "-trimpath", "-ldflags", ldflags, "-o", exe, "./cmd/s3select")
	},
	"clean": func(_ string) error {
		return rmrf("bin")
	},
}

func main() {
	task := os.Args[1]
	t, ok := tasks[task]
	if !ok {
		fmt.Fprintf(os.Stderr, "Don't know how to build task `%s`.\n", task)
		os.Exit(1)
	}

	err := t(task)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		fmt.Fprintf(os.Stderr, "building task `%s` failed.\n", task)
		os.Exit(1)
	}
}

func version() string {
	if versionEnv := os.Getenv("S3SELECT_VERSION"); versionEnv != "" {
		return versionEnv
	}
	if desc, err := cmdOutput("git", "describe", "--tags"); err == nil {
		return desc
	}
	rev, _ := cmdOutput("git", "rev-parse", "--short", "HEAD")
	return rev
}

func date() string {
	return time.Now().Format("2006-01-02")
}

func cmdOutput(args ...string) (string, error) {
	path, err := exec.LookPath(args[0])
	if err != nil {
		return "", err
	}
	cmd := exec.Command(path, args[1:]...)
	cmd.Stderr = io.Discard
	out, err := cmd.Output()
	return strings.TrimSuffix(string(out), "\n"), err
}

func rmrf(targets ...string) error {
	args := append([]string{"rm", "-rf"}, targets...)
	fmt.Printf("%v\n", args)
	for _, target := range targets {
		if err := os.RemoveAll(target); err != nil {
			return err
		}
	}
	return nil
}

func run(args ...string) error {
	path, err := exec.LookPath(args[0])
	if err != nil {
		return err
	}
	cmd := exec.Command(path, args[1:]...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}
