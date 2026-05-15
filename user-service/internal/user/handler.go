package user

import (
	"errors"
	"strings"
	"user-service/pkg/utils"

	"github.com/gofiber/fiber/v2"
)

type Handler struct {
	Service *Service
}

func NewHandler(service *Service) *Handler {
	return &Handler{Service: service}
}

// -------- REGISTER --------
func (h *Handler) Register(c *fiber.Ctx) error {
	var req RegisterRequest

	// parse request
	if err := c.BodyParser(&req); err != nil {
		return utils.ErrorResponse(c, 400, "invalid request body")
	}

	// call service
	user, err := h.Service.Register(req)
	if err != nil {
		return utils.ErrorResponse(c, 400, err.Error())
	}

	// success response
	return utils.SuccessResponse(
		c,
		200,
		"Hi "+user.Name+", successfully registered",
		nil,
	)
}

func (h *Handler) AddAddress(c *fiber.Ctx) error {
	userID, ok := c.Locals("user_id").(uint)
	if !ok {
		return utils.ErrorResponse(c, 401, "unauthorized")
	}

	var req AddAddressRequest

	// parse body
	if err := c.BodyParser(&req); err != nil {
		return utils.ErrorResponse(c, 400, "invalid request body")
	}

	// normalization FIRST
	req.Phone = utils.NormalizePhone(req.Phone)
	req.AddressType = strings.ToLower(strings.TrimSpace(req.AddressType))

	// struct validation (ONLY source of truth for validation)
	if err := utils.ValidateStruct(&req); err != nil {
		return utils.ErrorResponse(c, 400, err.Error())
	}

	address, err := h.Service.AddAddress(userID, req)
	if err != nil {
		return utils.ErrorResponse(c, 400, err.Error())
	}

	return utils.SuccessResponse(c, 200, "address added successfully", address)
}
func (h *Handler) GetAddresses(c *fiber.Ctx) error {

	userID, ok := c.Locals("user_id").(uint)

	if !ok {
		return utils.ErrorResponse(c, 401, "unauthorized")
	}

	addresses, err := h.Service.GetUserAddresses(userID)

	if err != nil {
		return utils.ErrorResponse(c, 500, "failed to fetch addresses")
	}

	return utils.SuccessResponse(
		c,
		200,
		"addresses fetched successfully",
		addresses,
	)
}

// -------- LOGIN --------
func (h *Handler) Login(c *fiber.Ctx) error {
	var req LoginRequest

	// parse request
	if err := c.BodyParser(&req); err != nil {
		return utils.ErrorResponse(c, 400, "invalid request body")
	}

	// call service
	user, err := h.Service.Login(req)
	if err != nil {
		return utils.ErrorResponse(c, 401, err.Error())
	}

	// generate jwt
	token, err := utils.GenerateJWT(
		user.ID,
		user.Role,
	)

	if err != nil {
		return utils.ErrorResponse(c, 500, "failed to generate token")
	}

	// success response
	return utils.SuccessResponse(
		c,
		200,
		"Hi "+user.Name+", logged in successfully",
		fiber.Map{
			"token": token,
			"user": fiber.Map{
				"id":          user.ID,
				"name":        user.Name,
				"email":       user.Email,
				"phone":       user.Phone,
				"role":        user.Role,
				"is_verified": user.IsVerified,
			},
		},
	)
}

func (h *Handler) RefreshToken(c *fiber.Ctx) error {
	userID, ok := c.Locals("user_id").(uint)
	if !ok {
		return utils.ErrorResponse(c, 401, "unauthorized")
	}

	user, err := h.Service.FindByID(userID)
	if err != nil {
		return utils.ErrorResponse(c, 401, err.Error())
	}

	token, err := utils.GenerateJWT(
		user.ID,
		user.Role,
	)

	if err != nil {
		return utils.ErrorResponse(c, 500, "failed to generate token")
	}

	return utils.SuccessResponse(
		c,
		200,
		"token refreshed successfully",
		fiber.Map{
			"token": token,
			"user": fiber.Map{
				"id":          user.ID,
				"name":        user.Name,
				"email":       user.Email,
				"phone":       user.Phone,
				"role":        user.Role,
				"is_verified": user.IsVerified,
			},
		},
	)
}

