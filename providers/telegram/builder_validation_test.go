package telegram_test

import (
	"testing"

	"github.com/shellvon/go-sender/core"
	"github.com/shellvon/go-sender/providers/telegram"
)

// TestBuilder_Validate ensures that messages built via builders pass their Validate() checks.
func TestBuilder_Validate(t *testing.T) {
	tests := []struct {
		name string
		msg  telegram.Message
	}{
		{
			name: "text",
			msg:  telegram.Text().Chat("123").Text("hello world").Build(),
		},
		{
			name: "photo",
			msg:  telegram.Photo().Chat("123").File("http://example.com/p.png").Build(),
		},
		{
			name: "audio",
			msg:  telegram.Audio().Chat("123").File("http://example.com/a.mp3").Build(),
		},
		{
			name: "voice",
			msg:  telegram.Voice().Chat("123").File("http://example.com/v.ogg").Build(),
		},
		{
			name: "document",
			msg:  telegram.Document().Chat("123").File("http://example.com/doc.pdf").Build(),
		},
		{
			name: "video",
			msg:  telegram.Video().Chat("123").File("http://example.com/video.mp4").Build(),
		},
		{
			name: "animation",
			msg:  telegram.Animation().Chat("123").File("http://example.com/anim.gif").Build(),
		},
		{
			name: "video_note",
			msg:  telegram.VideoNote().Chat("123").File("http://example.com/vn.mp4").Build(),
		},
		{
			name: "location",
			msg:  telegram.Location().Chat("123").Coordinates(12.34, 56.78).Build(),
		},
		{
			name: "contact",
			msg:  telegram.Contact().Chat("123").Phone("+1234567890").FirstName("Bob").Build(),
		},
		{
			name: "dice",
			msg:  telegram.Dice().Chat("123").Build(),
		},
	}

	for _, tt := range tests {
		// 安全地验证消息（如果支持）
		if validatable, ok := any(tt.msg).(core.Validatable); ok {
			if err := validatable.Validate(); err != nil {
				t.Errorf("Validate failed for %s: %v", tt.name, err)
			}
		}
	}
}
