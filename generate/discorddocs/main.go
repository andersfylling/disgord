package main

import (
	"bytes"
	"errors"
	"fmt"
	"go/format"
	"html/template"
	"io/ioutil"
	"path"
	"strings"
	"sync"
	"time"

	"github.com/andersfylling/disgord/internal/constant"

	"github.com/gocolly/colly"
)

func main() {
	fmt.Println(1)
	genPermissions()
}

func newScraper() *colly.Collector {
	c := colly.NewCollector(
		colly.Async(true),
	)
	c.UserAgent = constant.UserAgent
	return c
}

type Permission struct {
	Name         string
	Description  string
	ChannelTypes string
	Value        string
}

func genPermissions() {
	var permissions []Permission
	var err error
	var done bool
	httpTimeout := 3 * time.Second
	tableURL := "https://discordapp.com/developers/docs/topics/permissions"

	c := newScraper()

	wg := sync.WaitGroup{}
	wg.Add(1)
	mu := &sync.Mutex{}
	c.OnHTML("#permissions-bitwise-permission-flags", func(e *colly.HTMLElement) {
		mu.Lock()
		defer mu.Unlock()

		done = true
		wg.Done()
	})
	c.OnRequest(func(r *colly.Request) {
		fmt.Println("Visiting", r.URL)
	})
	c.OnResponse(func(r *colly.Response) {
		fmt.Println("status", r.StatusCode)
	})
	go func() {
		<-time.After(httpTimeout)
		mu.Lock()
		defer mu.Unlock()

		if done {
			return
		}
		err = errors.New("timeout")
		wg.Done()
	}()

	if err2 := c.Visit(tableURL); err2 != nil {
		panic(err2)
	}
	if wg.Wait(); err != nil {
		panic(err)
	}
	makeFile(permissions, "generate/discorddocs/permissions.gotpl", "permissions_gen.go")
}

func ensure(v ...interface{}) {
	for i := range v {
		if err, ok := v[i].(error); ok && err != nil {
			panic(err)
		}
	}
}

// ToCamelCase takes typical CONST names such as T_AS and converts them to TAs.
// TEST_EIN_TRES => TestEinTres
func ToCamelCase(s string) string {
	b := []byte(strings.ToLower(s))
	for i := range b {
		if b[i] == '_' && i < len(b)-1 {
			b[i+1] ^= 0x20
		}
	}
	s = strings.Replace(string(b), "-", "", -1)
	return s
}

func makeFile(data interface{}, tplFile, target string) {
	fMap := template.FuncMap{
		"ToUpper":      strings.ToUpper,
		"ToLower":      strings.ToLower,
		"Decapitalize": func(s string) string { return strings.ToLower(s[0:1]) + s[1:] },
		"ToCamelCase":  ToCamelCase,
	}

	// Open & parse our template
	tpl := template.Must(template.New(path.Base(tplFile)).Funcs(fMap).ParseFiles(tplFile))

	// Execute the template, inserting all the event information
	var b bytes.Buffer
	ensure(tpl.Execute(&b, data))

	// Format it according to gofmt standards
	formatted, err := format.Source(b.Bytes())
	ensure(err)

	// And write it.
	ensure(ioutil.WriteFile(target, formatted, 0644))
}