func (h *Handler) SendOTP(c *fiber.Ctx) error {

	var req LoginRequest

	if err := c.BodyParser(&req); err != nil {
		return utils.ErrorResponse(c, 400, "invalid request body")
	}

	req.Email = utils.NormalizeEmail(req.Email)
	req.Phone = utils.NormalizePhone(req.Phone)

	if req.Email == "" && req.Phone == "" {
		return utils.ErrorResponse(c, 400, "email or phone is required")
	}

	if req.Email != "" {
		if err := utils.ValidateEmail(req.Email); err != nil {
			return utils.ErrorResponse(c, 400, err.Error())
		}
	}

	if req.Phone != "" {
		if err := utils.ValidatePhone(req.Phone); err != nil {
			return utils.ErrorResponse(c, 400, err.Error())
		}
	}

	user, err := h.Service.FindByIdentifier(
		req.Email,
		req.Phone,
	)

	if err != nil {
		return utils.ErrorResponse(c, 400, "user not found")
	}

	if !user.IsVerified {
		return utils.ErrorResponse(
			c,
			403,
			"please verify your email before requesting otp",
		)
	}

	channel := "email"

	identifier := req.Email

	if req.Phone != "" {
		channel = "phone"
		identifier = req.Phone
	}

	err = h.Service.SendOTP(
		user,
		identifier,
		channel,
	)

	if err != nil {
		return utils.ErrorResponse(c, 500, "failed to send otp")
	}

	return utils.SuccessResponse(
		c,
		200,
		"OTP sent successfully",
		nil,
	)
}
func (h *Handler) LoginWithOTP(c *fiber.Ctx) error {
	var req OTPLoginRequest

	if err := c.BodyParser(&req); err != nil {
		return utils.ErrorResponse(c, 400, "invalid request body")
	}

	req.Email = utils.NormalizeEmail(req.Email)
	req.Phone = utils.NormalizePhone(req.Phone)

	if req.Email == "" && req.Phone == "" {
		return utils.ErrorResponse(c, 400, "email or phone is required")
	}

	if req.OTP == "" {
		return utils.ErrorResponse(c, 400, "otp is required")
	}

	user, err := h.Service.FindByIdentifier(req.Email, req.Phone)
	if err != nil {
		return utils.ErrorResponse(c, 400, "invalid request")
	}

	if !user.IsVerified {
		return utils.ErrorResponse(
			c,
			403,
			"please verify your email before login",
		)
	}

	identifier := req.Email

	if identifier == "" {
		identifier = req.Phone
	}

	if err := h.Service.VerifyOTP(identifier, req.OTP); err != nil {
		return utils.ErrorResponse(c, 400, err.Error())
	}

	// generate jwt
	token, err := utils.GenerateJWT(
		user.ID,
		user.Role,
	)

	if err != nil {
		return utils.ErrorResponse(c, 500, "failed to generate token")
	}

	return utils.SuccessResponse(
		c,
		200,
		"Hi "+user.Name+", logged in successfully",
		fiber.Map{
			"token": token,
			"user": fiber.Map{
				"id":          user.ID,
				"name":        user.Name,
				"email":       user.Email,
				"phone":       user.Phone,
				"role":        user.Role,
				"is_verified": user.IsVerified,
			},
		},
	)
}

func (h *Handler) VerifyEmail(c *fiber.Ctx) error {

	token := c.Query("token")

	if token == "" {
		return c.Status(400).Type("html", "utf-8").SendString(statusPageHTML(
			"Invalid Link",
			"Email verification link is missing a token.",
			false,
		))
	}

	err := h.Service.VerifyEmail(token)
	if err != nil {
		if errors.Is(err, ErrEmailAlreadyVerified) {
			return c.Type("html", "utf-8").SendString(statusPageHTML(
				"Already Verified",
				"Your email is already verified. You can log in.",
				true,
			))
		}

		return c.Status(400).Type("html", "utf-8").SendString(statusPageHTML(
			"Verification Failed",
			err.Error(),
			false,
		))
	}

	return c.Type("html", "utf-8").SendString(statusPageHTML(
		"Email Verified",
		"Your email has been verified successfully. You can now log in.",
		true,
	))
}

func (h *Handler) ForgotPassword(c *fiber.Ctx) error {
	var req ForgotPasswordRequest

	if err := c.BodyParser(&req); err != nil {
		return utils.ErrorResponse(c, 400, "invalid request body")
	}

	req.Email = utils.NormalizeEmail(req.Email)
	if req.Email == "" {
		return utils.ErrorResponse(c, 400, "email is required")
	}
	if err := utils.ValidateEmail(req.Email); err != nil {
		return utils.ErrorResponse(c, 400, err.Error())
	}

	if err := h.Service.ForgotPassword(req.Email); err != nil {
		return utils.ErrorResponse(c, 400, err.Error())
	}

	return utils.SuccessResponse(
		c,
		200,
		"password reset link sent to email",
		nil,
	)
}

