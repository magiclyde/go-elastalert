/**
 * Created by GoLand.
 * @author: clyde
 * @date: 2021/10/12 下午3:40
 * @refer: https://elastalert.readthedocs.io/en/latest/elastalert.html
 */

package elastalert

import (
	"fmt"
	"github.com/spf13/viper"
	"log"
)

type DurationStr string // eg "30s" "1h"

type Config struct {
	// EsUrl base URL of form http://ipaddr:port with no trailing slash
	EsUrl string `mapstructure:"es_url"`

	// VerifyCerts whether or not to verify TLS certificates
	VerifyCerts bool `mapstructure:"verify_certs"`

	// CertPem path to a PEM certificate to use as the client certificate
	CertPem string `mapstructure:"cert_pem"`

	// KeyPem path to a private key file to use as the client key
	KeyPem string `mapstructure:"key_pem"`

	// CaCert path to a ca cert
	CaCert string `mapstructure:"ca_cert"`

	// EsUsername basic-auth username for connecting to es_host
	EsUsername string `mapstructure:"es_username"`

	// EsPassword basic-auth password for connecting to es_host
	EsPassword string `mapstructure:"es_password"`

	// EsSendGetBodyAs Method for querying Elasticsearch - GET, POST or source. The default is GET
	EsSendGetBodyAs string `mapstructure:"es_send_get_body_as"`

	// EsConnTimeout timeout in seconds for connecting to and reading from es_host. The default is 20
	EsConnTimeout int `mapstructure:"es_conn_timeout"`

	// RulesLoader the loader to be used by ElastAlert to retrieve rules and hashes. The default is FileRulesLoader
	RulesLoader string `mapstructure:"rules_loader"`

	// RulesFolder the folder which contains rule configuration files
	RulesFolder string `mapstructure:"rules_folder"`

	// ScanSubdirectories whether or not ElastAlert should recursively descend the rules directory. The default is true
	ScanSubdirectories bool `mapstructure:"scan_subdirectories"`

	// BufferTime ElastAlert will continuously query against a window from the present to buffer_time ago.
	// This option is ignored for rules where use_count_query or use_terms_query is set to true. eg "1d2h3m4s"
	BufferTime DurationStr `mapstructure:"buffer_time"`

	// RunEvery How often ElastAlert should query Elasticsearch
	RunEvery DurationStr `mapstructure:"run_every"`

	// WritebackIndex The index on es_host to use, eg elastalert_status
	WritebackIndex string `mapstructure:"writeback_index"`

	WritebackAlias string `mapstructure:"writeback_alias"`

	// MaxQuerySize The maximum number of documents that will be downloaded from Elasticsearch in a single query.
	// The default is 10,000, and if you expect to get near this number, consider using use_count_query for the rule.
	// If this limit is reached, ElastAlert will scroll using the size of max_query_size through the set amount of pages,
	// when max_scrolling_count is set or until processing all results.
	MaxQuerySize int `mapstructure:"max_query_size"`

	// MaxScrollingCount The maximum amount of pages to scroll through. The default is 0, which means the scrolling has no limit
	MaxScrollingCount int `mapstructure:"max_scrolling_count"`

	// ScrollKeepalive The maximum time (formatted in Time Units) the scrolling context should be kept alive.
	// Avoid using high values as it abuses resources in Elasticsearch, but be mindful to allow sufficient time
	// to finish processing all the results. eg 30s
	ScrollKeepalive DurationStr `mapstructure:"scroll_keepalive"`

	// MaxAggregation The maximum number of alerts to aggregate together. The default is 10000
	// If a rule has aggregation set, all alerts occuring within a timeframe will be sent together.
	MaxAggregation int `mapstructure:"max_aggregation"`

	// OldQueryLimit The maximum time between queries for ElastAlert to start at the most recently run query. The default is one week.
	OldQueryLimit DurationStr `mapstructure:"old_query_limit"`

	// AlertTimeLimit the retry window for failed alerts.
	AlertTimeLimit DurationStr `mapstructure:"alert_time_limit"`

	// DisableRulesOnError  This defaults to True
	DisableRulesOnError bool `mapstructure:"disable_rules_on_error"`

	// ShowDisabledRules  This defaults to True.
	ShowDisabledRules bool `mapstructure:"show_disabled_rules"`

	NotifyEmail []string `mapstructure:"notify_email"`

	FromAddr string `mapstructure:"from_addr"`

	SmtpHost string `mapstructure:"smtp_host"`

	EmailReplyTo []string `mapstructure:"email_reply_to"`

	ReplaceDotsInFieldNames bool `mapstructure:"replace_dots_in_field_names"`

	// StringMultiFieldName If set, the suffix to use for the subfield for string multi-fields in Elasticsearch.
	// The default value is .raw for Elasticsearch 2 and .keyword for Elasticsearch 5
	StringMultiFieldName string `mapstructure:"string_multi_field_name"`

	AddMetadataAlert bool `mapstructure:"add_metadata_alert"`

	// SkipInvalid If True, skip invalid files instead of exiting.
	SkipInvalid bool `mapstructure:"skip_invalid"`
}

func NewConfig() *Config {
	load := func(v *viper.Viper, c *Config) error {
		if err := v.Unmarshal(&c); err != nil {
			return fmt.Errorf("unmarshal config err: %s", err.Error())
		}
		//log.Printf("loaded config: %+v", c)
		return nil
	}

	c := &Config{}

	v := viper.New()
	v.SetConfigName("config")
	v.SetConfigType("yaml")
	v.AddConfigPath("/etc/elastalert/")
	v.AddConfigPath(".")

	v.SetDefault("es_send_get_body_as", "GET")
	v.SetDefault("es_conn_timeout", 20)
	v.SetDefault("rules_loader", "FileRulesLoader")
	v.SetDefault("scan_subdirectories", true)
	v.SetDefault("max_query_size", 10000)
	v.SetDefault("max_aggregation", 10000)
	v.SetDefault("old_query_limit", "7d")
	v.SetDefault("disable_rules_on_error", true)

	if err := v.ReadInConfig(); err != nil {
		log.Printf("read config err: %s", err.Error())
	}

	if err := load(v, c); err != nil {
		log.Printf("load config err: %s", err.Error())
	}

	/*go func() {
		v.WatchConfig()
		v.OnConfigChange(func(e fsnotify.Event) {
			log.Printf("config file: %s changed", e.Name)
			err := load(v, c)
			if err != nil {
				log.Printf("reload config err: %s", err.Error())
			}
		})
	}()*/

	return c
}
