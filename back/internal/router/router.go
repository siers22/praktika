package router

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	chimiddleware "github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/siers22/praktika/back/internal/handler"
	"github.com/siers22/praktika/back/internal/middleware"
	"github.com/siers22/praktika/back/internal/model"
)

func Setup(
	auth *handler.AuthHandler,
	users *handler.UserHandler,
	equipment *handler.EquipmentHandler,
	categories *handler.CategoryHandler,
	departments *handler.DepartmentHandler,
	inventory *handler.InventoryHandler,
	movements *handler.MovementHandler,
	reports *handler.ReportHandler,
	auditLogs *handler.AuditHandler,
	authMw *middleware.AuthMiddleware,
	uploadDir string,
) http.Handler {
	r := chi.NewRouter()

	// Global middleware
	r.Use(chimiddleware.Recoverer)
	r.Use(middleware.Logger)
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type"},
		AllowCredentials: false,
		MaxAge:           300,
	}))

	// Static uploads
	r.Handle("/uploads/*", http.StripPrefix("/uploads/", http.FileServer(http.Dir(uploadDir))))

	r.Route("/api/v1", func(r chi.Router) {
		// Public routes
		r.Post("/auth/login", auth.Login)
		r.Post("/auth/refresh", auth.Refresh)

		// Authenticated routes
		r.Group(func(r chi.Router) {
			r.Use(authMw.Authenticate)

			r.Post("/auth/logout", auth.Logout)
			r.Put("/auth/password", auth.ChangePassword)

			// Users — admin only
			r.Route("/users", func(r chi.Router) {
				r.Use(middleware.RequireRole(model.RoleAdmin))
				r.Get("/", users.List)
				r.Post("/", users.Create)
				r.Get("/{id}", users.GetByID)
				r.Put("/{id}", users.Update)
				r.Patch("/{id}/status", users.UpdateStatus)
			})

			// Equipment
			r.Route("/equipment", func(r chi.Router) {
				r.Get("/", equipment.List)
				r.Get("/export/csv", equipment.ExportCSV)
				r.With(middleware.RequireRole(model.RoleAdmin, model.RoleInventory)).Post("/", equipment.Create)
				r.Get("/{id}", equipment.GetByID)
				r.With(middleware.RequireRole(model.RoleAdmin, model.RoleInventory)).Put("/{id}", equipment.Update)
				r.With(middleware.RequireRole(model.RoleAdmin)).Delete("/{id}", equipment.Archive)
				r.With(middleware.RequireRole(model.RoleAdmin, model.RoleInventory)).Post("/{id}/photos", equipment.UploadPhoto)
				r.With(middleware.RequireRole(model.RoleAdmin, model.RoleInventory)).Delete("/{id}/photos/{photoId}", equipment.DeletePhoto)
				r.Get("/{id}/movements", movements.ListByEquipment)
			})

			// Categories
			r.Route("/categories", func(r chi.Router) {
				r.Get("/", categories.List)
				r.With(middleware.RequireRole(model.RoleAdmin)).Post("/", categories.Create)
				r.With(middleware.RequireRole(model.RoleAdmin)).Put("/{id}", categories.Update)
				r.With(middleware.RequireRole(model.RoleAdmin)).Delete("/{id}", categories.Delete)
			})

			// Departments
			r.Route("/departments", func(r chi.Router) {
				r.Get("/", departments.List)
				r.With(middleware.RequireRole(model.RoleAdmin)).Post("/", departments.Create)
				r.With(middleware.RequireRole(model.RoleAdmin)).Put("/{id}", departments.Update)
				r.With(middleware.RequireRole(model.RoleAdmin)).Delete("/{id}", departments.Delete)
			})

			// Inventory
			r.Route("/inventories", func(r chi.Router) {
				r.Use(middleware.RequireRole(model.RoleAdmin, model.RoleInventory))
				r.Get("/", inventory.ListSessions)
				r.Post("/", inventory.CreateSession)
				r.Get("/{id}", inventory.GetSession)
				r.Post("/{id}/items", inventory.CheckItem)
				r.Put("/{id}/items/{itemId}", inventory.UpdateItem)
				r.Post("/{id}/complete", inventory.CompleteSession)
				r.Get("/{id}/export/csv", inventory.ExportCSV)
			})

			// Movements
			r.Route("/movements", func(r chi.Router) {
				r.With(middleware.RequireRole(model.RoleAdmin, model.RoleInventory)).Post("/", movements.Create)
				r.Get("/", movements.ListAll)
			})

			// Reports
			r.Route("/reports", func(r chi.Router) {
				r.Get("/summary", reports.Summary)
				r.Get("/by-department", reports.ByDepartment)
				r.Get("/dashboard", reports.Dashboard)
			})

			// Audit logs — admin only
			r.Route("/audit-logs", func(r chi.Router) {
				r.Use(middleware.RequireRole(model.RoleAdmin))
				r.Get("/", auditLogs.List)
			})
		})
	})

	return r
}
