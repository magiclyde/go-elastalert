/**
 * Created by GoLand.
 * @author: clyde
 * @date: 2021/10/20 下午1:51
 * @note:
 */

package elastalert

import (
	"os/exec"
	"testing"
	"time"
)

func TestFind(t *testing.T) {
	cmd := `time find /cargo/gos -name "*.go" | wc -l`
	out, _ := exec.Command("bash", "-c", cmd).Output()
	t.Logf("out: %s", out)
	// use: 0.23s found: 43293 files
}

func TestWalkDir(t *testing.T) {
	for _, descend := range []bool{false, true} {
		start := time.Now()
		out := WalkDir("/cargo/gos", "go", descend)
		n := 0
		for range out {
			n++
		}
		t.Logf("use: %fs found: %d files", time.Since(start).Seconds(), n)
		// use: 0.750630s found: 43293 files when descend
	}
}

func TestParallelWalkDir(t *testing.T) {
	for _, descend := range []bool{false, true} {
		start := time.Now()
		out := ParallelWalkDir("/cargo/gos", "go", descend)
		n := 0
		for range out {
			n++
		}
		t.Logf("use: %fs found: %d files", time.Since(start).Seconds(), n)
		// use: 0.385026s found: 43293 files when descend
	}
}
