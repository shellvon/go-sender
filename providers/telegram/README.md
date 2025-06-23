# Telegram Provider

This provider supports sending messages to Telegram chats, groups, and channels via the [Telegram Bot API](https://core.telegram.org/bots/api).

## Features

- **Multiple Account Support**: Configure multiple bot accounts with different strategies (round-robin, random, weighted)
- **Message Types**: Support for all Telegram Bot API message types including text, media, location, contact, poll, and dice
- **Rich Formatting**: Support for HTML and Markdown formatting in text messages
- **File Support**: Support for sending files via file_id or public HTTP URLs
- **Interactive Elements**: Support for polls, dice animations, and custom keyboards

## Configuration

```go
import (
    "github.com/shellvon/go-sender/core"
    "github.com/shellvon/go-sender/providers/telegram"
)

// Create Telegram configuration
config := telegram.Config{
    BaseConfig: core.BaseConfig{
        Strategy: core.StrategyRoundRobin,
    },
    Accounts: []core.Account{
        {
            Name:     "main-bot",
            Key:      "YOUR_BOT_TOKEN",
            Weight:   100,
            Disabled: false,
        },
        {
            Name:     "backup-bot",
            Key:      "YOUR_BACKUP_BOT_TOKEN",
            Weight:   80,
            Disabled: false,
        },
    },
}

// Create provider
provider, err := telegram.New(config)
if err != nil {
    log.Fatalf("Failed to create Telegram provider: %v", err)
}
```

## Message Types

### 1. Text Message

```go
// Simple text message
textMsg := telegram.NewTextMessage("@channel", "Hello from go-sender!")

// Text message with Markdown formatting
textMsg := telegram.NewTextMessage("@channel", "**Bold text** and *italic text*",
    telegram.WithParseMode("Markdown"),
    telegram.WithDisableWebPreview(true),
)

// Text message with entities
textMsg := telegram.NewTextMessage("@channel", "Hello @username",
    telegram.WithEntities([]telegram.MessageEntity{
        {Type: "mention", Offset: 6, Length: 9},
    }),
)
```

### 2. Photo Message

```go
// Photo from URL
photoMsg := telegram.NewPhotoMessage("@channel", "https://example.com/image.jpg",
    telegram.WithCaption("Beautiful image"),
    telegram.WithParseMode("HTML"),
)

// Photo from file_id
photoMsg := telegram.NewPhotoMessage("@channel", "AgACAgIAAxkBAAIB...",
    telegram.WithCaption("Reused image"),
    telegram.WithPhotoHasSpoiler(true),
)
```

### 3. Audio Message

```go
// Audio from URL
audioMsg := telegram.NewAudioMessage("@channel", "https://example.com/audio.mp3",
    telegram.WithAudioTitle("Song Title"),
    telegram.WithAudioPerformer("Artist Name"),
    telegram.WithAudioDuration(180),
)

// Audio from file_id
audioMsg := telegram.NewAudioMessage("@channel", "CQACAgIAAxkBAAIB...",
    telegram.WithCaption("Listen to this!"),
)
```

### 4. Document Message

```go
// Document from URL
docMsg := telegram.NewDocumentMessage("@channel", "https://example.com/document.pdf",
    telegram.WithCaption("Important document"),
    telegram.WithDocumentDisableContentTypeDetection(true),
)

// Document from file_id
docMsg := telegram.NewDocumentMessage("@channel", "BQACAgIAAxkBAAIB...",
    telegram.WithCaption("Shared document"),
)
```

### 5. Video Message

```go
// Video from URL
videoMsg := telegram.NewVideoMessage("@channel", "https://example.com/video.mp4",
    telegram.WithCaption("Amazing video"),
    telegram.WithVideoDuration(30),
    telegram.WithVideoWidth(1920),
    telegram.WithVideoHeight(1080),
    telegram.WithVideoSupportsStreaming(true),
)

// Video from file_id
videoMsg := telegram.NewVideoMessage("@channel", "BAACAgIAAxkBAAIB...",
    telegram.WithVideoHasSpoiler(true),
)
```

### 6. Animation Message

```go
// Animation (GIF) from URL
animationMsg := telegram.NewAnimationMessage("@channel", "https://example.com/animation.gif",
    telegram.WithCaption("Funny animation"),
    telegram.WithAnimationDuration(5),
    telegram.WithAnimationHasSpoiler(true),
)

// Animation from file_id
animationMsg := telegram.NewAnimationMessage("@channel", "CgACAgIAAxkBAAIB...",
    telegram.WithCaption("Reused animation"),
)
```

### 7. Voice Message

```go
// Voice from URL
voiceMsg := telegram.NewVoiceMessage("@channel", "https://example.com/voice.ogg",
    telegram.WithVoiceDuration(10),
    telegram.WithCaption("Voice message"),
)

// Voice from file_id
voiceMsg := telegram.NewVoiceMessage("@channel", "AwACAgIAAxkBAAIB...",
    telegram.WithCaption("Listen to this voice"),
)
```

### 8. Video Note Message

```go
// Video note from URL
videoNoteMsg := telegram.NewVideoNoteMessage("@channel", "https://example.com/videonote.mp4",
    telegram.WithVideoNoteDuration(15),
    telegram.WithVideoNoteLength(240),
)

// Video note from file_id
videoNoteMsg := telegram.NewVideoNoteMessage("@channel", "DQACAgIAAxkBAAIB...",
    telegram.WithVideoNoteDuration(10),
)
```

### 9. Location Message

```go
// Location message
locationMsg := telegram.NewLocationMessage("@channel", 40.7128, -74.0060,
    telegram.WithLocationHorizontalAccuracy(10),
    telegram.WithLocationLivePeriod(3600),
    telegram.WithLocationHeading(90),
    telegram.WithLocationProximityAlertRadius(1000),
)
```

### 10. Contact Message

```go
// Contact message
contactMsg := telegram.NewContactMessage("@channel", "+1234567890", "John Doe",
    telegram.WithContactLastName("Smith"),
    telegram.WithContactVCard("BEGIN:VCARD\nVERSION:3.0\nFN:John Doe\nTEL:+1234567890\nEND:VCARD"),
)
```

### 11. Poll Message

```go
// Regular poll
options := []telegram.InputPollOption{
    {Text: "Option 1"},
    {Text: "Option 2"},
    {Text: "Option 3"},
}
pollMsg := telegram.NewPollMessage("@channel", "What's your favorite color?", options,
    telegram.WithPollIsAnonymous(false),
    telegram.WithPollAllowsMultipleAnswers(true),
)

// Quiz poll
quizMsg := telegram.NewPollMessage("@channel", "What is 2+2?", []telegram.InputPollOption{
    {Text: "3"},
    {Text: "4"},
    {Text: "5"},
},
    telegram.WithPollType("quiz"),
    telegram.WithPollCorrectOptionID(1),
    telegram.WithPollExplanation("The correct answer is 4"),
)
```

### 12. Dice Message

```go
// Dice animation
diceMsg := telegram.NewDiceMessage("@channel",
    telegram.WithDiceEmoji("🎲"),
)

// Dart animation
dartMsg := telegram.NewDiceMessage("@channel",
    telegram.WithDiceEmoji("🎯"),
)
```

## Usage with Sender

```go
import (
    gosender "github.com/shellvon/go-sender"
    "github.com/shellvon/go-sender/core"
    "github.com/shellvon/go-sender/providers/telegram"
)

// Create sender
s := gosender.NewSender(nil)

// Register Telegram provider
telegramProvider, err := telegram.New(config)
if err != nil {
    log.Fatalf("Failed to create Telegram provider: %v", err)
}
s.RegisterProvider(core.ProviderTypeTelegram, telegramProvider, nil)

// Send message
ctx := context.Background()
textMsg := telegram.NewTextMessage("@channel", "Hello from go-sender!")
err = s.Send(ctx, textMsg)
if err != nil {
    log.Printf("Failed to send message: %v", err)
}
```

## Message Options

### Common Options (All Message Types)

- `WithSilent(silent bool)`: Send message silently (no notification sound)
- `WithProtectContent(protect bool)`: Protect message content from forwarding and saving
- `WithAllowPaidBroadcast(allow bool)`: Allow paid broadcast (up to 1000 messages/second)
- `WithMessageEffectID(effectID string)`: Add message effect (private chats only)
- `WithReplyParameters(params ReplyParameters)`: Reply to a specific message
- `WithReplyMarkup(markup interface{})`: Add inline keyboard or reply keyboard
- `WithBusinessConnectionID(id string)`: Send on behalf of business account
- `WithMessageThreadID(id int)`: Target specific message thread (forum topics)

### Media Message Options

- `WithCaption(caption string)`: Add caption to media
- `WithParseMode(mode string)`: Set caption parse mode (HTML, Markdown, MarkdownV2)
- `WithCaptionEntities(entities []MessageEntity)`: Set caption entities
- `WithShowCaptionAboveMedia(show bool)`: Show caption above media instead of below

### Text Message Options

- `WithParseMode(mode string)`: Set text parse mode
- `WithEntities(entities []MessageEntity)`: Set text entities
- `WithDisableWebPreview(disable bool)`: Disable web page preview
- `WithLinkPreviewOptions(options LinkPreviewOptions)`: Configure link preview

### Poll Message Options

- `WithPollQuestionParseMode(mode string)`: Set question parse mode
- `WithPollQuestionEntities(entities []MessageEntity)`: Set question entities
- `WithPollIsAnonymous(anonymous bool)`: Make poll anonymous
- `WithPollType(pollType string)`: Set poll type (regular/quiz)
- `WithPollAllowsMultipleAnswers(allow bool)`: Allow multiple answers
- `WithPollCorrectOptionID(id int)`: Set correct answer for quiz
- `WithPollExplanation(explanation string)`: Add explanation for quiz
- `WithPollExplanationParseMode(mode string)`: Set explanation parse mode
- `WithPollExplanationEntities(entities []MessageEntity)`: Set explanation entities
- `WithPollOpenPeriod(period int)`: Set poll open period in seconds
- `WithPollCloseDate(date int)`: Set poll close date (Unix timestamp)
- `WithPollIsClosed(closed bool)`: Close poll immediately

## File Sending Limitations

According to the [Telegram Bot API documentation](https://core.telegram.org/bots/api#sending-files), this provider has the following limitations:

### File Sources

- **file_id**: Existing file ID on Telegram servers
- **HTTP URL**: Publicly accessible file URL that Telegram can download
- **Local file upload**: Not supported (no multipart/form-data upload)

### File Size Limits

- Photos: 10 MB maximum
- Audio: 50 MB maximum
- Documents: 50 MB maximum
- Videos: 50 MB maximum
- Animations: 50 MB maximum
- Voice messages: 50 MB maximum
- Video notes: 50 MB maximum

### URL Requirements

- Must be HTTPS URLs
- Files must be publicly accessible
- Telegram servers must be able to download the file
- File format must be supported by Telegram

### Best Practices

1. **Prefer file_id**: Use existing file IDs when possible for better performance
2. **Use reliable URLs**: Ensure URLs are stable and accessible
3. **Check file formats**: Verify files are in supported formats
4. **Monitor file sizes**: Stay within Telegram's size limits

## Error Handling

The provider returns detailed error messages including:

- HTTP request errors
- Telegram API error codes and descriptions
- Parameter validation errors

```go
err := provider.Send(ctx, message)
if err != nil {
    // Error format: telegram API error 400: Bad Request: wrong type of the web page content
    log.Printf("Failed to send message: %v", err)
}
```

## Testing

Run unit tests:

```bash
go test ./providers/telegram/...
```

Run integration tests (requires environment variables):

```bash
export TELEGRAM_BOT_TOKEN="your-bot-token"
export TELEGRAM_CHAT_ID="your-chat-id"
go test ./providers/telegram/... -v -run TestTelegramProviderIntegration
```

## API Reference

### Config

- `BaseConfig`: Common configuration fields
  - `Disabled`: Whether the provider is disabled
  - `Strategy`: Selection strategy (round_robin, random, weighted)
- `Accounts`: Array of bot account configurations

### Message Types

All message types implement the `Message` interface and can be used with the provider's `Send` method.

## Related Links

- [Telegram Bot API Documentation](https://core.telegram.org/bots/api)
- [Sending Files Guide](https://core.telegram.org/bots/api#sending-files)
- [Bot Creation Guide](https://core.telegram.org/bots#how-do-i-create-a-bot)
