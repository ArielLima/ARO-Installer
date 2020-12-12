package audit

import (
	"testing"
	"time"

	"github.com/Azure/go-autorest/autorest/azure"
	"github.com/golang/mock/gomock"
	"github.com/onsi/gomega"
	"github.com/onsi/gomega/types"
	"github.com/sirupsen/logrus"
	"github.com/sirupsen/logrus/hooks/test"

	mock_env "github.com/Azure/ARO-RP/pkg/util/mocks/env"
	utillog "github.com/Azure/ARO-RP/test/util/log"
)

func TestAudit(t *testing.T) {
	logger, h := test.NewNullLogger()

	controller := gomock.NewController(t)
	defer controller.Finish()

	env := mock_env.NewMockInterface(controller)
	env.EXPECT().Environment().AnyTimes().Return(&azure.PublicCloud)
	env.EXPECT().Location().AnyTimes().Return("eastus")

	auditLog := NewEntry(env, logger)

	// maybe this goes into a helper function?
	auditLog.WithFields(logrus.Fields{
		PayloadKeyCategory:      CategoryAuthorization,
		PayloadKeyOperationName: "initializeAuthorizers",
		MetadataSource:          SourceRP,
	}).Print("see auditFullPayload field for full log data")

	if err := utillog.AssertLoggingOutput(h, []map[string]types.GomegaMatcher{
		{
			"level":         gomega.Equal(logrus.InfoLevel),
			"msg":           gomega.Equal("see auditFullPayload field for full log data"),
			MetadataSource:  gomega.Equal("aro-rp"),
			MetadataLogKind: gomega.Equal("ifxaudit"),
			MetadataCreatedTime: gomega.WithTransform(
				func(s string) time.Time {
					t, err := time.Parse(time.RFC3339, s)
					if err != nil {
						panic(err)
					}
					return t
				},
				gomega.BeTemporally("~", time.Now(), time.Second),
			),
			// auditFullPayload: "{"env_ver":2.1,"env_name":"#Ifx.AuditSchema","env_time":"2020-12-11T13:47:27Z","env_epoch":"ab34cafa-b047-4dc4-a8ff-e24b7f854b4d","env_seqNum":1,"env_popSample":0,"env_iKey":null,"env_flags":257,"env_cv":"","env_os":"linux","env_osVer":null,"env_appId":null,"env_appVer":null,"env_cloud_ver":1,"env_cloud_name":"AzurePublicCloud","env_cloud_role":"","env_cloud_roleVer":null,"env_cloud_roleInstance":"","env_cloud_environment":null,"env_cloud_location":"eastus","env_cloud_deploymentUnit":null,"CallerIdentities":[{"CallerDisplayName":"","CallerIdentityType":"PUID","CallerIdentityValue":"1261453","CallerIpAddress":""}],"Category":"Authorization","nCloud":"AzurePublicCloud","OperationName":"initializeAuthorizers","Result":{"ResultType":"","ResultDescription":""},"requestId":"","TargetResources":[{"TargetResourceType":"resource provider","TargetResourceName":"aro-rp"}]}"
		},
	}); err != nil {
		t.Error(err)
	}
}
