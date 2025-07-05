package providers

// Common message type string literals reused across multiple provider implementations.
// Keeping them in one place avoids magic strings scattered in different packages.
const (
	MsgTypeText         = "text"
	MsgTypeMarkdown     = "markdown"
	MsgTypeImage        = "image"
	MsgTypeLink         = "link"
	MsgTypeNews         = "news"
	MsgTypeFile         = "file"
	MsgTypeVoice        = "voice"
	MsgTypeActionCard   = "actionCard"
	MsgTypeFeedCard     = "feedCard"
	MsgTypeTemplateCard = "template_card"
)
