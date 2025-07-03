package telegram

// baseBuilder provides chainable methods common to all Telegram message builders.
// It follows the same self-type generic pattern used in the SMS module so that each
// concrete builder can embed *baseBuilder[*ConcreteBuilder] and every chain call
// returns the concrete builder type (no more casting or loss of type information).
//
// Fields map 1-1 to BaseMessage; the concrete builder will copy them into the
// resulting Telegram message during Build().

type baseBuilder[T any] struct {
	chatID string
	// Sends the message silently. Users will receive a notification with no sound.
	disableNotification bool
	protectContent      bool
	allowPaidBroadcast  bool

	// optional advanced fields
	businessConnectionID string
	messageThreadID      int

	messageEffectID string
	replyParameters *ReplyParameters
	replyMarkup     ReplyMarkup

	// self holds the concrete builder pointer so that each chain method can
	// return it, preserving fluent API ergonomics.
	self T
}

// Chat sets the chat ID (chat_id) for the message.
// chatID is Unique identifier for the target chat or username of the target channel (in the format @channelusername).
func (b *baseBuilder[T]) Chat(id string) T {
	b.chatID = id
	return b.self
}

// Silent toggles disable_notification.
func (b *baseBuilder[T]) Silent(disable bool) T {
	b.disableNotification = disable
	return b.self
}

// Protect toggles protect_content.
func (b *baseBuilder[T]) Protect(protect bool) T {
	b.protectContent = protect
	return b.self
}

// AllowPaidBroadcast toggles allow_paid_broadcast.
func (b *baseBuilder[T]) AllowPaidBroadcast(allow bool) T {
	b.allowPaidBroadcast = allow
	return b.self
}

// BusinessConnection sets business_connection_id for sending on behalf of a Business connection.
func (b *baseBuilder[T]) BusinessConnection(id string) T {
	b.businessConnectionID = id
	return b.self
}

// Thread sets message_thread_id for forum topic messages.
func (b *baseBuilder[T]) Thread(id int) T {
	b.messageThreadID = id
	return b.self
}

// Effect sets message_effect_id for special message effects (private chats only).
func (b *baseBuilder[T]) Effect(id string) T {
	b.messageEffectID = id
	return b.self
}

// Reply sets reply_parameters to reply to an existing message with advanced options.
func (b *baseBuilder[T]) Reply(params *ReplyParameters) T {
	b.replyParameters = params
	return b.self
}

// Markup attaches a ReplyMarkup (inline keyboard, reply keyboard, etc.).
func (b *baseBuilder[T]) Markup(markup ReplyMarkup) T {
	b.replyMarkup = markup
	return b.self
}

// DisableNotification toggles disable_notification.
// Sends the message silently. Users will receive a notification with no sound.
func (b *baseBuilder[T]) DisableNotification(disable bool) T {
	b.disableNotification = disable
	return b.self
}

// toBaseMessage converts builder state into a ready-to-use BaseMessage struct.
func (b *baseBuilder[T]) toBaseMessage(msgType MessageType) BaseMessage {
	return BaseMessage{
		MsgType:              msgType,
		ChatID:               b.chatID,
		DisableNotification:  b.disableNotification,
		ProtectContent:       b.protectContent,
		AllowPaidBroadcast:   b.allowPaidBroadcast,
		BusinessConnectionID: b.businessConnectionID,
		MessageThreadID:      b.messageThreadID,
		MessageEffectID:      b.messageEffectID,
		ReplyParameters:      b.replyParameters,
		ReplyMarkup:          b.replyMarkup,
	}
}
