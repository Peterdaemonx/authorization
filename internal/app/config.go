package app

import (
	"fmt"
	"os"
	"reflect"
	"time"

	"gitlab.cmpayments.local/creditcard/authorization/internal/infrastructure/connection"
	"gopkg.in/yaml.v3"
)

type Config struct {
	Listen     string `yaml:"listen"`
	CmPlatform struct {
		BaseDomain string `yaml:"base_domain"`
	} `yaml:"cm_platform"`
	GCP struct {
		ProjectID string `yaml:"project_id"`
		Spanner   struct {
			Instance     string        `yaml:"instance"`
			Database     string        `yaml:"database"`
			PoolSize     int           `yaml:"pool_size"`
			ReadTimeout  time.Duration `yaml:"read_timeout"`
			WriteTimeout time.Duration `yaml:"write_timeout"`
		} `yaml:"spanner"`
		PubSub struct {
			PresentmentClearedTopicID    string        `yaml:"presentment_cleared_topic_id"`
			AuthorizationCapturedTopicID string        `yaml:"authorization_captured_topic_id"`
			RefundCapturedTopicID        string        `yaml:"refund_captured_topic_id"`
			Timeout                      time.Duration `yaml:"timeout"`
		} `yaml:"pubsub"`
		Storage struct {
			BucketName string `yaml:"bucket_name"`
		} `yaml:"storage"`
	} `yaml:"gcp"`
	Development struct {
		MockCmPlatform       bool   `yaml:"mock_cm_platform"`
		MockPermissionStore  bool   `yaml:"mock_permission_store"`
		SpannerEmulatorAddr  string `yaml:"spanner_emulator_addr"`
		HumanReadableLogging bool   `yaml:"human_readable_logging"`
		MockData             bool   `yaml:"mock_data"`
		MockTokenization     bool   `yaml:"mock_tokenization"`
		MockCardInfo         bool   `yaml:"mock_card_info"`
		MockPublisher        bool   `yaml:"mock_publisher"`
		MockPubSub           bool   `yaml:"mock_pub_sub"`
		MockNonceStore       bool   `yaml:"mock_nonce_store"`
		MockStorageBucket    bool   `yaml:"mock_storage_bucket"`
	} `yaml:"development"`
	Auth struct {
		BaseURL   string `yaml:"baseURL"`
		PublicKey string `yaml:"publicKey"`
	} `yaml:"auth"`
	BlockedBins []string `yaml:"blocked_bins"`
	MasterCard  struct {
		ConnectionPool   connection.PoolConfiguration `yaml:"connection_pool"`
		BinrangeFiletype string                       `yaml:"binrange_filetype"`
	} `yaml:"mastercard"`
	Visa struct {
		ConnectionPool   connection.PoolConfiguration `yaml:"connection_pool"`
		SourceStationID  string                       `yaml:"source_station_id"`
		BinrangeFiletype string                       `yaml:"binrange_filetype"`
		AddTestPans      bool                         `yaml:"add_test_pans"`
	} `yaml:"visa"`
	Cors struct {
		AllowedOrigins string `yaml:"allowed_origins"`
	}
	Tokenization struct {
		BaseURL string `yaml:"base_url"`
	} `yaml:"tokenization"`
	Detokenization struct {
		BaseURL string `yaml:"base_url"`
	} `yaml:"detokenization"`
	CardInfoApi struct {
		BaseURL string `yaml:"base_url"`
	} `yaml:"card_info_api"`
	JWT                        string `yaml:"jwt"`
	MinLogLevel                string `yaml:"min_log_level"`
	AllowProductionCardNumbers bool   `yaml:"allow_production_card_numbers"`
}

func LoadConfig(path string, conf interface{}) error {
	if reflect.TypeOf(conf).Kind() != reflect.Ptr {
		// Panic, because this is effectively a compile-time error
		panic("LoadConfig: conf must be a pointer")
	}
	b, err := os.ReadFile(path)
	if err != nil {
		return fmt.Errorf("cannot read file: %w", err)
	}

	if err = yaml.Unmarshal(b, conf); err != nil {
		return fmt.Errorf("cannot parse yaml: %w", err)
	}

	return nil
}
