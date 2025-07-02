package telegram_test

import (
	"testing"

	"github.com/shellvon/go-sender/providers/telegram"
)

func TestNewVideoNoteMessageWithBuilder(t *testing.T) {
	chatID := "123456789"
	videoNote := "test_video_note.mp4"

	// Test basic creation
	msg := telegram.NewVideoNoteMessageWithBuilder(chatID, videoNote)
	if msg == nil {
		t.Fatal("NewVideoNoteMessageWithBuilder returned nil")
	}
	if msg.GetBase().MsgType != telegram.TypeVideoNote {
		t.Errorf("Expected MsgType to be TypeVideoNote, got %v", msg.GetBase().MsgType)
	}
	if msg.GetBase().ChatID != chatID {
		t.Errorf("Expected ChatID to be '%s', got '%s'", chatID, msg.GetBase().ChatID)
	}
	if msg.VideoNote != videoNote {
		t.Errorf("Expected VideoNote to be '%s', got '%s'", videoNote, msg.VideoNote)
	}

	// Test with options
	msg = telegram.NewVideoNoteMessageWithBuilder(chatID, videoNote,
		telegram.WithVideoNoteDuration(10),
		telegram.WithVideoNoteLength(320),
		telegram.WithVideoNoteThumbnail("thumb.jpg"),
	)
	if msg == nil {
		t.Fatal("NewVideoNoteMessageWithBuilder with options returned nil")
	}
	if msg.Duration != 10 {
		t.Errorf("Expected Duration to be 10, got %d", msg.Duration)
	}
	if msg.Length != 320 {
		t.Errorf("Expected Length to be 320, got %d", msg.Length)
	}
	if msg.Thumbnail != "thumb.jpg" {
		t.Errorf("Expected Thumbnail to be 'thumb.jpg', got '%s'", msg.Thumbnail)
	}
}

func TestMediaMessageBuilder_BuildMediaMessage_VideoNote(t *testing.T) {
	chatID := "123456789"
	videoNote := "test_video_note.mp4"

	builder := telegram.NewMediaMessageBuilder(telegram.TypeVideoNote, chatID, videoNote)
	msg := builder.BuildMediaMessage()

	// Should return VideoNoteMessage
	if videoNoteMsg, ok := msg.(*telegram.VideoNoteMessage); ok {
		if videoNoteMsg.GetBase().MsgType != telegram.TypeVideoNote {
			t.Errorf("Expected MsgType to be TypeVideoNote, got %v", videoNoteMsg.GetBase().MsgType)
		}
		if videoNoteMsg.GetBase().ChatID != chatID {
			t.Errorf("Expected ChatID to be '%s', got '%s'", chatID, videoNoteMsg.GetBase().ChatID)
		}
		if videoNoteMsg.VideoNote != videoNote {
			t.Errorf("Expected VideoNote to be '%s', got '%s'", videoNote, videoNoteMsg.VideoNote)
		}
	} else {
		t.Fatal("Expected VideoNoteMessage, got different type")
	}
}
