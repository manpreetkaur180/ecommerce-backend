package templates

import "fmt"

// VerifyEmailTemplate renders email verification message
func VerifyEmailTemplate(name, link string) string {
	return fmt.Sprintf(`
Hello %s,

Please verify your email by clicking the link below:

%s

This link will expire in 24 hours.

If you did not create this account, you can safely ignore this email.

Thanks,  
Team
`, name, link)
}