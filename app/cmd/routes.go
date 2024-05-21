package cmd

import (
	"time"

	"github.com/getfider/fider/app/handlers"
	"github.com/getfider/fider/app/handlers/apiv1"
	"github.com/getfider/fider/app/handlers/webhooks"
	"github.com/getfider/fider/app/middlewares"
	"github.com/getfider/fider/app/models/enum"
	"github.com/getfider/fider/app/pkg/env"
	"github.com/getfider/fider/app/pkg/web"
)

func routes(r *web.Engine) *web.Engine {
	r.Worker().Use(middlewares.WorkerSetup())

	r.Get("/_health", handlers.Health())

	r.Use(middlewares.CatchPanic())
	r.Use(middlewares.Instrumentation())

	r.NotFound(func(c *web.Context) error {
		mw := middlewares.Chain(
			middlewares.WebSetup(),
			middlewares.Tenant(),
			middlewares.User(),
		)
		next := mw(func(c *web.Context) error {
			return c.NotFound()
		})
		return next(c)
	})

	r.Use(middlewares.Secure())
	r.Use(middlewares.Compress())

	assets := r.Group()
	{
		assets.Use(middlewares.CORS())
		assets.Use(middlewares.ClientCache(365 * 24 * time.Hour))
		assets.Get("/static/favicon", handlers.Favicon())
		assets.Static("/assets/*filepath", "dist")
	}

	r.Use(middlewares.Session())

	r.Get("/robots.txt", handlers.RobotsTXT())
	r.Post("/_api/log-error", handlers.LogError())

	r.Use(middlewares.Maintenance())
	r.Use(middlewares.WebSetup())
	r.Use(middlewares.Tenant())
	r.Use(middlewares.User())

	r.Get("/", handlers.Index())
	r.Get("/privacy", handlers.LegalPage("Privacy Policy", "privacy.md"))
	r.Get("/terms", handlers.LegalPage("Terms of Service", "terms.md"))

	r.Post("/_api/tenants", handlers.CreateTenant())
	r.Get("/_api/tenants/:subdomain/availability", handlers.CheckAvailability())
	r.Get("/signup", handlers.SignUp())
	r.Get("/oauth/:provider", handlers.SignInByOAuth())
	r.Get("/oauth/:provider/callback", handlers.OAuthCallback())

	if env.IsBillingEnabled() {
		wh := r.Group()
		{
			wh.Post("/webhooks/paddle", webhooks.IncomingPaddleWebhook())
		}
	}

	r.Get("/sitemap.xml", handlers.Sitemap())

	tenantAssets := r.Group()
	{
		tenantAssets.Use(middlewares.ClientCache(5 * 24 * time.Hour))
		tenantAssets.Get("/static/avatars/letter/:id/:name", handlers.LetterAvatar())
		tenantAssets.Get("/static/avatars/gravatar/:id/:name", handlers.Gravatar())

		tenantAssets.Use(middlewares.ClientCache(30 * 24 * time.Hour))
		tenantAssets.Get("/static/favicon/*bkey", handlers.Favicon())
		tenantAssets.Get("/static/images/*bkey", handlers.ViewUploadedImage())
		/*
			tenantAssets.Get("/static/custom/:md5.css", func(c *web.Context) error {
				return c.Blob(http.StatusOK, "text/css", []byte(c.Tenant().CustomCSS))
			})
		*/
	}

	r.Get("/_design", handlers.Page("Design System", "A preview of Fider UI elements", "DesignSystem/DesignSystem.page"))
	r.Get("/signup/verify", handlers.VerifySignUpKey())
	r.Get("/signout", handlers.SignOut())
	r.Get("/oauth/:provider/token", handlers.OAuthToken())
	r.Get("/oauth/:provider/echo", handlers.OAuthEcho())

	r.Get("/signin", handlers.SignInPage())
	r.Get("/not-invited", handlers.NotInvitedPage())
	r.Get("/signin/verify", handlers.VerifySignInKey(enum.EmailVerificationKindSignIn))
	r.Get("/invite/verify", handlers.VerifySignInKey(enum.EmailVerificationKindUserInvitation))
	r.Post("/_api/signin/complete", handlers.CompleteSignInProfile())
	r.Post("/_api/signin", handlers.SignInByEmail())

	ui := r.Group()
	{
		//From this step, a User is required
		ui.Use(middlewares.IsAuthenticated())

		ui.Get("/settings", handlers.UserSettings())
		ui.Get("/notifications", handlers.Notifications())
		ui.Get("/notifications/:id", handlers.ReadNotification())
		ui.Get("/change-email/verify", handlers.VerifyChangeEmailKey())

		ui.Delete("/_api/user", handlers.DeleteUser())
		ui.Post("/_api/user/regenerate-apikey", handlers.RegenerateAPIKey())
		ui.Post("/_api/user/settings", handlers.UpdateUserSettings())
		ui.Post("/_api/user/change-email", handlers.ChangeUserEmail())
		ui.Post("/_api/notifications/read-all", handlers.ReadAllNotifications())
		ui.Get("/_api/notifications/unread/total", handlers.TotalUnreadNotifications())

		// From this step, only Collaborators and Administrators are allowed
		ui.Use(middlewares.IsAuthorized(enum.RoleCollaborator, enum.RoleAdministrator))

		// locale is forced to English for administrative pages.
		// This is meant to be removed when all pages are translated.
		ui.Use(middlewares.SetLocale("en"))

		ui.Post("/api/v1/board", handlers.CreateTenant())

		ui.Get("/board/:board/admin", handlers.GeneralSettingsPage())
		ui.Get("/board/:board/admin/advanced", handlers.AdvancedSettingsPage())
		ui.Get("/board/:board/admin/privacy", handlers.Page("Privacy · Site Settings", "", "Administration/pages/PrivacySettings.page"))
		ui.Get("/board/:board/admin/invitations", handlers.Page("Invitations · Site Settings", "", "Administration/pages/Invitations.page"))
		ui.Get("/board/:board/admin/members", handlers.ManageMembers())
		ui.Get("/board/:board/admin/tags", handlers.ManageTags())
		// ui.Get("/board/:board/admin/authentication", handlers.ManageAuthentication())
		ui.Get("/_api/admin/oauth/:provider", handlers.GetOAuthConfig())

		//From this step, only Administrators are allowed
		ui.Use(middlewares.IsAuthorized(enum.RoleAdministrator))

		ui.Get("/board/:board/admin/export", handlers.Page("Export · Site Settings", "", "Administration/pages/Export.page"))
		ui.Get("/board/:board/admin/export/posts.csv", handlers.ExportPostsToCSV())
		ui.Get("/board/:board/admin/export/backup.zip", handlers.ExportBackupZip())
		ui.Get("/board/:board/admin/webhooks", handlers.ManageWebhooks())
		ui.Post("/_api/board/:board/admin/webhook", handlers.CreateWebhook())
		ui.Put("/_api/board/:board/admin/webhook/:id", handlers.UpdateWebhook())
		ui.Delete("/_api/board/:board/admin/webhook/:id", handlers.DeleteWebhook())
		ui.Get("/_api/board/:board/admin/webhook/test/:id", handlers.TestWebhook())
		ui.Post("/_api/board/:board/admin/webhook/preview", handlers.PreviewWebhook())
		ui.Get("/_api/board/:board/admin/webhook/props/:type", handlers.GetWebhookProps())
		ui.Post("/_api/board/:board/admin/settings/general", handlers.UpdateSettings())
		ui.Post("/_api/board/:board/admin/settings/advanced", handlers.UpdateAdvancedSettings())
		ui.Post("/_api/board/:board/admin/settings/privacy", handlers.UpdatePrivacy())
		ui.Post("/_api/board/:board/admin/settings/emailauth", handlers.UpdateEmailAuthAllowed())
		ui.Post("/_api/board/:board/admin/oauth", handlers.SaveOAuthConfig())
		ui.Post("/_api/board/:board/admin/roles/:role/users", handlers.ChangeUserRole())
		ui.Put("/_api/board/:board/admin/users/:userID/block", handlers.BlockUser())
		ui.Delete("/_api/board/:board/admin/users/:userID/block", handlers.UnblockUser())

		if env.IsBillingEnabled() {
			ui.Get("/board/:board/admin/billing", handlers.ManageBilling())
			ui.Post("/_api/billing/checkout-link", handlers.GenerateCheckoutLink())
		}
	}

	//Starting from this step, a Tenant is required
	r.Use(middlewares.RequireTenant())
	//If tenant is pending, block it from using any other route
	r.Use(middlewares.BlockPendingTenants())

	//Block if it's private tenant with unauthenticated user
	r.Use(middlewares.CheckTenantPrivacy())

	r.Get("/board/:board", handlers.BoardDetails())
	r.Get("/board/:board/posts/:number", handlers.PostDetails())
	r.Get("/board/:board/posts/:number/:slug", handlers.PostDetails())

	// Public operations
	// Does not require authentication
	publicApi := r.Group()
	{
		publicApi.Get("/api/v1/board/:board/posts", apiv1.SearchPosts())
		publicApi.Get("/api/v1/board/:board/tags", apiv1.ListTags())
		publicApi.Get("/api/v1/board/:board/posts/:number", apiv1.GetPost())
		publicApi.Get("/api/v1/board/:board/posts/:number/comments", apiv1.ListComments())
		publicApi.Get("/api/v1/board/:board/posts/:number/comments/:id", apiv1.GetComment())
	}

	// Operations used to manage the content of a site
	// Available to any authenticated user
	membersApi := r.Group()
	{
		membersApi.Use(middlewares.IsAuthenticated())
		membersApi.Use(middlewares.AddUserToBoard(enum.RoleVisitor))
		membersApi.Use(middlewares.BlockLockedTenants())

		membersApi.Post("/api/v1/board/:board/posts", apiv1.CreatePost())
		membersApi.Put("/api/v1/board/:board/posts/:number", apiv1.UpdatePost())
		membersApi.Post("/api/v1/board/:board/posts/:number/comments", apiv1.PostComment())
		membersApi.Put("/api/v1/board/:board/posts/:number/comments/:id", apiv1.UpdateComment())
		membersApi.Delete("/api/v1/board/:board/posts/:number/comments/:id", apiv1.DeleteComment())
		membersApi.Post("/api/v1/board/:board/posts/:number/votes", apiv1.AddVote())
		membersApi.Delete("/api/v1/board/:board/posts/:number/votes", apiv1.RemoveVote())
		membersApi.Post("/api/v1/board/:board/posts/:number/subscription", apiv1.Subscribe())
		membersApi.Delete("/api/v1/board/:board/posts/:number/subscription", apiv1.Unsubscribe())

		membersApi.Use(middlewares.IsAuthorized(enum.RoleCollaborator, enum.RoleAdministrator))
		membersApi.Put("/api/v1/board/:board/posts/:number/status", apiv1.SetResponse())

		// membersApi.Get("/api/v1/boards", apiv1.ListBoards())
	}

	// Operations used to manage a site
	// Available to both collaborators and administrators
	staffApi := r.Group()
	{
		staffApi.Use(middlewares.SetLocale("en"))
		staffApi.Use(middlewares.IsAuthenticated())
		staffApi.Use(middlewares.IsAuthorized(enum.RoleCollaborator, enum.RoleAdministrator))

		staffApi.Get("/api/v1/users", apiv1.ListUsers())
		staffApi.Get("/api/v1/board/:board/posts/:number/votes", apiv1.ListVotes())
		staffApi.Post("/api/v1/invitations/send", apiv1.SendInvites())
		staffApi.Post("/api/v1/invitations/sample", apiv1.SendSampleInvite())

		staffApi.Use(middlewares.BlockLockedTenants())
		staffApi.Post("/api/v1/board/:board/posts/:number/tags/:slug", apiv1.AssignTag())
		staffApi.Delete("/api/v1/board/:board/posts/:number/tags/:slug", apiv1.UnassignTag())
	}

	// Operations used to manage a site
	// Only available to administrators
	adminApi := r.Group()
	{
		adminApi.Use(middlewares.SetLocale("en"))
		adminApi.Use(middlewares.IsAuthenticated())
		adminApi.Use(middlewares.IsAuthorized(enum.RoleAdministrator))

		adminApi.Post("/api/v1/users", apiv1.CreateUser())
		adminApi.Post("/api/v1/board/:board/tags", apiv1.CreateEditTag())
		adminApi.Put("/api/v1/board/:board/tags/:slug", apiv1.CreateEditTag())
		adminApi.Delete("/api/v1/board/:board/tags/:slug", apiv1.DeleteTag())

		// Multiple boards
		// adminApi.Post("/api/v1/boards", apiv1.CreateBoard())
		// adminApi.Put("/api/v1/boards/:number", apiv1.UpdateBoard())
		// adminApi.Delete("/api/v1/boards/:number", apiv1.DeleteBoard())

		adminApi.Use(middlewares.BlockLockedTenants())
		adminApi.Delete("/api/v1/board/:board/posts/:number", apiv1.DeletePost())
	}

	return r
}
