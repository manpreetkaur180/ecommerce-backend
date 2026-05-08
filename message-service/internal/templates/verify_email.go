package templates

import "fmt"

func VerifyEmailTemplate(
	name string,
	link string,
) string {

	return fmt.Sprintf(`
Hello %s,

Please verify your email by clicking the link below:

%s

This link will expire in 24 hours.

Thanks,
Team
`, name, link)
}