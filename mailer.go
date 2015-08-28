/*

Package mailer implements utility function to send mail.

*/
package mailer

import (
	"bytes"
	"encoding/json"
	"fmt"
	"text/template"

	"github.com/sourcegraph/go-ses"
)

var templates = map[string]*Template{}

type Template struct {
	From    string
	Name    string
	Body    string
	Subject string
	tpl     *template.Template
}

func RegisterTemplate(tpl *Template) {
	tpl.tpl = template.Must(template.New(tpl.Name).Parse(tpl.Body))
	templates[tpl.Name] = tpl
}

type Mail struct {
	payload map[string]interface{}
}

func New() *Mail {
	return &Mail{
		payload: map[string]interface{}{},
	}
}

func (m *Mail) To(email string) *Mail {
	m.payload["to"] = email
	return m
}

func (m *Mail) Tpl(name string, data map[string]interface{}) *Mail {
	m.payload["tpl"] = name
	if data == nil {
		return m
	}
	for k, v := range data {
		m.payload[k] = v
	}
	return m
}

// Payload serialize the payload in JSON,
// can be sent to a message queue for processing
func (m *Mail) Payload() []byte {
	js, err := json.Marshal(m.payload)
	if err != nil {
		panic(err)
	}
	return js
}

// Send actually send the mail
func Send(m *Mail) error {
	//payload := map[string]interface{}{}
	//if err := json.Unmarshal(m.Body, &payload); err != nil {
	//	return err
	//}
	payload := m.payload
	ntpl, ok := payload["tpl"].(string)
	if !ok {
		return fmt.Errorf("missing tpl name")
	}
	tpl, ok := templates[ntpl]
	if !ok {
		return fmt.Errorf("unknown tpl: %v", ntpl)
	}
	sub := tpl.Subject
	buf := bytes.NewBufferString("")
	if err := tpl.tpl.Execute(buf, payload); err != nil {
		return err
	}
	to, ok := payload["to"].(string)
	if !ok {
		return fmt.Errorf("missing destination address to")
	}
	res, err := ses.EnvConfig.SendEmail(tpl.From, to, sub, buf.String())
	if err != nil {
		return fmt.Errorf("failed to send email: %v / %v", err, res)
	}
	return nil
}
