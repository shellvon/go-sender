package telegram

import "github.com/shellvon/go-sender/core"

// MessageType represents the type of Telegram message.
type MessageType string

const (
	TypeText      MessageType = "text"
	TypePhoto     MessageType = "photo"
	TypeDocument  MessageType = "document"
	TypeLocation  MessageType = "location"
	TypeContact   MessageType = "contact"
	TypePoll      MessageType = "poll"
	TypeAudio     MessageType = "audio"
	TypeVideo     MessageType = "video"
	TypeAnimation MessageType = "animation"
	TypeVoice     MessageType = "voice"
	TypeVideoNote MessageType = "video_note"
	TypeVenue     MessageType = "venue"
	TypeDice      MessageType = "dice"
)

// MessageEntityType represents the type of a message entity
// Based on Telegram Bot API MessageEntity
//   - https://core.telegram.org/bots/api#messageentity
//
// Type of the entity.
// Can be one of:
//   - "mention" (@username)
//   - "hashtag" (#hashtag)
//   - "cashtag" ($USD)
//   - "bot_command" (/start@jobs_bot)
//   - "url" (https://telegram.org)
//   - "email" (do-not-reply@telegram.org)
//   - "phone_number" (+1-212-555-0123)
//   - "bold" (bold text)
//   - "italic" (italic text)
//   - "underline" (underlined text)
//   - "strikethrough" (strikethrough text)
//   - "spoiler" (spoiler message)
//   - "blockquote" (block quotation)
//   - "expandable_blockquote" (collapsed-by-default block quotation)
//   - "code" (monowidth string)
//   - "pre" (monowidth block)
//   - "text_link" (for clickable text URLs)
//   - "text_mention" (for users without usernames)
//   - "custom_emoji" (for inline custom emoji stickers)
//
// See https://core.telegram.org/bots/api#messageentity for more details.
type MessageEntityType string

const (
	EntityTypeMention              MessageEntityType = "mention"
	EntityTypeHashtag              MessageEntityType = "hashtag"
	EntityTypeCashtag              MessageEntityType = "cashtag"
	EntityTypeBotCommand           MessageEntityType = "bot_command"
	EntityTypeURL                  MessageEntityType = "url"
	EntityTypeEmail                MessageEntityType = "email"
	EntityTypePhoneNumber          MessageEntityType = "phone_number"
	EntityTypeBold                 MessageEntityType = "bold"
	EntityTypeItalic               MessageEntityType = "italic"
	EntityTypeUnderline            MessageEntityType = "underline"
	EntityTypeStrikethrough        MessageEntityType = "strikethrough"
	EntityTypeSpoiler              MessageEntityType = "spoiler"
	EntityTypeBlockquote           MessageEntityType = "blockquote"
	EntityTypeExpandableBlockquote MessageEntityType = "expandable_blockquote"
	EntityTypeCode                 MessageEntityType = "code"
	EntityTypePre                  MessageEntityType = "pre"
	EntityTypeTextLink             MessageEntityType = "text_link"
	EntityTypeTextMention          MessageEntityType = "text_mention"
	EntityTypeCustomEmoji          MessageEntityType = "custom_emoji"
)

// String implements the Stringer interface.
func (t MessageEntityType) String() string {
	return string(t)
}

// IsValid checks if the entity type is valid.
func (t MessageEntityType) IsValid() bool {
	switch t {
	case EntityTypeMention, EntityTypeHashtag, EntityTypeCashtag, EntityTypeBotCommand,
		EntityTypeURL, EntityTypeEmail, EntityTypePhoneNumber, EntityTypeBold,
		EntityTypeItalic, EntityTypeUnderline, EntityTypeStrikethrough, EntityTypeSpoiler,
		EntityTypeBlockquote, EntityTypeExpandableBlockquote, EntityTypeCode,
		EntityTypePre, EntityTypeTextLink, EntityTypeTextMention, EntityTypeCustomEmoji:
		return true
	default:
		return false
	}
}

