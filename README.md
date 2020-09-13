# emailnum

Take a list of domains and check their SPF and DMARC records.

## Install

``` go get -u gitlab.com/mokytis/emailnum
```

## Basic Usage

`emailnum` accepts domains from stdin

```
$ cat domains
example.com
example.org
example.net
gitlab.com
github.com

$ cat domains | emailnum
xample.com "v=spf1 -all"
github.com "v=spf1 ip4:192.30.252.0/22 ip4:208.74.204.0/22 ip4:46.19.168.0/23 include:_spf.google.com include:esp.github.com include:_spf.createsend.com include:servers.mcsv.net ~all"
github.com "v=DMARC1; p=none; rua=mailto:dmarc@github.com,mailto:d@rua.agari.com"
gitlab.com "v=spf1 include:mail.zendesk.com include:_spf.google.com include:mktomail.com include:_spf.salesforce.com include:_spf-ip.gitlab.com a:zgateway.zuora.com -all"
gitlab.com "v=DMARC1; p=reject; pct=100"
example.org "v=spf1 -all"
example.net "v=spf1 -all"
```

## Flags

### `-c` concurrently level

Sets the amount of workers to use

### `-dmarc` dmarc filter

| Value | Description                                                                                      |
|------:|--------------------------------------------------------------------------------------------------|
|    -1 | Check for dmarc records. Only show output if a dmarc record doesn't exists                       |
|     0 | Don't check for dmarc records                                                                    |
|     1 | Check for dmarc records. If other filters match show output regardless of dmarc record existence |
|     2 | Check for dmarc records. Only show output if a dmarc record exists                               |

### `-spf` spf filter

| Value | Description                                                                                  |
|------:|----------------------------------------------------------------------------------------------|
|    -1 | Check for spf records. Only show output if a spf record doesn't exists                       |
|     0 | Don't check for spf records                                                                  |
|     1 | Check for spf records. If other filters match show output regardless of spf record existence |
|     2 | Check for spf records. Only show output if a spf record exists                               |

### `-x` show domain only

* If `-x` is set, output domains that match the given filters, but don't output
  the records themselves.
* If `-x` is not set, output domains and matching records for domains that
  match the filters. If both filters are set to `-1` then there will be no
  output even if the domain matches.

## More Complex Examples

### Find SPF records for domains that have no DMARC

```
$ cat domains | emailnum -spf 2 -dmarc -1
example.com "v=spf1 -all"
example.com "v=spf1 -all"
example.org "v=spf1 -all"
example.net "v=spf1 -all"
```

### Find domains that have no DMARC, but don't check SPF

```
$ cat domains | emailnum -dmarc -1 -spf 0 -x
example.org
example.com
example.net
```

### Find domains with SPF records ending in `?all`

```
$ cat domains | emailnum -dmarc 0 | grep ~all
github.com "v=spf1 ip4:192.30.252.0/22 ip4:208.74.204.0/22 ip4:46.19.168.0/23 include:_spf.google.com include:esp.github.com include:_spf.createsend.com include:servers.mcsv.net ~all"
```

### Find domains with SPF records ending in `?all` that don't have DMARC

```
$ cat domains | emailnum -dmarc -1 | grep ~all
```

### Find domains with SPF records ending in `?all` that do have DMARC

```
$ cat domains | emailnum -dmarc 2 | grep ~all
github.com "v=spf1 ip4:192.30.252.0/22 ip4:208.74.204.0/22 ip4:46.19.168.0/23 include:_spf.google.com include:esp.github.com include:_spf.createsend.com include:servers.mcsv.net ~all"
```

## Thanks

This tool was created to automate some of my workflow, but was inspired by
tools made my [tomnomnom](https://github.com/tomnomnom).
