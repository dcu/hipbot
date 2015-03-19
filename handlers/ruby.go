package handlers

import (
	"github.com/daneharrigan/hipchat"
	"os/exec"
	"strings"
)

type RubyHandler struct {
}

func (ruby *RubyHandler) Matches(message *hipchat.Message) bool {
	return strings.HasPrefix(message.Body, "ruby:")
}

func (ruby *RubyHandler) Process(client *hipchat.Client, roomId string, message *hipchat.Message) {
	code := strings.Replace(message.Body, "ruby:", "", 1)
	code = strings.Replace(code, `\`, `\\`, -1)
	code = strings.Replace(code, `'`, `\'`, -1)

	script := `Timeout.timeout(1) { puts Object.module_eval('$SAFE=3;` + code + `') }`

	println("Executing command:", script)
	commandName := "ruby"
	args := []string{"-rtimeout", "-e", script}

	cmd := exec.Command(commandName, args...)
	output, err := cmd.CombinedOutput()
	if err != nil {
		client.Say(roomId, "Ruby", "Error: "+string(output))
	} else {
		client.Say(roomId, "Ruby", string(output))
	}
}