// DiceEmoji represents the emoji for dice messages
// Based on SendDiceParams from Telegram Bot API
//   - https://core.telegram.org/bots/api#senddice
//
// Emoji on which the dice throw animation is based.
// Must be one of
//   - 🎲 (default)
//   - 🎯
//   - 🏀
//   - ⚽
//   - 🎳
//   - 🎰
type DiceEmoji string

const (
	DiceEmojiDice    DiceEmoji = "🎲" // default
	DiceEmojiTarget  DiceEmoji = "🎯"
	DiceEmojiBasket  DiceEmoji = "🏀"
	DiceEmojiSoccer  DiceEmoji = "⚽"
	DiceEmojiBowling DiceEmoji = "🎳"
	DiceEmojiSlot    DiceEmoji = "🎰"
)

// String implements the Stringer interface.
func (e DiceEmoji) String() string {
	return string(e)
}

// IsValid checks if the dice emoji is valid.
func (e DiceEmoji) IsValid() bool {
	switch e {
	case DiceEmojiDice, DiceEmojiTarget, DiceEmojiBasket, DiceEmojiSoccer, DiceEmojiBowling, DiceEmojiSlot:
		return true
	default:
		return false
	}
}

// MessageEntity represents a special entity that appears in message text
// Based on Telegram Bot API MessageEntity
// https://core.telegram.org/bots/api#messageentity
type MessageEntity struct {
	// Type - Type of the entity.
	// Required. Can be one of the following:
	// - "mention" (@username)
	// - "hashtag" (#hashtag)
	// - "cashtag" ($USD)
	// - "bot_command" (/start@jobs_bot)
	// - "url" (https://telegram.org)
	// - "email" (do-not-reply@telegram.org)
	// - "phone_number" (+1-212-555-0123)
	// - "bold" (bold text)
	// - "italic" (italic text)
	// - "underline" (underlined text)
	// - "strikethrough" (strikethrough text)
	// - "spoiler" (spoiler message)
	// - "blockquote" (block quotation)
	// - "expandable_blockquote" (collapsed-by-default block quotation)
	// - "code" (monowidth string)
	// - "pre" (monowidth block)
	// - "text_link" (for clickable text URLs)
	// - "text_mention" (for users without usernames)
	// - "custom_emoji" (for inline custom emoji stickers)
	Type MessageEntityType `json:"type"`

	// Offset - Offset in UTF-16 code units to the start of the entity.
	// Required. Zero-based position where the entity starts in the message text.
	// Must be within the bounds of the message text.
	Offset int `json:"offset"`

	// Length - Length of the entity in UTF-16 code units.
	// Required. Number of characters that the entity spans.
	// Must be positive and not exceed the message length.
	Length int `json:"length"`

	// URL - For "text_link" entities only, URL that will be opened after user taps on the text.
	// Optional. Only used when Type is "text_link".
	// Must be a valid URL.
	URL string `json:"url,omitempty"`

	// User - For "text_mention" entities only, the mentioned user.
	// Optional. Only used when Type is "text_mention".
	// Contains information about the mentioned user.
	User *User `json:"user,omitempty"`

	// Language - For "pre" entities only, the programming language of the entity text.
	// Optional. Only used when Type is "pre".
	// Used for syntax highlighting.
	Language string `json:"language,omitempty"`

	// CustomEmojiID - For "custom_emoji" entities only, unique identifier of the custom emoji.
	// Optional. Only used when Type is "custom_emoji".
	// Must be a valid custom emoji ID.
	CustomEmojiID string `json:"custom_emoji_id,omitempty"`
}

// Validate validates the MessageEntity
func (e *MessageEntity) Validate() error {
	if !e.Type.IsValid() {
		return core.NewParamError("invalid entity type: " + string(e.Type))
	}
	if e.Offset < 0 {
		return core.NewParamError("offset cannot be negative")
	}
	if e.Length <= 0 {
		return core.NewParamError("length must be positive")
	}

	// Validate type-specific fields
	switch e.Type {
	case EntityTypeTextLink:
		if e.URL == "" {
			return core.NewParamError("URL is required for text_link entities")
		}
	case EntityTypeTextMention:
		if e.User == nil {
			return core.NewParamError("User is required for text_mention entities")
		}
	case EntityTypeCustomEmoji:
		if e.CustomEmojiID == "" {
			return core.NewParamError("CustomEmojiID is required for custom_emoji entities")
		}
	}

	return nil
}

