package dingtalk

import "github.com/shellvon/go-sender/core"

// Account represents a DingTalk bot account.
// It follows the three-tier design: AccountMeta + Credentials + extra
//   - AccountMeta: Name, Weight, Disabled (from core.BaseAccount)
//   - Credentials: APIKey (access token), APISecret (signing secret), AppID (optional) (from core.BaseAccount)
//   - Extra: No additional fields needed for DingTalk bot
//
// It embeds core.BaseAccount so it automatically satisfies core.BasicAccount
// and core.Selectable interfaces.
type Account struct {
	core.BaseAccount
}
