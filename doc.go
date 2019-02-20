/*
Package macgen lets you generate local adminstered mac address prefixes and calculate unique mac sufix per ip address

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

*/
package macgen