// User represents a Telegram user
// Based on Telegram Bot API User
// https://core.telegram.org/bots/api#user
type User struct {
	// ID - Unique identifier for this user or bot.
	// Required. This number may have more than 32 significant bits.
	ID int64 `json:"id"`

	// IsBot - True, if this user is a bot.
	// Required. Indicates whether the user is a bot or a human user.
	IsBot bool `json:"is_bot"`

	// FirstName - User's or bot's first name.
	// Required. The user's first name as they have set it.
	FirstName string `json:"first_name"`

	// LastName - User's or bot's last name.
	// Optional. The user's last name as they have set it.
	LastName string `json:"last_name,omitempty"`

	// Username - User's or bot's username.
	// Optional. The user's username, if they have one.
	// Username is unique across all users and bots.
	Username string `json:"username,omitempty"`

	// LanguageCode - IETF language tag of the user's language.
	// Optional. The user's language code (e.g., "en", "zh-CN").
	// Used for localization purposes.
	LanguageCode string `json:"language_code,omitempty"`

	// IsPremium - True, if this user is a Telegram Premium user.
	// Optional. Indicates whether the user has Telegram Premium subscription.
	IsPremium bool `json:"is_premium,omitempty"`

	// AddedToAttachmentMenu - True, if this user added the bot to the attachment menu.
	// Optional. Indicates whether the user has added the bot to their attachment menu.
	AddedToAttachmentMenu bool `json:"added_to_attachment_menu,omitempty"`

	// CanJoinGroups - True, if the bot can be invited to groups.
	// Optional. Indicates whether the bot can be added to groups.
	// Returned only in getMe.
	CanJoinGroups bool `json:"can_join_groups,omitempty"`

	// CanReadAllGroupMessages - True, if privacy mode is disabled for the bot.
	// Optional. Indicates whether the bot can read all group messages.
	// Returned only in getMe.
	CanReadAllGroupMessages bool `json:"can_read_all_group_messages,omitempty"`

	// SupportsInlineQueries - True, if the bot supports inline queries.
	// Optional. Indicates whether the bot supports inline mode.
	// Returned only in getMe.
	SupportsInlineQueries bool `json:"supports_inline_queries,omitempty"`
}

// LinkPreviewOptions represents link preview generation options
// Based on Telegram Bot API LinkPreviewOptions
// https://core.telegram.org/bots/api#linkpreviewoptions
type LinkPreviewOptions struct {
	// IsDisabled - True, if the link preview is disabled.
	// Optional. If true, no link preview will be shown for links in the message.
	// Default is false (link previews are enabled).
	IsDisabled bool `json:"is_disabled,omitempty"`

	// URL - URL to use for the link preview.
	// Optional. If specified, this URL will be used for the link preview instead of the first URL found in the message.
	// Must be a valid URL.
	URL string `json:"url,omitempty"`

	// PreferSmallMedia - True, if the media in the link preview is suppposed to be shrunk.
	// Optional. If true, small media will be preferred over large media in link previews.
	// Default is false (large media is preferred).
	PreferSmallMedia bool `json:"prefer_small_media,omitempty"`

	// PreferLargeMedia - True, if the media in the link preview is suppposed to be enlarged.
	// Optional. If true, large media will be preferred over small media in link previews.
	// Default is false (small media is preferred).
	PreferLargeMedia bool `json:"prefer_large_media,omitempty"`

	// ShowAboveText - True, if the link preview must be shown above the message text.
	// Optional. If true, the link preview will be displayed above the message text.
	// Default is false (link preview is shown below the text).
	ShowAboveText bool `json:"show_above_text,omitempty"`
}

