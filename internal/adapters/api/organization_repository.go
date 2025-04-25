package api

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"time"

	"github.com/greysquirr3l/bishoujo-huntress/internal/adapters/http/client"
	"github.com/greysquirr3l/bishoujo-huntress/internal/domain/organization"
	"github.com/greysquirr3l/bishoujo-huntress/internal/ports/repository"
)

// OrganizationRepository implements the repository.OrganizationRepository interface
type OrganizationRepository struct {
	httpClient client.Client
	baseURL    string
}

// NewOrganizationRepository creates a new OrganizationRepository instance
func NewOrganizationRepository(httpClient client.Client, baseURL string) *OrganizationRepository {
	return &OrganizationRepository{
		httpClient: httpClient,
		baseURL:    baseURL,
	}
}

// Get retrieves a specific organization by its ID
func (r *OrganizationRepository) Get(ctx context.Context, id string) (*organization.Organization, error) {
	endpoint := fmt.Sprintf("%s/organizations/%s", r.baseURL, url.PathEscape(id))

	var orgDTO organizationDTO
	if err := r.httpClient.Get(ctx, endpoint, &orgDTO); err != nil {
		return nil, fmt.Errorf("failed to get organization: %w", err)
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
	endpoint := fmt.Sprintf("%s/organizations", r.baseURL)

	// Add query parameters for filtering
	queryParams := url.Values{}
	if params != nil {
		if params.Page > 0 {
			queryParams.Set("page", fmt.Sprintf("%d", params.Page))
		}
		if params.Limit > 0 {
			queryParams.Set("limit", fmt.Sprintf("%d", params.Limit))
		}
		if params.AccountID > 0 {
			queryParams.Set("account_id", fmt.Sprintf("%d", params.AccountID))
		}
		if params.Status != "" {
			queryParams.Set("status", params.Status)
		}
		if params.Search != "" {
			queryParams.Set("search", params.Search)
		}
		if params.Industry != "" {
			queryParams.Set("industry", params.Industry)
		}
		// Add tags if present
		for _, tag := range params.Tags {
			queryParams.Add("tags", tag)
		}
	}

	if len(queryParams) > 0 {
		endpoint = fmt.Sprintf("%s?%s", endpoint, queryParams.Encode())
	}

	var response struct {
		Organizations []organizationDTO `json:"organizations"`
		Pagination    struct {
			CurrentPage  int `json:"current_page"`
			TotalPages   int `json:"total_pages"`
			TotalItems   int `json:"total_items"`
			ItemsPerPage int `json:"items_per_page"`
		} `json:"pagination"`
	}

	if err := r.httpClient.Get(ctx, endpoint, &response); err != nil {
		return nil, nil, fmt.Errorf("failed to list organizations: %w", err)
	}

	// Convert DTOs to domain entities
	organizations := make([]*organization.Organization, 0, len(response.Organizations))
	for _, dto := range response.Organizations {
		org, err := r.convertDTOToOrganization(dto)
		if err != nil {
			return nil, nil, fmt.Errorf("invalid organization data: %w", err)
		}
		organizations = append(organizations, org)
	}

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
	endpoint := fmt.Sprintf("%s/organizations", r.baseURL)

	// Validate the organization before sending
	if err := org.Validate(); err != nil {
		return nil, fmt.Errorf("invalid organization: %w", err)
	}

	// Convert domain entity to DTO for API request
	dto := r.convertOrganizationToDTO(org)

	var responseDTO organizationDTO
	if err := r.httpClient.Post(ctx, endpoint, dto, &responseDTO, http.StatusCreated); err != nil {
		return nil, fmt.Errorf("failed to create organization: %w", err)
	}

	// Convert response DTO back to domain entity
	createdOrg, err := r.convertDTOToOrganization(responseDTO)
	if err != nil {
		return nil, fmt.Errorf("invalid organization data in response: %w", err)
	}

	return createdOrg, nil
}

// Update updates an existing organization
func (r *OrganizationRepository) Update(ctx context.Context, org *organization.Organization) (*organization.Organization, error) {
	if org.ID == "" {
		return nil, fmt.Errorf("organization ID is required for update")
	}

	// Validate the organization before sending
	if err := org.Validate(); err != nil {
		return nil, fmt.Errorf("invalid organization: %w", err)
	}

	endpoint := fmt.Sprintf("%s/organizations/%s", r.baseURL, url.PathEscape(org.ID))

	// Convert domain entity to DTO for API request
	dto := r.convertOrganizationToDTO(org)

	var responseDTO organizationDTO
	if err := r.httpClient.Put(ctx, endpoint, dto, &responseDTO); err != nil {
		return nil, fmt.Errorf("failed to update organization: %w", err)
	}

	// Convert response DTO back to domain entity
	updatedOrg, err := r.convertDTOToOrganization(responseDTO)
	if err != nil {
		return nil, fmt.Errorf("invalid organization data in response: %w", err)
	}

	return updatedOrg, nil
}

// Delete removes an organization by its ID
func (r *OrganizationRepository) Delete(ctx context.Context, id string) error {
	endpoint := fmt.Sprintf("%s/organizations/%s", r.baseURL, url.PathEscape(id))

	if err := r.httpClient.Delete(ctx, endpoint, nil); err != nil {
		return fmt.Errorf("failed to delete organization: %w", err)
	}

	return nil
}

// GetUsers retrieves all users associated with an organization
func (r *OrganizationRepository) GetUsers(ctx context.Context, organizationID string) ([]*organization.User, *repository.Pagination, error) {
	endpoint := fmt.Sprintf("%s/organizations/%s/users", r.baseURL, url.PathEscape(organizationID))

	var response struct {
		Users      []userDTO `json:"users"`
		Pagination struct {
			CurrentPage  int `json:"current_page"`
			TotalPages   int `json:"total_pages"`
			TotalItems   int `json:"total_items"`
			ItemsPerPage int `json:"items_per_page"`
		} `json:"pagination"`
	}

	if err := r.httpClient.Get(ctx, endpoint, &response); err != nil {
		return nil, nil, fmt.Errorf("failed to get organization users: %w", err)
	}

	// Convert DTOs to domain entities
	users := make([]*organization.User, 0, len(response.Users))
	for _, dto := range response.Users {
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
	endpoint := fmt.Sprintf("%s/organizations/%s/users", r.baseURL, url.PathEscape(organizationID))

	// Convert domain entity to DTO for API request
	dto := r.convertUserToDTO(user)

	if err := r.httpClient.Post(ctx, endpoint, dto, nil, http.StatusCreated); err != nil {
		return fmt.Errorf("failed to add user to organization: %w", err)
	}

	return nil
}

// RemoveUser removes a user from an organization
func (r *OrganizationRepository) RemoveUser(ctx context.Context, organizationID string, userID string) error {
	endpoint := fmt.Sprintf("%s/organizations/%s/users/%s", r.baseURL, url.PathEscape(organizationID), url.PathEscape(userID))

	if err := r.httpClient.Delete(ctx, endpoint, nil); err != nil {
		return fmt.Errorf("failed to remove user from organization: %w", err)
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
	return time.Parse(time.RFC3339, timeStr)
}
