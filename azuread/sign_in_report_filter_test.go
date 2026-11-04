package azuread

import (
	"context"
	"errors"
	"slices"
	"testing"
	"time"

	"github.com/microsoftgraph/msgraph-sdk-go/models"
	"github.com/turbot/steampipe-plugin-sdk/v6/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v6/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v6/plugin/quals"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func timeQual(op string, at time.Time) *quals.Qual {
	return &quals.Qual{Column: "created_date_time", Operator: op, Value: &proto.QualValue{Value: &proto.QualValue_TimestampValue{TimestampValue: timestamppb.New(at)}}}
}

func listQual(column, op string) *quals.Qual {
	return &quals.Qual{Column: column, Operator: op, Value: &proto.QualValue{Value: &proto.QualValue_ListValue{ListValue: &proto.QualValueList{}}}}
}

func stringQual(column, op, value string) *quals.Qual {
	return &quals.Qual{Column: column, Operator: op, Value: &proto.QualValue{Value: &proto.QualValue_StringValue{StringValue: value}}}
}

func qualMap(qs ...*quals.Qual) plugin.KeyColumnQualMap {
	m := plugin.KeyColumnQualMap{}
	for _, q := range qs {
		if m[q.Column] == nil {
			m[q.Column] = &plugin.KeyColumnQuals{Name: q.Column}
		}
		m[q.Column].Quals = append(m[q.Column].Quals, q)
	}
	return m
}

func signIn(id string, created *time.Time, clientApp string) models.SignInable {
	s := models.NewSignIn()
	s.SetId(&id)
	s.SetCreatedDateTime(created)
	if clientApp != "" {
		s.SetClientAppUsed(&clientApp)
	}
	return s
}

func TestSignInGraphFilter(t *testing.T) {
	whole := time.Date(2026, 9, 20, 10, 0, 0, 0, time.UTC)
	fractional := time.Date(2026, 9, 20, 10, 0, 0, 500_000_000, time.UTC)
	offset := time.Date(2026, 9, 20, 15, 30, 0, 0, time.FixedZone("IST", 5*3600+1800))

	cases := []struct {
		name  string
		quals plugin.KeyColumnQualMap
		want  []string
	}{
		{"no quals", qualMap(), []string{}},
		{">= whole second", qualMap(timeQual(">=", whole)), []string{"createdDateTime ge 2026-09-20T10:00:00Z"}},
		{"> widens to ge", qualMap(timeQual(">", whole)), []string{"createdDateTime ge 2026-09-20T10:00:00Z"}},
		{">= fraction floors", qualMap(timeQual(">=", fractional)), []string{"createdDateTime ge 2026-09-20T10:00:00Z"}},
		{"<= whole second still covers its sub-second sign-ins", qualMap(timeQual("<=", whole)), []string{"createdDateTime le 2026-09-20T10:00:01Z"}},
		{"<= fraction ceils", qualMap(timeQual("<=", fractional)), []string{"createdDateTime le 2026-09-20T10:00:01Z"}},
		{"< widens to le", qualMap(timeQual("<", whole)), []string{"createdDateTime le 2026-09-20T10:00:01Z"}},
		{"= becomes a one-second range", qualMap(timeQual("=", whole)), []string{"createdDateTime ge 2026-09-20T10:00:00Z", "createdDateTime le 2026-09-20T10:00:01Z"}},
		{"offset renders as UTC", qualMap(timeQual(">=", offset)), []string{"createdDateTime ge 2026-09-20T10:00:00Z"}},
		{"client app", qualMap(stringQual("client_app_used", "=", "Authenticated SMTP")), []string{"clientAppUsed eq 'Authenticated SMTP'"}},
		{"client app quote escaped", qualMap(stringQual("client_app_used", "=", "O'Brien")), []string{"clientAppUsed eq 'O''Brien'"}},
		{"blank client app is still pushed", qualMap(stringQual("client_app_used", "=", "")), []string{"clientAppUsed eq ''"}},
		{"an infinite bound is left to Postgres", qualMap(timeQual("<=", time.Date(294276, 12, 31, 0, 0, 0, 0, time.UTC)), timeQual(">", time.Date(-4713, 1, 1, 0, 0, 0, 0, time.UTC))), []string{}},
		{
			"unsplit IN lists are left to Postgres, never read as the epoch",
			qualMap(listQual("created_date_time", "="), listQual("client_app_used", "="), timeQual(">=", whole)),
			[]string{"createdDateTime ge 2026-09-20T10:00:00Z"},
		},
		{
			"window and client app together",
			qualMap(timeQual(">=", whole), stringQual("client_app_used", "=", "IMAP4")),
			[]string{"createdDateTime ge 2026-09-20T10:00:00Z", "clientAppUsed eq 'IMAP4'"},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got := parseSignInQuals(tc.quals).graphFilter()
			if !slices.Equal(got, tc.want) {
				t.Errorf("filter = %q, want %q", got, tc.want)
			}
		})
	}
}

