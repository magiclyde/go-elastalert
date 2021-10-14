/**
 * Created by GoLand.
 * @author: clyde
 * @date: 2021/10/13 下午4:03
 * @note:
 */

package elastalert

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

// WalkDir
func WalkDir(dirPth, suffix string, descend bool) (files []string, err error) {
	suffix = strings.ToLower(suffix)

	if descend {
		err = filepath.Walk(dirPth, func(filename string, fi os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			if fi.IsDir() {
				return nil
			}
			if strings.HasSuffix(strings.ToLower(fi.Name()), suffix) {
				files = append(files, filename)
			}
			return nil
		})
		return
	}

	fis, err := ioutil.ReadDir(dirPth)
	if err != nil {
		return nil, err
	}
	for _, fi := range fis {
		if fi.IsDir() {
			continue
		}
		if strings.HasSuffix(strings.ToLower(fi.Name()), suffix) {
			files = append(files, filepath.Join(dirPth, fi.Name()))
		}
	}
	return
}
