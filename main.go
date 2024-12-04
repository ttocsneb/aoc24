package main

import (
	"bufio"
	_ "embed"
	"flag"
	"fmt"
	"io"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"os"
	"os/exec"
	"reflect"
	"slices"
	"time"

	"github.com/joho/godotenv"
	_ "github.com/joho/godotenv"
)

//go:embed d1.py
var d1 []byte

//go:embed d2.py
var d2 []byte

type Day struct {
	Input *bufio.Reader
}

func (self *Day) D1() error {
	python, err := exec.LookPath("python3")
	if err != nil {
		return err
	}

	cmd := exec.Command(python)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	w, err := cmd.StdinPipe()
	go func() {
		w.Write(d1)
		w.Close()
	}()

	return cmd.Run()
}

func (self *Day) D2() error {
	python, err := exec.LookPath("python3")
	if err != nil {
		return err
	}

	cmd := exec.Command(python)
	w, err := cmd.StdinPipe()
	go func() {
		w.Write(d2)
		w.Close()
	}()
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func main() {
	flag.Usage = func() {
		fmt.Printf("Usage of %s:\n", os.Args[0])

		flag.PrintDefaults()
	}

	godotenv.Load(".env")

	if sess := os.Getenv("SESSION"); sess != "" {
		u, _ := url.Parse("https://adventofcode.com")
		var err error
		http.DefaultClient.Jar, err = cookiejar.New(&cookiejar.Options{})
		if err != nil {
			panic(err)
		}
		http.DefaultClient.Jar.SetCookies(u, []*http.Cookie{
			{
				Name:     "session",
				Value:    sess,
				Path:     "/",
				Domain:   ".adventofcode.com",
				Expires:  time.Now().Add(time.Second * 300),
				Secure:   true,
				HttpOnly: true,
				SameSite: http.SameSiteNoneMode,
			},
		})

	}

	file := flag.String("file", "", "file to read")
	day := flag.Int("day", 0, "Day to run")
	all := flag.Bool("all", false, "Run all days")
	flag.StringVar(file, "f", "", "file to read")
	flag.IntVar(day, "d", 0, "Day to run")
	flag.BoolVar(all, "a", false, "Run all days")

	vars := []string{"-redirects", "-upload", "-key", "-file", "-expires", "-cert"}

	flags := []string{}
	args := []string{}
	nextIsFlag := false
	for _, arg := range os.Args[1:] {
		if nextIsFlag {
			flags = append(flags, arg)
			nextIsFlag = false
			continue
		}
		if len(arg) > 0 && arg[0] == '-' {
			if slices.Contains(vars, arg) {
				nextIsFlag = true
			}
			flags = append(flags, arg)
		} else {
			args = append(args, arg)
		}
	}

	flag.Parse()

	if *all {
		*day = 0
	}

	var f *os.File

	runtime := &Day{}
	readFile := func() {
		var err error
		if *file == "" {
			*file = fmt.Sprintf("input-d%d", *day)
			if _, err := os.Stat(*file); err != nil {
				resp, err := http.Get(fmt.Sprintf("https://adventofcode.com/2024/day/%d/input", *day))
				if err != nil {
					panic(err)
				}
				if resp.StatusCode != 200 {
					data, _ := io.ReadAll(resp.Body)
					panic(string(data))
				}
				defer resp.Body.Close()
				f, err := os.OpenFile(*file, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0o644)
				if err != nil {
					panic(err)
				}
				_, err = io.Copy(f, resp.Body)
				if err != nil {
					panic(err)
				}
				f.Close()
			}
		}
		f, err = os.Open(*file)
		if err != nil {
			panic(err)
		}

		runtime.Input = bufio.NewReader(f)
	}

	v := reflect.ValueOf(runtime)

	getFunc := func(day int) func() error {
		m := v.MethodByName(fmt.Sprintf("D%d", day))
		if !m.IsValid() || m.IsZero() || m.IsNil() {
			return nil
		}
		return func() error {
			readFile()
			res := m.Call([]reflect.Value{})
			f.Close()
			err := res[0].Interface()
			if err != nil {
				return err.(error)
			}
			return nil
		}
	}

	if *day == 0 {
		var fn func() error
		for i := 1; true; i++ {
			f := getFunc(i)
			if f == nil {
				break
			}
			*day = i
			fn = f
			if *all {
				*file = ""
				fmt.Printf("=== Day %d ===\n", i)
				err := fn()
				if err != nil {
					panic(err)
				}
			}
		}
		if !*all {
			fmt.Printf("=== Day %d ===\n", *day)
			err := fn()
			if err != nil {
				panic(err)
			}
		}
	} else {
		fn := getFunc(*day)
		if fn == nil {
			panic("Could not find day")
		}
		fmt.Printf("=== Day %d ===\n", *day)
		err := fn()
		if err != nil {
			panic(err)
		}
	}

}