// ReplyParameters represents description of the message to reply to
// Based on Telegram Bot API ReplyParameters
// https://core.telegram.org/bots/api#replyparameters
type ReplyParameters struct {
	// MessageID - Identifier of the message that will be replied to in the current chat.
	// Required. The ID of the message to reply to.
	// Must be a valid message ID in the current chat.
	MessageID int64 `json:"message_id"`

	// ChatID - If the message to be replied to is from a different chat.
	// Optional. Unique identifier for the target chat or username of the target channel.
	// If not specified, the message will be replied to in the same chat.
	ChatID string `json:"chat_id,omitempty"`

	// AllowSendingWithoutReply - Pass True if the message should be sent even if the specified message to reply to is not found.
	// Optional. If true, the message will be sent even if the replied message doesn't exist.
	// Default is false (message won't be sent if replied message is not found).
	AllowSendingWithoutReply bool `json:"allow_sending_without_reply,omitempty"`

	// Quote - Quoted part of the message to be replied to.
	// Optional. Text to be quoted in the reply.
	// 0-1024 characters.
	Quote string `json:"quote,omitempty"`

	// QuoteParseMode - Mode for parsing entities in the quote.
	// Optional. Can be "HTML", "Markdown", or "MarkdownV2".
	// If not specified, the quote will be sent as plain text.
	QuoteParseMode string `json:"quote_parse_mode,omitempty"`

	// QuoteEntities - A JSON-serialized list of special entities that appear in the quote.
	// Optional. Can be specified instead of quote_parse_mode.
	// Allows for precise control over quote formatting.
	QuoteEntities []MessageEntity `json:"quote_entities,omitempty"`

	// QuotePosition - Position of the quote in the original message in UTF-16 code units.
	// Optional. Zero-based position where the quote starts in the original message.
	// Must be within the bounds of the original message.
	QuotePosition int `json:"quote_position,omitempty"`
}

// ReplyMarkup represents additional interface options
// Based on Telegram Bot API ReplyMarkup (union type)
// https://core.telegram.org/bots/api#replymarkup
// This is a union type that can be one of:
// - InlineKeyboardMarkup
// - ReplyKeyboardMarkup
// - ReplyKeyboardRemove
// - ForceReply.
type ReplyMarkup interface {
	// This is a union type, so we'll use interface{}
	// In practice, this will be one of the concrete types below
}

// InlineKeyboardMarkup represents an inline keyboard
// Based on Telegram Bot API InlineKeyboardMarkup
// https://core.telegram.org/bots/api#inlinekeyboardmarkup
type InlineKeyboardMarkup struct {
	// InlineKeyboard - Array of button rows, each represented by an Array of InlineKeyboardButton objects.
	// Required. Each row is an array of buttons that will be displayed horizontally.
	// Maximum 8 rows, each row can have up to 8 buttons.
	InlineKeyboard [][]InlineKeyboardButton `json:"inline_keyboard"`
}

