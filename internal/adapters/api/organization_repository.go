// Package api provides API client implementations for Huntress resources.
package api

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"time"

	// Correct the import path for the HTTP client
	"github.com/greysquirr3l/bishoujo-huntress/internal/domain/organization"
	httpClient "github.com/greysquirr3l/bishoujo-huntress/internal/infrastructure/http"
	"github.com/greysquirr3l/bishoujo-huntress/internal/ports/repository"
)

// OrganizationRepository implements the repository.OrganizationRepository interface
type OrganizationRepository struct {
	// Use the imported http client type
	httpClient *httpClient.Client
	baseURL    string
}

// NewOrganizationRepository creates a new OrganizationRepository instance
// Update the function signature to use the correct client type
func NewOrganizationRepository(client *httpClient.Client, baseURL string) *OrganizationRepository {
	return &OrganizationRepository{
		httpClient: client,
		baseURL:    baseURL,
	}
}

// Get retrieves a specific organization by its ID
func (r *OrganizationRepository) Get(ctx context.Context, id string) (*organization.Organization, error) {
	// Construct the full URL using the base URL and path
	path := fmt.Sprintf("/organizations/%s", url.PathEscape(id))
	// endpoint := fmt.Sprintf("%s%s", r.baseURL, path) // BaseURL is handled by the client's Do method

	var orgDTO organizationDTO
	// Use the client's Do method for the request
	resp, err := r.httpClient.Do(ctx, http.MethodGet, path, nil, &orgDTO, nil) // Pass nil for RequestOptions if none
	if err != nil {
		return nil, fmt.Errorf("failed to get organization: %w", err)
	}
	if resp != nil {
		defer func() {
			if err := resp.Body.Close(); err != nil {
				fmt.Fprintf(os.Stderr, "error closing response body: %v\n", err)
			}
		}()
	}
	// Check response status code if needed, though Do might handle it
	if resp.StatusCode != http.StatusOK {
		// Handle non-OK status, potentially using APIError from the client package
		return nil, fmt.Errorf("unexpected status code %d when getting organization", resp.StatusCode)
	}

	// Convert DTO to domain entity
	org, err := r.convertDTOToOrganization(orgDTO)
	if err != nil {
		return nil, fmt.Errorf("invalid organization data: %w", err)
	}

	return org, nil
}

// List retrieves all organizations with optional filtering
func (r *OrganizationRepository) List(ctx context.Context, params *organization.ListParams) ([]*organization.Organization, *repository.Pagination, error) {
	path := "/organizations"

	// Add query parameters for filtering
	queryParams := url.Values{}
	if params != nil {
		if params.Page > 0 {
			queryParams.Set("page", fmt.Sprintf("%d", params.Page))
		}
		if params.Limit > 0 {
			// Use "per_page" as per Huntress API docs
			queryParams.Set("per_page", fmt.Sprintf("%d", params.Limit))
		}
		// AccountID filtering might not be supported directly, check API docs
		// if params.AccountID > 0 {
		// 	queryParams.Set("account_id", fmt.Sprintf("%d", params.AccountID))
		// }
		if params.Status != "" {
			queryParams.Set("status", params.Status)
		}
		if params.Search != "" {
			queryParams.Set("search", params.Search)
		}
		// Industry filtering might not be supported directly, check API docs
		// if params.Industry != "" {
		// 	queryParams.Set("industry", params.Industry)
		// }
		// Tags filtering might not be supported directly, check API docs
		// for _, tag := range params.Tags {
		// 	queryParams.Add("tags", tag)
		// }
	}

	// Prepare RequestOptions
	reqOpts := &httpClient.RequestOptions{
		Query: queryParams,
	}

	// Define the expected response structure based on Huntress API
	// The API likely returns a list directly or a structure containing the list and pagination
	var response struct {
		// Assuming the API returns a structure like this based on common patterns
		Data       []organizationDTO     `json:"data"`       // Check actual API response key
		Pagination httpClient.Pagination `json:"pagination"` // Use pagination struct from http client
	}

	// Use the client's Do method
	resp, err := r.httpClient.Do(ctx, http.MethodGet, path, nil, &response, reqOpts)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to list organizations: %w", err)
	}
	if resp != nil {
		defer func() {
			if err := resp.Body.Close(); err != nil {
				fmt.Fprintf(os.Stderr, "error closing response body: %v\n", err)
			}
		}()
	}
	if resp.StatusCode != http.StatusOK {
		return nil, nil, fmt.Errorf("unexpected status code %d when listing organizations", resp.StatusCode)
	}

	// Convert DTOs to domain entities
	organizations := make([]*organization.Organization, 0, len(response.Data))
	for _, dto := range response.Data {
		org, err := r.convertDTOToOrganization(dto)
		if err != nil {
			// Consider logging the error and skipping the problematic entry
			// or returning the error immediately depending on requirements.
			return nil, nil, fmt.Errorf("invalid organization data in list response: %w", err)
		}
		organizations = append(organizations, org)
	}

	// Map the client's pagination struct to the repository's pagination struct
	pagination := &repository.Pagination{
		Page:       response.Pagination.CurrentPage,
		PerPage:    response.Pagination.ItemsPerPage,
		TotalItems: response.Pagination.TotalItems,
		TotalPages: response.Pagination.TotalPages,
	}

	return organizations, pagination, nil
}

