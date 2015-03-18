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
	code = strings.Replace(code, "'", `\'`, -1)

	script := `puts Object.module_eval('$SAFE=3;` + code + `')`

	println("Executing command:", script)
	commandName := "ruby"
	args := []string{"-e", script}

	cmd := exec.Command(commandName, args...)
	output, err := cmd.Output()
	if err != nil {
		client.Say(roomId, "Ruby", "Error: "+err.Error())
	} else {
		client.Say(roomId, "Ruby", string(output))
	}
}
