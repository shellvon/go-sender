package wecombot

import "github.com/shellvon/go-sender/core"

// Account represents a WeCom Bot account.
// It follows the three-tier design: AccountMeta + Credentials + extra
//   - AccountMeta: Name, Weight, Disabled (from core.BaseAccount)
//   - Credentials: APIKey (bot key), APISecret (optional), AppID (optional) (from core.BaseAccount)
//   - Extra: No additional fields needed for WeCom Bot
//
// It embeds core.BaseAccount so it automatically satisfies core.BasicAccount
// and core.Selectable interfaces.
type Account struct {
	core.BaseAccount
}
