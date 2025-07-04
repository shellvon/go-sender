package telegram_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/shellvon/go-sender/core"
	"github.com/shellvon/go-sender/providers/telegram"
)

// rewriteRoundTripper rewrites request scheme/host to the test server while preserving the original path.
// This allows us to verify the exact endpoint chosen by the transformer without hitting real Telegram servers.
type rewriteRoundTripper struct {
	base         http.RoundTripper
	targetHost   string
	targetScheme string
}

func (rt rewriteRoundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	req.URL.Host = rt.targetHost
	req.URL.Scheme = rt.targetScheme
	return rt.base.RoundTrip(req)
}

// TestProvider_Send_AllMessageTypes verifies that the Telegram provider can transform and send
// every supported message type, and that the correct Telegram Bot API endpoint is selected.
func TestProvider_Send_AllMessageTypes(t *testing.T) {
	// Prepare provider configuration.
	cfg := telegram.Config{
		Accounts: []*telegram.Account{
			{
				BaseAccount: core.BaseAccount{
					AccountMeta: core.AccountMeta{Name: "primary"},
					Credentials: core.Credentials{APIKey: "bot123:token"},
				},
			},
		},
	}
	provider, err := telegram.New(cfg)
	if err != nil {
		t.Fatalf("failed to create telegram provider: %v", err)
	}

	// Table-driven definition for each message type we want to cover.
	cases := []struct {
		name     string
		buildMsg func() telegram.Message
		endpoint string // expected suffix in request path
	}{
		{
			name: "text",
			buildMsg: func() telegram.Message {
				return telegram.NewTextMessage("100", "hello")
			},
			endpoint: "/sendMessage",
		},
		{
			name: "photo",
			buildMsg: func() telegram.Message {
				return telegram.NewPhotoMessage("100", "http://example.com/p.png")
			},
			endpoint: "/sendPhoto",
		},
		{
			name: "audio",
			buildMsg: func() telegram.Message {
				return telegram.NewAudioMessage("100", "http://example.com/a.mp3")
			},
			endpoint: "/sendAudio",
		},
		{
			name: "voice",
			buildMsg: func() telegram.Message {
				return telegram.NewVoiceMessage("100", "http://example.com/v.ogg")
			},
			endpoint: "/sendVoice",
		},
		{
			name: "document",
			buildMsg: func() telegram.Message {
				return telegram.NewDocumentMessage("100", "http://example.com/doc.pdf")
			},
			endpoint: "/sendDocument",
		},
		{
			name: "video",
			buildMsg: func() telegram.Message {
				return telegram.NewVideoMessage("100", "http://example.com/video.mp4")
			},
			endpoint: "/sendVideo",
		},
		{
			name: "animation",
			buildMsg: func() telegram.Message {
				return telegram.NewAnimationMessage("100", "http://example.com/anim.gif")
			},
			endpoint: "/sendAnimation",
		},
		{
			name: "video_note",
			buildMsg: func() telegram.Message {
				return telegram.NewVideoNoteMessage("100", "http://example.com/vn.mp4")
			},
			endpoint: "/sendVideoNote",
		},
		{
			name: "location",
			buildMsg: func() telegram.Message {
				return telegram.NewLocationMessage("100", 12.34, 56.78)
			},
			endpoint: "/sendLocation",
		},
		{
			name: "contact",
			buildMsg: func() telegram.Message {
				return telegram.NewContactMessage("100", "+1234567890", "Bob")
			},
			endpoint: "/sendContact",
		},
		{
			name: "poll",
			buildMsg: func() telegram.Message {
				opts := []telegram.InputPollOption{{Text: "Yes"}, {Text: "No"}}
				return telegram.NewPollMessage("100", "Are you OK?", opts)
			},
			endpoint: "/sendPoll",
		},
		{
			name: "dice",
			buildMsg: func() telegram.Message {
				return telegram.NewDiceMessage("100")
			},
			endpoint: "/sendDice",
		},
		{
			name: "venue",
			buildMsg: func() telegram.Message {
				return telegram.NewVenueMessage("100", 10.0, 11.0, "A place", "Somewhere")
			},
			endpoint: "/sendVenue",
		},
	}

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			var gotPath string
			// Spin up an HTTP test server to capture requests.
			ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				gotPath = r.URL.Path
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusOK)
				_, _ = w.Write([]byte(`{"ok": true}`))
			}))
			defer ts.Close()

			// Rewrite outgoing requests to our test server.
			parts := strings.SplitN(ts.URL, "://", 2)
			scheme, host := parts[0], parts[1]
			client := &http.Client{
				Transport: rewriteRoundTripper{
					base:         http.DefaultTransport,
					targetHost:   host,
					targetScheme: scheme,
				},
			}

			// Execute send.
			if sendErr := provider.Send(context.Background(), tt.buildMsg(), &core.ProviderSendOptions{HTTPClient: client}); sendErr != nil {
				t.Fatalf("Send() returned error: %v", sendErr)
			}

			// Validate endpoint mapping.
			if !strings.HasSuffix(gotPath, tt.endpoint) {
				t.Errorf("expected request to end with %s, got %s", tt.endpoint, gotPath)
			}
		})
	}
}
