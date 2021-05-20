package gkafka

import (
	"crypto/tls"
	"crypto/x509"
	"errors"
	"fmt"
	"github.com/Shopify/sarama"
	"io/ioutil"
)

type Engine struct {
	config   *KafkaConfig
	dataChan chan []byte
}

type KafkaConfig struct {
	Brokers      []string `json:"brokers"`
	Topic        string   `json:"topic"`
	Group        string   `json:"group"`
	User         string   `json:"user"`
	Pwd          string   `json:"pwd"`
	CertFile     string   `json:"cert_file"`
	KeyFile      string   `json:"key_file"`
	CaFile       string   `json:"ca_file"`
	KafkaVersion string   `json:"kafka_version"`
	VerifySsl    bool     `json:"verify_ssl"`
}

var Default = &Engine{}

func getSaramaConfig(kc *KafkaConfig) (*sarama.Config, error) {
	var err error
	config := sarama.NewConfig()
	if kc.KafkaVersion != "" {
		config.Version, err = sarama.ParseKafkaVersion(kc.KafkaVersion)
		if nil != err {
			return nil, errors.New(fmt.Sprintf("ParseKafkaVersion, err=%s", err.Error()))
		}
	}

	if kc.VerifySsl {
		config.Net.SASL.Enable = true
		config.Net.SASL.User = kc.User
		config.Net.SASL.Password = kc.Pwd
		config.Net.SASL.Handshake = true
		config.Net.TLS.Enable = true
		var tlsConfig = &tls.Config{
			InsecureSkipVerify: true,
		}

		if kc.KeyFile != "" && kc.CertFile != "" {
			cert, err := tls.LoadX509KeyPair(kc.CertFile, kc.KeyFile)
			if err != nil {
				return nil, errors.New(fmt.Sprintf("CertFile or KeyFile is fail, err=%s", err.Error()))
			}
			tlsConfig.Certificates = []tls.Certificate{cert}
		}

		certBytes, err := ioutil.ReadFile(kc.CaFile)
		if nil != err {
			return nil, errors.New(fmt.Sprintf("CaFile is fail, err=%s", err.Error()))
		}
		clientCertPool := x509.NewCertPool()
		ok := clientCertPool.AppendCertsFromPEM(certBytes)
		if !ok {
			return nil, errors.New("AppendCertsFromPEM fail")
		}
		tlsConfig.RootCAs = clientCertPool
		config.Net.TLS.Config = tlsConfig
	}

	config.Producer.Partitioner = sarama.NewRandomPartitioner
	config.Producer.Return.Errors = true
	config.Producer.RequiredAcks = sarama.WaitForAll
	config.Producer.Return.Successes = true
	config.Consumer.Group.Rebalance.Strategy = sarama.BalanceStrategyRoundRobin
	config.Consumer.Offsets.Initial = sarama.OffsetOldest
	if err = config.Validate(); nil != err {
		return nil, errors.New(fmt.Sprintf("Default config.Validate(), err=%s", err.Error()))
	}
	return config, nil
}

func InitKafka(kc *KafkaConfig) *Engine {
	Default.config = kc
	return Default
}
