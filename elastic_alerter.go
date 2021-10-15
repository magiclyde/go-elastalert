/**
 * Created by GoLand.
 * @author: clyde
 * @date: 2021/10/12 下午3:45
 * @note:
 */

package elastalert

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"github.com/olivere/elastic/v7"
	"github.com/xhit/go-str2duration/v2"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

type ElasticAlerter struct {
	cfg         *Config
	esClient    *elastic.Client
	rulesLoader RulesLoader
	rules       []Rule
	startTime   time.Time
}

func NewElasticAlerter(cfg *Config) *ElasticAlerter {
	e := &ElasticAlerter{cfg: cfg, startTime: time.Now()}
	e.init()

	return e
}

func (e *ElasticAlerter) init() {
	e.initEsClient()
	e.initRules()
}

func (e *ElasticAlerter) getHttpClient() *http.Client {
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}

	if e.cfg.VerifyCerts {
		certPEMBlock, err := ioutil.ReadFile(e.cfg.CertPem)
		if err != nil {
			log.Fatalf("read cert err: %s", err.Error())
		}
		keyPEMBlock, err := ioutil.ReadFile(e.cfg.KeyPem)
		if err != nil {
			log.Fatalf("read cert key err: %s", err.Error())
		}
		cert, err := tls.X509KeyPair(certPEMBlock, keyPEMBlock)
		if err != nil {
			log.Fatalf("tls.X509KeyPair err: %s", err.Error())
		}

		caCert, err := ioutil.ReadFile(e.cfg.CaCert)
		if err != nil {
			log.Fatalf("read ca cert err: %s", err.Error())
		}
		caCertPool := x509.NewCertPool()
		caCertPool.AppendCertsFromPEM(caCert)

		tr.TLSClientConfig = &tls.Config{
			Certificates: []tls.Certificate{cert},
			RootCAs:      caCertPool,
		}
	}

	hc := &http.Client{
		Transport: tr,
		Timeout:   time.Second * time.Duration(e.cfg.EsConnTimeout),
	}
	return hc
}

func (e *ElasticAlerter) initEsClient() {
	opts := []elastic.ClientOptionFunc{
		elastic.SetURL(e.cfg.EsUrl),
		elastic.SetSniff(false),
		elastic.SetHttpClient(e.getHttpClient()),
		elastic.SetSendGetBodyAs(e.cfg.EsSendGetBodyAs),
	}
	if e.cfg.EsUsername != "" || e.cfg.EsPassword != "" {
		opts = append(opts, elastic.SetBasicAuth(e.cfg.EsUsername, e.cfg.EsPassword))
	}

	client, err := elastic.NewClient(opts...)
	if err != nil {
		log.Fatalf("elastic.NewClient err: %s", err.Error())
	}
	e.esClient = client
}

func (e *ElasticAlerter) initRules() {
	switch e.cfg.RulesLoader {
	case "FileRulesLoader":
		e.rulesLoader = NewFileRulesLoader(e.cfg.RulesFolder, SetDescend(e.cfg.ScanSubdirectories))
	default:
		log.Fatalf("rules loader: %s not supported", e.cfg.RulesLoader)
	}

	e.rules = e.rulesLoader.Load()
	log.Printf("%d rules loaded", len(e.rules))
}

func (e *ElasticAlerter) Run(ctx context.Context) {
	duration, err := str2duration.ParseDuration(string(e.cfg.RunEvery))
	if err != nil {
		log.Fatalf("str2duration.ParseDuration err: %s", err.Error())
	}
	log.Printf("run every %+v", duration)

	ticker := time.NewTicker(duration)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			log.Fatalf("run cancelled")

		case <-ticker.C:
			for _, rule := range e.rules {
				rule.SetInitialStartTime(e.startTime)

				switch rule.GetType() {
				case "cardinality":
					cardinality, ok := rule.(RuleCardinality)
					if ok {
						e.runCardinality(ctx, cardinality)
					}

				case "change":
					change, ok := rule.(RuleChange)
					if ok {
						e.runChange(ctx, change)
					}

				case "frequency":
					frequency, ok := rule.(RuleFrequency)
					if ok {
						e.runFrequency(ctx, frequency)
					}

				case "new_term":
					newTerm, ok := rule.(RuleNewTerm)
					if ok {
						e.runNewTerm(ctx, newTerm)
					}

				case "percentage_match":
					percentageMatch, ok := rule.(RulePercentageMatch)
					if ok {
						e.runPercentageMatch(ctx, percentageMatch)
					}

				case "metric_aggregation":
					metricAggregation, ok := rule.(RuleMetricAggregation)
					if ok {
						e.runMetricAggregation(ctx, metricAggregation)
					}

				case "spike_aggregation":
					spikeAggregation, ok := rule.(RuleSpikeAggregation)
					if ok {
						e.runSpikeAggregation(ctx, spikeAggregation)
					}

				case "spike":
					spike, ok := rule.(RuleSpike)
					if ok {
						e.runSpike(ctx, spike)
					}
				}

			}
		}
	}
}

func (e *ElasticAlerter) runCardinality(ctx context.Context, rl RuleCardinality) {
	log.Println("runCardinality...")
}

func (e *ElasticAlerter) runChange(ctx context.Context, rl RuleChange) {
	log.Println("runChange...")
}

func (e *ElasticAlerter) runFrequency(ctx context.Context, rl RuleFrequency) {
	log.Println("runFrequency...")
}

func (e *ElasticAlerter) runNewTerm(ctx context.Context, rl RuleNewTerm) {
	log.Println("runNewTerm...")
}

func (e *ElasticAlerter) runPercentageMatch(ctx context.Context, rl RulePercentageMatch) {
	log.Println("runPercentageMatch...")
}

func (e *ElasticAlerter) runMetricAggregation(ctx context.Context, rl RuleMetricAggregation) {
	log.Println("runMetricAggregation...")
}

func (e *ElasticAlerter) runSpikeAggregation(ctx context.Context, rl RuleSpikeAggregation) {
	log.Println("runSpikeAggregation...")
}

func (e *ElasticAlerter) runSpike(ctx context.Context, rl RuleSpike) {
	log.Println("runSpike...")
}