// Create creates a new organization
func (r *OrganizationRepository) Create(ctx context.Context, org *organization.Organization) (*organization.Organization, error) {
	path := "/organizations"

	// Validate the organization before sending
	// Domain validation should ideally happen in the application layer or domain service
	// if err := org.Validate(); err != nil {
	// 	return nil, fmt.Errorf("invalid organization: %w", err)
	// }

	// Convert domain entity to DTO for API request
	dto := r.convertOrganizationToDTO(org)

	var responseDTO organizationDTO
	// Use the client's Do method
	resp, err := r.httpClient.Do(ctx, http.MethodPost, path, dto, &responseDTO, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create organization: %w", err)
	}
	if resp != nil {
		defer func() {
			if err := resp.Body.Close(); err != nil {
				fmt.Fprintf(os.Stderr, "error closing response body: %v\n", err)
			}
		}()
	}
	// Check for expected status code (e.g., 201 Created)
	if resp.StatusCode != http.StatusCreated {
		return nil, fmt.Errorf("unexpected status code %d when creating organization", resp.StatusCode)
	}

	// Convert response DTO back to domain entity
	createdOrg, err := r.convertDTOToOrganization(responseDTO)
	if err != nil {
		return nil, fmt.Errorf("invalid organization data in create response: %w", err)
	}

	return createdOrg, nil
}

// Update updates an existing organization
func (r *OrganizationRepository) Update(ctx context.Context, org *organization.Organization) (*organization.Organization, error) {
	if org.ID == "" {
		return nil, fmt.Errorf("organization ID is required for update")
	}

	// Domain validation should ideally happen in the application layer or domain service
	// if err := org.Validate(); err != nil {
	// 	return nil, fmt.Errorf("invalid organization: %w", err)
	// }

	path := fmt.Sprintf("/organizations/%s", url.PathEscape(org.ID))

	// Convert domain entity to DTO for API request
	dto := r.convertOrganizationToDTO(org) // Consider sending only changed fields (PATCH) if API supports it

	var responseDTO organizationDTO
	// Use PUT or PATCH depending on API design (PUT usually replaces, PATCH updates)
	resp, err := r.httpClient.Do(ctx, http.MethodPut, path, dto, &responseDTO, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to update organization: %w", err)
	}
	if resp != nil {
		defer func() {
			if err := resp.Body.Close(); err != nil {
				fmt.Fprintf(os.Stderr, "error closing response body: %v\n", err)
			}
		}()
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code %d when updating organization", resp.StatusCode)
	}

	// Convert response DTO back to domain entity
	updatedOrg, err := r.convertDTOToOrganization(responseDTO)
	if err != nil {
		return nil, fmt.Errorf("invalid organization data in update response: %w", err)
	}

	return updatedOrg, nil
}

// Delete removes an organization by its ID
func (r *OrganizationRepository) Delete(ctx context.Context, id string) error {
	path := fmt.Sprintf("/organizations/%s", url.PathEscape(id))

	// Use the client's Do method, expecting no response body on success (nil for result)
	resp, err := r.httpClient.Do(ctx, http.MethodDelete, path, nil, nil, nil)
	if err != nil {
		return fmt.Errorf("failed to delete organization: %w", err)
	}
	if resp != nil {
		defer func() {
			if err := resp.Body.Close(); err != nil {
				fmt.Fprintf(os.Stderr, "error closing response body: %v\n", err)
			}
		}()
	}
	// Check for expected status code (e.g., 204 No Content)
	if resp.StatusCode != http.StatusNoContent {
		return fmt.Errorf("unexpected status code %d when deleting organization", resp.StatusCode)
	}

	return nil
}

