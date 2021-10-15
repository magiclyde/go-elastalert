/**
 * Created by GoLand.
 * @author: clyde
 * @date: 2021/10/13 下午1:42
 * @note:
 */

package elastalert

import (
	"github.com/spf13/viper"
	"log"
	"sync"
)

type RulesLoader interface {
	GetRules() []Rule
}

type FileRulesLoaderOption func(*FileRulesLoader)

type FileRulesLoader struct {
	Path    string
	Suffix  string
	Descend bool
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

func (l *FileRulesLoader) GetRules() []Rule {
	files, err := WalkDir(l.Path, l.Suffix, l.Descend)
	if err != nil {
		log.Printf("WalkDir err: %s", err.Error())
		return nil
	}

	rules := make([]Rule, len(files))

	maxGoroutines := 16
	guard := make(chan struct{}, maxGoroutines)

	wg := sync.WaitGroup{}
	wg.Add(len(files))

	for i, file := range files {
		guard <- struct{}{} // would block if guard channel is already filled
		go func(i int, file string) {
			defer func() {
				wg.Done()
				<-guard
			}()

			runtimeViper := viper.New()
			runtimeViper.SetConfigFile(file)

			if err := runtimeViper.ReadInConfig(); err != nil {
				log.Printf("ReadInConfig err: %s from file: %s", err.Error(), file)
				return
			}

			switch runtimeViper.Get("type") {
			case "cardinality":
				var rule RuleCardinality
				if err := runtimeViper.Unmarshal(&rule); err != nil {
					log.Printf("Unmarshal err: %s from file: %s", err.Error(), file)
					return
				}
				rules[i] = rule

			case "change":
				var rule RuleChange
				if err := runtimeViper.Unmarshal(&rule); err != nil {
					log.Printf("Unmarshal err: %s from file: %s", err.Error(), file)
					return
				}
				rules[i] = rule

			case "frequency":
				var rule RuleFrequency
				if err := runtimeViper.Unmarshal(&rule); err != nil {
					log.Printf("Unmarshal err: %s from file: %s", err.Error(), file)
					return
				}
				rules[i] = rule

			case "new_term":
				var rule RuleNewTerm
				if err := runtimeViper.Unmarshal(&rule); err != nil {
					log.Printf("Unmarshal err: %s from file: %s", err.Error(), file)
					return
				}
				rules[i] = rule

			case "percentage_match":
				var rule RulePercentageMatch
				if err := runtimeViper.Unmarshal(&rule); err != nil {
					log.Printf("Unmarshal err: %s from file: %s", err.Error(), file)
					return
				}
				rules[i] = rule

			case "metric_aggregation":
				var rule RuleMetricAggregation
				if err := runtimeViper.Unmarshal(&rule); err != nil {
					log.Printf("Unmarshal err: %s from file: %s", err.Error(), file)
					return
				}
				rules[i] = rule

			case "spike_aggregation":
				var rule RuleSpikeAggregation
				if err := runtimeViper.Unmarshal(&rule); err != nil {
					log.Printf("Unmarshal err: %s from file: %s", err.Error(), file)
					return
				}
				rules[i] = rule

			case "spike":
				var rule RuleSpike
				if err := runtimeViper.Unmarshal(&rule); err != nil {
					log.Printf("Unmarshal err: %s from file: %s", err.Error(), file)
					return
				}
				rules[i] = rule
			}

		}(i, file)
	}

	wg.Wait()

	// rm nil value from slice without allocating a new slice
	for i := 0; i < len(rules); {
		if rules[i] != nil {
			i++
			continue
		}
		if i < len(rules)-1 {
			copy(rules[i:], rules[i+1:])
		}
		rules[len(rules)-1] = nil
		rules = rules[:len(rules)-1]
	}

	return rules
}
