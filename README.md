# Mailer

Send mail with golang via [AWS SES](https://aws.amazon.com/ses/).

The package expects `AWS_ACCESS_KEY_ID` and `AWS_SECRET_KEY` to be set.

Templates are parsed with `text/template`.

```go
// First we register a template
welcometpl := `Welcome {{.to}}`
tpl := &mailer.Template{
	Name: "welcome",
	Subject: "Welcome",
	Body: welcometpl,
	From: "contact@acme.com",
}
mailer.RegisterTemplate(tpl)

// Create a `mailer.Mail` object
m := mailer.New().Tpl("welcome", nil).To("thomas.sileo@gmail.com")
// m.Payload() returns JSON if you want to send the payload in a message queue

// Actually send the mail
if err := mailer.Send(m); err != nil {
	panic(err)
}
```

## TODO

 - Add a basic worker (`chan *Mail`)
 - HTML template support?
 - Loads email template from a YAML file?

