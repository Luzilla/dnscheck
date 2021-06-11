package main

import (
	"fmt"
	"log"
	"net"
	"os"

	"github.com/miekg/dns"
	"github.com/miekg/dns/dnsutil"
	"github.com/olekukonko/tablewriter"
	"github.com/urfave/cli/v2"
)

var version string
var date string
var commit string

func main() {
	app := cli.NewApp()

	app.Name = "dnscheck"
	app.Usage = "A cli tool to check records on nameservers"
	app.Version = fmt.Sprintf(
		"%s (build: %s, commit: %s)",
		version,
		date,
		commit)

	app.Authors = []*cli.Author{
		{
			Name:  "Till Klampaeckel",
			Email: "till@luzilla-capital.com",
		},
	}
	app.Commands = []*cli.Command{
		{
			Name:    "check",
			Aliases: []string{"c"},
			Usage:   "Check the host",
			Flags: []cli.Flag{
				&cli.StringFlag{
					Name:     "host",
					Required: true,
					Usage:    "The host to check",
					EnvVars:  []string{"DNSCHECK_HOST"},
				},
				&cli.StringFlag{
					Name:    "type",
					Value:   "A",
					Usage:   "Record type (A, MX, ...)",
					EnvVars: []string{"DNSCHECK_TYPE"},
				},
			},
			Action: func(c *cli.Context) error {
				host := c.String("host")

				config := getConfig()
				client := getClient()

				log.Printf("Discovered local resolver: %s:%s", config.Servers[0], config.Port)

				log.Printf("Probing for NS records (nameservers) for %s\n", host)

				response, err := fetchDNS(host, client, config)
				if err != nil {
					return err
				}

				var answers []dns.RR

				if len(response.Answer) == 0 {
					// need to dig deeper
					var soaRecord *dns.SOA

					soaRecord, ok := response.Ns[0].(*dns.SOA)
					if !ok {
						return fmt.Errorf("Couldn't find authority for %s", host)
					}

					log.Printf("Found SOA record: %s", soaRecord.Header().Name)

					zone := dnsutil.TrimDomainName(soaRecord.Header().Name, ".")

					response, err := fetchDNS(zone, client, config)
					if err != nil {
						return err
					}

					if len(response.Answer) == 0 {
						return fmt.Errorf("Couldn't find DNS server: %s", soaRecord.Hdr.Name)
					}
					answers = response.Answer
				} else {
					answers = response.Answer
				}

				log.Println(fmt.Sprintf("Found %d nameservers for %s.", len(answers), host))

				table := tablewriter.NewWriter(os.Stdout)
				table.SetHeader([]string{"NS", "TTL", "TYPE", "Data"})

				foundError := false

				// Stuff must be in the answer section
				for _, a := range answers {
					nsRecord, ok := a.(*dns.NS)

					if !ok {
						continue
					}

					var row []string
					server := dnsutil.TrimDomainName(nsRecord.Ns, ".")

					aMessage := createMessage(host, dns.TypeA)
					aResponse, _, err := client.Exchange(aMessage, net.JoinHostPort(server, "53"))
					if err != nil {
						return err
					}

					if response.Rcode != dns.RcodeSuccess {
						row = append(row, server)
						row = append(row, "XXX")
						row = append(row, "XXX")
						row = append(row, "XXX")
						table.Append(row)

						foundError = true
						continue
					}

					if len(aResponse.Answer) == 0 {
						row = append(row, server)
						row = append(row, "XXX")
						row = append(row, "XXX")
						row = append(row, "XXX")
						table.Append(row)

						foundError = true
						continue
					}

					for _, b := range aResponse.Answer {
						if aRecord, ok := b.(*dns.A); ok {
							row = append(row, server)
							row = append(row, fmt.Sprint(aRecord.Header().Ttl))
							row = append(row, "A")
							row = append(row, aRecord.A.String())
							table.Append(row)
						}
					}
				}
				table.Render()

				if foundError {
					return fmt.Errorf("One or more errors were discovered during this check.")
				}

				return nil
			},
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}

func getConfig() *dns.ClientConfig {
	config, _ := dns.ClientConfigFromFile("/etc/resolv.conf")
	return config
}

func getClient() *dns.Client {
	return new(dns.Client)
}

func createMessage(host string, record uint16) *dns.Msg {
	m := new(dns.Msg)
	m.SetQuestion(dns.Fqdn(host), record)
	m.RecursionDesired = true

	return m
}

func fetchDNS(host string, client *dns.Client, config *dns.ClientConfig) (*dns.Msg, error) {
	log.Printf("=> host: %s", host)
	nsMessage := createMessage(host, dns.TypeNS)

	var response *dns.Msg

	response, _, err := client.Exchange(
		nsMessage,
		net.JoinHostPort(config.Servers[0], config.Port))
	if err != nil {
		return response, err
	}

	if response.Rcode != dns.RcodeSuccess {
		return response, fmt.Errorf("Invalid answer to NS query for %s", host)
	}

	return response, nil
}
