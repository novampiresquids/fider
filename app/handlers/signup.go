package handlers

import (
	"net/http"
	"time"

	"github.com/getfider/fider/app/metrics"
	"github.com/getfider/fider/app/models/entity"
	"github.com/getfider/fider/app/models/enum"
	"github.com/getfider/fider/app/models/query"

	"github.com/getfider/fider/app/models/cmd"
	"github.com/getfider/fider/app/pkg/bus"

	webutil "github.com/getfider/fider/app/pkg/web/util"

	"github.com/getfider/fider/app/tasks"

	"strings"

	"github.com/getfider/fider/app"
	"github.com/getfider/fider/app/actions"
	"github.com/getfider/fider/app/pkg/env"
	"github.com/getfider/fider/app/pkg/errors"
	"github.com/getfider/fider/app/pkg/validate"
	"github.com/getfider/fider/app/pkg/web"
)

// CheckAvailability checks if given domain is available to be used
func CheckAvailability() web.HandlerFunc {
	return func(c *web.Context) error {
		subdomain := strings.ToLower(c.Param("subdomain"))

		messages, _ := validate.Subdomain(c, subdomain)
		if len(messages) > 0 {
			return c.Ok(web.Map{
				"message": strings.Join(messages, ","),
			})
		}

		return c.Ok(web.Map{})
	}
}

// CreateTenant creates a new tenant
func CreateTenant() web.HandlerFunc {
	return func(c *web.Context) error {
		if env.Config.SignUpDisabled {
			return c.NotFound()
		}

		action := actions.NewCreateTenant()
		action.Email = c.User().Email
		action.Name = c.User().Name
		action.LegalAgreement = true
		if result := c.BindTo(action); !result.Ok {
			return c.HandleValidation(result)
		}

		socialSignUp := action.Token != ""

		/*
			status := enum.TenantPending
			if socialSignUp {
				status = enum.TenantActive
			}
		*/
		status := enum.TenantActive

		createTenant := &cmd.CreateTenant{
			Name:           action.TenantName,
			Subdomain:      action.Subdomain,
			Status:         status,
			WelcomeMessage: action.WelcomeMessage,
		}
		err := bus.Dispatch(c, createTenant)
		if err != nil {
			return c.Failure(err)
		}

		metrics.TotalTenants.Inc()
		c.SetTenant(createTenant.Result)

		user := &entity.User{
			Membership: []*entity.UserMembership{
				&entity.UserMembership{
					Board: createTenant.Result,
					Role:  enum.RoleAdministrator,
				},
			},
		}

		// siteURL := web.TenantBaseURL(c, c.Tenant())

		if socialSignUp {
			user.Name = action.UserClaims.OAuthName
			user.Email = action.UserClaims.OAuthEmail
			user.Providers = []*entity.UserProvider{
				{
					UID:  action.UserClaims.OAuthID,
					Name: action.UserClaims.OAuthProvider,
				},
			}

			if err := bus.Dispatch(c, &cmd.RegisterUser{User: user}); err != nil {
				return c.Failure(err)
			}

			if env.IsSingleHostMode() {
				webutil.AddAuthUserCookie(c, user)
			} else {
				webutil.SetSignUpAuthCookie(c, user)
			}

		} else {
			user.Name = action.Name
			user.Email = action.Email

			err := bus.Dispatch(c, &cmd.SaveVerificationKey{
				Key:      action.VerificationKey,
				Duration: 48 * time.Hour,
				Request:  action,
			})
			if err != nil {
				return c.Failure(err)
			}

			// c.Enqueue(tasks.SendSignUpEmail(action, siteURL))
		}
		changeRole := &cmd.ChangeUserRole{
			UserID:   c.User().ID,
			Role:     enum.RoleAdministrator,
			TenantId: createTenant.Result.ID,
		}
		err = bus.Dispatch(c, changeRole)
		if err != nil {
			return c.HandleValidation(validate.Failed("Could not add admin to board"))
		}

		// if status == enum.TenantActive && user.Email != "" {
		// c.Enqueue(tasks.SendWelcomeEmail(user.Name, user.Email, siteURL))
		// }

		return c.Ok(web.Map{
			"id": createTenant.Result.ID,
		})
	}
}

// SignUp is the entry point for installation / signup
func SignUp() web.HandlerFunc {
	return func(c *web.Context) error {
		if env.Config.SignUpDisabled {
			return c.NotFound()
		}

		if env.IsSingleHostMode() {
			firstTenant := &query.GetFirstTenant{}
			err := bus.Dispatch(c, firstTenant)
			if err != nil && errors.Cause(err) != app.ErrNotFound {
				return c.Failure(err)
			}

			if firstTenant.Result != nil {
				return c.Redirect("/")
			}
		} else {
			baseURL := web.OAuthBaseURL(c)
			if !strings.HasPrefix(c.Request.URL.String(), baseURL) {
				return c.Redirect(baseURL + "/signup")
			}
		}

		return c.Page(http.StatusOK, web.Props{
			Page:        "SignUp/SignUp.page",
			Title:       "Sign up",
			Description: "Sign up for Fider and let your customers share, vote and discuss on suggestions they have to make your product even better.",
		})
	}
}

// VerifySignUpKey checks if verify key is correct, activate the tenant and sign in user
func VerifySignUpKey() web.HandlerFunc {
	return func(c *web.Context) error {
		if c.Tenant().Status != enum.TenantPending {
			return c.NotFound()
		}

		key := c.QueryParam("k")
		result, err := validateKey(enum.EmailVerificationKindSignUp, key, c)
		if result == nil {
			return err
		}

		if err = bus.Dispatch(c, &cmd.ActivateTenant{TenantID: c.Tenant().ID}); err != nil {
			return c.Failure(err)
		}

		user := &entity.User{
			Name:  result.Name,
			Email: result.Email,
			Membership: []*entity.UserMembership{
				&entity.UserMembership{
					Board: c.Tenant(),
					Role:  enum.RoleAdministrator,
				},
			},
		}

		if err = bus.Dispatch(c, &cmd.RegisterUser{User: user}); err != nil {
			return c.Failure(err)
		}

		err = bus.Dispatch(c, &cmd.SetKeyAsVerified{Key: key})
		if err != nil {
			return c.Failure(err)
		}

		webutil.AddAuthUserCookie(c, user)

		c.Enqueue(tasks.SendWelcomeEmail(user.Name, user.Email, c.BaseURL()))

		return c.Redirect(c.BaseURL())
	}
}
