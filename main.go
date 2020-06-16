package main

import (
    "awesomeProject/controller"
    "flag"
    "fmt"
    "k8s.io/client-go/kubernetes"
    "k8s.io/client-go/tools/clientcmd"
    "k8s.io/klog/v2"
    "time"
)

var (
    masterURL  string
    kubeconfig string
)

func main() {
    klog.InitFlags(nil)
    flag.Parse()

    // cfg, err := clientcmd.BuildConfigFromFlags(masterURL, kubeconfig)
    cfg, err := clientcmd.BuildConfigFromFlags("", "")
    if err != nil {
        fmt.Print("Error to build config")
        return
    }

    client, err := kubernetes.NewForConfig(cfg)
    if err != nil {
        fmt.Print("failed to build client from config")
        return
    }

    controllerObject := controller.NewController(client)

    controllerObject.Run()


    var _ = client
    time.Sleep(1 * time.Hour)
}

func init() {
    flag.StringVar(&kubeconfig, "kubeconfig", "/Users/harish.bohara/.kube/config", "Path to a kubeconfig. Only required if out-of-cluster.")
    flag.StringVar(&masterURL, "master", "", "The address of the Kubernetes API server. Overrides any value in kubeconfig. Only required if out-of-cluster.")
}
