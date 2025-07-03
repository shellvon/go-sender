package dingtalk

// LinkBuilder provides a fluent API to construct DingTalk link messages.
//
// Example:
//
//	msg := dingtalk.NewLinkBuilder().
//	         Title("GitHub").
//	         Text("Repository link").
//	         MessageURL("https://github.com").
//	         PicURL("https://github.githubassets.com/images/modules/logos_page/GitHub-Mark.png").
//	         Build()
type LinkBuilder struct {
	title      string
	text       string
	messageURL string
	picURL     string
}

// Link creates a new LinkBuilder instance.
func Link() *LinkBuilder { return &LinkBuilder{} }

// Title sets the link title.
func (b *LinkBuilder) Title(t string) *LinkBuilder { b.title = t; return b }

// Text sets the link text description.
func (b *LinkBuilder) Text(desc string) *LinkBuilder { b.text = desc; return b }

// MessageURL sets the target URL of the link.
func (b *LinkBuilder) MessageURL(url string) *LinkBuilder { b.messageURL = url; return b }

// PicURL sets the optional picture URL.
func (b *LinkBuilder) PicURL(url string) *LinkBuilder { b.picURL = url; return b }

// Build assembles a *LinkMessage.
func (b *LinkBuilder) Build() *LinkMessage {
	return &LinkMessage{
		BaseMessage: BaseMessage{MsgType: TypeLink},
		Link: LinkContent{
			Title:      b.title,
			Text:       b.text,
			MessageURL: b.messageURL,
			PicURL:     b.picURL,
		},
	}
}