func TestSignInAdmitsOnlyRowsPostgresKeeps(t *testing.T) {
	bound := time.Date(2026, 9, 20, 10, 0, 0, 0, time.UTC)
	at := func(d time.Duration) *time.Time { v := bound.Add(d); return &v }

	cases := []struct {
		name  string
		quals plugin.KeyColumnQualMap
		row   models.SignInable
		want  bool
	}{
		{"no quals keeps a row without a time", qualMap(), signIn("a", nil, ""), true},
		{"< drops the extra second Graph returned", qualMap(timeQual("<", bound)), signIn("a", at(400*time.Millisecond), ""), false},
		{"< keeps an earlier row", qualMap(timeQual("<", bound)), signIn("a", at(-time.Microsecond), ""), true},
		{"< drops a row Postgres rounds up onto the bound", qualMap(timeQual("<", bound)), signIn("a", at(-400*time.Nanosecond), ""), false},
		{"<= keeps a row Postgres rounds down onto the bound", qualMap(timeQual("<=", bound)), signIn("a", at(400*time.Nanosecond), ""), true},
		{"a half-microsecond tie is kept either way", qualMap(timeQual("<=", bound)), signIn("a", at(500*time.Nanosecond), ""), true},
		{"<= drops a row past the bound", qualMap(timeQual("<=", bound)), signIn("a", at(2*time.Microsecond), ""), false},
		{">= drops the extra second Graph returned", qualMap(timeQual(">=", bound.Add(500*time.Millisecond))), signIn("a", at(100*time.Millisecond), ""), false},
		{"range keeps a row inside it", qualMap(timeQual(">=", bound), timeQual("<", bound.Add(time.Hour))), signIn("a", at(time.Minute), ""), true},
		{"a time qual drops a row without a time", qualMap(timeQual(">=", bound)), signIn("a", nil, ""), false},
		{"an unsplit list qual is left to Postgres", qualMap(listQual("created_date_time", "=")), signIn("a", at(time.Hour), ""), true},
		{"client app matches exactly", qualMap(stringQual("client_app_used", "=", "IMAP4")), signIn("a", nil, "IMAP4"), true},
		{"Graph's case-insensitive match is dropped", qualMap(stringQual("client_app_used", "=", "imap4")), signIn("a", nil, "IMAP4"), false},
		{"a client app qual drops a row without one", qualMap(stringQual("client_app_used", "=", "IMAP4")), signIn("a", nil, ""), false},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			if got := parseSignInQuals(tc.quals).admits(tc.row); got != tc.want {
				t.Errorf("admits = %v, want %v", got, tc.want)
			}
		})
	}
}

