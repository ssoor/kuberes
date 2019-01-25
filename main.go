package main

import (
	"github.com/ssoor/kuberes/pkg/log"

	"github.com/ssoor/kuberes/cmds"
	"sigs.k8s.io/kustomize/k8sdeps"
)

func main() {
	if err := cmds.New(k8sdeps.NewFactory()).Execute(); nil != err {
		log.Warn("execute cmd failed, ", err)
	}
}
