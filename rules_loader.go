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
	rules       []Rule
	loaded      bool
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

// SetSuffix set the suffix used by the FileRulesLoader.
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
	if l.loaded {
		return l.rules
	}

	files := WalkDir(l.Path, l.Suffix, l.Descend)
	rules := l.loadRule(files)
	for rule := range rules {
		l.rules = append(l.rules, rule)
	}
	l.loaded = true

	return l.rules
}

func (l *FileRulesLoader) loadRule(in <-chan string) <-chan Rule {
	out := make(chan Rule, cap(in))
	go func() {
		for path := range in {
			runtimeViper := viper.New()
			runtimeViper.SetConfigFile(path)
			if err := runtimeViper.ReadInConfig(); err != nil {
				log.Printf("ReadInConfig err: %s from : %s", err.Error(), path)
				continue
			}
			typ := runtimeViper.GetString("type")
			ruleObj, ok := l.ruleTypeMap[typ]
			if !ok {
				log.Printf("unsupported type: %s", typ)
				continue
			}
			val := reflect.New(reflect.TypeOf(ruleObj))
			if err := runtimeViper.Unmarshal(val.Interface()); err != nil {
				log.Printf("Unmarshal err: %s from file: %s", err.Error(), path)
				continue
			}
			out <- val.Elem().Interface().(Rule)
		}
		close(out)
	}()
	return out
}
