package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

type Version struct {
	Tag  string `json:"tag"`
	Date string `json:"date"`
}

type ProjectInfo struct {
	Name             string `json:"name"`
	FrameworkVersion string `json:"frameworkVersion"`
	FrameworkDate    string `json:"frameworkDate"`
}

const (
	frameworkRepo = "git@github.com:AgileExecutives/ae-framework.git"

	versionFile = "version.json"

	projectMetaDir  = ".ae"
	projectMetaFile = "project.json"
)

func main() {
	baseDir := flag.String("b", ".", "base directory")
	flag.Parse()

	args := flag.Args()

	if len(args) < 1 {
		usage()
		os.Exit(1)
	}

	command := args[0]
	project := ""
	if len(args) > 1 {
		project = args[1]
	}

	switch command {
	case "init":
		if project == "" {
			usage()
			os.Exit(1)
		}
		if err := initProject(*baseDir, project); err != nil {
			fmt.Println("ERROR:", err)
			os.Exit(1)
		}
	case "version":
		tmp, err := CloneFramework()
		if err != nil {
			fmt.Println("ERROR:", err)
			os.Exit(1)
		}
		defer os.RemoveAll(tmp)

		version, err := ReadVersion(tmp)
		if err != nil {
			fmt.Println("ERROR:", err)
			os.Exit(1)
		}

		fmt.Printf("%s (%s)\n", version.Tag, version.Date)

	case "update":
		if project == "" {
			usage()
			os.Exit(1)
		}
		projectDir := filepath.Join(*baseDir, project)
		if err := updateProject(projectDir); err != nil {
			fmt.Println("ERROR:", err)
			os.Exit(1)
		}

	default:
		usage()
	}
}

func usage() {
	fmt.Println(`
Usage:

	ae -b <basedir> <command> <project>

Commands:

	init     create new project
	version  show current framework version
	update   update an existing project to latest framework
`)
}

func initProject(baseDir, project string) error {

	target := filepath.Join(baseDir, project)

	if _, err := os.Stat(target); !os.IsNotExist(err) {
		return fmt.Errorf("project already exists: %s", target)
	}

	tmp, err := CloneFramework()
	if err != nil {
		return err
	}
	defer os.RemoveAll(tmp)

	fmt.Println("Creating project:", target)

	if err := os.MkdirAll(target, 0755); err != nil {
		return err
	}

	//
	// Copy application skeleton
	//

	skeleton := filepath.Join(
		tmp,
		"ae-cli",
		"app-skeleton",
	)

	fmt.Println("Copy skeleton")

	if err := copyDir(
		skeleton,
		target,
	); err != nil {
		return err
	}

	//
	// Copy frontend framework
	//

	srcFramework := filepath.Join(
		tmp,
		"core-frontend",
		"src",
		"framework",
	)

	appFramework := filepath.Join(
		target,
		"<name>-app",
		"src",
		"framework",
	)

	fmt.Println("Copy frontend framework")

	if err := copyDir(
		srcFramework,
		appFramework,
	); err != nil {
		return err
	}

	//
	// Replace placeholders
	//

	fmt.Println("Replacing placeholders")

	err = replaceRecursive(
		target,
		"<name>",
		project,
	)

	if err != nil {
		return err
	}

	//
	// Rename files/directories
	//

	fmt.Println("Renaming files")

	if err := renameRecursive(
		target,
		"<name>",
		project,
	); err != nil {
		return err
	}

	//
	// Initialize git
	//

	fmt.Println("Initialize git")

	if err := runInDir(
		target,
		"git",
		"init",
	); err != nil {
		return err
	}

	fmt.Println()
	fmt.Println("Project created successfully:")
	fmt.Println(target)

	// write project metadata with framework version
	version, err := ReadVersion(tmp)
	if err == nil {
		_ = WriteProjectInfo(target, &ProjectInfo{
			Name:             project,
			FrameworkVersion: version.Tag,
			FrameworkDate:    version.Date,
		})
	}

	return nil
}

func copyDir(src, dst string) error {

	return filepath.Walk(
		src,
		func(path string, info os.FileInfo, err error) error {

			if err != nil {
				return err
			}

			rel, err := filepath.Rel(src, path)
			if err != nil {
				return err
			}

			target := filepath.Join(dst, rel)

			if info.IsDir() {
				return os.MkdirAll(target, info.Mode())
			}

			return copyFile(path, target)
		},
	)
}

