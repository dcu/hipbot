package handlers

import (
	"github.com/dcu/hipbot/xmpp"
	"os/exec"
	"strings"
)

type RubyHandler struct {
}

func (ruby *RubyHandler) Matches(message *xmpp.Chat) bool {
	return strings.HasPrefix(message.Text, "ruby:")
}

func (ruby *RubyHandler) Process(client *xmpp.Client, roomId string, message *xmpp.Chat) {
	code := strings.Replace(message.Text, "ruby:", "", 1)
	code = strings.Replace(code, `\`, `\\`, -1)
	code = strings.Replace(code, `'`, `\'`, -1)

	script := `Timeout.timeout(1) { puts Object.module_eval('$SAFE=3;` + code + `') }`

	println("Executing command:", script)
	commandName := "ruby"
	args := []string{"-rtimeout", "-e", script}

	cmd := exec.Command(commandName, args...)
	output, err := cmd.CombinedOutput()
	if err != nil {
		client.Send(xmpp.Chat{
			Remote: roomId,
			Text:   "Error: " + string(output),
			Type:   "groupchat",
		})
	} else {
		client.Send(xmpp.Chat{
			Remote: roomId,
			Text:   string(output),
			Type:   "groupchat",
		})
	}
}
