package email

import "github.com/shellvon/go-sender/core"

// Account represents a single SMTP account configuration.
// It follows the three-tier design: AccountMeta + Credentials + extra
//   - AccountMeta: Name, Weight, Disabled (from core.BaseAccount)
//   - Credentials: APIKey (username), APISecret (password) (from core.BaseAccount)
//   - Extra: Host, Port, From (SMTP-specific configuration)
//
// It embeds core.BaseAccount so it automatically satisfies core.BasicAccount
// and core.Selectable interfaces.
type Account struct {
	core.BaseAccount

	Host string `json:"host"` // SMTP server host
	Port int    `json:"port"` // SMTP port (25/465/587)
	From string `json:"from"` // Default "From" address
}

func (a *Account) Username() string { return a.GetCredentials().APIKey }
func (a *Account) Password() string { return a.GetCredentials().APISecret }
