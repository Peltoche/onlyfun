package auth

import "github.com/Peltoche/onlyfun/internal/tools/secret"

type LoginPageTmpl struct {
	Username      string
	UsernameError string
	PasswordError string
}

func (t *LoginPageTmpl) Template() string { return "auth/page_login" }

type BootstrapPageTmpl struct {
	Username      string
	Password      secret.Text
	UsernameError string
	PasswordError string
	ConfirmError  string
}

func (t *BootstrapPageTmpl) Template() string {
	return "auth/page_bootstrap"
}