// InlineKeyboardButton represents one button of an inline keyboard
// Based on Telegram Bot API InlineKeyboardButton
// https://core.telegram.org/bots/api#inlinekeyboardbutton
type InlineKeyboardButton struct {
	// Text - Label text on the button.
	// Required. The text that will be displayed on the button.
	// 1-64 characters.
	Text string `json:"text"`

	// URL - HTTP or tg:// URL to be opened when the button is pressed.
	// Optional. Opens the specified URL when the button is pressed.
	// Must be a valid HTTP URL or tg:// URL.
	URL string `json:"url,omitempty"`

	// CallbackData - Data to be sent in a callback query to the bot when button is pressed.
	// Optional. 1-64 bytes. Used for handling button presses in the bot.
	// The bot will receive a callback query with this data.
	CallbackData string `json:"callback_data,omitempty"`

	// WebApp - Description of the Web App that will be launched when the user presses the button.
	// Optional. The Web App will be able to send an arbitrary message on behalf of the user.
	// Used for integrating web applications with Telegram.
	WebApp *WebAppInfo `json:"web_app,omitempty"`

	// LoginURL - An HTTP URL used to automatically authorize the user.
	// Optional. Can be used as a replacement for the Telegram Login Widget.
	// Used for OAuth-like authentication flows.
	LoginURL *LoginURL `json:"login_url,omitempty"`

	// SwitchInlineQuery - If set, pressing the button will prompt the user to select one of their chats.
	// Optional. Opens a list of suitable chats. Tapping on a chat will send its chat identifier.
	// Used for switching to inline mode in a specific chat.
	SwitchInlineQuery string `json:"switch_inline_query,omitempty"`

	// SwitchInlineQueryCurrentChat - If set, pressing the button will insert the bot's username and the specified inline query.
	// Optional. Opens a list of suitable chats. Tapping on a chat will send its chat identifier.
	// Used for switching to inline mode in the current chat.
	SwitchInlineQueryCurrentChat string `json:"switch_inline_query_current_chat,omitempty"`

	// SwitchInlineQueryChosenChat - If set, pressing the button will prompt the user to select one of their chats.
	// Optional. Opens a list of suitable chats. Tapping on a chat will send its chat identifier.
	// Used for switching to inline mode in a chosen chat with specific parameters.
	SwitchInlineQueryChosenChat *SwitchInlineQueryChosenChat `json:"switch_inline_query_chosen_chat,omitempty"`

	// CallbackGame - Description of the game that will be launched when the user presses the button.
	// Optional. Used for launching games from inline keyboards.
	// NOTE: This type of button must always be the first button in the first row.
	CallbackGame *CallbackGame `json:"callback_game,omitempty"`

	// Pay - Specify True, to send a Pay button.
	// Optional. Used for payments. NOTE: This type of button must always be the first button in the first row.
	// Used for initiating payments through Telegram.
	Pay bool `json:"pay,omitempty"`
}

// WebAppInfo represents information about a Web App
// Based on Telegram Bot API WebAppInfo
// https://core.telegram.org/bots/api#webappinfo
type WebAppInfo struct {
	// URL - An HTTPS URL of a Web App to be opened with additional data as specified in Initializing Web Apps.
	// Required. The URL of the web application that will be opened.
	// Must be a valid HTTPS URL.
	URL string `json:"url"`
}

// LoginURL represents a parameter of the inline keyboard button used to automatically authorize a user
// Based on Telegram Bot API LoginURL
// https://core.telegram.org/bots/api#loginurl
type LoginURL struct {
	// URL - An HTTP URL to be opened with user authorization data added to the query string.
	// Required. The URL where the user will be redirected after authorization.
	// Must be a valid HTTP URL.
	URL string `json:"url"`

	// ForwardText - New text of the button in forwarded messages.
	// Optional. Text to be displayed on the button when the message is forwarded.
	// 1-64 characters.
	ForwardText string `json:"forward_text,omitempty"`

	// BotUsername - Username of a bot, which will be used for user authorization.
	// Optional. If not specified, the current bot's username will be assumed.
	// The URL domain must be the same as the domain linked with the bot.
	BotUsername string `json:"bot_username,omitempty"`

	// RequestWriteAccess - Pass True to request the permission for your bot to send messages to the user.
	// Optional. If true, the bot will request permission to send messages to the user.
	// Default is false.
	RequestWriteAccess bool `json:"request_write_access,omitempty"`
}

// SwitchInlineQueryChosenChat represents an inline keyboard button that switches the current user to inline mode in a chosen chat
// Based on Telegram Bot API SwitchInlineQueryChosenChat
// https://core.telegram.org/bots/api#switchinlinequerychosenchat
type SwitchInlineQueryChosenChat struct {
	// Query - The default inline query to be inserted in the input field.
	// Optional. If left empty, only the bot's username will be inserted.
	// 1-64 characters.
	Query string `json:"query,omitempty"`

	// AllowUserChats - True, if private chats with users can be chosen.
	// Optional. If true, private chats with users will be included in the chat selection.
	// Default is false.
	AllowUserChats bool `json:"allow_user_chats,omitempty"`

	// AllowBotChats - True, if private chats with bots can be chosen.
	// Optional. If true, private chats with bots will be included in the chat selection.
	// Default is false.
	AllowBotChats bool `json:"allow_bot_chats,omitempty"`

	// AllowGroupChats - True, if group and supergroup chats can be chosen.
	// Optional. If true, group and supergroup chats will be included in the chat selection.
	// Default is false.
	AllowGroupChats bool `json:"allow_group_chats,omitempty"`

	// AllowChannelChats - True, if channel chats can be chosen.
	// Optional. If true, channel chats will be included in the chat selection.
	// Default is false.
	AllowChannelChats bool `json:"allow_channel_chats,omitempty"`
}

