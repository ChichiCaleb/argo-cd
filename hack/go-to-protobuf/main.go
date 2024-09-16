// go-to-protobuf generates a Protobuf IDL from a Go struct, respecting any
// existing IDL tags on the Go struct.
package main

import (
	goflag "flag"

	"github.com/argoproj/argo-cd/v2/hack/go-to-protobuf/protobuf"
	flag "github.com/spf13/pflag"
	"k8s.io/klog/v2"
)

var g = protobuf.New()

func init() {
	klog.InitFlags(nil)
	g.BindFlags(flag.CommandLine)
	goflag.Set("logtostderr", "true")
	flag.CommandLine.AddGoFlagSet(goflag.CommandLine)
}

func main() {
	flag.Parse()
	protobuf.Run(g)
}
