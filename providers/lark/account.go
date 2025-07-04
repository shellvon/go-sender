package lark

import "github.com/shellvon/go-sender/core"

// Account represents a Lark/Feishu bot account.
// It follows the three-tier design: AccountMeta + Credentials + extra
//   - AccountMeta: Name, Weight, Disabled (from core.BaseAccount)
//   - Credentials: APIKey (bot access token), APISecret (optional secret), AppID (optional) (from core.BaseAccount)
//   - Extra: No additional fields needed for Lark bot
//
// It embeds core.BaseAccount so it automatically satisfies core.BasicAccount
// and core.Selectable interfaces.
type Account struct {
	core.BaseAccount
}