// CallbackGame represents a placeholder, currently holds no information
// Based on Telegram Bot API CallbackGame
// https://core.telegram.org/bots/api#callbackgame
// A placeholder, currently holds no information. Use BotFather to set up your game.
type CallbackGame struct{}

// ReplyKeyboardMarkup represents a custom keyboard with reply options
// Based on Telegram Bot API ReplyKeyboardMarkup.
type ReplyKeyboardMarkup struct {
	Keyboard              [][]KeyboardButton `json:"keyboard"`
	IsPersistent          bool               `json:"is_persistent,omitempty"`
	ResizeKeyboard        bool               `json:"resize_keyboard,omitempty"`
	OneTimeKeyboard       bool               `json:"one_time_keyboard,omitempty"`
	InputFieldPlaceholder string             `json:"input_field_placeholder,omitempty"`
	Selective             bool               `json:"selective,omitempty"`
}

// KeyboardButton represents one button of the reply keyboard
// Based on Telegram Bot API KeyboardButton.
type KeyboardButton struct {
	Text            string                     `json:"text"`
	RequestUser     *KeyboardButtonRequestUser `json:"request_user,omitempty"`
	RequestChat     *KeyboardButtonRequestChat `json:"request_chat,omitempty"`
	RequestContact  bool                       `json:"request_contact,omitempty"`
	RequestLocation bool                       `json:"request_location,omitempty"`
	RequestPoll     *KeyboardButtonPollType    `json:"request_poll,omitempty"`
	WebApp          *WebAppInfo                `json:"web_app,omitempty"`
}

// KeyboardButtonRequestUser represents a button that requests the user to create and send a KeyboardButtonRequestUser
// Based on Telegram Bot API KeyboardButtonRequestUser.
type KeyboardButtonRequestUser struct {
	RequestID     int  `json:"request_id"`
	UserIsBot     bool `json:"user_is_bot,omitempty"`
	UserIsPremium bool `json:"user_is_premium,omitempty"`
}

// KeyboardButtonRequestChat represents a button that requests the user to create and send a KeyboardButtonRequestChat
// Based on Telegram Bot API KeyboardButtonRequestChat.
type KeyboardButtonRequestChat struct {
	RequestID               int                      `json:"request_id"`
	ChatIsChannel           bool                     `json:"chat_is_channel"`
	ChatIsForum             bool                     `json:"chat_is_forum,omitempty"`
	ChatHasUsername         bool                     `json:"chat_has_username,omitempty"`
	ChatIsCreated           bool                     `json:"chat_is_created,omitempty"`
	UserAdministratorRights *ChatAdministratorRights `json:"user_administrator_rights,omitempty"`
	BotAdministratorRights  *ChatAdministratorRights `json:"bot_administrator_rights,omitempty"`
	BotIsMember             bool                     `json:"bot_is_member,omitempty"`
}

