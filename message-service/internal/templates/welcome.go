package templates

func WelcomeTemplate(name string) string {

	return `<!doctype html>
<html>
  <body style="margin:0;padding:24px;background:#0b1220;font-family:system-ui,-apple-system,Segoe UI,Roboto,Arial,sans-serif;color:#e8eefc;">
    <div style="max-width:560px;margin:0 auto;background:#121a2b;border:1px solid rgba(255,255,255,.08);border-radius:16px;padding:24px;box-shadow:0 10px 30px rgba(0,0,0,.28);">
      <div style="font-size:13px;letter-spacing:.08em;text-transform:uppercase;color:#7aa2ff;margin-bottom:10px;">Welcome</div>
      <h2 style="margin:0 0 12px;font-size:22px;color:#ffffff;">Hello ` + name + `</h2>
      <p style="margin:0 0 10px;color:#b8c6e6;line-height:1.55;">Your account has been created successfully.</p>
      <p style="margin:0;color:#b8c6e6;line-height:1.55;">You can now explore products.</p>
      <div style="height:1px;background:rgba(255,255,255,.08);margin:20px 0;"></div>
      <p style="margin:0;color:#9fb0d6;font-size:12px;">Thanks,<br/>Wits Innovation Lab</p>
    </div>
  </body>
</html>`
}
