package macgen

import (
	"fmt"
	"math/rand"
	"net"
	"strconv"
	"strings"
	"time"
)

// GenerateRandomLocalMacPrefix generates a random LAA mac address
// unicast true sets the unicast bit
// unicast false sets the multicast bit
func GenerateRandomLocalMacPrefix(unicast bool) string {

	// randomize things
	r := rand.New(rand.NewSource(time.Now().UnixNano()))

	// [s]ignificant [b]its [oc]tant [1] [a]llowed
	var sboc1a []uint8
	if unicast {
		// 2, 6, A, E
		sboc1a = []uint8{2, 6, 10, 14}
	} else {
		// 3, 7, B, F
		sboc1a = []uint8{3, 7, 11, 15}
	}

	excludedCombinations := []string{"02", "AA"}

	left := r.Intn(15)
	right := sboc1a[r.Intn(len(sboc1a))]

	oc1 := fmt.Sprintf("%x%x", left, right)

	for _, ex := range excludedCombinations {
		if oc1 == ex {
			oc1 = GenerateRandomLocalMacPrefix(unicast)
		}
	}

	oc2 := fmt.Sprintf("%x%x", r.Intn(15), r.Intn(15))
	oc3 := fmt.Sprintf("%x%x", r.Intn(15), r.Intn(15))
	return fmt.Sprintf("%s:%s:%s", oc1, oc2, oc3)
}

// CalculateNICSufix calculates the mac address by given IP
func CalculateNICSufix(ip net.IP) (string, error) {
	if ip.To4() == nil {
		return "", fmt.Errorf("ip is not v4")
	}

	split := strings.Split(ip.String(), ".")

	n := make(map[int]int, 3)
	for i := 1; i <= 3; i++ {
		num, err := strconv.Atoi(split[i])
		n[i] = num
		if err != nil {
			return "", err
		}
	}
	// since we can generate 2^24  addresses
	// we can map all addresses we can generate to
	// all addresses of the /8 net of the given IP address, which fits perfectly !
	// so we just need to multiply the 2nd 3rd and 4th byte to get the actuall number
	// inside of our 2^24 possible addresses
	// we actually need to generate the (n[1]+1)*(n[2]+1)*(n[3]+1)-1 mac address
	// why the +1, simply because we do count X.0.0.0 as the first address which is
	// 1*1*1 and X.255.255.255 the last which is 256*256*256 = 16^8 = 2^24
	// but now we have to substract by one since we count from 0 to (2^24)-1 and not from
	// 1 to 2^24

	hex := fmt.Sprintf("%06x\n", (n[1]+1)*(n[2]+1)*(n[3]+1)-1)

	return fmt.Sprintf("%s:%s:%s", hex[0:2], hex[2:4], hex[4:6]), nil
}
