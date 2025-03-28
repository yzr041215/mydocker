package runc

import (
	"testing"
)

// sudo mount -t proc proc "/home/yzr/Documents/mydocker/data/overlay2/e518df1ad4972f366a0a74b1a2e0859e5262ed03825a700e6a06b4a5d782daa9/merged/proc"
// sudo /usr/local/go/bin/go test -v -count=1 -run ^TestRunA$ engine/internal/pkg/runc
func TestRunA(t *testing.T) {
	//ass
	 rootfs, err := Mount("mysql")
	 if err != nil {
	 	t.Fatal(err)
	 }
	//rootfs := "/home/yzr/Documents/mydocker/data/overlay2/e518df1ad4972f366a0a74b1a2e0859e5262ed03825a700e6a06b4a5d782daa9/merged"
	if err := Run(rootfs); err != nil {
		t.Fatal(err)
	}
}
