package components

import (
	"context"
	"fmt"
	ytv1 "github.com/ytsaurus/yt-k8s-operator/api/v1"
	corev1 "k8s.io/api/core/v1"
	"strings"

	"github.com/ytsaurus/yt-k8s-operator/pkg/apiproxy"
	"github.com/ytsaurus/yt-k8s-operator/pkg/consts"
	"github.com/ytsaurus/yt-k8s-operator/pkg/labeller"
	"github.com/ytsaurus/yt-k8s-operator/pkg/resources"
	"github.com/ytsaurus/yt-k8s-operator/pkg/ytconfig"
)

type Spyt struct {
	labeller *labeller.Labeller
	spyt     *apiproxy.Spyt
	cfgen    *ytconfig.Generator
	ytsaurus *ytv1.Ytsaurus

	secret *resources.StringSecret

	initUser        *InitJob
	initEnvironment *InitJob
}

func NewSpyt(
	cfgen *ytconfig.Generator,
	spyt *apiproxy.Spyt,
	ytsaurus *ytv1.Ytsaurus) *Spyt {

	spytSpec := spyt.GetResource().Spec

	l := labeller.Labeller{
		ObjectMeta:     &spyt.GetResource().ObjectMeta,
		APIProxy:       spyt.APIProxy(),
		ComponentLabel: fmt.Sprintf("ytsaurus-spyt-%s", spytSpec.Name),
		ComponentName:  fmt.Sprintf("SPYT-%s", spytSpec.Name),
	}

	return &Spyt{
		labeller: &l,
		spyt:     spyt,
		cfgen:    cfgen,
		ytsaurus: ytsaurus,
		initUser: NewInitJob(
			&l,
			spyt.APIProxy(),
			spyt,
			ytsaurus.Spec.ImagePullSecrets,
			"user",
			consts.ClientConfigFileName,
			ytsaurus.Spec.CoreImage,
			cfgen.GetNativeClientConfig),
		initEnvironment: NewInitJob(
			&l,
			spyt.APIProxy(),
			spyt,
			ytsaurus.Spec.ImagePullSecrets,
			"spyt-environment",
			consts.ClientConfigFileName,
			spytSpec.Image,
			cfgen.GetNativeClientConfig),
		secret: resources.NewStringSecret(
			l.GetSecretName(),
			&l,
			spyt.APIProxy()),
	}
}

func (s *Spyt) createInitUserScript() string {
	token, _ := s.secret.GetValue(consts.TokenSecretKey)
	commands := createUserCommand("spyt_releaser", "", token, true)
	script := []string{
		initJobWithNativeDriverPrologue(),
	}
	script = append(script, commands...)

	return strings.Join(script, "\n")
}

func (s *Spyt) createInitScript() string {
	script := []string{
		"/entrypoint.sh",
	}

	return strings.Join(script, "\n")
}

func (s *Spyt) doSync(ctx context.Context, dry bool) (SyncStatus, error) {
	var err error

	if s.ytsaurus.Status.State != ytv1.ClusterStateRunning {
		return SyncStatusBlocked, err
	}

	// Create user for spyt initialization.
	if s.secret.NeedSync(consts.TokenSecretKey, "") {
		if !dry {
			secretSpec := s.secret.Build()
			secretSpec.StringData = map[string]string{
				consts.TokenSecretKey: ytconfig.RandString(30),
			}
			err = s.secret.Sync(ctx)
		}
		return SyncStatusPending, err
	}
	if !dry {
		s.initUser.SetInitScript(s.createInitUserScript())
	}
	status, err := s.initUser.Sync(ctx, dry)
	if status != SyncStatusReady {
		return status, err
	}

	if !dry {
		s.initEnvironment.SetInitScript(s.createInitScript())
		job := s.initEnvironment.Build()
		container := &job.Spec.Template.Spec.Containers[0]
		token, _ := s.secret.GetValue(consts.TokenSecretKey)
		container.Env = []corev1.EnvVar{
			{
				Name:  "YT_PROXY",
				Value: s.cfgen.GetHTTPProxiesAddress(consts.DefaultHTTPProxyRole),
			},
			{
				Name:  "YT_TOKEN",
				Value: token,
			},
			{
				Name:  "EXTRA_PUBLISH_CLUSTER_OPTIONS",
				Value: "--ignore-existing",
			},
		}
	}

	return s.initEnvironment.Sync(ctx, dry)
}

func (s *Spyt) Fetch(ctx context.Context) error {
	return resources.Fetch(ctx, []resources.Fetchable{
		s.initUser,
		s.initEnvironment,
		s.secret,
	})
}

func (s *Spyt) Status(ctx context.Context) SyncStatus {
	status, err := s.doSync(ctx, true)
	if err != nil {
		panic(err)
	}

	return status
}

func (s *Spyt) Sync(ctx context.Context) error {
	_, err := s.doSync(ctx, false)
	return err
}