// ChatAdministratorRights represents the rights of an administrator in a chat
// Based on Telegram Bot API ChatAdministratorRights.
type ChatAdministratorRights struct {
	IsAnonymous         bool `json:"is_anonymous,omitempty"`
	CanManageChat       bool `json:"can_manage_chat,omitempty"`
	CanDeleteMessages   bool `json:"can_delete_messages,omitempty"`
	CanManageVideoChats bool `json:"can_manage_video_chats,omitempty"`
	CanRestrictMembers  bool `json:"can_restrict_members,omitempty"`
	CanPromoteMembers   bool `json:"can_promote_members,omitempty"`
	CanChangeInfo       bool `json:"can_change_info,omitempty"`
	CanInviteUsers      bool `json:"can_invite_users,omitempty"`
	CanPostMessages     bool `json:"can_post_messages,omitempty"`
	CanEditMessages     bool `json:"can_edit_messages,omitempty"`
	CanPinMessages      bool `json:"can_pin_messages,omitempty"`
	CanPostStories      bool `json:"can_post_stories,omitempty"`
	CanEditStories      bool `json:"can_edit_stories,omitempty"`
	CanDeleteStories    bool `json:"can_delete_stories,omitempty"`
	CanManageTopics     bool `json:"can_manage_topics,omitempty"`
}

// KeyboardButtonPollType represents the type of poll that is requested to be created
// Based on Telegram Bot API KeyboardButtonPollType.
type KeyboardButtonPollType struct {
	Type PollType `json:"type,omitempty"`
}

// ReplyKeyboardRemove represents a reply keyboard that is requested to be removed
// Based on Telegram Bot API ReplyKeyboardRemove.
type ReplyKeyboardRemove struct {
	RemoveKeyboard bool `json:"remove_keyboard"`
	Selective      bool `json:"selective,omitempty"`
}

// ForceReply represents a reply keyboard that is requested to be shown to the user
// Based on Telegram Bot API ForceReply.
type ForceReply struct {
	ForceReply            bool   `json:"force_reply"`
	InputFieldPlaceholder string `json:"input_field_placeholder,omitempty"`
	Selective             bool   `json:"selective,omitempty"`
}

// InputPollOption represents an answer option in a poll
// Based on Telegram Bot API InputPollOption
// https://core.telegram.org/bots/api#inputpolloption
type InputPollOption struct {
	// Option text, 1-100 characters
	Text string `json:"text"`
	// Optional. Mode for parsing entities in the text. See formatting options for more details. Currently, only custom emoji entities are allowed
	TextParseMode string `json:"text_parse_mode,omitempty"`
	// Optional. A JSON-serialized list of special entities that appear in the poll option text. It can be specified instead of text_parse_mode
	TextEntities []MessageEntity `json:"text_entities,omitempty"`
}

// Message interface definition.
type Message interface {
	core.Message
	GetMsgType() MessageType
}

// ParseMode defines the supported parse modes for Telegram messages.
//   - https://core.telegram.org/bots/api#formatting-options
//
// Can be one of:
//   - "Markdown"
//   - "MarkdownV2"
//   - "HTML"
type ParseMode string

const (
	// ParseModeMarkdown is the legacy Markdown format (version 1).
	ParseModeMarkdown ParseMode = "Markdown"
	// ParseModeMarkdownV2 is the modern Markdown format (version 2).
	ParseModeMarkdownV2 ParseMode = "MarkdownV2"
	// ParseModeHTML is the HTML format.
	ParseModeHTML ParseMode = "HTML"
)

// String implements the Stringer interface.
func (p ParseMode) String() string {
	return string(p)
}

// IsValid checks if the parse mode is valid.
func (p ParseMode) IsValid() bool {
	switch p {
	case ParseModeMarkdown, ParseModeMarkdownV2, ParseModeHTML:
		return true
	default:
		return false
	}
}

// PollType represents the type of a poll
// Based on Telegram Bot API PollType
// https://core.telegram.org/bots/api#sendpoll
//
// Can be one of:
//   - "regular" (default)
//   - "quiz"
//
// See https://core.telegram.org/bots/api#sendpoll for more details.
type PollType string

const (
	PollTypeRegular PollType = "regular"
	PollTypeQuiz    PollType = "quiz"
)

// String implements the Stringer interface.
func (p PollType) String() string {
	return string(p)
}

// IsValid checks if the poll type is valid.
func (p PollType) IsValid() bool {
	switch p {
	case PollTypeRegular, PollTypeQuiz:
		return true
	default:
		return false
	}
}
