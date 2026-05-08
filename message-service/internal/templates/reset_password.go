package templates

func ResetPasswordTemplate(
	name string,
	link string,
) string {

	return `
	<h2>Hello ` + name + `</h2>

	<p>Your email has been verified successfully.</p>

	<p>Click below to reset your password:</p>

	<a href="` + link + `">
		Reset Password
	</a>

	<p>This link will expire soon.</p>
	`
}