package collector

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
)

func getNetDevInfo() (map[string]map[string]string, error) {

	netDev, err := getNetDev()
	if err != nil {
		return nil, fmt.Errorf("couldn't get network devices: %s", err)
	}

	// /proc/net/route stores the kernel's routing table
	// The interface whose destination is 00000000 is the interface of the default gateway
	file, err := os.Open("/proc/net/route")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	scanner.Scan()
	for scanner.Scan() {
		if scanner.Text() == "" {
			break
		}

		s := strings.FieldsFunc(scanner.Text(), Split)
		if s[1] == "00000000" {
			netDev[s[0]]["default"] = "true"
		}
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
		return nil, err
	}

	// ls -l /sys/class/net/
	// docker0 -> ../../devices/virtual/net/docker0
	// eth0 -> ../../devices/pci0000:00/0000:00:01.0/virtio0/net/eth0
	for iface, info := range netDev {
		devLink, _ := os.Readlink("/sys/class/net/" + iface)
		if strings.Contains(devLink, "virtual") {
			info["type"] = "virtual"
		} else {
			info["type"] = "physical"

			// extract pci id
			folderPath := strings.Split(devLink, "/")
			info["pci_id"] = folderPath[4]
		}
	}

	return netDev, nil
}
