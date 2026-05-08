package templates

func ForgotPasswordVerifyTemplate(
	name string,
	link string,
) string {

	return `<!doctype html>
<html>
  <body style="margin:0;padding:24px;background:#0b1220;font-family:system-ui,-apple-system,Segoe UI,Roboto,Arial,sans-serif;color:#e8eefc;">
    <div style="max-width:520px;margin:0 auto;background:#121a2b;border:1px solid rgba(255,255,255,.08);border-radius:16px;padding:22px;">
      <h2 style="margin:0 0 12px;font-size:20px;">Hello ` + name + `</h2>
      <p style="margin:0 0 10px;color:#b8c6e6;">You requested to reset your password.</p>
      <p style="margin:0 0 14px;color:#b8c6e6;">Please verify your email first by clicking below:</p>
      <a href="` + link + `" style="display:inline-block;background:#4f7cff;color:#fff;text-decoration:none;padding:10px 14px;border-radius:12px;font-weight:600;">
        Verify Email
      </a>
      <p style="margin:14px 0 0;color:#9fb0d6;font-size:12px;">After verification, you will receive another email with the password reset link.</p>
    </div>
  </body>
</html>`
}