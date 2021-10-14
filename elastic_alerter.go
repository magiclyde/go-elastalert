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
}

func NewElasticAlerter(cfg *Config) *ElasticAlerter {
	e := &ElasticAlerter{cfg: cfg}
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

	if err := e.rulesLoader.Load(); err != nil {
		log.Fatalf("e.rulesLoader.Load() err: %s", err.Error())
	}
	e.rules = e.rulesLoader.GetRules()
}

func (e *ElasticAlerter) Run(ctx context.Context) {
	duration, err := str2duration.ParseDuration(e.cfg.RunEvery)
	if err != nil {
		log.Fatalf("str2duration.ParseDuration err: %s, run_every: %s", err.Error(), e.cfg.RunEvery)
	}
	ticker := time.NewTicker(duration)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			log.Println("tick...")
			for _, rule := range e.rules {
				log.Printf("apply rule, name: %s, type: %s, index: %s", rule.Name, rule.Typ, rule.Index)
				// todo...
			}

		case <-ctx.Done():
			log.Fatalf("run cancelled")
		}
	}
}
