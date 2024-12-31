package entities

type Bot struct {
	Name     string
	Script   string
	UserName string
}

func NewBot(name string, script string, userName string) *Bot {
	return &Bot{
		Name:     name,
		Script:   script,
		UserName: userName,
	}
}
