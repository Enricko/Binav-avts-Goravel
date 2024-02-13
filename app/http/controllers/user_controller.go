package controllers

import (
	"fmt"
	extensions "goravel/app"
	"goravel/app/http/requests/user_request"
	"goravel/app/models"
	"time"

	"github.com/goravel/framework/contracts/http"
	"github.com/goravel/framework/contracts/mail"
	"github.com/goravel/framework/facades"
)

// UserController handles user-related HTTP requests.
type UserController struct {
	// No dependencies for now
}

// NewUserController creates a new instance of UserController.
func NewUserController() *UserController {
	return &UserController{}
}

// Show is a sample method to demonstrate returning a JSON response.
func (r *UserController) Show(ctx http.Context) http.Response {
	header := ctx.Request().Header("Authorization", "")
	payload, err := facades.Auth().Parse(ctx, header)
	if err != nil {
		return extensions.HandleBadRequestError(ctx, err)
	}

	if payload.Guard != "user" {
		return ctx.Response().Json(http.StatusForbidden, http.Json{
			"message": "The token is not user token.",
			"status":  http.StatusBadRequest,
		})
	}

	var existingUser models.User
	if err := facades.Orm().Query().Model(&existingUser).Where("id_user = ?", payload.Key).First(&existingUser); err != nil {
		return ctx.Response().Json(http.StatusBadRequest, http.Json{
			"message": err.Error(),
			"status":  http.StatusBadRequest,
		})
	}

	return ctx.Response().Json(http.StatusOK, http.Json{
		"message": "Successfully get data.",
		"user":    existingUser,
		"status":  http.StatusOK,
	})
}

// Login handles user login requests.
func (r *UserController) Login(ctx http.Context) http.Response {
	// Parse the incoming request
	var request user_request.LoginUserRequest

	// Validate the request
	if shouldReturn, returnValue := extensions.RequestValidation(ctx, &request); shouldReturn {
		return returnValue
	}

	// Find the user by email
	var existingUser models.User
	if err := facades.Orm().Query().Model(&existingUser).Where("email = ?", request.Email).FirstOrFail(&existingUser); err != nil {
		return ctx.Response().Json(http.StatusNotFound, http.Json{
			"message": "Data not found.",
			"status":  http.StatusNotFound,
		})
	}

	// Check if the password matches
	if !facades.Hash().Check(request.Password, existingUser.Password) {
		return ctx.Response().Json(http.StatusForbidden, http.Json{
			"message": "Email or Password incorrect.",
			"status":  http.StatusForbidden,
		})
	}

	// Generate JWT token and return response
	token, err := facades.Auth().Login(ctx, &existingUser)
	if err != nil {
		return ctx.Response().Json(http.StatusBadRequest, http.Json{
			"message": err.Error(),
			"status":  http.StatusBadRequest,
		})
	}

	return ctx.Response().Json(http.StatusAccepted, http.Json{
		"message": "Welcome back " + existingUser.Name,
		"token":   token,
		"status":  http.StatusAccepted,
	})
}
func (r *UserController) Logout(ctx http.Context) http.Response {
	header := ctx.Request().Header("Authorization", "")
	_, err := facades.Auth().Parse(ctx, header)
	if err != nil {
		return extensions.HandleBadRequestError(ctx, err)
	}

	if err := facades.Auth().Logout(ctx); err != nil {
		return extensions.HandleBadRequestError(ctx, err)
	}

	return ctx.Response().Success().Json(http.Json{
		"message": "Logout successfully.",
		"status":  http.StatusOK,
	})
}

func (r *UserController) ForgotPassword(ctx http.Context) http.Response {
	var request user_request.ForgotUserRequest

	// Validate the request
	if shouldReturn, returnValue := extensions.RequestValidation(ctx, &request); shouldReturn {
		return returnValue
	}

	var existingUser models.User
	if err := facades.Orm().Query().Model(&existingUser).Where("email = ?", request.Email).FirstOrFail(&existingUser); err != nil {
		return ctx.Response().Json(http.StatusNotFound, http.Json{
			"message": "Data not found.",
			"status":  http.StatusNotFound,
		})
	}
	fmt.Printf("existingUser: %v\n")

	otpCode, err := extensions.GenerateOTP(6)
	if err != nil {
		return extensions.HandleBadRequestError(ctx, err)
	}

	var resetModel models.ResetCodePassword
	if err := facades.Orm().Query().UpdateOrCreate(
		&resetModel,
		models.ResetCodePassword{
			Email: request.Email,
		},
		models.ResetCodePassword{
			Code: otpCode,
		}); err != nil {
		return extensions.HandleBadRequestError(ctx, err)
	}

	if err := facades.Mail().To([]string{request.Email}).
		Cc([]string{request.Email}).
		Bcc([]string{request.Email}).
		Content(mail.Content{Subject: "Subject", Html: MailOtpUI(existingUser.Name, otpCode)}).
		Send(); err != nil {
		return extensions.HandleBadRequestError(ctx, err)
	}
	return ctx.Response().Success().Json(http.Json{
		"message": "Otp has been mailed." + otpCode,
		"status":  http.StatusOK,
	})
}

