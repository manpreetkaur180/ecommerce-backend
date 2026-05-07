package templates

import "fmt"

func WelcomeTemplate(name string) string {

	return fmt.Sprintf(`
Hello %s,

Welcome to our platform.

We are excited to have you onboard.

Thanks,
Wits Innovation Lab
`, name)
}