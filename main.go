package main

import (
	"flag"
	"fmt"
	"time"

	"github.com/bigdatadev/goryman"
	"github.com/golang/glog"
	"github.com/google/cadvisor/client"
	"github.com/google/cadvisor/info"
)

var riemannAddress = flag.String("riemann_address", "localhost:5555", "specify the riemann server location")
var cadvisorAddress = flag.String("cadvisor_address", "http://localhost:8080", "specify the cadvisor API server location")
var sampleInterval = flag.Int("sample_interval", 1000, "specify the sampling interval")

func pushToRiemann(r *goryman.GorymanClient, service string, metric int, tags []string) {
	err := r.SendEvent(&goryman.Event{
		Service: service,
		Metric:  metric,
		Tags:    tags,
	})
	if err != nil {
		glog.Fatalf("unable to write to riemann: %s", err)
	}
}

func main() {
	defer glog.Flush()
	flag.Parse()

	// Setting up the Riemann client
	r := goryman.NewGorymanClient(*riemannAddress)
	err := r.Connect()
	if err != nil {
		glog.Fatalf("unable to connect to riemann: %s", err)
	}
	//defer r.Close()

	// Setting up the cadvisor client
	c, err := client.NewClient(*cadvisorAddress)
	if err != nil {
		glog.Fatalf("unable to setup cadvisor client: %s", err)
	}

	// Setting up the 1 second ticker
	ticker := time.NewTicker(1 * time.Second).C
	for {
		select {
		case <-ticker:
			// Make the call to get all the possible data points
			request := info.ContainerInfoRequest{10}
			returned, err := c.AllDockerContainers(&request)
			if err != nil {
				glog.Fatalf("unable to retrieve machine data: %s", err)
			}
			// Start dumping data into riemann
			// Loop into each ContainerInfo
			// Get stats
			// Push into riemann
			for _, container := range returned {
				pushToRiemann(r, fmt.Sprintf("Load %s", container.Name), int(container.Stats[0].Cpu.Load), []string{})
			}
		}
	}

}