// GetUsers retrieves all users associated with an organization
func (r *OrganizationRepository) GetUsers(ctx context.Context, organizationID string) ([]*organization.User, *repository.Pagination, error) {
	path := fmt.Sprintf("/organizations/%s/users", url.PathEscape(organizationID))

	// Define expected response structure
	var response struct {
		Data       []userDTO             `json:"data"` // Check actual API response key
		Pagination httpClient.Pagination `json:"pagination"`
	}

	// Use the client's Do method
	resp, err := r.httpClient.Do(ctx, http.MethodGet, path, nil, &response, nil)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to get organization users: %w", err)
	}
	if resp != nil {
		defer func() {
			if err := resp.Body.Close(); err != nil {
				fmt.Fprintf(os.Stderr, "error closing response body: %v\n", err)
			}
		}()
	}
	if resp.StatusCode != http.StatusOK {
		return nil, nil, fmt.Errorf("unexpected status code %d when getting organization users", resp.StatusCode)
	}

	// Convert DTOs to domain entities
	users := make([]*organization.User, 0, len(response.Data))
	for _, dto := range response.Data {
		user := r.convertDTOToUser(dto)
		users = append(users, user)
	}

	pagination := &repository.Pagination{
		Page:       response.Pagination.CurrentPage,
		PerPage:    response.Pagination.ItemsPerPage,
		TotalItems: response.Pagination.TotalItems,
		TotalPages: response.Pagination.TotalPages,
	}

	return users, pagination, nil
}

// AddUser adds a user to an organization
func (r *OrganizationRepository) AddUser(ctx context.Context, organizationID string, user *organization.User) error {
	path := fmt.Sprintf("/organizations/%s/users", url.PathEscape(organizationID))

	// Convert domain entity to DTO for API request
	dto := r.convertUserToDTO(user)

	// Use the client's Do method, potentially expecting a user DTO or no body on success
	resp, err := r.httpClient.Do(ctx, http.MethodPost, path, dto, nil, nil) // Assuming no response body needed
	if err != nil {
		return fmt.Errorf("failed to add user to organization: %w", err)
	}
	if resp != nil {
		defer func() {
			if err := resp.Body.Close(); err != nil {
				fmt.Fprintf(os.Stderr, "error closing response body: %v\n", err)
			}
		}()
	}
	// Check for expected status code (e.g., 201 Created or 200 OK)
	if resp.StatusCode != http.StatusCreated && resp.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected status code %d when adding user to organization", resp.StatusCode)
	}

	return nil
}

// RemoveUser removes a user from an organization
func (r *OrganizationRepository) RemoveUser(ctx context.Context, organizationID string, userID string) error {
	path := fmt.Sprintf("/organizations/%s/users/%s", url.PathEscape(organizationID), url.PathEscape(userID))

	// Use the client's Do method, expecting no response body on success
	resp, err := r.httpClient.Do(ctx, http.MethodDelete, path, nil, nil, nil)
	if err != nil {
		return fmt.Errorf("failed to remove user from organization: %w", err)
	}
	if resp != nil {
		defer func() {
			if err := resp.Body.Close(); err != nil {
				fmt.Fprintf(os.Stderr, "error closing response body: %v\n", err)
			}
		}()
	}
	// Check for expected status code (e.g., 204 No Content)
	if resp.StatusCode != http.StatusNoContent {
		return fmt.Errorf("unexpected status code %d when removing user from organization", resp.StatusCode)
	}

	return nil
}

// Data transfer objects (DTOs) for API communication
type organizationDTO struct {
	ID          string                 `json:"id"`
	AccountID   int                    `json:"account_id"`
	Name        string                 `json:"name"`
	Description string                 `json:"description,omitempty"`
	CreatedAt   string                 `json:"created_at"`
	UpdatedAt   string                 `json:"updated_at"`
	Status      string                 `json:"status"`
	Address     addressDTO             `json:"address,omitempty"`
	ContactInfo contactInfoDTO         `json:"contact_info,omitempty"`
	Settings    map[string]interface{} `json:"settings,omitempty"`
	Tags        []string               `json:"tags,omitempty"`
	Industry    string                 `json:"industry,omitempty"`
	AgentCount  int                    `json:"agent_count,omitempty"`
}

