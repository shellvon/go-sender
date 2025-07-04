package telegram

import "github.com/shellvon/go-sender/core"

// Account represents a Telegram bot account.
// It follows the three-tier design: AccountMeta + Credentials + extra
//   - AccountMeta: Name, Weight, Disabled (from core.BaseAccount)
//   - Credentials: APIKey (bot token), APISecret (optional), AppID (optional) (from core.BaseAccount)
//   - Extra: No additional fields needed for Telegram bot
//
// It embeds core.BaseAccount so it automatically satisfies core.BasicAccount
// and core.Selectable interfaces.
type Account struct {
	core.BaseAccount
}