func copyFile(src, dst string) error {

	in, err := os.Open(src)

	if err != nil {
		return err
	}

	defer in.Close()

	if err := os.MkdirAll(
		filepath.Dir(dst),
		0755,
	); err != nil {
		return err
	}

	out, err := os.Create(dst)

	if err != nil {
		return err
	}

	defer out.Close()

	_, err = io.Copy(out, in)

	return err
}

func replaceRecursive(root, old, new string) error {

	return filepath.Walk(
		root,
		func(path string, info os.FileInfo, err error) error {

			if err != nil {
				return err
			}

			if info.IsDir() {
				return nil
			}

			data, err := os.ReadFile(path)

			if err != nil {
				return err
			}

			if !strings.Contains(
				string(data),
				old,
			) {
				return nil
			}

			data = []byte(
				strings.ReplaceAll(
					string(data),
					old,
					new,
				),
			)

			return os.WriteFile(
				path,
				data,
				info.Mode(),
			)
		},
	)
}

func renameRecursive(root, old, new string) error {

	var paths []string

	err := filepath.Walk(
		root,
		func(path string, info os.FileInfo, err error) error {

			if err == nil {
				paths = append(paths, path)
			}

			return err
		},
	)

	if err != nil {
		return err
	}

	for i := len(paths) - 1; i >= 0; i-- {

		oldPath := paths[i]

		if strings.Contains(
			filepath.Base(oldPath),
			old,
		) {

			newPath := filepath.Join(
				filepath.Dir(oldPath),
				strings.ReplaceAll(
					filepath.Base(oldPath),
					old,
					new,
				),
			)

			if oldPath != newPath {
				os.Rename(oldPath, newPath)
			}
		}
	}

	return nil
}

func run(cmd string, args ...string) error {

	c := exec.Command(cmd, args...)

	c.Stdout = os.Stdout
	c.Stderr = os.Stderr

	return c.Run()
}

func runInDir(dir, cmd string, args ...string) error {

	c := exec.Command(cmd, args...)

	c.Dir = dir
	c.Stdout = os.Stdout
	c.Stderr = os.Stderr

	return c.Run()
}

func CloneFramework() (string, error) {
	tmp, err := os.MkdirTemp("", "ae-framework-*")
	if err != nil {
		return "", err
	}

	err = run(
		"git",
		"clone",
		frameworkRepo,
		tmp,
	)

	if err != nil {
		os.RemoveAll(tmp)
		return "", err
	}

	return tmp, nil
}

func ReadVersion(repo string) (*Version, error) {
	file := filepath.Join(repo, versionFile)
	b, err := os.ReadFile(file)
	if err != nil {
		return nil, err
	}

	var v Version
	if err := json.Unmarshal(b, &v); err != nil {
		return nil, err
	}

	return &v, nil
}

func WriteProjectInfo(projectDir string, info *ProjectInfo) error {
	metaDir := filepath.Join(projectDir, projectMetaDir)
	if err := os.MkdirAll(metaDir, 0755); err != nil {
		return err
	}

	file := filepath.Join(metaDir, projectMetaFile)
	b, err := json.MarshalIndent(info, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(file, b, 0644)
}

func ReadProjectInfo(projectDir string) (*ProjectInfo, error) {
	file := filepath.Join(projectDir, projectMetaDir, projectMetaFile)
	b, err := os.ReadFile(file)
	if err != nil {
		return nil, err
	}

	var info ProjectInfo
	if err := json.Unmarshal(b, &info); err != nil {
		return nil, err
	}

	return &info, nil
}

func updateProject(projectDir string) error {
	info, err := ReadProjectInfo(projectDir)
	if err != nil {
		return err
	}

	tmp, err := CloneFramework()
	if err != nil {
		return err
	}
	defer os.RemoveAll(tmp)

	version, err := ReadVersion(tmp)
	if err != nil {
		return err
	}

	if info.FrameworkVersion == version.Tag {
		fmt.Println("Already up to date.")
		return nil
	}

	fmt.Printf("Updating %s -> %s\n", info.FrameworkVersion, version.Tag)

	frameworkSrc := filepath.Join(tmp, "core-frontend", "src", "framework")

	frameworkDst := filepath.Join(projectDir, info.Name+"-app", "src", "framework")

	os.RemoveAll(frameworkDst)

	if err := copyDir(frameworkSrc, frameworkDst); err != nil {
		return err
	}

	info.FrameworkVersion = version.Tag
	info.FrameworkDate = version.Date

	return WriteProjectInfo(projectDir, info)
}
