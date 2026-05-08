package templates

func ResetPasswordTemplate(
	name string,
	link string,
) string {

	return `<!doctype html>
<html>
  <body style="margin:0;padding:24px;background:#0b1220;font-family:system-ui,-apple-system,Segoe UI,Roboto,Arial,sans-serif;color:#e8eefc;">
    <div style="max-width:520px;margin:0 auto;background:#121a2b;border:1px solid rgba(255,255,255,.08);border-radius:16px;padding:22px;">
      <h2 style="margin:0 0 12px;font-size:20px;">Hello ` + name + `</h2>
      <p style="margin:0 0 10px;color:#b8c6e6;">Your email has been verified successfully.</p>
      <p style="margin:0 0 14px;color:#b8c6e6;">Click below to reset your password:</p>
      <a href="` + link + `" style="display:inline-block;background:#4f7cff;color:#fff;text-decoration:none;padding:10px 14px;border-radius:12px;font-weight:600;">
        Reset Password
      </a>
      <p style="margin:14px 0 0;color:#9fb0d6;font-size:12px;">This link will expire soon.</p>
    </div>
  </body>
</html>`
}