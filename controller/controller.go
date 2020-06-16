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

    // KubeInformerFactory - this will callback us every 1 sec
    // workQueue - it is just a queue to store the updates and process them one-by-one
    controller := &Controller{
        KubeClientSet:       KubeClientSet,
        KubeInformerFactory: informers.NewSharedInformerFactory(KubeClientSet, 1*time.Second),
        workQueue:           workqueue.NewNamedRateLimitingQueue(workqueue.DefaultControllerRateLimiter(), "my-controller-rate-limit"),
    }

    // ************************************ This is the real thing *****************************************************
    // This will register a listener which is called when we have update
    controller.KubeInformerFactory.Apps().V1().Deployments().Informer().AddEventHandler(
        cache.ResourceEventHandlerFuncs{
            AddFunc: func(obj interface{}) {
                // This will be called first time with all services, deployment etc
            },
            UpdateFunc: func(oldObj, newObj interface{}) {
                // Called when we see a update
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

    // Start informer (hack to put local channel to close)
    stop := make(chan struct{})
    controller.KubeInformerFactory.Start(stop)

    return controller;
}

// This is the method which will run to do the job
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

    obj, shutdown := c.workQueue.Get()
    if shutdown {
        return false
    }

    defer c.workQueue.Done(obj)

    var key string
    var ok bool
    if key, ok = obj.(string); ok {
        namespace, name, _ := cache.SplitMetaNamespaceKey(key)
        var _ = name

        /* var ctx = context.Background()
         var opts metav1.GetOptions
         service, _ := c.KubeClientSet.CoreV1().Services(namespace).Get(
             ctx,
             name,
             opts,
         )
         var _ = service
         // fmt.Println("--> ", service)*/

        // We will work with only default namespace - only on "my-app"
        if namespace == "default" {

            var opts1 metav1.GetOptions
            d, _ := c.KubeClientSet.AppsV1().Deployments(namespace).Get(
                context.Background(),
                "my-app",
                opts1,
            )
            if d != nil {
                fmt.Println("Got Deployment - ", d)

                // Update replica set to 10
                var v int32 = 10;
                d.Spec.Replicas = &v

                // Apply changes
                var err error
                d, err = c.KubeClientSet.AppsV1().Deployments(namespace).Update(context.TODO(), d, metav1.UpdateOptions{})

                if err != nil {
                    fmt.Println("Something went wrong when we tried to scale my-app to 10")
                } else {
                    fmt.Println("Done - see your my-app should scaled to 10 nodes")
                }
            } else {
                fmt.Println("Something is wrong - we did not get my-app")
            }
        }

        time.Sleep(1 * time.Second)
    }

    c.workQueue.Forget(obj)

    return true
}
