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
	"reflect"
	"sync"
)

type RulesLoader interface {
	Load() []Rule
}

type FileRulesLoaderOption func(*FileRulesLoader)

type FileRulesLoader struct {
	Path        string
	Suffix      string
	Descend     bool
	ruleTypeMap map[string]Rule
}

func NewFileRulesLoader(path string, options ...FileRulesLoaderOption) *FileRulesLoader {
	l := &FileRulesLoader{
		Path:        path,
		Suffix:      "yaml",
		Descend:     true,
		ruleTypeMap: make(map[string]Rule),
	}

	for _, f := range options {
		f(l)
	}

	l.ruleTypeMap = map[string]Rule{
		"cardinality":        RuleCardinality{},
		"change":             RuleChange{},
		"frequency":          RuleFrequency{},
		"new_term":           RuleNewTerm{},
		"percentage_match":   RulePercentageMatch{},
		"metric_aggregation": RuleMetricAggregation{},
		"spike_aggregation":  RuleSpikeAggregation{},
		"spike":              RuleSpike{},
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

func (l *FileRulesLoader) Load() []Rule {
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

			typ := runtimeViper.GetString("type")
			ruleObj, ok := l.ruleTypeMap[typ]
			if !ok {
				log.Printf("unsupported type: %s", typ)
				return
			}
			val := reflect.New(reflect.TypeOf(ruleObj))
			if err := runtimeViper.Unmarshal(val.Interface()); err != nil {
				log.Printf("Unmarshal err: %s from file: %s", err.Error(), file)
				return
			}
			rules[i] = val.Elem().Interface().(Rule)

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
