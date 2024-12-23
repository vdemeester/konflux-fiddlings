package main

import (
	"context"
	"encoding/json"
	"math/rand"
	"time"

	pipelinev1 "github.com/tektoncd/pipeline/pkg/apis/pipeline/v1"
	clientset "github.com/tektoncd/pipeline/pkg/client/clientset/versioned"
	pipelineclient "github.com/tektoncd/pipeline/pkg/client/injection/client"
	pipelineruninformer "github.com/tektoncd/pipeline/pkg/client/injection/informers/pipeline/v1/pipelinerun"
	pipelinerunreconciler "github.com/tektoncd/pipeline/pkg/client/injection/reconciler/pipeline/v1/pipelinerun"
	listers "github.com/tektoncd/pipeline/pkg/client/listers/pipeline/v1"
	jsonpatch "gomodules.xyz/jsonpatch/v2"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	// corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/utils/clock"
	kubeclient "knative.dev/pkg/client/injection/kube/client"

	// filteredinformerfactory "knative.dev/pkg/client/injection/kube/informers/factory/filtered"
	"knative.dev/pkg/configmap"
	"knative.dev/pkg/controller"
	"knative.dev/pkg/injection"
	"knative.dev/pkg/injection/sharedmain"
	"knative.dev/pkg/logging"
	"knative.dev/pkg/reconciler"
	"knative.dev/pkg/signals"
)

func main() {
	// This parses flags.
	cfg := injection.ParseAndGetRESTConfigOrDie()
	if cfg.QPS == 0 {
		cfg.QPS = 2 * rest.DefaultQPS
	}
	if cfg.Burst == 0 {
		cfg.Burst = rest.DefaultBurst
	}
	ctx := injection.WithNamespaceScope(signals.NewContext(), "vdemeest-tenant")
	ctx = sharedmain.WithHADisabled(ctx)

	sharedmain.MainWithConfig(ctx, "fiddlings-controller", cfg,
		newController(clock.RealClock{}),
	)
}

// newController creates a new controller.
func newController(clock clock.PassiveClock) func(context.Context, configmap.Watcher) *controller.Impl {
	return func(ctx context.Context, watcher configmap.Watcher) *controller.Impl {
		kubeclientset := kubeclient.Get(ctx)
		pipelineclientset := pipelineclient.Get(ctx)
		pipelineRunInformer := pipelineruninformer.Get(ctx)

		c := &Reconciler{
			kubeclient:        kubeclientset,
			pipelineclientset: pipelineclientset,
			pipelineRunLister: pipelineRunInformer.Lister(),
		}
		impl := pipelinerunreconciler.NewImpl(ctx, c, func(impl *controller.Impl) controller.Options {
			return controller.Options{
				AgentName: "fiddlings-controller",
			}
		})

		if _, err := pipelineRunInformer.Informer().AddEventHandler(controller.HandleAll(impl.Enqueue)); err != nil {
			logging.FromContext(ctx).Panicf("Couldn't register PipelineRun informer event handler: %w", err)
		}

		return impl
	}
}

type Reconciler struct {
	kubeclient        kubernetes.Interface
	pipelineclientset clientset.Interface
	pipelineRunLister listers.PipelineRunLister
}

func (r *Reconciler) ReconcileKind(ctx context.Context, pr *pipelinev1.PipelineRun) reconciler.Event {
	logger := logging.FromContext(ctx)

	// no-op on non-pending PipelineRuns
	if pr.Spec.Status != pipelinev1.PipelineRunSpecStatusPending {
		return nil
	}
	logger.Infof("Reconciling PipelineRun %s", pr.Name)

	// Randomly sleep for a given amount of time
	v := rand.Intn(60) + 10

	logger.Infof("PipelineRun %s is pending, starting it in %d seconds", pr.Name, v)
	pr.Spec.Status = ""

	time.Sleep(time.Duration(v) * time.Second)

	logger.Infof("Patching PipelineRun %s", pr.Name)

	patches := []jsonpatch.JsonPatchOperation{{
		Operation: "remove",
		Path:      "/spec/status",
	}}
	patchBytes, err := json.Marshal(patches)
	if err != nil {
		logger.Errorf("Failed to marshal patch: %v", err)
		return err
	}

	// Patch the PipelineRun to remove the status
	if _, err := r.pipelineclientset.TektonV1().PipelineRuns(pr.Namespace).Patch(ctx, pr.Name, types.JSONPatchType, patchBytes, metav1.PatchOptions{}, ""); err != nil {
		logger.Errorf("Failed to patch PipelineRun %s: %v", pr.Name, err)
		return err
	}

	logger.Infof("PipelineRun %s is now going to run", pr.Name)

	return nil
}
