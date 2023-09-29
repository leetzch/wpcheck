package main

import (
	"bufio"
	"fmt"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"
)

var (
	green  = "\033[32m"
	red    = "\033[31m"
	yellow = "\033[33m"
	reset  = "\033[0m"
)

var zner = `
     █     █░ ██▓███      ▄████▄   ██░ ██ ▓█████  ▄████▄   ██ ▄█▀
    ▓█░ █ ░█░▓██░  ██▒   ▒██▀ ▀█  ▓██░ ██▒▓█   ▀ ▒██▀ ▀█   ██▄█▒ 
    ▒█░ █ ░█ ▓██░ ██▓▒   ▒▓█    ▄ ▒██▀▀██░▒███   ▒▓█    ▄ ▓███▄░ 
    ░█░ █ ░█ ▒██▄█▓▒ ▒   ▒▓▓▄ ▄██▒░▓█ ░██ ▒▓█  ▄ ▒▓▓▄ ▄██▒▓██ █▄ 
    ░░██▒██▓ ▒██▒ ░  ░   ▒ ▓███▀ ░░▓█▒░██▓░▒████▒▒ ▓███▀ ░▒██▒ █▄
    ░ ▓░▒ ▒  ▒▓▒░ ░  ░   ░ ░▒ ▒  ░ ▒ ░░▒░▒░░ ▒░ ░░ ░▒ ▒  ░▒ ▒▒ ▓▒
      ▒ ░ ░  ░▒ ░          ░  ▒    ▒ ░▒░ ░ ░ ░  ░  ░  ▒   ░ ░▒ ▒░
      ░   ░  ░░          ░         ░  ░░ ░   ░   ░        ░ ░░ ░ 
        ░                ░ ░       ░  ░  ░   ░  ░░ ░      ░  ░   
                         ░                       ░                      
                  Mass Wordpress Checker
              Super Fast - Good Quality - FREE
         Leet'z - t.me/leetzch - github.com/leetzch

                       Copyright Z
`

func zwp(domain []string, thread int, wg *sync.WaitGroup) {
	defer wg.Done()

	client := &http.Client{Timeout: 5 * time.Second}
	for _, site := range domain {
		if !strings.HasPrefix(site, "http") && !strings.HasPrefix(site, "https") {
			url := "http://" + site + "/wp-login.php"
			resp, err := client.Head(url)
			if err != nil || resp.StatusCode != http.StatusOK {
				url = "https://" + site + "/wp-login.php"
				resp, err = client.Head(url)
			}
			if err == nil && resp.StatusCode == http.StatusOK {
				fmt.Printf("%s--: %s => [Wordpress]%s\n", green, url, reset)
				file, _ := os.OpenFile("result.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
				defer file.Close()
				file.WriteString(url + "\n")
			} else {
				fmt.Printf("%s--| %s >> [Other CMS]%s\n", red, site, reset)
			}
		}
	}
}

func main() {
	fmt.Print(zner)

	reader := bufio.NewReader(os.Stdin)
	fmt.Print("--| root@z[list]: ")
	list, _ := reader.ReadString('\n')
	list = strings.TrimSpace(list)
	fmt.Print("--| root@z[thread]: ")
	var threads int
	_, err := fmt.Scanf("%d", &threads)
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println()

	domainList, err := readFile(list)
	if err != nil {
		fmt.Println(err)
		return
	}

	threadDomains := make([][]string, threads)
	for i := 0; i < threads; i++ {
		threadDomains[i] = getThreadDomains(domainList, threads, i)
	}

	var wg sync.WaitGroup
	for i, domain := range threadDomains {
		wg.Add(1)
		go zwp(domain, i+1, &wg)
	}

	wg.Wait()

	fmt.Println("\nPress Enter to exit...")
	reader.ReadString('\n')
}

func readFile(filename string) ([]string, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var lines []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return lines, nil
}

func getThreadDomains(domains []string, totalThreads, threadIndex int) []string {
	var threadDomains []string
	for i := threadIndex; i < len(domains); i += totalThreads {
		threadDomains = append(threadDomains, domains[i])
	}
	return threadDomains
}
