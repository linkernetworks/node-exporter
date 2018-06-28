//phstsai
package collector

import (
	"fmt"
	"net"

	"github.com/prometheus/client_golang/prometheus"
)

type netIfaceCollector struct {
	subsystem   string
	metricDescs map[string]*prometheus.Desc
}

func init() {
	registerCollector("iface", defaultEnabled, NewNetIfaceCollector)
}

// NewNetIfaceCollector returns a new Collector exposing network device stats.
func NewNetIfaceCollector() (Collector, error) {
	return &netIfaceCollector{
		subsystem:   "network",
		metricDescs: map[string]*prometheus.Desc{},
	}, nil
}

func (c *netIfaceCollector) Update(ch chan<- prometheus.Metric) error {
	netDevDefault, err := getNetDevInfo()
	if err != nil {
		return fmt.Errorf("couldn't get default network devices: %s", err)
	}

	for key, info := range netDevDefault {
		desc, ok := c.metricDescs[key]

		if !ok {
			desc = prometheus.NewDesc(
				prometheus.BuildFQName(namespace, c.subsystem, "interface"),
				fmt.Sprintf("List network devices and label the default one."),
				[]string{"device", "default", "ip_address", "type"},
				nil,
			)
			c.metricDescs[key] = desc
		}

		ch <- prometheus.MustNewConstMetric(desc, prometheus.GaugeValue, float64(1), key, info["default"], info["addr"], info["type"])
	}
	return nil
}

func Split(r rune) bool {
	return r == ' ' || r == '\t'
}

func getNetDev() (map[string]map[string]string, error) {
	netDev := map[string]map[string]string{}

	ifaces, _ := net.Interfaces()
	for _, iface := range ifaces {
		info := map[string]string{}
		netDev[iface.Name] = info

		addr, _ := iface.Addrs()
		if addr != nil && len(addr) > 0 {
			netDev[iface.Name]["addr"] = addr[0].String()
		}
		netDev[iface.Name]["default"] = "false"
	}

	return netDev, nil
}
