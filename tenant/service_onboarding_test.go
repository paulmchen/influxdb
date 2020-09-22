package tenant_test

import (
	"context"
	"testing"

	"github.com/google/go-cmp/cmp"

	influxdb "github.com/influxdata/influxdb/v2"
	icontext "github.com/influxdata/influxdb/v2/context"
	"github.com/influxdata/influxdb/v2/kv"
	"github.com/influxdata/influxdb/v2/tenant"
	influxdbtesting "github.com/influxdata/influxdb/v2/testing"
	"go.uber.org/zap/zaptest"
)

func TestBoltOnboardingService(t *testing.T) {
	influxdbtesting.OnboardInitialUser(initBoltOnboardingService, t)
}

func initBoltOnboardingService(f influxdbtesting.OnboardingFields, t *testing.T) (influxdb.OnboardingService, func()) {
	s, closeStore, err := NewTestInmemStore(t)
	if err != nil {
		t.Fatalf("failed to create new bolt kv store: %v", err)
	}

	svc := initOnboardingService(s, f, t)
	return svc, func() {
		closeStore()
	}
}

func initOnboardingService(s kv.Store, f influxdbtesting.OnboardingFields, t *testing.T) influxdb.OnboardingService {
	storage := tenant.NewStore(s)
	ten := tenant.NewService(storage)

	// we will need an auth service as well
	svc := tenant.NewOnboardService(ten, kv.NewService(zaptest.NewLogger(t), s))

	ctx := context.Background()

	t.Logf("Onboarding: %v", f.IsOnboarding)
	if !f.IsOnboarding {
		// create a dummy so so we can no longer onboard
		err := ten.CreateUser(ctx, &influxdb.User{Name: "dummy", Status: influxdb.Active})
		if err != nil {
			t.Fatal(err)
		}
	}

	return svc
}

func TestOnboardURM(t *testing.T) {
	s, _, _ := NewTestInmemStore(t)
	storage := tenant.NewStore(s)
	ten := tenant.NewService(storage)

	// we will need an auth service as well
	svc := tenant.NewOnboardService(ten, kv.NewService(zaptest.NewLogger(t), s))

	ctx := icontext.SetAuthorizer(context.Background(), &influxdb.Authorization{
		UserID: 123,
	})

	onboard, err := svc.OnboardUser(ctx, &influxdb.OnboardingRequest{
		User:   "name",
		Org:    "name",
		Bucket: "name",
	})

	if err != nil {
		t.Fatal(err)
	}

	urms, _, err := ten.FindUserResourceMappings(ctx, influxdb.UserResourceMappingFilter{ResourceID: onboard.Org.ID})
	if err != nil {
		t.Fatal(err)
	}

	if len(urms) > 1 {
		t.Fatal("additional URMs created")
	}
	if urms[0].UserID != onboard.User.ID {
		t.Fatal("org assigned to the wrong user")
	}
}

