package postgres_test

import (
	"slices"
	"testing"

	"github.com/getfider/fider/app/models/cmd"
	"github.com/getfider/fider/app/models/entity"
	"github.com/getfider/fider/app/models/enum"
	"github.com/getfider/fider/app/models/query"
	"github.com/getfider/fider/app/pkg/bus"
)

func TestVisitorMentionPreferences(t *testing.T) {
	for _, tc := range []struct {
		name    string
		setting string
		web     bool
		email   bool
	}{
		{name: "default", web: true, email: true},
		{name: "opted out", setting: "0"},
		{name: "web only", setting: "1", web: true},
		{name: "email only", setting: "2", email: true},
		{name: "both", setting: "3", web: true, email: true},
	} {
		t.Run(tc.name, func(t *testing.T) {
			SetupDatabaseTest(t)
			defer TeardownDatabaseTest()
			if aryaStark.Role != enum.RoleVisitor {
				t.Fatal("expected a visitor fixture")
			}
			if tc.setting != "" {
				if err := bus.Dispatch(aryaStarkCtx, &cmd.UpdateCurrentUserSettings{Settings: map[string]string{
					enum.NotificationEventMention.UserSettingsKeyName: tc.setting,
				}}); err != nil {
					t.Fatal(err)
				}
			}
			settings := &query.GetCurrentUserSettings{}
			if err := bus.Dispatch(aryaStarkCtx, settings); err != nil {
				t.Fatal(err)
			}
			want := tc.setting
			if want == "" {
				want = "3"
			}
			if got := settings.Result[enum.NotificationEventMention.UserSettingsKeyName]; got != want {
				t.Errorf("mention setting = %q, want %q", got, want)
			}
			for channel, enabled := range map[enum.NotificationChannel]bool{
				enum.NotificationChannelWeb: tc.web, enum.NotificationChannelEmail: tc.email,
			} {
				subscribers := &query.GetActiveSubscribers{Channel: channel, Event: enum.NotificationEventMention}
				if err := bus.Dispatch(aryaStarkCtx, subscribers); err != nil {
					t.Fatal(err)
				}
				got := slices.ContainsFunc(subscribers.Result, func(user *entity.User) bool { return user.ID == aryaStark.ID })
				if got != enabled {
					t.Errorf("visitor eligible on channel %d = %t, want %t", channel, got, enabled)
				}
			}
		})
	}
}
