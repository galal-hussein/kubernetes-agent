package kubernetesclient

import (
	"crypto/tls"
	"crypto/x509"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"k8s.io/client-go/tools/clientcmd"
)

const (
	caLocation         = "/etc/kubernetes/ssl/ca.pem"
	kubeconfigLocation = "/etc/kubernetes/ssl/kubeconfig"
)

var (
	token  string
	caData []byte
)

func Init() error {
	bytes, err := ioutil.ReadAll(os.Stdin)
	if err != nil {
		return fmt.Errorf("Failed to read token from stdin: %v", err)
	}
	token = strings.TrimSpace(string(bytes))
	if token == "" {
		return errors.New("No token passed in from stdin")
	}

	caData, err = ioutil.ReadFile(caLocation)
	if err != nil {
		return fmt.Errorf("Failed to read CA cert %s: %v", caLocation, err)
	}

	//starting converting
	var kubeconfig string
	kubeconfig = kubeconfigLocation

	// use the current context in kubeconfig
	config, err := clientcmd.BuildConfigFromFlags("", *kubeconfig)
	if err != nil {
		panic(err.Error())
	}
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(err.Error())
	}
	return nil
}

func GetAuthorizationHeader() string {
	return fmt.Sprintf("Bearer %s", token)
}

func GetTLSClientConfig() *tls.Config {
	certPool := x509.NewCertPool()
	certPool.AppendCertsFromPEM(caData)
	return &tls.Config{
		RootCAs: certPool,
	}
}
