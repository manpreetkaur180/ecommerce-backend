package templates

import "fmt"

func OTPTemplate(name string, otp string) string {

	return fmt.Sprintf(`
Hello %s,

Your OTP code is: %s

This OTP is valid for 5 minutes.

Do not share it with anyone.

Thanks,
Your App Team
`, name, otp)
}