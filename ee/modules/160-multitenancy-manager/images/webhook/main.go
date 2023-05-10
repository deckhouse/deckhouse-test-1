/*
Copyright 2023 Flant JSC
Licensed under the Deckhouse Platform Enterprise Edition (EE) license. See https://github.com/deckhouse/deckhouse/blob/main/ee/LICENSE
*/

package main

import (
	"flag"
	"os"

	"k8s.io/apimachinery/pkg/runtime"
	utilruntime "k8s.io/apimachinery/pkg/util/runtime"
	clientgoscheme "k8s.io/client-go/kubernetes/scheme"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/healthz"
	"sigs.k8s.io/controller-runtime/pkg/log/zap"

	"github.com/deckhouse/deckhouse/ee/modules/160-multitenancy-manager/hooks/apis/deckhouse.io/v1alpha1"
)

var (
	scheme = runtime.NewScheme()
)

func init() {
	utilruntime.Must(clientgoscheme.AddToScheme(scheme))

	utilruntime.Must(v1alpha1.AddToScheme(scheme))
}

func main() {
	var certsDir string
	var metricsAddr string
	var healtzAddr string
	var webhookHost string
	var webhookPort int

	flag.StringVar(&certsDir, "certs.dir", "", "The directory that contains the server key and certificate for webhook.")
	flag.StringVar(&metricsAddr, "metrics.address", "127.0.0.1:8080", "The address the metric endpoint binds to.")
	flag.StringVar(&healtzAddr, "healthz.address", ":8081", "The address the probe endpoint binds to.")
	flag.StringVar(&webhookHost, "webhook.host", "0.0.0.0", "The hist the webhook endpoints bind to.")
	flag.IntVar(&webhookPort, "webhook.port", 9443, "The port the webhook endpoints bind to.")
	opts := zap.Options{
		Development: true,
	}
	opts.BindFlags(flag.CommandLine)
	flag.Parse()

	setupLog := zap.New(zap.UseFlagOptions(&opts))
	ctrl.SetLogger(setupLog.WithName("controller-runtime"))

	mgr, err := ctrl.NewManager(ctrl.GetConfigOrDie(), ctrl.Options{
		Scheme:                 scheme,
		HealthProbeBindAddress: healtzAddr,
		MetricsBindAddress:     metricsAddr,
		Port:                   webhookPort,
		Host:                   webhookHost,
		// Logger:                 ctrl.Log.WithName("webhook"),
		CertDir: certsDir,
	})
	if err != nil {
		setupLog.Error(err, "unable to start manager")
		os.Exit(1)
	}

	if err = (&v1alpha1.ProjectType{}).SetupWebhookWithManager(mgr); err != nil {
		setupLog.Error(err, "unable to create webhook", "webhook", "ProjectType")
		os.Exit(1)
	}

	if err := mgr.AddHealthzCheck("healthz", healthz.Ping); err != nil {
		setupLog.Error(err, "unable to set up health check")
		os.Exit(1)
	}

	setupLog.Info("starting manager")
	if err := mgr.Start(ctrl.SetupSignalHandler()); err != nil {
		setupLog.Error(err, "problem running manager")
		os.Exit(1)
	}
}
