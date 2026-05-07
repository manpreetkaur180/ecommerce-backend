package templates

import "fmt"

func OTPSMSTemplate(otp string) string {

	return fmt.Sprintf(
		"Your OTP code is %s. Valid for 5 minutes.",
		otp,
	)
}