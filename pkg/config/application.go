package config

import (
	"github.com/Coinoner/trader_provisioner/pkg/models"
	"gorm.io/gorm"
	"k8s.io/client-go/kubernetes"
	"log"
	"os"
	"strconv"
)

var ApplicationConfig Config

type Config struct {
	traderImage        string
	traderPort         int32
	traderPrefix       string
	traderIngressName  string
	imagePullSecrets   string
	maxTraderPerUser   int64
	clusterCertificate string
	hostname            string
	db                  DB
	KubernetesClientSet *kubernetes.Clientset
}

func InitializeConfig() {
	log.Printf("configuration initialized")
	port, _ := strconv.ParseInt(os.Getenv("TRADER_PORT"), 10, 32)
	ApplicationConfig = Config{
		traderImage:       os.Getenv("TRADER_IMAGE"),
		traderPrefix:      os.Getenv("TRADER_PREFIX"),
		traderPort:        int32(port),
		traderIngressName: os.Getenv("TRADER_INGRESS_NAME"),
		imagePullSecrets:  os.Getenv("IMAGE_PULL_SECRETS"),
		clusterCertificate: os.Getenv("CLUSTER_CERTIFICATE"),
		hostname:          os.Getenv("HOSTNAME"),
		maxTraderPerUser:  100,
	}

	ApplicationConfig.db = DatabaseInitialize()
	_ = ApplicationConfig.db.instance.AutoMigrate(&models.Trader{})
}

func (c *Config) GetDbInstance() *gorm.DB {
	return c.db.instance
}

func (c *Config) GetTraderPrefix() string {
	return c.traderPrefix
}

func (c *Config) GetTraderImage() string {
	return c.traderPrefix
}

func (c *Config) GetTraderPort() int32 {
	return c.traderPort
}

func (c *Config) GetMaxTraderPerUser() int64 {
	return c.maxTraderPerUser
}

func (c *Config) GetImagePullSecrets() string {
	return c.imagePullSecrets
}

func (c *Config) GetClusterCertificate() string {
	return c.clusterCertificate
}

func (c *Config) GetHostname() string {
	return c.hostname
}