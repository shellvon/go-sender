package dingtalk

// ActionCardBuilder provides a fluent API to construct DingTalk action-card messages.
// Both single-button and multi-button cards are supported.
//
// Example (single button):
//
//	msg := dingtalk.NewActionCardBuilder().
//	         Title("Deploy Finished").
//	         Text("### deployment done").
//	         SingleButton("View", "https://ci.example.com").
//	         Build()
//
// Example (multi buttons):
//
//	msg := dingtalk.NewActionCardBuilder().
//	         Title("Choose").
//	         Text("select a link").
//	         BtnOrientation("1").
//	         AddButton("Google", "https://google.com").
//	         AddButton("Bing", "https://bing.com").
//	         Build()
type ActionCardBuilder struct {
	title          string
	text           string
	btnOrientation string
	singleTitle    string
	singleURL      string
	buttons        []ActionCardButton
}

// ActionCard creates a new ActionCardBuilder instance.
func ActionCard() *ActionCardBuilder { return &ActionCardBuilder{} }

// Title sets the card title.
func (b *ActionCardBuilder) Title(t string) *ActionCardBuilder { b.title = t; return b }

// Text sets the markdown body of the card.
func (b *ActionCardBuilder) Text(markdown string) *ActionCardBuilder {
	b.text = markdown
	return b
}

// BtnOrientation sets button orientation ("0" = vertical, "1" = horizontal).
func (b *ActionCardBuilder) BtnOrientation(o string) *ActionCardBuilder {
	b.btnOrientation = o
	return b
}

// SingleButton sets the single-button variant.
// Based on SendActionCardParams from DingTalk API
// https://open.dingtalk.com/document/orgapp/custom-bot-send-message-type
//   - singleTitle is required.
//   - singleURL is required.
//   - btnOrientation is optional, default is "0".
func (b *ActionCardBuilder) SingleButton(title, url string) *ActionCardBuilder {
	b.singleTitle = title
	b.singleURL = url
	return b
}

// AddButton appends one button for the multi-button variant.
func (b *ActionCardBuilder) AddButton(title, url string) *ActionCardBuilder {
	b.buttons = append(b.buttons, ActionCardButton{Title: title, ActionURL: url})
	return b
}

// Buttons sets all buttons at once (overwrites additions).
func (b *ActionCardBuilder) Buttons(btns []ActionCardButton) *ActionCardBuilder {
	b.buttons = btns
	return b
}

// Build assembles a *ActionCardMessage.
func (b *ActionCardBuilder) Build() *ActionCardMessage {
	return &ActionCardMessage{
		BaseMessage: BaseMessage{MsgType: TypeActionCard},
		ActionCard: ActionCardContent{
			Title:          b.title,
			Text:           b.text,
			BtnOrientation: b.btnOrientation,
			SingleTitle:    b.singleTitle,
			SingleURL:      b.singleURL,
			Btns:           b.buttons,
		},
	}
}
