package templates

func ResetPasswordTemplate(
	name string,
	link string,
) string {

	return `<!doctype html>
<html>
  <body style="margin:0;padding:24px;background:#0b1220;font-family:system-ui,-apple-system,Segoe UI,Roboto,Arial,sans-serif;color:#e8eefc;">
    <div style="max-width:560px;margin:0 auto;background:#121a2b;border:1px solid rgba(255,255,255,.08);border-radius:16px;padding:24px;box-shadow:0 10px 30px rgba(0,0,0,.28);">
      <div style="font-size:13px;letter-spacing:.08em;text-transform:uppercase;color:#7aa2ff;margin-bottom:10px;">Password Reset</div>
      <h2 style="margin:0 0 12px;font-size:22px;color:#ffffff;">Hello ` + name + `</h2>
      <p style="margin:0 0 14px;color:#b8c6e6;line-height:1.55;">We received a request to reset your password. Click below to choose a new one.</p>
      <a href="` + link + `" style="display:inline-block;background:#4f7cff;color:#fff;text-decoration:none;padding:12px 16px;border-radius:12px;font-weight:700;">
        Reset Password
      </a>
      <p style="margin:16px 0 0;color:#9fb0d6;font-size:12px;line-height:1.5;">This link expires in 15 minutes. If you did not request this, you can safely ignore this email.</p>
    </div>
  </body>
</html>`
}
