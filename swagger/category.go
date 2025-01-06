package swagger

import (
	"backend/controllers"
	"backend/models"
	"net/http"

	"github.com/zc2638/swag"
	"github.com/zc2638/swag/endpoint"
)

func SetupCategorySwagger(api *swag.API) {
	categoryController := controllers.NewCategoryController()

	// Endpoint: Create Category
	api.AddEndpoint(
		endpoint.New(
			http.MethodPost, "/categories",
			endpoint.Handler(categoryController.CreateCategory),
			endpoint.Summary("Create a new category"),
			endpoint.Description("Allows a user to create a new category"),
			endpoint.Body(models.Category{}, "Category object to create", true),
			endpoint.Response(http.StatusCreated, "Successfully created category", endpoint.SchemaResponseOption(models.Category{})),
			endpoint.Response(http.StatusBadRequest, "Invalid category data"),
			endpoint.Response(http.StatusConflict, "Category already exists"),
			endpoint.Response(http.StatusInternalServerError, "Internal server error"),
			endpoint.Security("bearer_auth"),
			endpoint.Tags("Categories"),
		),
	)

	// Endpoint: Get All Categories
	api.AddEndpoint(
		endpoint.New(
			http.MethodGet, "/categories",
			endpoint.Handler(categoryController.GetCategories),
			endpoint.Summary("Retrieve all categories"),
			endpoint.Description("Fetches a paginated list of all categories"),
			endpoint.Query("search", "string", "Optional search keyword", false),
			endpoint.Query("page", "integer", "Page number for pagination", false),
			endpoint.Query("pageSize", "integer", "Number of items per page", false),
			endpoint.Response(http.StatusInternalServerError, "Internal server error"),
			endpoint.Security("bearer_auth"),
			endpoint.Tags("Categories"),
		),
	)

	// Endpoint: Get Category By ID
	api.AddEndpoint(
		endpoint.New(
			http.MethodGet, "/categories/{id}",
			endpoint.Handler(categoryController.GetCategoryById),
			endpoint.Summary("Retrieve a category by ID"),
			endpoint.Description("Fetches the details of a specific category using its ID"),
			endpoint.Path("id", "string", "ID of the category to retrieve", true),
			endpoint.Response(http.StatusOK, "Details of the category", endpoint.SchemaResponseOption(models.Category{})),
			endpoint.Response(http.StatusNotFound, "Category not found"),
			endpoint.Response(http.StatusInternalServerError, "Internal server error"),
			endpoint.Security("bearer_auth"),
			endpoint.Tags("Categories"),
		),
	)

	// Endpoint: Update Category
	api.AddEndpoint(
		endpoint.New(
			http.MethodPut, "/categories/{id}",
			endpoint.Handler(categoryController.UpdateCategory),
			endpoint.Summary("Update a category"),
			endpoint.Description("Allows a user to update the details of an existing category"),
			endpoint.Path("id", "string", "ID of the category to update", true),
			endpoint.Body(models.Category{}, "Updated category data", true),
			endpoint.Response(http.StatusOK, "Successfully updated category", endpoint.SchemaResponseOption(models.Category{})),
			endpoint.Response(http.StatusBadRequest, "Invalid category data"),
			endpoint.Response(http.StatusNotFound, "Category not found"),
			endpoint.Response(http.StatusConflict, "Category update conflict"),
			endpoint.Response(http.StatusInternalServerError, "Internal server error"),
			endpoint.Security("bearer_auth"),
			endpoint.Tags("Categories"),
		),
	)

	// Endpoint: Delete Category
	api.AddEndpoint(
		endpoint.New(
			http.MethodDelete, "/categories/{id}",
			endpoint.Handler(categoryController.DeleteCategory),
			endpoint.Summary("Delete a category"),
			endpoint.Description("Allows a user to delete a category by ID"),
			endpoint.Path("id", "string", "ID of the category to delete", true),
			endpoint.Response(http.StatusNoContent, "Successfully deleted category"),
			endpoint.Response(http.StatusBadRequest, "Invalid ID"),
			endpoint.Response(http.StatusNotFound, "Category not found"),
			endpoint.Response(http.StatusInternalServerError, "Internal server error"),
			endpoint.Security("bearer_auth"),
			endpoint.Tags("Categories"),
		),
	)
}
