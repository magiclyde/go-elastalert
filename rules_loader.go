/**
 * Created by GoLand.
 * @author: clyde
 * @date: 2021/10/13 下午1:42
 * @note:
 */

package elastalert

import (
	"bytes"
	"github.com/spf13/viper"
	"io/ioutil"
	"log"
)

type RulesLoader interface {
	Load() error
	GetRules() []Rule
}

type FileRulesLoaderOption func(*FileRulesLoader)

type FileRulesLoader struct {
	Path    string
	Suffix  string
	Descend bool
	rules   []Rule
}

func NewFileRulesLoader(path string, options ...FileRulesLoaderOption) *FileRulesLoader {
	l := &FileRulesLoader{
		Path:    path,
		Suffix:  "yaml",
		Descend: true,
	}

	for _, f := range options {
		f(l)
	}

	return l
}

// SetSuffix sets the suffix used by the FileRulesLoader.
func SetSuffix(s string) FileRulesLoaderOption {
	return func(l *FileRulesLoader) {
		l.Suffix = s
	}
}

// SetDescend recursively descend the rules directory
func SetDescend(d bool) FileRulesLoaderOption {
	return func(l *FileRulesLoader) {
		l.Descend = d
	}
}

func (l *FileRulesLoader) Load() error {
	files, err := WalkDir(l.Path, l.Suffix, l.Descend)
	if err != nil {
		return err
	}

	viper.SetConfigType(l.Suffix)

	for _, file := range files {
		block, err := ioutil.ReadFile(file)
		if err != nil {
			log.Printf("ioutil.ReadFile err: %s", err.Error())
			continue
		}
		if err := viper.ReadConfig(bytes.NewBuffer(block)); err != nil {
			log.Printf("viper.ReadConfig err: %s", err.Error())
			continue
		}

		switch viper.Get("type") {
		case "cardinality":
			var rule RuleCardinality
			if err := viper.Unmarshal(&rule); err != nil {
				log.Printf("viper.Unmarshal err: %s", err.Error())
				continue
			}
			l.rules = append(l.rules, rule)

		case "change":
			var rule RuleChange
			if err := viper.Unmarshal(&rule); err != nil {
				log.Printf("viper.Unmarshal err: %s", err.Error())
				continue
			}
			l.rules = append(l.rules, rule)

		case "frequency":
			var rule RuleFrequency
			if err := viper.Unmarshal(&rule); err != nil {
				log.Printf("viper.Unmarshal err: %s", err.Error())
				continue
			}
			l.rules = append(l.rules, rule)

		case "new_term":
			var rule RuleNewTerm
			if err := viper.Unmarshal(&rule); err != nil {
				log.Printf("viper.Unmarshal err: %s", err.Error())
				continue
			}
			l.rules = append(l.rules, rule)

		case "percentage_match":
			var rule RulePercentageMatch
			if err := viper.Unmarshal(&rule); err != nil {
				log.Printf("viper.Unmarshal err: %s", err.Error())
				continue
			}
			l.rules = append(l.rules, rule)

		case "metric_aggregation":
			var rule RuleMetricAggregation
			if err := viper.Unmarshal(&rule); err != nil {
				log.Printf("viper.Unmarshal err: %s", err.Error())
				continue
			}
			l.rules = append(l.rules, rule)

		case "spike_aggregation":
			var rule RuleSpikeAggregation
			if err := viper.Unmarshal(&rule); err != nil {
				log.Printf("viper.Unmarshal err: %s", err.Error())
				continue
			}
			l.rules = append(l.rules, rule)

		case "spike":
			var rule RuleSpike
			if err := viper.Unmarshal(&rule); err != nil {
				log.Printf("viper.Unmarshal err: %s", err.Error())
				continue
			}
			l.rules = append(l.rules, rule)
		}
	}

	return nil
}

func (l *FileRulesLoader) GetRules() []Rule {
	return l.rules
}
