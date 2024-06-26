package main

import (
	"os"
	"strings"
	"time"

	"github.com/fsnotify/fsnotify"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

var (
	lg = logrus.New()
)

func setupConfig() *viper.Viper {
	cfg := viper.New()
	cfg.AddConfigPath(".")
	cfg.AddConfigPath("$HOME/.config")
	cfg.AddConfigPath("/etc/ca-injector")
	cfg.SetConfigName("config")

	cfg.AutomaticEnv()
	cfg.SetEnvPrefix("cainjector") // will be uppercased automatically
	cfg.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	cfg.SetDefault("tls.key", "tls.key")
	cfg.SetDefault("tls.crt", "tls.crt")
	cfg.SetDefault("caBundle.annotation", "ca-injector.zeiss.com/inject-ca-from")
	cfg.SetDefault("caBundle.configMap", "")
	cfg.SetDefault("caBundle.crt", "ca.crt")
	cfg.SetDefault("admissionWebhook.labelSelector", "ca-injector.zeiss.com/inject")
	cfg.SetDefault("admissionWebhook.enableNamespacesByDefault", false)
	cfg.SetDefault("admissionWebhook.ignoreNamespaces", []string{"kube-system", "kube-public", "kube-node-lease", "local-path-storage"})

	cfg.SetDefault("shutdown.timeout", 10*time.Second)

	if err := cfg.ReadInConfig(); err != nil {
		lg.WithError(err).Error("could not read initial config")
	}

	cfg.OnConfigChange(func(_ fsnotify.Event) {
		if err := cfg.ReadInConfig(); err != nil {
			lg.WithError(err).Warn("could not reload config")
		}
	})
	if os.Getenv("KUBERNETES_SERVICE_PORT") != "" {
		lg.SetFormatter(&logrus.JSONFormatter{})
	}

	go cfg.WatchConfig()

	return cfg
}
