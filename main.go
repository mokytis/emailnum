package main

import (
	"bufio"
	"flag"
	"fmt"
	"net"
	"os"
	"strings"
	"sync"
)

type Domain struct {
	name          string
	spf_records   []string
	dmarc_records []string
}

func get_spf(domain string) []string {
	potential_spf, _ := net.LookupTXT(domain)
	var spf_records []string
	for _, record := range potential_spf {
		if strings.HasPrefix(record, "v=spf1") {
			spf_records = append(spf_records, record)
		}
	}
	return spf_records
}

func get_dmarc(domain string) []string {
	potential_dmarc, _ := net.LookupTXT("_dmarc." + domain)
	var dmarc_records []string
	for _, record := range potential_dmarc {
		if strings.HasPrefix(record, "v=DMARC1") {
			dmarc_records = append(dmarc_records, record)
		}
	}
	return dmarc_records
}

func main() {
	var concurrency int
	flag.IntVar(&concurrency, "c", 20, "concurrency level")

	var spf int
	flag.IntVar(&spf, "spf", 1, "-1 = must not exist\n 0 = don't check\n 1 = can exist\n 2 = must exist\n")

	var dmarc int
	flag.IntVar(&dmarc, "dmarc", 1, "-1 = must not exist\n 0 = don't check\n 1 = can exist\n 2 = must exist\n")

	var domonly bool
	flag.BoolVar(&domonly, "x", false, "just output domain names that match the rules")

	flag.Parse()

	if spf < -1 || spf > 2 {
		fmt.Fprintf(os.Stderr, "spf should be one of (-1, 0, 1, 2) not %d", spf)
		os.Exit(1)
	}
	if dmarc < -1 || dmarc > 2 {
		fmt.Fprintf(os.Stderr, "dmarc should be one of (-1, 0, 1, 2) not %d", dmarc)
		os.Exit(1)
	}

	domains := make(chan string)
	output := make(chan Domain)

	var dnswg sync.WaitGroup
	for i := 0; i < concurrency; i++ {
		dnswg.Add(1)
		go func() {
			defer dnswg.Done()
			for dom := range domains {
				var spf_records []string
				var dmarc_records []string
				if spf != 0 {
					spf_records = get_spf(dom)
				}
				if dmarc != 0 {
					dmarc_records = get_dmarc(dom)
				}
				domain := Domain{dom, spf_records, dmarc_records}
				output <- domain
			}
		}()

	}

	var outputWG sync.WaitGroup
	outputWG.Add(1)
	go func() {
		defer outputWG.Done()
		for o := range output {
			var spf_good bool
			var dmarc_good bool

			if spf == 1 || spf == 0 {
				spf_good = true
			} else if spf == -1 {
				spf_good = len(o.spf_records) == 0
			} else if spf == 2 {
				spf_good = len(o.spf_records) > 0
			}

			if dmarc == 1 || dmarc == 0 {
				dmarc_good = true
			} else if dmarc == -1 {
				dmarc_good = len(o.dmarc_records) == 0
			} else if dmarc == 2 {
				dmarc_good = len(o.dmarc_records) > 0
			}

			if spf_good && dmarc_good {
				if domonly {
					fmt.Println(o.name)
				} else {
					for _, record := range o.spf_records {
						fmt.Printf("%s \"%s\"\n", o.name, record)
					}
					for _, record := range o.dmarc_records {
						fmt.Printf("%s \"%s\"\n", o.name, record)
					}
				}
			}
		}
	}()

	go func() {
		dnswg.Wait()
		close(output)
	}()

	sc := bufio.NewScanner(os.Stdin)
	for sc.Scan() {
		domain := strings.ToLower(sc.Text())
		domains <- domain
	}
	close(domains)

	outputWG.Wait()
}
