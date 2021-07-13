package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"regexp"
	"sort"
	"strconv"
)

const defaultPattern = `release-(\d+)\.(\d+)\.(\d+).*`

var (
	m = make(map[int]string)
)

func main() {
	if len(os.Args) < 2 {
		fmt.Fprintf(os.Stderr, "usage %s <path to release dir> [optional pattern]\n", os.Args[0])
		os.Exit(1)
	}

	var (
		path    = os.Args[1]
		pattern string
	)
	// read optional pattern
	if len(os.Args) > 2 {
		pattern = os.Args[2]
	} else {
		pattern = defaultPattern
	}

	d, err := ioutil.ReadDir(path)
	if err != nil {
		panic(err)
	}
	regex := regexp.MustCompile(pattern)

	for _, f := range d {
		name := f.Name()
		if f.IsDir() {
			match := regex.FindAllStringSubmatch(name, -1)

			if len(match) > 0 && len(match[0]) >= 4 {
				m1, m2, m3 := match[0][1], match[0][2], match[0][3]

				n, _ := strconv.Atoi(fmt.Sprintf("%03s%03s%03s", m1, m2, m3))
				m[n] = name
			}
		}
	}

	if len(m) == 0 { // no matches found
		fmt.Fprintln(os.Stderr, "no releases found")
		if os.Getenv("DEBUG") != "" {
			for _, dd := range d {
				fmt.Fprintln(os.Stderr, "name:", dd.Name(), "is dir:", dd.IsDir())
			}
		}
		os.Exit(1)
	}

	keys := make([]int, len(m))
	i := 0
	for k := range m {
		keys[i] = k
		i++
	}
	sort.Ints(keys)

	if os.Getenv("DEBUG") != "" {
		for _, k := range keys {
			fmt.Fprintln(os.Stderr, k, m[k])
		}
	}

	fmt.Println(m[keys[len(keys)-1]])
}