func TestSignInReportListDeclaresFilterKeyColumns(t *testing.T) {
	keyColumns := tableAzureAdSignInReport(context.Background()).List.KeyColumns

	created := keyColumns.Find("created_date_time")
	if created == nil || created.Require != plugin.Optional {
		t.Fatalf("created_date_time must be an optional key column, got %+v", created)
	}
	if !slices.Equal(created.Operators, []string{">", ">=", "=", "<", "<="}) {
		t.Errorf("created_date_time operators = %v", created.Operators)
	}

	clientApp := keyColumns.Find("client_app_used")
	if clientApp == nil || clientApp.Require != plugin.Optional {
		t.Fatalf("client_app_used must be an optional key column, got %+v", clientApp)
	}
}

func TestSignInReportListLimiterTargetsOnlySignInLists(t *testing.T) {
	limiter := signInReportListLimiter
	if err := limiter.Initialise(); err != nil {
		t.Fatalf("limiter where clause does not parse: %v", err)
	}
	if errs := limiter.Validate(); len(errs) > 0 {
		t.Fatalf("limiter invalid: %v", errs)
	}

	tags := tableAzureAdSignInReport(context.Background()).List.Tags
	if !limiter.SatisfiesFilters(tags) {
		t.Errorf("limiter does not match the sign-in list call's tags %v", tags)
	}
	if limiter.SatisfiesFilters(map[string]string{"service": "auditLogs", "action": "directoryAudits.list"}) {
		t.Error("limiter must not throttle other audit-log tables")
	}
	if !slices.Contains(Plugin(context.Background()).RateLimiters, limiter) {
		t.Error("plugin does not register the sign-in list limiter")
	}
}

func signInPage(nextLink string, rows ...models.SignInable) models.SignInCollectionResponseable {
	page := models.NewSignInCollectionResponse()
	page.SetValue(rows)
	if nextLink != "" {
		page.SetOdataNextLink(&nextLink)
	}
	return page
}

func TestIterateSignInPagesWaitsBeforeEveryNextPage(t *testing.T) {
	pages := map[string]models.SignInCollectionResponseable{
		"page2": signInPage("page3", signIn("b", nil, ""), nil, signIn("c", nil, "")),
		"page3": signInPage("", signIn("d", nil, "")),
	}
	var events, ids []string
	wait := func(context.Context) { events = append(events, "wait") }
	fetchPage := func(_ context.Context, link string) (models.SignInCollectionResponseable, error) {
		events = append(events, link)
		return pages[link], nil
	}

	err := iterateSignInPages(context.Background(), signInPage("page2", signIn("a", nil, "")), wait, fetchPage, func(s models.SignInable) bool {
		ids = append(ids, *s.GetId())
		return true
	})
	if err != nil {
		t.Fatalf("iterate: %v", err)
	}
	if !slices.Equal(ids, []string{"a", "b", "c", "d"}) {
		t.Errorf("ids = %v (a null entry must be skipped)", ids)
	}
	if !slices.Equal(events, []string{"wait", "page2", "wait", "page3"}) {
		t.Errorf("events = %v, want a wait before each fetch", events)
	}
}

func TestIterateSignInPagesStops(t *testing.T) {
	fetches := 0
	noWait := func(context.Context) {}
	failing := func(context.Context, string) (models.SignInCollectionResponseable, error) {
		fetches++
		return nil, errors.New("throttled")
	}
	first := signInPage("page2", signIn("a", nil, ""))

	if err := iterateSignInPages(context.Background(), first, noWait, failing, func(models.SignInable) bool { return false }); err != nil || fetches != 0 {
		t.Errorf("a stopped emit must not fetch: err=%v fetches=%d", err, fetches)
	}

	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	if err := iterateSignInPages(ctx, first, noWait, failing, func(models.SignInable) bool { return true }); err != nil || fetches != 0 {
		t.Errorf("a cancelled query must stop quietly after the wait: err=%v fetches=%d", err, fetches)
	}

	if err := iterateSignInPages(context.Background(), first, noWait, failing, func(models.SignInable) bool { return true }); err == nil || err.Error() != "throttled" {
		t.Errorf("err = %v, want the page error", err)
	}
}
