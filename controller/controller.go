package controller

import (
    "context"
    "fmt"
    metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
    "k8s.io/apimachinery/pkg/util/runtime"
    "k8s.io/client-go/informers"
    "k8s.io/client-go/kubernetes"
    "k8s.io/client-go/tools/cache"
    "k8s.io/client-go/util/workqueue"
    "time"
)

type Controller struct {
    workQueue           workqueue.RateLimitingInterface
    KubeClientSet       kubernetes.Interface
    KubeInformerFactory informers.SharedInformerFactory
}

func NewController(KubeClientSet kubernetes.Interface) *Controller {

    controller := &Controller{
        KubeClientSet:       KubeClientSet,
        KubeInformerFactory: informers.NewSharedInformerFactory(KubeClientSet, 1*time.Second),
        workQueue:           workqueue.NewNamedRateLimitingQueue(workqueue.DefaultControllerRateLimiter(), "my-controller-rate-limit"),
    }

    controller.KubeInformerFactory.Apps().V1().Deployments().Informer().AddEventHandler(
        cache.ResourceEventHandlerFuncs{
            AddFunc: func(obj interface{}) {

            },
            UpdateFunc: func(oldObj, newObj interface{}) {
                // fmt.Print("UpdateFunc - ", newObj)
                var key string
                var err error
                if key, err = cache.MetaNamespaceKeyFunc(newObj); err != nil {
                    runtime.HandleError(err)
                    return
                }
                controller.workQueue.Add(key)
            },
        },
    );
    stop := make(chan struct{})
    controller.KubeInformerFactory.Start(stop)

    return controller;
}

func (c *Controller) Run() {
    go func() {
        c.runWorker()
    }()
}

func (c *Controller) runWorker() {
    defer runtime.HandleCrash()
    defer c.workQueue.ShutDown()

    for c.processNextWorkItem() {
    }
}

func (c *Controller) processNextWorkItem() bool {
    // fmt.Println("processNextWorkItem -->> ")

    obj, shutdown := c.workQueue.Get()
    if shutdown {
        return false
    }

    defer c.workQueue.Done(obj)

    //fmt.Println(obj)

    var key string
    var ok bool
    if key, ok = obj.(string); ok {
        namespace, name, _ := cache.SplitMetaNamespaceKey(key)
        // fmt.Println(namespace, "-", name, "-", err)

        var ctx = context.Background()
        var opts metav1.GetOptions
        service, _ := c.KubeClientSet.CoreV1().Services(namespace).Get(
            ctx,
            name,
            opts,
        )
        var _ = service
        // fmt.Println("--> ", service)

        if namespace == "default" {

            var opts1 metav1.GetOptions
            d, _ := c.KubeClientSet.AppsV1().Deployments(namespace).Get(
                context.Background(),
                "my-app",
                opts1,
            )
            if d != nil {
                fmt.Println("-->d ", d)
                var v int32 = 10;
                d.Spec.Replicas = &v

                var err error
                d, err = c.KubeClientSet.AppsV1().Deployments(namespace).Update(context.TODO(), d, metav1.UpdateOptions{})
                if err != nil {
                    fmt.Println("Bad ")
                } else {
                    fmt.Println("good ")
                }
            }
        }

        time.Sleep(1 * time.Second)

        //deployment, err = c.KubeClientSet.AppsV1().Deployments(namespace).Update(context.TODO(), newDeployment(foo), metav1.UpdateOptions{})
    }

    /*var object metav1.Object
    if object, ok = obj.(metav1.Object); !ok {
        fmt.Println("Not good - 1")

    }
    if ownerRef := metav1.GetControllerOf(object); ownerRef != nil {
        fmt.Println("Not good")
    }*/
    c.workQueue.Forget(obj)

    return true
}