func (h *Handler) ResetPasswordForm(c *fiber.Ctx) error {
	token := c.Query("token")
	if token == "" {
		return c.Status(400).SendString("invalid reset link")
	}

	c.Type("html", "utf-8")
	return c.SendString(`<!doctype html>
<html lang="en">
<head>
  <meta charset="utf-8" />
  <meta name="viewport" content="width=device-width, initial-scale=1" />
  <title>Reset Password</title>
  <style>
    body{font-family:system-ui,-apple-system,Segoe UI,Roboto,Arial,sans-serif;background:#0b1220;color:#e8eefc;margin:0;display:flex;min-height:100vh;align-items:center;justify-content:center;padding:24px}
    .card{width:100%;max-width:420px;background:#121a2b;border:1px solid rgba(255,255,255,.08);border-radius:16px;padding:22px;box-shadow:0 10px 30px rgba(0,0,0,.35)}
    h1{font-size:20px;margin:0 0 14px}
    label{display:block;font-size:13px;margin:12px 0 6px;color:#b8c6e6}
    input{width:100%;padding:12px 12px;border-radius:12px;border:1px solid rgba(255,255,255,.12);background:#0b1220;color:#e8eefc;outline:none}
    button{margin-top:16px;width:100%;padding:12px;border:0;border-radius:12px;background:#4f7cff;color:white;font-weight:600;cursor:pointer}
    .hint{margin-top:10px;font-size:12px;color:#9fb0d6}
  </style>
</head>
<body>
  <div class="card">
    <h1>Reset your password</h1>
    <form method="POST" action="/api/v1/user/reset-password">
      <input type="hidden" name="token" value="` + token + `" />
      <label>New password</label>
      <input type="password" name="new_password" placeholder="Enter new password" required />
      <label>Confirm new password</label>
      <input type="password" name="confirm_password" placeholder="Confirm new password" required />
      <button type="submit">Reset password</button>
      <div class="hint">Password must be at least 8 characters and include uppercase, lowercase, and a number.</div>
    </form>
  </div>
</body>
</html>`)
}

func (h *Handler) ResetPassword(c *fiber.Ctx) error {
	// support both JSON and form submissions
	type resetReq struct {
		Token           string `json:"token" form:"token"`
		NewPassword     string `json:"new_password" form:"new_password"`
		ConfirmPassword string `json:"confirm_password" form:"confirm_password"`
	}

	var req resetReq
	if err := c.BodyParser(&req); err != nil {
		return utils.ErrorResponse(c, 400, "invalid request body")
	}

	if err := h.Service.ResetPassword(req.Token, req.NewPassword, req.ConfirmPassword); err != nil {
		// if this came from the HTML form, show a friendly message
		if !strings.Contains(strings.ToLower(c.Get("Content-Type")), "application/json") {
			return c.Status(400).Type("html", "utf-8").SendString(statusPageHTML(
				"Reset Failed",
				err.Error(),
				false,
			))
		}
		return utils.ErrorResponse(c, 400, err.Error())
	}

	// HTML form success
	if !strings.Contains(strings.ToLower(c.Get("Content-Type")), "application/json") {
		return c.Type("html", "utf-8").SendString(statusPageHTML(
			"Password Reset",
			"Your password has been reset successfully. You can now log in with your new password.",
			true,
		))
	}

	return utils.SuccessResponse(
		c,
		200,
		"password reset successfully",
		nil,
	)
}

func statusPageHTML(title string, message string, success bool) string {
	accent := "#4f7cff"
	if !success {
		accent = "#ff6b6b"
	}

	return `<!doctype html>
<html lang="en">
<head>
  <meta charset="utf-8" />
  <meta name="viewport" content="width=device-width, initial-scale=1" />
  <title>` + title + `</title>
  <style>
    body{font-family:system-ui,-apple-system,Segoe UI,Roboto,Arial,sans-serif;background:#0b1220;color:#e8eefc;margin:0;display:flex;min-height:100vh;align-items:center;justify-content:center;padding:24px}
    .card{width:100%;max-width:440px;background:#121a2b;border:1px solid rgba(255,255,255,.08);border-radius:16px;padding:24px;box-shadow:0 10px 30px rgba(0,0,0,.35)}
    .mark{width:42px;height:42px;border-radius:14px;background:` + accent + `;margin-bottom:14px}
    h1{font-size:22px;margin:0 0 10px}
    p{margin:0;color:#b8c6e6;line-height:1.55}
  </style>
</head>
<body>
  <div class="card">
    <div class="mark"></div>
    <h1>` + title + `</h1>
    <p>` + message + `</p>
  </div>
</body>
</html>`
}

