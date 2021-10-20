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
