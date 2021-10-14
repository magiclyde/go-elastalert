/**
 * Created by GoLand.
 * @author: clyde
 * @date: 2021/10/12 下午3:51
 * @note:
 */

package elastalert

import "time"

type Rule interface {
	GetName() string
	GetType() string
	GetIndex() string
}

type RuleBase struct {
	Name        string      `mapstructure:"name"`
	Typ         string      `mapstructure:"type"`
	Index       string      `mapstructure:"index"`
	Description string      `mapstructure:"description"`
	NumEvents   int         `mapstructure:"num_events"`
	TimeFrame   DurationStr `mapstructure:"timeframe"`
	Filter      interface{} `mapstructure:"filter"`
	Alert       []string    `mapstructure:"alert"`
	Email       []string    `mapstructure:"email"`

	InitialStartTime time.Time `mapstructure:"-"`
}

func (r RuleBase) GetName() string {
	return r.Name
}

func (r RuleBase) GetType() string {
	return r.Typ
}

func (r RuleBase) GetIndex() string {
	return r.Index
}

type RuleCardinality struct {
	RuleBase         `mapstructure:",squash"`
	CardinalityField string `mapstructure:"cardinality_field"` // Count the number of unique values for this field
	MinCardinality   int    `mapstructure:"min_cardinality"`
}

type RuleChange struct {
	RuleBase   `mapstructure:",squash"`
	CompareKey string `mapstructure:"compare_key"` // The field to look for changes in
	IgnoreNull bool   `mapstructure:"ignore_null"` // Ignore documents without the compare_key (country_name) field
	QueryKey   string `mapstructure:"query_key"`   // The change must occur in two documents with the same query_key
}

type RuleFrequency struct {
	RuleBase `mapstructure:",squash"`
}

type RuleNewTerm struct {
	RuleBase        `mapstructure:",squash"`
	Fields          []string    `mapstructure:"fields"`
	TermsWindowSize DurationStr `mapstructure:"terms_window_size"`
}

type RulePercentageMatch struct {
	RuleBase               `mapstructure:",squash"`
	BufferTime             DurationStr `mapstructure:"buffer_time"`
	QueryKey               string      `mapstructure:"query_key"`
	DocType                string      `mapstructure:"doc_type"`
	MinPercentage          int         `mapstructure:"min_percentage"`
	MaxPercentage          int         `mapstructure:"max_percentage"`
	BucketInterval         DurationStr `mapstructure:"bucket_interval"`
	SyncBucketInterval     bool        `mapstructure:"sync_bucket_interval"`
	AllowBufferTimeOverlap bool        `mapstructure:"allow_buffer_time_overlap"`
	UseRunEveryQuerySize   bool        `mapstructure:"use_run_every_query_size"`
}

type RuleMetricAggregation struct {
	RuleBase               `mapstructure:",squash"`
	BufferTime             DurationStr `mapstructure:"buffer_time"`
	MetricAggKey           string      `mapstructure:"metric_agg_key"`
	MetricAggType          string      `mapstructure:"metric_agg_type"`
	QueryKey               string      `mapstructure:"query_key"`
	DocType                string      `mapstructure:"doc_type"`
	BucketInterval         DurationStr `mapstructure:"bucket_interval"`
	SyncBucketInterval     bool        `mapstructure:"sync_bucket_interval"`
	AllowBufferTimeOverlap bool        `mapstructure:"allow_buffer_time_overlap"`
	UseRunEveryQuerySize   bool        `mapstructure:"use_run_every_query_size"`
	MinThreshold           float64     `mapstructure:"min_threshold"`
	MaxThreshold           float64     `mapstructure:"max_threshold"`
}

type RuleSpike struct {
	RuleBase     `mapstructure:",squash"`
	ThresholdCur int    `mapstructure:"threshold_cur"`
	ThresholdRef int    `mapstructure:"threshold_ref"`
	SpikeHeight  int    `mapstructure:"spike_height"`
	SpikeType    string `mapstructure:"spike_type"`
}

type RuleSpikeAggregation struct {
	RuleBase      `mapstructure:",squash"`
	BufferTime    DurationStr `mapstructure:"buffer_time"`
	MetricAggKey  string      `mapstructure:"metric_agg_key"`
	MetricAggType string      `mapstructure:"metric_agg_type"`
	QueryKey      string      `mapstructure:"query_key"`
	DocType       string      `mapstructure:"doc_type"`
	ThresholdCur  int         `mapstructure:"threshold_cur"`
	ThresholdRef  int         `mapstructure:"threshold_ref"`
	SpikeHeight   int         `mapstructure:"spike_height"`
	SpikeType     string      `mapstructure:"spike_type"`
}
