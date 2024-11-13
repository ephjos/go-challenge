package handlers_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"slices"
	"testing"

	"github.com/gloveboxhq/glovebox-go-code-challenge/comms/email"
	"github.com/gloveboxhq/glovebox-go-code-challenge/comms/email/mockemail"
	"github.com/gloveboxhq/glovebox-go-code-challenge/comms/handlers"
)

func TestAddPolicyVehicle(t *testing.T) {

	t.Parallel()

	type testCase struct {
		method       string
		payload      handlers.AddPolicyVehicleReq
		expectTplID  email.TplID
		expectStatus int
	}

	testCases := map[string]testCase{
		"pass": {
			method: http.MethodPost,
			payload: handlers.AddPolicyVehicleReq{
				EmailTo: "foo@bar.com",
				Message: json.RawMessage(`{"foo":"bar"}`),
			},
			expectTplID:  email.TplAddPolicyVehicle,
			expectStatus: http.StatusOK,
		},
		"fail invalid method": {
			method:       http.MethodGet,
			payload:      handlers.AddPolicyVehicleReq{},
			expectStatus: http.StatusMethodNotAllowed,
		},
	}

	testFactory := func(tc testCase) func(*testing.T) {
		return func(t *testing.T) {

			payload, err := json.Marshal(tc.payload)
			if err != nil {
				t.Fatalf("could nor marshal payload to json: %v", err)
			}

			testEmail := mockemail.NewClient()
			defer testEmail.FlushSendLogs()

			w := httptest.NewRecorder()
			req := httptest.NewRequest(tc.method, "/api/comms/add-policy-vehicle", bytes.NewReader(payload))

			handlers.AddPolicyVehicle(testEmail)(w, req)

			resp := w.Result()

			if resp.StatusCode != tc.expectStatus {
				t.Fatalf("expected status %v but got %v", tc.expectStatus, resp.StatusCode)
			}

			if resp.StatusCode == http.StatusOK {

				if testEmail.SendLogs().IsEmpty() {
					t.Fatalf("expected email log but got empty")
				}

				lastEmail := testEmail.SendLogs().Last()

				if !slices.Equal(lastEmail.ExtractTos(), []string{tc.payload.EmailTo}) {
					t.Fatalf("expected to %v but got %v", lastEmail.ExtractTos(), tc.payload.EmailTo)
				}

				if string(lastEmail.ExtractMessage()) != string(tc.payload.Message) {
					t.Fatalf("expected message %v but got %v", tc.payload.Message, lastEmail.ExtractMessage())
				}

				if lastEmail.ExtractTplID() != tc.expectTplID {
					t.Fatalf("expected tpl %v but got %v", tc.expectTplID, lastEmail.ExtractTplID())
				}
			}
		}
	}

	for name, tc := range testCases {
		t.Run(name, testFactory(tc))
	}
}

func TestAddPolicyDriver(t *testing.T) {

	t.Parallel()

	type testCase struct {
		method       string
		payload      handlers.AddPolicyDriverReq
		expectTplID  email.TplID
		expectStatus int
	}

	testCases := map[string]testCase{
		"pass": {
			method: http.MethodPost,
			payload: handlers.AddPolicyDriverReq{
				EmailTo: "foo@bar.com",
				Message: json.RawMessage(`{"foo":"bar"}`),
			},
			expectTplID:  email.TplAddPolicyDriver,
			expectStatus: http.StatusOK,
		},
		"fail invalid method": {
			method:       http.MethodGet,
			payload:      handlers.AddPolicyDriverReq{},
			expectStatus: http.StatusMethodNotAllowed,
		},
	}

	testFactory := func(tc testCase) func(*testing.T) {
		return func(t *testing.T) {

			payload, err := json.Marshal(tc.payload)
			if err != nil {
				t.Fatalf("could nor marshal payload to json: %v", err)
			}

			testEmail := mockemail.NewClient()
			defer testEmail.FlushSendLogs()

			w := httptest.NewRecorder()
			req := httptest.NewRequest(tc.method, "/api/comms/add-policy-driver", bytes.NewReader(payload))

			handlers.AddPolicyDriver(testEmail)(w, req)

			resp := w.Result()

			if resp.StatusCode != tc.expectStatus {
				t.Fatalf("expected status %v but got %v", tc.expectStatus, resp.StatusCode)
			}

			if resp.StatusCode == http.StatusOK {

				if testEmail.SendLogs().IsEmpty() {
					t.Fatalf("expected email log but got empty")
				}

				lastEmail := testEmail.SendLogs().Last()

				if !slices.Equal(lastEmail.ExtractTos(), []string{tc.payload.EmailTo}) {
					t.Fatalf("expected to %v but got %v", lastEmail.ExtractTos(), tc.payload.EmailTo)
				}

				if string(lastEmail.ExtractMessage()) != string(tc.payload.Message) {
					t.Fatalf("expected message %v but got %v", tc.payload.Message, lastEmail.ExtractMessage())
				}

				if lastEmail.ExtractTplID() != tc.expectTplID {
					t.Fatalf("expected tpl %v but got %v", tc.expectTplID, lastEmail.ExtractTplID())
				}
			}
		}
	}

	for name, tc := range testCases {
		t.Run(name, testFactory(tc))
	}
}

