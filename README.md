[![Built with Mage](https://magefile.org/badge.svg)](https://magefile.org)
[![GoDoc](https://img.shields.io/badge/godoc-reference-green.svg)](https://godoc.org/github.com/cseeger-epages/mac-gen-go)
[![Go Report Card](https://goreportcard.com/badge/github.com/cseeger-epages/mac-gen-go)](https://goreportcard.com/report/github.com/cseeger-epages/mac-gen-go)


# A IEEE 802 MAC address generator

## Purpose

This generator can be used to create IEEE 802 conform local administered MAC addresses

## get it

```
go get github.com/cseeger-epages/mac-gen-go
```

## Usage

```
package main

import (
  "fmt"
  "net"

  gm "github.com/cseeger-epages/mac-gen-go"
)

func main() {

  // generate a random local administered unicast mac prefix
  prefix := gm.GenerateRandomLocalMacPrefix(true)

  // calculates the NIC Sufix by ip address
  sufix, err := gm.CalculateNICSufix(net.ParseIP("129.168.12.127"))
  if err != nil {
    // your error handling here
  }
  mac := fmt.Sprintf("%s:%s", prefix, sufix)

  fmt.Println(mac)
}

```
see [examples](https://github.com/cseeger-epages/mac-gen-go/examples) for more information

### testing and coverage

```
mage test
mage coverage
```

## Background

### general stuff

- a /8 net has 16,777,216 addresses
- a mac address has 6 octets (2^48 addresses)
- the first 3 octets are the OUI (Organisationally Unique Identifier)
- and the last 3 octets are the NIC (Network interface controler specific)
- every octet has 8 bit
- `[b7][b6][b5][b4][b3][b2][b1][b0]`
- the first octet has special meaning for b0 and b1
- b0 -> 0: unicast, 1: multicast
- b1 -> 0: globally unique (OUI), 1: locally administered
- to generate a mac address without having a OUI
- b0 ->0
- using only the last 3 octets (16^6 = 2^24 addresses = 16777216 = /8 net)
- is enough to create mac addresses for all ip addresses
- per net the first 3 octets need to be generated as a new prefix

### prefix block generation

generating a prefix block of 3 octets

- first octet in binary xxxxxx10 ( or xxxxxx11 if you want a multicast address)
- second and 3rd octet can be generated freely
- 8 bit block in binary
- `[0|1] [0|1] [0|1] [0|1] [0|1] [0|1] [0|1] [0|1]`
- left 4 bits does not matter since they do not effect the b0 and b1 bit
- to have b0 -> 0 and b1 -> 1 (a LA unicast address)
- `[0|1] [0|1] 1 0`
- so we have sth + 2
- where sth is ether 0 (0 0) or 8 (1 0) or 4 (0 1) or 12 (1 1)
- so all valid combinations are 2, 6, 10 and 14 which are in hex
- 2, 6, A, E
- the first octet has to be sth like
- `[0..F][2|6|A|E]`
- if the address should be multicast
- b0 -> 1 and b1 -> 1
- so sth + 3
- which results in 3, 7, B, F
- and the first octet has to be sth like
- `[0..F][3|7|B|F]`
- there are 6 free bits in this octet so we have
- 2^6 permutations for this octet
- if we include the other 2 octests (sum up to additonal 16 bits)
- we have 2^22 permutations for free local administratered mac addresses (either uni or multicast)

- we should exclude some of them because of some registries failed in correctly checking if an OUI has local administraited address (LAA) bit set or not -.-
- here is a list https://gist.github.com/aallan/b4bb86db86079509e6159810ae9bd3e4
- and checked against http:-standards-oui.ieee.org/oui/oui.txt

- to exclude
- 02 since here are a few registraited ones
- AA due to DEC LAA

- so we now have
- `0[3|6|7|A|B|E|F]` - 7
- `[1..9][2|3|6|7|A|B|E|F]` - 9\*8
- A[2|3|6|7|B|E|F] - 7
- and
- `[B..F][2|3|6|7|A|B|E|F]` - 5\*8
- we will lose some of our actual permutations
- without the first octet we have 2^16 permutations (remember 2 octest = 8 bit per octet -> 16 bits)
- for our first octet we now have 126 possibile permutations left
- so we get 2^16\*126 = 8257536 possible permutations (this includes multicast)
- this is about 49,21% of the full addresspace (2^24) so we lose about half of the addresses
- in a perfect world where no registrie failed we would exactly lose 50% (2^16\*128 / 2^24) = (2^16 \* 2^7 = 2^23) / 2^24 = 0.5
- so we lose about 0.79% addresses due to registration failure exclusions
- but yeah we do have over 8 million left - show me someone who has more than 8 millions local networks
- the registrated private networks are 10.0.0.0/8 (16.77.216) 172.16.0.0/12 (1.048.576) and 192.168.0.0/16 (65.536)
- so in the worst case you can cut all these networks into /32 networks with only one ip address so you can have
- up to 17891328 networks so our 2^16\*126 permutations fit about 46.15% of this worst case szenario or 23.44% when only using unicast

