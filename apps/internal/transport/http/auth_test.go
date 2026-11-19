package http

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"reflect"
	"strings"
	"testing"

	"github.com/google/uuid"
	"github.com/nawfal-btw/CE1-Group-3/apps/internal/application"
	"github.com/nawfal-btw/CE1-Group-3/apps/internal/auth"
)

// Local-only test adapter; no database, OIDC discovery or production wiring.
type incidentAuthenticatorFunc func(*http.Request) (application.Identity, error)

func (f incidentAuthenticatorFunc) AuthenticateIncident(r *http.Request) (application.Identity, error) {
	return f(r)
}

func TestIncidentAuthenticationFailures(t *testing.T) {
	const sensitive = "sensitive internal diagnostic"
	for _, tc := range []struct {
		name   string
		err    error
		status int
	}{
		{"unauthenticated", auth.ErrUnauthenticated, 401},
		{"unavailable", auth.ErrAuthUnavailable, 503},
		{"wrapped unauthenticated", fmt.Errorf("%s: %w", sensitive, auth.ErrUnauthenticated), 401},
		{"wrapped unavailable", fmt.Errorf("%s: %w", sensitive, auth.ErrAuthUnavailable), 503},
		{"unclassified", errors.New(sensitive), 401},
	} {
		t.Run(tc.name, func(t *testing.T) {
			for _, target := range []string{"helper", "GET", "POST", "PATCH"} {
				t.Run(target, func(t *testing.T) {
					var logs bytes.Buffer
					logger := slog.New(slog.NewJSONHandler(&logs, nil))
					calls := 0
					authenticator := incidentAuthenticatorFunc(func(*http.Request) (application.Identity, error) {
						calls++
						// Even a partial identity supplied alongside an error must be discarded.
						return application.Identity{SubjectID: uuid.New(), Roles: []string{"dispatcher"}}, tc.err
					})
					// The uninitialized application service cannot run: its nil dependencies
					// would panic. Exercising all three handlers verifies they stop first.
					h := NewIncidentHandler(application.IncidentService{}, authenticator, logger)
					method := target
					if target == "helper" {
						method = "GET"
					}
					r := httptest.NewRequest(method, "https://bff.example.test/v1/incident", strings.NewReader("{}"))
					r.SetPathValue("id", uuid.NewString())
					w := httptest.NewRecorder()
					switch target {
					case "helper":
						identity, ok := h.authenticateIncident(w, r)
						if ok || identity.SubjectID != uuid.Nil || identity.Roles != nil {
							t.Fatal("failed authentication returned an identity")
						}
					case "GET":
						h.GetIncident(w, r)
					case "POST":
						h.PostIncident(w, r)
					case "PATCH":
						h.PatchIncident(w, r)
					}
					if calls != 1 || w.Code != tc.status || w.Body.String() != http.StatusText(tc.status)+"\n" {
						t.Fatal("incorrect authentication call count, status or generic body")
					}
					if strings.Contains(w.Body.String(), sensitive) || strings.Contains(logs.String(), sensitive) {
						t.Fatal("authentication diagnostics leaked")
					}
					lines := strings.Split(strings.TrimSpace(logs.String()), "\n")
					var event map[string]any
					if err := json.Unmarshal([]byte(lines[len(lines)-1]), &event); err != nil {
						t.Fatal("invalid authentication log")
					}
					if event["status"] != float64(tc.status) || event["error"] != nil {
						t.Fatal("authentication log omitted status or included raw error")
					}
				})
			}
		})
	}
}

func TestIncidentAuthenticationSuccess(t *testing.T) {
	want := application.Identity{SubjectID: uuid.New(), Roles: []string{"dispatcher", "operator"}}
	var logs bytes.Buffer
	logger := slog.New(slog.NewJSONHandler(&logs, nil))
	calls := 0
	authenticator := incidentAuthenticatorFunc(func(*http.Request) (application.Identity, error) {
		calls++
		return want, nil
	})
	h := NewIncidentHandler(application.IncidentService{}, authenticator, logger)
	w := httptest.NewRecorder()
	got, ok := h.authenticateIncident(w, httptest.NewRequest("GET", "https://bff.example.test/v1/incident", nil))
	if !ok || calls != 1 || !reflect.DeepEqual(got, want) {
		t.Fatal("successful authentication changed application UUID or roles")
	}
	if w.Body.Len() != 0 || len(w.Header()) != 0 || logs.Len() != 0 {
		t.Fatal("successful authentication unexpectedly wrote a response or failure log")
	}
}
