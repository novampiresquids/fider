package middlewares

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/getfider/fider/app/models/cmd"
	"github.com/getfider/fider/app/models/entity"
	"github.com/getfider/fider/app/models/enum"
	"github.com/getfider/fider/app/models/query"
	"github.com/getfider/fider/app/pkg/bus"
	"github.com/getfider/fider/app/pkg/log"

	"github.com/getfider/fider/app/pkg/validate"

	"github.com/getfider/fider/app"
	"github.com/getfider/fider/app/pkg/errors"
	"github.com/getfider/fider/app/pkg/jwt"
	"github.com/getfider/fider/app/pkg/web"
	webutil "github.com/getfider/fider/app/pkg/web/util"
)

// User gets JWT Auth token from cookie and insert into context
func User() web.MiddlewareFunc {
	return func(next web.HandlerFunc) web.HandlerFunc {
		return func(c *web.Context) error {
			var (
				token string
				user  *entity.User
			)

			cookie, err := c.Request.Cookie(web.CookieAuthName)
			if err == nil {
				token = cookie.Value
			} else {
				token = webutil.GetSignUpAuthCookie(c)
				if token != "" {
					webutil.AddAuthTokenCookie(c, token)
				}
			}

			if token != "" {
				claims, err := jwt.DecodeFiderClaims(token)
				if err != nil {
					c.RemoveCookie(web.CookieAuthName)
					return next(c)
				}

				userByClaimsID := &query.GetUserByID{UserID: claims.UserID}
				err = bus.Dispatch(c, userByClaimsID)
				user = userByClaimsID.Result
				if err != nil {
					if errors.Cause(err) == app.ErrNotFound {
						c.RemoveCookie(web.CookieAuthName)
						return next(c)
					}
					return err
				}
			} else if c.Request.IsAPI() {
				authHeader := c.Request.GetHeader("Authorization")
				parts := strings.Split(authHeader, "Bearer")
				if len(parts) == 2 {
					apiKey := strings.TrimSpace(parts[1])
					getUserByAPIKey := &query.GetUserByAPIKey{APIKey: apiKey}
					err = bus.Dispatch(c, getUserByAPIKey)
					if err != nil {
						if errors.Cause(err) == app.ErrNotFound {
							return c.HandleValidation(validate.Failed("API Key is invalid"))
						}
						return err
					}
					user = getUserByAPIKey.Result

					if !user.IsCollaborator(c.Tenant()) {
						return c.HandleValidation(validate.Failed("API Key is invalid"))
					}

					if impersonateUserIDStr := c.Request.GetHeader("X-Fider-UserID"); impersonateUserIDStr != "" {
						if !user.IsAdministrator(c.Tenant()) {
							return c.HandleValidation(validate.Failed("Only Administrators are allowed to impersonate another user"))
						}
						impersonateUserID, err := strconv.Atoi(impersonateUserIDStr)
						if err != nil {
							return c.HandleValidation(validate.Failed(fmt.Sprintf("User not found for given impersonate UserID '%s'", impersonateUserIDStr)))
						}
						userByImpersonateID := &query.GetUserByID{UserID: impersonateUserID}
						err = bus.Dispatch(c, userByImpersonateID)
						user = userByImpersonateID.Result
						if err != nil {
							if errors.Cause(err) == app.ErrNotFound {
								return c.HandleValidation(validate.Failed(fmt.Sprintf("User not found for given impersonate UserID '%s'", impersonateUserIDStr)))
							}
							return err
						}
					}
				}
			}

			// if user != nil && c.Tenant() != nil && user.Tenant.ID == c.Tenant().ID {
			if user != nil {
				// blocked users are unable to sign in
				if user.Status == enum.UserBlocked {
					c.RemoveCookie(web.CookieAuthName)
					return c.Unauthorized()
				}

				// Determine role of user in the current board
				if user.BoardRole == 0 && c.Tenant() != nil {
					for _, r := range user.Membership {
						if r.Board.ID == c.Tenant().ID {
							user.BoardRole = r.Role
						}
					}
				}

				c.SetUser(user)
			}

			return next(c)
		}
	}
}

func AddUserToBoard(role enum.Role) web.MiddlewareFunc {

	return func(next web.HandlerFunc) web.HandlerFunc {
		return func(c *web.Context) error {
			user := c.User()
			if user == nil {
				return c.HandleValidation(validate.Failed("User can not be set to be a member of this board"))
			}
			log.Info(c, "Checking user role")
			log.Info(c, user.BoardRole.String())
			if user.BoardRole == 0 {
				changeRole := &cmd.ChangeUserRole{
					UserID:   user.ID,
					Role:     enum.RoleVisitor,
					TenantId: c.Tenant().ID,
				}
				err := bus.Dispatch(c, changeRole)
				if err != nil {
					return c.HandleValidation(validate.Failed("Could not add member to board"))
				}
			}
			return next(c)
		}
	}
}
