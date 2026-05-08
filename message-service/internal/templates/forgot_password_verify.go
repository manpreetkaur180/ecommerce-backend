package templates

func ForgotPasswordVerifyTemplate(
	name string,
	link string,
) string {

	return `
	<h2>Hello ` + name + `</h2>

	<p>You requested to reset your password.</p>

	<p>Please verify your email first by clicking below:</p>

	<a href="` + link + `">
		Verify Email
	</a>

	<p>After verification, you will receive another email with the password reset link.</p>
	`
}