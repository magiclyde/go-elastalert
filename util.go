/**
 * Created by GoLand.
 * @author: clyde
 * @date: 2021/10/13 下午4:03
 * @note:
 */

package elastalert

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"sync"
)

func Concat(values ...string) string {
	var buffer bytes.Buffer
	for _, s := range values {
		buffer.WriteString(s)
	}
	return buffer.String()
}

func WalkDir(dir, suffix string, descend bool) <-chan string {
	ext := Concat(".", strings.ToLower(suffix))
	out := make(chan string)
	go func() {
		filepath.Walk(dir, func(path string, fi os.FileInfo, _ error) (err error) {
			if fi.IsDir() && path != dir {
				if descend {
					return
				}
				return filepath.SkipDir
			}
			//filter file by extension
			if strings.ToLower(filepath.Ext(path)) == ext {
				out <- path
			}
			return
		})
		defer close(out)
	}()
	return out
}

func ParallelWalkDir(dir, suffix string, descend bool) <-chan string {
	out := make(chan string)
	go func() {
		maxGoroutines := 1000
		guard := make(chan struct{}, maxGoroutines)

		wg := &sync.WaitGroup{}
		wg.Add(1)
		parallelWalkDir(dir, suffix, descend, wg, out, guard)
		wg.Wait()

		close(out)
	}()
	return out
}

func parallelWalkDir(dir, suffix string, descend bool, wg *sync.WaitGroup, out chan string, guard chan struct{}) {
	defer func() {
		wg.Done()
		<-guard
	}()

	guard <- struct{}{} // would block if guard channel is already filled

	ext := Concat(".", strings.ToLower(suffix))
	visit := func(path string, fi os.FileInfo, err error) error {
		if fi.IsDir() && path != dir {
			if descend {
				wg.Add(1)
				go parallelWalkDir(path, suffix, descend, wg, out, guard)
			}
			return filepath.SkipDir
		}
		//filter file by extension
		if strings.ToLower(filepath.Ext(path)) == ext {
			out <- path
		}
		return nil
	}
	filepath.Walk(dir, visit)
}
