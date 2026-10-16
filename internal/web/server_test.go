package web

import (
	"io/fs"
	"net/http"
	"net/http/httptest"
	"testing"
)

// Every page must parse together with the shared layout and render without errors.
func TestPagesRender(t *testing.T) {
	templateFS, err := fs.Sub(content, "templates")
	if err != nil {
		t.Fatal(err)
	}

	renderer, err := newTemplateRenderer(templateFS)
	if err != nil {
		t.Fatalf("templates do not parse: %v", err)
	}

	data := pageData{Title: "Test", Version: "v1.0.0", IsLoggedIn: true, Username: "admin", RefreshSeconds: 60, ScanRadius: 30}
	for _, page := range allPages {
		t.Run(page, func(t *testing.T) {
			data.Page = page
			recorder := httptest.NewRecorder()
			if err := renderer.Instance(page, data).Render(recorder); err != nil {
				t.Fatalf("render failed: %v", err)
			}
			if recorder.Body.Len() == 0 {
				t.Fatal("page is empty")
			}
		})
	}
}

func TestMaskSecrets(t *testing.T) {
	config := map[string]interface{}{
		"DiscordWebHookURL": "https://discord.example/secret",
		"SlackWebHookURL":   "",
		"NtfyTopic":         "my-topic",
	}

	maskSecrets(config)

	if config["DiscordWebHookURL"] != "configured" {
		t.Errorf("secret was not masked: %v", config["DiscordWebHookURL"])
	}
	if config["SlackWebHookURL"] != "" {
		t.Errorf("an empty secret should stay empty: %v", config["SlackWebHookURL"])
	}
	if config["NtfyTopic"] != "my-topic" {
		t.Errorf("a normal setting must not change: %v", config["NtfyTopic"])
	}
}

// A session cookie that cannot be verified (for example after a restart) must not break the page.
func TestInvalidSessionCookieIsIgnored(t *testing.T) {
	server := NewServer(Config{})
	templateFS, err := fs.Sub(content, "templates")
	if err != nil {
		t.Fatal(err)
	}
	renderer, err := newTemplateRenderer(templateFS)
	if err != nil {
		t.Fatal(err)
	}
	server.engine.HTMLRender = renderer

	request := httptest.NewRequest("GET", "/login", nil)
	request.AddCookie(&http.Cookie{Name: "jetspotter_session", Value: "not-a-valid-cookie"})
	recorder := httptest.NewRecorder()
	server.engine.ServeHTTP(recorder, request)

	if recorder.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", recorder.Code)
	}
}
