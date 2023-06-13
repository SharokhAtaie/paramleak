package main

import (
	"bufio"
	"context"
	"fmt"
	"github.com/projectdiscovery/goflags"
	"github.com/projectdiscovery/gologger"
	fUtils "github.com/projectdiscovery/utils/file"
	httpUtil "github.com/projectdiscovery/utils/http"
	sUtils "github.com/projectdiscovery/utils/slice"
	"github.com/SharokhAtaie/paramleak/regex"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"regexp"
	"strings"
	"time"
)

type options struct {
	url     string
	list    string
	delay   time.Duration
	silent  bool
	method  string
	body    string
	header  string
	verbose bool
}

var (
	opt = &options{}
)

func main() {
	flagSet := goflags.NewFlagSet()
	flagSet.SetDescription(`Paramleak is a tool for get all parameters from input`)
	flagSet.StringVarP(&opt.url, "url", "u", "", "url for get parameters")
	flagSet.StringVarP(&opt.list, "list", "l", "", "List of Url for Get All parameters")
	flagSet.StringVarP(&opt.method, "method", "X", "GET", "Http Method for request")
	flagSet.StringVarP(&opt.body, "body", "d", "", "Body data for Post Request")
	flagSet.StringVarP(&opt.header, "header", "H", "", "Custom Header (You can set only one custom header)")
	flagSet.DurationVarP(&opt.delay, "delay", "p", 0, "time for delay example: 1000 Millisecond (1 second)")
	flagSet.BoolVarP(&opt.verbose, "verbose", "v", false, "Verbose mode")
	flagSet.BoolVarP(&opt.silent, "silent", "s", false, "Silent mode")

	if err := flagSet.Parse(); err != nil {
		gologger.Error().Msgf("Could not parse flags: %s", err)
	}

	if opt.url == "" && opt.list == "" && !fUtils.HasStdin() {
		printUsage()
		return
	}

	if !isValidHTTPMethod(opt.method) {
		gologger.Error().Msgf("Invalid HTTP method!\nAllowed methods: GET, POST, PUT, PATCH, DELETE")
		return
	}

	var parts = strings.SplitN(opt.header, ":", 2)
	if opt.header != "" && len(parts) != 2 {
		gologger.Error().Msgf("invalid custom header! (Your input must be two part separate with ':')\n\t -H 'Cookie: test=test;'")
		return
	}

	if !opt.silent {
		Banner()
	}

	input := getInput()
	if len(input) > 0 {
		results := run(input, opt.method, opt.body, opt.delay)
		for _, v := range results {
			fmt.Println(v)
		}
	}
}

func getInput() []string {
	input := make([]string, 0)
	switch {
	case fUtils.HasStdin():
		scanner := bufio.NewScanner(os.Stdin)
		for scanner.Scan() {
			input = append(input, scanner.Text())
		}
	case opt.url != "":
		input = append(input, opt.url)
	case opt.list != "":
		lists, err := fUtils.ReadFile(opt.list)
		if err != nil {
			gologger.Error().Msgf("%v", err)
			return nil
		}
		for content := range lists {
			input = append(input, content)
		}
	}
	return input
}

func run(lists []string, method, body string, delay time.Duration) []string {
	output := make([]string, 0)
	c := make(chan []string)

	for _, Url := range lists {
		go func(URL string) {
			time.Sleep(delay * time.Millisecond)
			data, err := Request(method, URL, body)
			if err != nil {
				gologger.Error().Msgf("%v", err)
				return
			}
			result := regex.Regex(data)
			c <- result
		}(Url)
	}

	for range lists {
		result := <-c
		for _, s := range result {
			reg := regexp.MustCompile(" ")
			s = reg.ReplaceAllString(s, "")
			output = append(output, s)
		}
	}

	return sUtils.Dedupe(output)
}

var client = &http.Client{
	Timeout: 10 * time.Second,
}

func Request(method, urlStr, bodyStr string) (string, error) {
	u, err := url.Parse(urlStr)
	if err != nil {
		return "", err
	}

	if u.Scheme == "" {
		urlStr = "https://" + u.Host + u.Path
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	reqBody := strings.NewReader(bodyStr)
	req, err := http.NewRequestWithContext(ctx, method, urlStr, reqBody)
	if err != nil {
		return "", err
	}

	req.Header.Set("User-Agent", "Mozilla/5.0 (X11; Linux x86_64; rv:109.0) Gecko/20100101 Firefox/112.0")
	req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,*/*;q=0.8")
	req.Header.Set("Origin", "https://"+u.Host)

	var parts = strings.SplitN(opt.header, ":", 2)
	if opt.header != "" {
		key := strings.TrimSpace(parts[0])
		value := strings.TrimSpace(parts[1])
		// Set custom header
		req.Header.Set(key, value)
	}

	if opt.verbose {
		dump, err := httpUtil.DumpRequest(req)
		if err != nil {
			return "", err
		}
		gologger.Print().Msgf(dump)
	}

	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	return string(respBody), nil
}

func isValidHTTPMethod(method string) bool {
	validMethods := []string{"GET", "POST", "PUT", "PATCH", "DELETE"}
	for _, m := range validMethods {
		if method == m {
			return true
		}
	}
	return false
}

func printUsage() {
	Banner()
	gologger.Print().Msgf("Flags:\n")
	gologger.Print().Msgf("\t-url,      -u       URL for getting all parameters")
	gologger.Print().Msgf("\t-list,     -l       List of URLs for getting all parameters")
	gologger.Print().Msgf("\t-method,   -X       HTTP method for requests")
	gologger.Print().Msgf("\t-body,     -d       Body data for POST/PATCH requests")
	gologger.Print().Msgf("\t-header,   -H       Custom header (You can set only 1 custom header)")
	gologger.Print().Msgf("\t-delay,    -p       Delay time example: 1000 milliseconds (1 second)")
	gologger.Print().Msgf("\t-verbose,  -v       Verbose mode")
	gologger.Print().Msgf("\t-silent,   -s       Silent mode")
}

func Banner() {
	gologger.Print().Msgf(`
██████╗  █████╗ ██████╗  █████╗ ███╗   ███╗██╗     ███████╗ █████╗ ██╗  ██╗
██╔══██╗██╔══██╗██╔══██╗██╔══██╗████╗ ████║██║     ██╔════╝██╔══██╗██║ ██╔╝
██████╔╝███████║██████╔╝███████║██╔████╔██║██║     █████╗  ███████║█████╔╝ 
██╔═══╝ ██╔══██║██╔══██╗██╔══██║██║╚██╔╝██║██║     ██╔══╝  ██╔══██║██╔═██╗ 
██║     ██║  ██║██║  ██║██║  ██║██║ ╚═╝ ██║███████╗███████╗██║  ██║██║  ██╗
╚═╝     ╚═╝  ╚═╝╚═╝  ╚═╝╚═╝  ╚═╝╚═╝     ╚═╝╚══════╝╚══════╝╚═╝  ╚═╝╚═╝  ╚═╝
`)
	gologger.Print().Msgf("\t\tCreated by Sharo_k_h :)\n\n")
}
