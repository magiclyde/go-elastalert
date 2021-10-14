/**
 * Created by GoLand.
 * @author: clyde
 * @date: 2021/10/12 下午3:51
 * @note:
 */

package elastalert

type Rule struct {
	Name      string `mapstructure:"name"`
	Typ       string `mapstructure:"type"`
	Index     string `mapstructure:"index"`
	NumEvents int    `mapstructure:"num_events"`
}