type addressDTO struct {
	Street1 string `json:"street1,omitempty"`
	Street2 string `json:"street2,omitempty"`
	City    string `json:"city,omitempty"`
	State   string `json:"state,omitempty"`
	ZipCode string `json:"zip_code,omitempty"`
	Country string `json:"country,omitempty"`
}

type contactInfoDTO struct {
	Name        string `json:"name,omitempty"`
	Email       string `json:"email,omitempty"`
	PhoneNumber string `json:"phone_number,omitempty"`
	Title       string `json:"title,omitempty"`
}

type userDTO struct {
	ID        string `json:"id"`
	Email     string `json:"email"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Role      string `json:"role"`
	Status    string `json:"status"`
}

// Helper functions to convert between DTOs and domain entities
func (r *OrganizationRepository) convertDTOToOrganization(dto organizationDTO) (*organization.Organization, error) {
	createdAt, err := parseTime(dto.CreatedAt)
	if err != nil {
		return nil, fmt.Errorf("invalid created_at timestamp: %w", err)
	}

	updatedAt, err := parseTime(dto.UpdatedAt)
	if err != nil {
		return nil, fmt.Errorf("invalid updated_at timestamp: %w", err)
	}

	org := &organization.Organization{
		ID:          dto.ID,
		AccountID:   dto.AccountID,
		Name:        dto.Name,
		Description: dto.Description,
		Status:      dto.Status,
		CreatedAt:   createdAt,
		UpdatedAt:   updatedAt,
		Settings:    dto.Settings,
		Tags:        dto.Tags,
		Industry:    dto.Industry,
		AgentCount:  dto.AgentCount,
		Address: organization.Address{
			Street1: dto.Address.Street1,
			Street2: dto.Address.Street2,
			City:    dto.Address.City,
			State:   dto.Address.State,
			ZipCode: dto.Address.ZipCode,
			Country: dto.Address.Country,
		},
		ContactInfo: organization.ContactInfo{
			Name:        dto.ContactInfo.Name,
			Email:       dto.ContactInfo.Email,
			PhoneNumber: dto.ContactInfo.PhoneNumber,
			Title:       dto.ContactInfo.Title,
		},
	}

	return org, nil
}

func (r *OrganizationRepository) convertOrganizationToDTO(org *organization.Organization) organizationDTO {
	return organizationDTO{
		ID:          org.ID,
		AccountID:   org.AccountID,
		Name:        org.Name,
		Description: org.Description,
		Status:      org.Status,
		CreatedAt:   org.CreatedAt.Format(time.RFC3339),
		UpdatedAt:   org.UpdatedAt.Format(time.RFC3339),
		Settings:    org.Settings,
		Tags:        org.Tags,
		Industry:    org.Industry,
		AgentCount:  org.AgentCount,
		Address: addressDTO{
			Street1: org.Address.Street1,
			Street2: org.Address.Street2,
			City:    org.Address.City,
			State:   org.Address.State,
			ZipCode: org.Address.ZipCode,
			Country: org.Address.Country,
		},
		ContactInfo: contactInfoDTO{
			Name:        org.ContactInfo.Name,
			Email:       org.ContactInfo.Email,
			PhoneNumber: org.ContactInfo.PhoneNumber,
			Title:       org.ContactInfo.Title,
		},
	}
}

func (r *OrganizationRepository) convertDTOToUser(dto userDTO) *organization.User {
	return &organization.User{
		ID:        dto.ID,
		Email:     dto.Email,
		FirstName: dto.FirstName,
		LastName:  dto.LastName,
		Role:      dto.Role,
		Status:    dto.Status,
	}
}

func (r *OrganizationRepository) convertUserToDTO(user *organization.User) userDTO {
	return userDTO{
		ID:        user.ID,
		Email:     user.Email,
		FirstName: user.FirstName,
		LastName:  user.LastName,
		Role:      user.Role,
		Status:    user.Status,
	}
}

// Helper function to parse time strings from the API
func parseTime(timeStr string) (time.Time, error) {
	if timeStr == "" {
		return time.Time{}, nil
	}
	t, err := time.Parse(time.RFC3339, timeStr)
	if err != nil {
		return time.Time{}, fmt.Errorf("failed to parse time: %w", err)
	}
	return t, nil
}