func TestAddPolicyCoverage(t *testing.T) {

	t.Parallel()

	type testCase struct {
		method       string
		payload      handlers.AddPolicyCoverageReq
		expectTplID  email.TplID
		expectStatus int
	}

	testCases := map[string]testCase{
		"pass": {
			method: http.MethodPost,
			payload: handlers.AddPolicyCoverageReq{
				EmailTo: []string{"foo@bar.com"},
				EmailCC: []string{"baz@bar.com"},
				Message: json.RawMessage(`{"foo":"bar"}`),
			},
			expectTplID:  email.TplAddPolicyCoverage,
			expectStatus: http.StatusOK,
		},
		// NOTE: this passes in this faked sendgrid scenario, but I imagine sengrid
		// would error if passed no tos and we'd either have to handle the error or
		// catch empty tos at the route handler level
		"pass empty to": {
			method: http.MethodPost,
			payload: handlers.AddPolicyCoverageReq{
				EmailTo: []string{},
				EmailCC: []string{"baz@bar.com"},
				Message: json.RawMessage(`{"foo":"bar"}`),
			},
			expectTplID:  email.TplAddPolicyCoverage,
			expectStatus: http.StatusOK,
		},
		"pass empty cc": {
			method: http.MethodPost,
			payload: handlers.AddPolicyCoverageReq{
				EmailTo: []string{"foo@bar.com"},
				EmailCC: []string{},
				Message: json.RawMessage(`{"foo":"bar"}`),
			},
			expectTplID:  email.TplAddPolicyCoverage,
			expectStatus: http.StatusOK,
		},
		"pass 2 to": {
			method: http.MethodPost,
			payload: handlers.AddPolicyCoverageReq{
				EmailTo: []string{"foo@bar.com", "foo@baz.com"},
				EmailCC: []string{"baz@bar.com"},
				Message: json.RawMessage(`{"foo":"bar"}`),
			},
			expectTplID:  email.TplAddPolicyCoverage,
			expectStatus: http.StatusOK,
		},
		"pass 2 cc": {
			method: http.MethodPost,
			payload: handlers.AddPolicyCoverageReq{
				EmailTo: []string{"baz@bar.com"},
				EmailCC: []string{"foo@bar.com", "foo@baz.com"},
				Message: json.RawMessage(`{"foo":"bar"}`),
			},
			expectTplID:  email.TplAddPolicyCoverage,
			expectStatus: http.StatusOK,
		},
		"pass 2 to and 2 cc": {
			method: http.MethodPost,
			payload: handlers.AddPolicyCoverageReq{
				EmailTo: []string{"baz@bar.com", "fizz@buzz.com"},
				EmailCC: []string{"foo@bar.com", "foo@baz.com"},
				Message: json.RawMessage(`{"foo":"bar"}`),
			},
			expectTplID:  email.TplAddPolicyCoverage,
			expectStatus: http.StatusOK,
		},
		"pass many": {
			method: http.MethodPost,
			payload: handlers.AddPolicyCoverageReq{
				EmailTo: []string{"baz@bar.com", "fizz@buzz.com", "beep@boop.com"},
				EmailCC: []string{"foo@bar.com", "foo@baz.com", "zig@zag.com"},
				Message: json.RawMessage(`{"foo":"bar"}`),
			},
			expectTplID:  email.TplAddPolicyCoverage,
			expectStatus: http.StatusOK,
		},
		"fail invalid method": {
			method:       http.MethodGet,
			payload:      handlers.AddPolicyCoverageReq{},
			expectStatus: http.StatusMethodNotAllowed,
		},
	}

	testFactory := func(tc testCase) func(*testing.T) {
		return func(t *testing.T) {

			payload, err := json.Marshal(tc.payload)
			if err != nil {
				t.Fatalf("could nor marshal payload to json: %v", err)
			}

			testEmail := mockemail.NewClient()
			defer testEmail.FlushSendLogs()

			w := httptest.NewRecorder()
			req := httptest.NewRequest(tc.method, "/api/comms/add-policy-coverage", bytes.NewReader(payload))

			handlers.AddPolicyCoverage(testEmail)(w, req)

			resp := w.Result()

			if resp.StatusCode != tc.expectStatus {
				t.Fatalf("expected status %v but got %v", tc.expectStatus, resp.StatusCode)
			}

			if resp.StatusCode == http.StatusOK {

				if testEmail.SendLogs().IsEmpty() {
					t.Fatalf("expected email log but got empty")
				}

				lastEmail := testEmail.SendLogs().Last()

				if !slices.Equal(lastEmail.ExtractTos(), tc.payload.EmailTo) {
					t.Fatalf("expected to %v but got %v", lastEmail.ExtractTos(), tc.payload.EmailTo)
				}

				if !slices.Equal(lastEmail.ExtractCCs(), tc.payload.EmailCC) {
					t.Fatalf("expected cc %v but got %v", lastEmail.ExtractCCs(), tc.payload.EmailCC)
				}

				if string(lastEmail.ExtractMessage()) != string(tc.payload.Message) {
					t.Fatalf("expected message %v but got %v", tc.payload.Message, lastEmail.ExtractMessage())
				}

				if lastEmail.ExtractTplID() != tc.expectTplID {
					t.Fatalf("expected tpl %v but got %v", tc.expectTplID, lastEmail.ExtractTplID())
				}
			}
		}
	}

	for name, tc := range testCases {
		t.Run(name, testFactory(tc))
	}
}