func TestOnboardAuth(t *testing.T) {
	s, _, _ := NewTestInmemStore(t)
	storage := tenant.NewStore(s)
	ten := tenant.NewService(storage)

	// we will need an auth service as well
	svc := tenant.NewOnboardService(ten, kv.NewService(zaptest.NewLogger(t), s))

	ctx := icontext.SetAuthorizer(context.Background(), &influxdb.Authorization{
		UserID: 123,
	})

	onboard, err := svc.OnboardUser(ctx, &influxdb.OnboardingRequest{
		User:   "name",
		Org:    "name",
		Bucket: "name",
	})

	if err != nil {
		t.Fatal(err)
	}

	auth := onboard.Auth
	expectedPerm := []influxdb.Permission{
		influxdb.Permission{Action: influxdb.ReadAction, Resource: influxdb.Resource{OrgID: &onboard.Org.ID, Type: influxdb.AuthorizationsResourceType}},
		influxdb.Permission{Action: influxdb.WriteAction, Resource: influxdb.Resource{OrgID: &onboard.Org.ID, Type: influxdb.AuthorizationsResourceType}},
		influxdb.Permission{Action: influxdb.ReadAction, Resource: influxdb.Resource{OrgID: &onboard.Org.ID, Type: influxdb.BucketsResourceType}},
		influxdb.Permission{Action: influxdb.WriteAction, Resource: influxdb.Resource{OrgID: &onboard.Org.ID, Type: influxdb.BucketsResourceType}},
		influxdb.Permission{Action: influxdb.ReadAction, Resource: influxdb.Resource{OrgID: &onboard.Org.ID, Type: influxdb.DashboardsResourceType}},
		influxdb.Permission{Action: influxdb.WriteAction, Resource: influxdb.Resource{OrgID: &onboard.Org.ID, Type: influxdb.DashboardsResourceType}},
		influxdb.Permission{Action: influxdb.ReadAction, Resource: influxdb.Resource{ID: &onboard.Org.ID, Type: influxdb.OrgsResourceType}},
		influxdb.Permission{Action: influxdb.WriteAction, Resource: influxdb.Resource{ID: &onboard.Org.ID, Type: influxdb.OrgsResourceType}},
		influxdb.Permission{Action: influxdb.ReadAction, Resource: influxdb.Resource{OrgID: &onboard.Org.ID, Type: influxdb.SourcesResourceType}},
		influxdb.Permission{Action: influxdb.WriteAction, Resource: influxdb.Resource{OrgID: &onboard.Org.ID, Type: influxdb.SourcesResourceType}},
		influxdb.Permission{Action: influxdb.ReadAction, Resource: influxdb.Resource{OrgID: &onboard.Org.ID, Type: influxdb.TasksResourceType}},
		influxdb.Permission{Action: influxdb.WriteAction, Resource: influxdb.Resource{OrgID: &onboard.Org.ID, Type: influxdb.TasksResourceType}},
		influxdb.Permission{Action: influxdb.ReadAction, Resource: influxdb.Resource{OrgID: &onboard.Org.ID, Type: influxdb.TelegrafsResourceType}},
		influxdb.Permission{Action: influxdb.WriteAction, Resource: influxdb.Resource{OrgID: &onboard.Org.ID, Type: influxdb.TelegrafsResourceType}},
		influxdb.Permission{Action: influxdb.ReadAction, Resource: influxdb.Resource{OrgID: &onboard.Org.ID, Type: influxdb.UsersResourceType}},
		influxdb.Permission{Action: influxdb.WriteAction, Resource: influxdb.Resource{OrgID: &onboard.Org.ID, Type: influxdb.UsersResourceType}},
		influxdb.Permission{Action: influxdb.ReadAction, Resource: influxdb.Resource{OrgID: &onboard.Org.ID, Type: influxdb.VariablesResourceType}},
		influxdb.Permission{Action: influxdb.WriteAction, Resource: influxdb.Resource{OrgID: &onboard.Org.ID, Type: influxdb.VariablesResourceType}},
		influxdb.Permission{Action: influxdb.ReadAction, Resource: influxdb.Resource{OrgID: &onboard.Org.ID, Type: influxdb.ScraperResourceType}},
		influxdb.Permission{Action: influxdb.WriteAction, Resource: influxdb.Resource{OrgID: &onboard.Org.ID, Type: influxdb.ScraperResourceType}},
		influxdb.Permission{Action: influxdb.ReadAction, Resource: influxdb.Resource{OrgID: &onboard.Org.ID, Type: influxdb.SecretsResourceType}},
		influxdb.Permission{Action: influxdb.WriteAction, Resource: influxdb.Resource{OrgID: &onboard.Org.ID, Type: influxdb.SecretsResourceType}},
		influxdb.Permission{Action: influxdb.ReadAction, Resource: influxdb.Resource{OrgID: &onboard.Org.ID, Type: influxdb.LabelsResourceType}},
		influxdb.Permission{Action: influxdb.WriteAction, Resource: influxdb.Resource{OrgID: &onboard.Org.ID, Type: influxdb.LabelsResourceType}},
		influxdb.Permission{Action: influxdb.ReadAction, Resource: influxdb.Resource{OrgID: &onboard.Org.ID, Type: influxdb.ViewsResourceType}},
		influxdb.Permission{Action: influxdb.WriteAction, Resource: influxdb.Resource{OrgID: &onboard.Org.ID, Type: influxdb.ViewsResourceType}},
		influxdb.Permission{Action: influxdb.ReadAction, Resource: influxdb.Resource{OrgID: &onboard.Org.ID, Type: influxdb.DocumentsResourceType}},
		influxdb.Permission{Action: influxdb.WriteAction, Resource: influxdb.Resource{OrgID: &onboard.Org.ID, Type: influxdb.DocumentsResourceType}},
		influxdb.Permission{Action: influxdb.ReadAction, Resource: influxdb.Resource{OrgID: &onboard.Org.ID, Type: influxdb.NotificationRuleResourceType}},
		influxdb.Permission{Action: influxdb.WriteAction, Resource: influxdb.Resource{OrgID: &onboard.Org.ID, Type: influxdb.NotificationRuleResourceType}},
		influxdb.Permission{Action: influxdb.ReadAction, Resource: influxdb.Resource{OrgID: &onboard.Org.ID, Type: influxdb.NotificationEndpointResourceType}},
		influxdb.Permission{Action: influxdb.WriteAction, Resource: influxdb.Resource{OrgID: &onboard.Org.ID, Type: influxdb.NotificationEndpointResourceType}},
		influxdb.Permission{Action: influxdb.ReadAction, Resource: influxdb.Resource{OrgID: &onboard.Org.ID, Type: influxdb.ChecksResourceType}},
		influxdb.Permission{Action: influxdb.WriteAction, Resource: influxdb.Resource{OrgID: &onboard.Org.ID, Type: influxdb.ChecksResourceType}},
		influxdb.Permission{Action: influxdb.ReadAction, Resource: influxdb.Resource{OrgID: &onboard.Org.ID, Type: influxdb.DBRPResourceType}},
		influxdb.Permission{Action: influxdb.WriteAction, Resource: influxdb.Resource{OrgID: &onboard.Org.ID, Type: influxdb.DBRPResourceType}},
		influxdb.Permission{Action: influxdb.ReadAction, Resource: influxdb.Resource{ID: &onboard.User.ID, Type: influxdb.UsersResourceType}},
		influxdb.Permission{Action: influxdb.WriteAction, Resource: influxdb.Resource{ID: &onboard.User.ID, Type: influxdb.UsersResourceType}},
	}
	if !cmp.Equal(auth.Permissions, expectedPerm) {
		t.Fatalf("unequal permissions: \n %+v", cmp.Diff(auth.Permissions, expectedPerm))
	}

}
