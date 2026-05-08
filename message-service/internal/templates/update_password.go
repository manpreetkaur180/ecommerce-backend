package templates

func UpdatePasswordTemplate(
	name string,
	link string,
) string {
	return `
	<h2>Hello ` + name + `</h2>

	<p>Click below to update your password:</p>

	<a href="` + link + `">
		Update Password
	</a>

	<p>This link will expire soon.</p>
	`
}

