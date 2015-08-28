package mailer

import "testing"

func TestMail(t *testing.T) {
	t.Logf("Testing payload generation...")
	p := New().Tpl("welcome", nil).To("thomas.sileo@gmail.com").Payload()
	t.Logf("payload: %v", string(p))
	p2 := New().Tpl("welcome", map[string]interface{}{"count": 100}).To("thomas.sileo@gmail.com").Payload()
	t.Logf("payload2: %v", string(p2))
}