func TestAddPolicyAddress(t *testing.T) {

	t.Parallel()

	type testCase struct {
		method       string
		payload      handlers.AddPolicyAddressReq
		expectTplID  email.TplID
		expectStatus int
	}

	testCases := map[string]testCase{
		"pass": {
			method: http.MethodPost,
			payload: handlers.AddPolicyAddressReq{
				EmailTo: "foo@bar.com",
				Message: json.RawMessage(`{"foo":"bar"}`),
			},
			expectTplID:  email.TplAddPolicyAddress,
			expectStatus: http.StatusOK,
		},
		"fail invalid method": {
			method:       http.MethodGet,
			payload:      handlers.AddPolicyAddressReq{},
			expectStatus: http.StatusMethodNotAllowed,
		},
	}

	testFactory := func(tc testCase) func(*testing.T) {
		return func(t *testing.T) {

			payload, err := json.Marshal(tc.payload)
			if err != nil {
				t.Fatalf("could nor marshal payload to json: %v", err)
			}

			testEmail := mockemail.NewClient()
			defer testEmail.FlushSendLogs()

			w := httptest.NewRecorder()
			req := httptest.NewRequest(tc.method, "/api/comms/add-policy-address", bytes.NewReader(payload))

			handlers.AddPolicyAddress(testEmail)(w, req)

			resp := w.Result()

			if resp.StatusCode != tc.expectStatus {
				t.Fatalf("expected status %v but got %v", tc.expectStatus, resp.StatusCode)
			}

			if resp.StatusCode == http.StatusOK {

				if testEmail.SendLogs().IsEmpty() {
					t.Fatalf("expected email log but got empty")
				}

				lastEmail := testEmail.SendLogs().Last()

				if !slices.Equal(lastEmail.ExtractTos(), []string{tc.payload.EmailTo}) {
					t.Fatalf("expected to %v but got %v", lastEmail.ExtractTos(), tc.payload.EmailTo)
				}

				if string(lastEmail.ExtractMessage()) != string(tc.payload.Message) {
					t.Fatalf("expected message %v but got %v", tc.payload.Message, lastEmail.ExtractMessage())
				}

				if lastEmail.ExtractTplID() != tc.expectTplID {
					t.Fatalf("expected tpl %v but got %v", tc.expectTplID, lastEmail.ExtractTplID())
				}
			}
		}
	}

	for name, tc := range testCases {
		t.Run(name, testFactory(tc))
	}
}
