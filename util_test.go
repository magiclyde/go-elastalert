/**
 * Created by GoLand.
 * @author: clyde
 * @date: 2021/10/13 下午4:24
 * @note:
 */

package elastalert

import "testing"

func TestWalkDir(t *testing.T) {
	files, err := WalkDir("/tmp/", "yaml", true)
	t.Logf("file: %+v, err: %+v", files, err)
}
