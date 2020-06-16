package main

import (
    "awesomeProject/controller"
    "flag"
    "fmt"
    "k8s.io/client-go/kubernetes"
    "k8s.io/client-go/tools/clientcmd"
    "k8s.io/klog/v2"
    "os"
)

var (
    masterURL  string
    kubeconfig string
)

func main() {
    klog.InitFlags(nil)
    flag.Parse()

    stop := make(chan struct{})
    go func() {
        <-stop
        close(stop)
        os.Exit(1)
    }()

    // If you are running from command line then uncomment line with (masterURL, kubeconfig)
    // it and comment then line with ("", "")
    // >> go run main.go -kubeconfig=<your home dir>/.kube/config
    cfg, err := clientcmd.BuildConfigFromFlags(masterURL, kubeconfig)

    // cfg, err := clientcmd.BuildConfigFromFlags("", "")
    if err != nil {
        fmt.Print("Error to build config")
        return
    }

    client, err := kubernetes.NewForConfig(cfg)
    if err != nil {
        fmt.Print("failed to build client from config")
        return
    }

    // Make a controller and run ti
    controllerObject := controller.NewController(client, stop)
    controllerObject.Run(stop)

    // Should use a channel to kill this controller - just a easy hack for now using sleep
    // time.Sleep(1 * time.Hour)
}

func init() {
    flag.StringVar(&kubeconfig, "kubeconfig", "/Users/harish.bohara/.kube/config", "Path to a kubeconfig. Only required if out-of-cluster.")
    flag.StringVar(&masterURL, "master", "", "The address of the Kubernetes API server. Overrides any value in kubeconfig. Only required if out-of-cluster.")
}