func (h *Handler) RequestUpdatePassword(c *fiber.Ctx) error {
	var req UpdatePasswordRequest

	if err := c.BodyParser(&req); err != nil {
		return utils.ErrorResponse(c, 400, "invalid request body")
	}

	req.Email = utils.NormalizeEmail(req.Email)
	if req.Email == "" {
		return utils.ErrorResponse(c, 400, "email is required")
	}
	if err := utils.ValidateEmail(req.Email); err != nil {
		return utils.ErrorResponse(c, 400, err.Error())
	}

	if err := h.Service.RequestUpdatePassword(req.Email); err != nil {
		return utils.ErrorResponse(c, 400, err.Error())
	}

	return utils.SuccessResponse(
		c,
		200,
		"update password link sent to email",
		nil,
	)
}

func (h *Handler) UpdatePasswordForm(c *fiber.Ctx) error {
	token := c.Query("token")
	if token == "" {
		return c.Status(400).Type("html", "utf-8").SendString(statusPageHTML(
			"Invalid Link",
			"Password update link is missing a token.",
			false,
		))
	}

	c.Type("html", "utf-8")
	return c.SendString(`<!doctype html>
<html lang="en">
<head>
  <meta charset="utf-8" />
  <meta name="viewport" content="width=device-width, initial-scale=1" />
  <title>Update Password</title>
  <style>
    body{font-family:system-ui,-apple-system,Segoe UI,Roboto,Arial,sans-serif;background:#0b1220;color:#e8eefc;margin:0;display:flex;min-height:100vh;align-items:center;justify-content:center;padding:24px}
    .card{width:100%;max-width:420px;background:#121a2b;border:1px solid rgba(255,255,255,.08);border-radius:16px;padding:22px;box-shadow:0 10px 30px rgba(0,0,0,.35)}
    h1{font-size:20px;margin:0 0 14px}
    label{display:block;font-size:13px;margin:12px 0 6px;color:#b8c6e6}
    input{width:100%;padding:12px 12px;border-radius:12px;border:1px solid rgba(255,255,255,.12);background:#0b1220;color:#e8eefc;outline:none}
    button{margin-top:16px;width:100%;padding:12px;border:0;border-radius:12px;background:#4f7cff;color:white;font-weight:600;cursor:pointer}
    .hint{margin-top:10px;font-size:12px;color:#9fb0d6}
  </style>
</head>
<body>
  <div class="card">
    <h1>Update your password</h1>
    <form method="POST" action="/api/v1/user/update-password/confirm">
      <input type="hidden" name="token" value="` + token + `" />
      <label>Old password</label>
      <input type="password" name="old_password" placeholder="Enter old password" required />
      <label>New password</label>
      <input type="password" name="new_password" placeholder="Enter new password" required />
      <label>Confirm new password</label>
      <input type="password" name="confirm_new_password" placeholder="Confirm new password" required />
      <button type="submit">Update password</button>
      <div class="hint">Password must be at least 8 characters and include uppercase, lowercase, and a number.</div>
    </form>
  </div>
</body>
</html>`)
}

func (h *Handler) ConfirmUpdatePassword(c *fiber.Ctx) error {
	type confirmReq struct {
		Token              string `json:"token" form:"token"`
		OldPassword        string `json:"old_password" form:"old_password"`
		NewPassword        string `json:"new_password" form:"new_password"`
		ConfirmNewPassword string `json:"confirm_new_password" form:"confirm_new_password"`
	}

	var req confirmReq
	if err := c.BodyParser(&req); err != nil {
		return utils.ErrorResponse(c, 400, "invalid request body")
	}

	if err := h.Service.ConfirmUpdatePassword(
		req.Token,
		req.OldPassword,
		req.NewPassword,
		req.ConfirmNewPassword,
	); err != nil {
		if !strings.Contains(strings.ToLower(c.Get("Content-Type")), "application/json") {
			return c.Status(400).Type("html", "utf-8").SendString(statusPageHTML(
				"Update Failed",
				err.Error(),
				false,
			))
		}
		return utils.ErrorResponse(c, 400, err.Error())
	}

	if !strings.Contains(strings.ToLower(c.Get("Content-Type")), "application/json") {
		return c.Type("html", "utf-8").SendString(statusPageHTML(
			"Password Updated",
			"Your password has been updated successfully. Use the new password the next time you log in.",
			true,
		))
	}

	return utils.SuccessResponse(
		c,
		200,
		"password updated successfully",
		nil,
	)
}
