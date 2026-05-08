package user

import (
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

	// success response
	return utils.SuccessResponse(
		c,
		200,
		"Hi "+user.Name+", logged in successfully",
		nil,
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

	identifier := req.Email

	if identifier == "" {
		identifier = req.Phone
	}

	if err := h.Service.VerifyOTP(identifier, req.OTP); err != nil {
		return utils.ErrorResponse(c, 400, err.Error())
	}

	user, err := h.Service.FindByIdentifier(req.Email, req.Phone)
	if err != nil {
		return utils.ErrorResponse(c, 400, "invalid request")
	}

	return utils.SuccessResponse(
		c,
		200,
		"Hi "+user.Name+", logged in successfully",
		nil,
	)
}


func (h *Handler) VerifyEmail(c *fiber.Ctx) error {

	token := c.Query("token")

	if token == "" {
		return utils.ErrorResponse(c, 400, "token required")
	}

	err := h.Service.VerifyEmail(token)
	if err != nil {
		return utils.ErrorResponse(c, 400, err.Error())
	}

	return utils.SuccessResponse(
		c,
		200,
		"email verified successfully",
		nil,
	)
<<<<<<< Updated upstream
=======
}
func (h *Handler) ForgotPassword(c *fiber.Ctx) error {

	var req ForgotPasswordRequest

	if err := c.BodyParser(&req); err != nil {
		return utils.ErrorResponse(
			c,
			400,
			"invalid request body",
		)
	}

	req.Email = utils.NormalizeEmail(req.Email)

	if req.Email == "" {
		return utils.ErrorResponse(
			c,
			400,
			"email is required",
		)
	}

	err := h.Service.ForgotPassword(req.Email)

	if err != nil {
		return utils.ErrorResponse(
			c,
			400,
			err.Error(),
		)
	}

	return utils.SuccessResponse(
		c,
		200,
		"verification email sent",
		nil,
	)
}
func (h *Handler) VerifyResetRequest(
	c *fiber.Ctx,
) error {

	token := c.Query("token")

	if token == "" {
		return c.SendString(
			"invalid verification link",
		)
	}

	err := h.Service.VerifyResetRequest(token)

	if err != nil {
		return c.SendString(
			"verification failed: " + err.Error(),
		)
	}

	return c.SendString(
		"Email verified successfully.\n\nPlease check your email for the password reset link.",
	)
}
func (h *Handler) ResetPasswordPage(
	c *fiber.Ctx,
) error {

	token := c.Query("token")

	html := `
<!DOCTYPE html>
<html>
<head>
<title>Reset Password</title>
<meta name="viewport" content="width=device-width, initial-scale=1.0">
<style>
body{
	font-family:Arial;
	background:#f4f4f4;
	padding:20px;
}
.container{
	max-width:400px;
	margin:auto;
	background:white;
	padding:20px;
	border-radius:10px;
	box-shadow:0 0 10px rgba(0,0,0,0.1);
}
input{
	width:100%;
	padding:12px;
	margin-top:10px;
	border:1px solid #ccc;
	border-radius:6px;
}
button{
	width:100%;
	padding:12px;
	margin-top:20px;
	background:#000;
	color:white;
	border:none;
	border-radius:6px;
	font-size:16px;
}
</style>
</head>
<body>

<div class="container">
<h2>Reset Password</h2>

<form method="POST" action="/api/v1/user/reset-password">

<input type="hidden" name="token" value="` + token + `">

<input type="password" name="new_password" placeholder="New Password" required>

<input type="password" name="confirm_password" placeholder="Confirm Password" required>

<button type="submit">
Update Password
</button>

</form>
</div>

</body>
</html>
`

	return c.Type("html").SendString(html)
}
func (h *Handler) ResetPassword(c *fiber.Ctx) error {

	var req ResetPasswordRequest

	if err := c.BodyParser(&req); err != nil {
		return c.SendString(
			"invalid request",
		)
	}

	err := h.Service.ResetPassword(req)

	if err != nil {
		return c.SendString(
			"failed: " + err.Error(),
		)
	}

	return c.SendString(
		"Password updated successfully. You can now login.",
	)
}
func (h *Handler) UpdatePasswordRequest(
	c *fiber.Ctx,
) error {

	var req UpdatePasswordRequest

	if err := c.BodyParser(&req); err != nil {
		return utils.ErrorResponse(
			c,
			400,
			"invalid request body",
		)
	}

	req.Email = utils.NormalizeEmail(req.Email)

	if req.Email == "" {
		return utils.ErrorResponse(
			c,
			400,
			"email is required",
		)
	}

	err := h.Service.SendUpdatePasswordLink(
		req.Email,
	)

	if err != nil {
		return utils.ErrorResponse(
			c,
			400,
			err.Error(),
		)
	}

	return utils.SuccessResponse(
		c,
		200,
		"update password link sent",
		nil,
	)
}
func (h *Handler) UpdatePasswordPage(
	c *fiber.Ctx,
) error {

	token := c.Query("token")

	html := `
<!DOCTYPE html>
<html>
<head>
<title>Update Password</title>
<meta name="viewport" content="width=device-width, initial-scale=1.0">
<style>
body{
	font-family:Arial;
	background:#f4f4f4;
	padding:20px;
}
.container{
	max-width:400px;
	margin:auto;
	background:white;
	padding:20px;
	border-radius:10px;
	box-shadow:0 0 10px rgba(0,0,0,0.1);
}
input{
	width:100%;
	padding:12px;
	margin-top:10px;
	border:1px solid #ccc;
	border-radius:6px;
}
button{
	width:100%;
	padding:12px;
	margin-top:20px;
	background:#000;
	color:white;
	border:none;
	border-radius:6px;
	font-size:16px;
}
</style>
</head>
<body>

<div class="container">
<h2>Update Password</h2>

<form method="POST" action="/api/v1/user/update-password">

<input type="hidden" name="token" value="` + token + `">

<input type="password" name="old_password" placeholder="Old Password" required>

<input type="password" name="new_password" placeholder="New Password" required>

<input type="password" name="confirm_new_password" placeholder="Confirm New Password" required>

<button type="submit">
Update Password
</button>

</form>
</div>

</body>
</html>
`

	return c.Type("html").SendString(html)
}
func (h *Handler) UpdatePassword(
	c *fiber.Ctx,
) error {

	var req UpdatePasswordConfirmRequest

	if err := c.BodyParser(&req); err != nil {
		return c.SendString(
			"invalid request",
		)
	}

	err := h.Service.UpdatePassword(req)

	if err != nil {
		return c.SendString(
			"failed: " + err.Error(),
		)
	}

	return c.SendString(`
<!DOCTYPE html>
<html>
<head>
<title>Password Updated</title>
<meta name="viewport" content="width=device-width, initial-scale=1.0">
<style>
body{
	font-family:Arial;
	background:#f4f4f4;
	padding:20px;
	display:flex;
	align-items:center;
	justify-content:center;
	height:100vh;
}
.container{
	max-width:400px;
	background:white;
	padding:30px;
	border-radius:12px;
	text-align:center;
	box-shadow:0 0 15px rgba(0,0,0,0.1);
}
h2{
	color:green;
}
p{
	color:#444;
	margin-top:10px;
}
</style>
</head>
<body>

<div class="container">
<h2>Password Updated Successfully</h2>

<p>
Your password has been updated successfully.
</p>

<p>
You can now login from your mobile app.
</p>
</div>

</body>
</html>
`)
>>>>>>> Stashed changes
}