func (r *UserController) CheckCode(ctx http.Context) http.Response {
	var request user_request.ForgotUserRequest

	// Validate the request
	if shouldReturn, returnValue := extensions.RequestValidation(ctx, &request); shouldReturn {
		return returnValue
	}

	var existingOtp models.ResetCodePassword
	shouldReturn, returnValue := CheckOTPValidation(existingOtp, request, ctx)
	if shouldReturn {
		return returnValue
	}

	return ctx.Response().Success().Json(http.Json{
		"message": "OTP verification successful.",
		"status":  http.StatusOK,
	})
}

func (r *UserController) ResetPassword(ctx http.Context) http.Response {
	var request user_request.ForgotUserRequest

	// Validate the request
	if shouldReturn, returnValue := extensions.RequestValidation(ctx, &request); shouldReturn {
		return returnValue
	}

	var existingOtp models.ResetCodePassword
	shouldReturn, returnValue := CheckOTPValidation(existingOtp, request, ctx)
	if shouldReturn {
		return returnValue
	}

	// Password Hashing / Dcrypt
	password, shouldReturn, returnValue := extensions.PasswordHash(request.Password, ctx)
	if shouldReturn {
		return returnValue
	}

	var existingUser models.User
	if err := facades.Orm().
		Query().
		Model(&existingUser).
		Where("email = ?", request.Email).FirstOrFail(&existingUser); err != nil {
		return extensions.HandleInternalServerError(ctx, err)
	}
	existingUser.Password = password
	existingUser.PasswordString = existingUser.Name + request.Password + request.Email
	if err := facades.Orm().
		Query().
		Save(&existingUser); err != nil {
		return extensions.HandleInternalServerError(ctx, err)
	}

	if _, err := facades.Orm().Query().Where("email = ?", request.Email).Delete(&existingOtp); err != nil {
		return extensions.HandleInternalServerError(ctx, err)
	}

	return ctx.Response().Success().Json(http.Json{
		"message": "Password reset successful.",
		"status":  http.StatusOK,
	})
}

func CheckOTPValidation(existingOtp models.ResetCodePassword, request user_request.ForgotUserRequest, ctx http.Context) (bool, http.Response) {
	if err := facades.Orm().Query().Model(&existingOtp).Where("email = ?", request.Email).FirstOrFail(&existingOtp); err != nil {
		return true, ctx.Response().Json(http.StatusNotFound, http.Json{
			"message": "Data not found.",
			"status":  http.StatusNotFound,
		})
	}

	if !extensions.VerifyOTP(existingOtp.Code, request.OtpCode) {
		return true, ctx.Response().Json(http.StatusBadRequest, http.Json{
			"message": "OTP code invalid.",
			"status":  http.StatusNotFound,
		})
	}

	if time.Since(existingOtp.CreatedAt) > 15*time.Minute {

		if _, err := facades.Orm().Query().Where("email = ?", request.Email).Delete(&existingOtp); err != nil {

			return true, extensions.HandleInternalServerError(ctx, err)
		}
		return true, ctx.Response().Json(http.StatusBadRequest, http.Json{
			"message": "OTP expired.",
			"status":  http.StatusOK,
		})
	}
	return false, nil
}

func MailOtpUI(name string, otpCode string) string {
	fmt.Printf("otpCode: %v\n", otpCode)
	return fmt.Sprintf(`<!DOCTYPE html>
	<html lang="en">
	
	<head>
		<meta charset="UTF-8">
		<meta name="viewport" content="width=device-width, initial-scale=1.0">
		<meta http-equiv="X-UA-Compatible" content="ie=edge">
		<title>Binav AVTS</title>
		<link href="https://cdn.jsdelivr.net/npm/bootstrap@5.3.2/dist/css/bootstrap.min.css" rel="stylesheet"
			integrity="sha384-T3c6CoIi6uLrA9TneNEoa7RxnatzjcDSCmG1MXxSR1GAsXEV/Dwwykc2MPK8M2HN" crossorigin="anonymous">
		<script src="https://cdn.jsdelivr.net/npm/bootstrap@5.3.2/dist/js/bootstrap.bundle.min.js"
			integrity="sha384-C6RzsynM9kWDrMNeT87bh95OGNyZPhcTNXj1NW7RuBCsyN/o0jlpcV8Qyq46cDfL" crossorigin="anonymous">
		</script>
	</head>
	
	<body>
		<div class="card m-auto" style="width: 40rem; margin:auto">
			<div style="border: 1px solid #868383; border-radius: 5px 5px 0px 0px; padding: 20px;">
				<img src="https://encrypted-tbn0.gstatic.com/images?q=tbn:ANd9GcSmlPQ2XMJ8hQT05TyZbN_hD_QXAOvjI-79c_b98A9h&s"
					class="card-img-top" alt="logo">
			</div>
			<div class="container center" style="border: 1px solid #ccc; border-radius: 0px 0px 5px 5px; padding: 20px;">
	
				<div class="card-body">
					<p>Hey <span class="w-900">%s</span><br></p>
					<p>There was a request to change your password on Binav AVTS !!!<br></p>
					<p style="color:red">If you did not make this request then please ignore this email.<br></p>
					<p>Otherwise, Here the Verification Code: <h4>%s</h4><br></p>
					<p>Thank you for choosing us. We are hopeful for your success and your growth!<br></p>
					<p>---------<br></p>
					<p>PT. Binav Maju Sejahtera<br>https://binav-avts.id</p>
				</div>
			</div>
		</div>
	</body>

	</html>
	`, name, otpCode)
}
