package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	git "github.com/go-git/go-git/v5"
)

var (
	_base = "."
)

func main() {
	infos, err := ioutil.ReadDir(_base)
	if err != nil {
		panic(err)
	}

	for _, info := range infos {
		if !info.IsDir() {
			continue
		}

		src, err := filepath.Abs(filepath.Join(_base, info.Name()))
		if err != nil {
			panic(err)
		}

		r, err := git.PlainOpen(src)
		if err != nil {
			fmt.Println("not git:", src)
			continue
		}

		origin, err := r.Remote("origin")
		if err != nil {
			continue
		}

		if len(origin.Config().URLs) != 1 {
			fmt.Println("unexpected origin urls:", origin.Config().URLs)
			continue
		}

		url := origin.Config().URLs[0]
		url = rmPrefix(url, "https://")
		url = rmPrefix(url, "git://")
		url = rmPrefix(url, "git@")
		url = rmSuffix(url, ".git")
		url = strings.Replace(url, "github.com:", "github.com/", 1)

		dst, err := filepath.Abs(filepath.Join(_base, url))
		if err != nil {
			panic(err)
		}

		_, err = os.Stat(dst)
		if err != nil {
			if os.IsNotExist(err) {
				if yn(fmt.Sprintf("move %s to %s?", src, dst)) {
					if err := os.MkdirAll(filepath.Dir(dst), os.ModePerm); err != nil {
						panic(err)
					}
					if err := os.Rename(src, dst); err != nil {
						panic(err)
					}
				}
				continue
			} else {
				panic(err)
			}
		} else {
			fmt.Println("already exists:", dst)
		}
	}
}

func yn(prompt string) bool {
	reader := bufio.NewReader(os.Stdin)
	fmt.Printf("%s (y/n): ", prompt)
	text, _ := reader.ReadString('\n')
	return strings.TrimSpace(text) == "y"
}
func rmPrefix(s, p string) string {
	if strings.HasPrefix(s, p) {
		return s[len(p):]
	}
	return s
}
func rmSuffix(s, p string) string {
	if strings.HasSuffix(s, p) {
		return s[:len(s)-len(p)]
	}
	return s
}
