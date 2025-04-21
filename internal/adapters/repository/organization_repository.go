// Package repository provides repository implementations for various domain entities.
package repository

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"time"

	"github.com/greysquirr3l/bishoujo-huntress/internal/domain/organization"
	"github.com/greysquirr3l/bishoujo-huntress/internal/ports/repository"
)

// OrganizationRepositoryImpl implements the organization repository interface
type OrganizationRepositoryImpl struct {
	httpClient  HTTPClient
	baseURL     string
	authHeaders map[string]string
}

// NewOrganizationRepository creates a new organization repository
func NewOrganizationRepository(httpClient HTTPClient, baseURL string, authHeaders map[string]string) repository.OrganizationRepository {
	return &OrganizationRepositoryImpl{
		httpClient:  httpClient,
		baseURL:     baseURL,
		authHeaders: authHeaders,
	}
}

// Get retrieves an organization by ID
func (r *OrganizationRepositoryImpl) Get(ctx context.Context, id int) (*organization.Organization, error) {
	path := fmt.Sprintf("/organizations/%d", id)
	req, err := createRequest(ctx, http.MethodGet, r.baseURL+path, nil, r.authHeaders)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	resp, err := r.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to execute request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, handleErrorResponse(resp)
	}

	var orgDTO organizationDTO
	if err := json.NewDecoder(resp.Body).Decode(&orgDTO); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return orgDTO.toDomain(), nil
}

// List retrieves multiple organizations based on filters
func (r *OrganizationRepositoryImpl) List(ctx context.Context, filters repository.OrganizationFilters) ([]*organization.Organization, repository.Pagination, error) {
	// Construct query parameters
	query := url.Values{}
	if filters.Page > 0 {
		query.Set("page", strconv.Itoa(filters.Page))
	}
	if filters.Limit > 0 {
		query.Set("limit", strconv.Itoa(filters.Limit))
	}
	if filters.AccountID != nil {
		query.Set("account_id", strconv.Itoa(*filters.AccountID))
	}
	if filters.Status != nil {
		query.Set("status", string(*filters.Status))
	}
	if filters.Search != "" {
		query.Set("search", filters.Search)
	}
	if filters.Industry != "" {
		query.Set("industry", filters.Industry)
	}
	// Add tags if present
	for _, tag := range filters.Tags {
		query.Add("tags", tag)
	}
	// Add ordering
	for _, order := range filters.OrderBy {
		direction := "asc"
		if order.Direction == repository.OrderDesc {
			direction = "desc"
		}
		query.Add("sort", fmt.Sprintf("%s:%s", order.Field, direction))
	}
	// Add time range if present
	if filters.TimeRange != nil {
		if filters.TimeRange.Start != "" {
			query.Set("from", filters.TimeRange.Start)
		}
		if filters.TimeRange.End != "" {
			query.Set("to", filters.TimeRange.End)
		}
	}

	path := "/organizations"
	if len(query) > 0 {
		path += "?" + query.Encode()
	}

	req, err := createRequest(ctx, http.MethodGet, r.baseURL+path, nil, r.authHeaders)
	if err != nil {
		return nil, repository.Pagination{}, fmt.Errorf("failed to create request: %w", err)
	}

	resp, err := r.httpClient.Do(req)
	if err != nil {
		return nil, repository.Pagination{}, fmt.Errorf("failed to execute request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, repository.Pagination{}, handleErrorResponse(resp)
	}

	var orgDTOs []organizationDTO
	if err := json.NewDecoder(resp.Body).Decode(&orgDTOs); err != nil {
		return nil, repository.Pagination{}, fmt.Errorf("failed to decode response: %w", err)
	}

	// Extract pagination information from headers
	pagination := extractPagination(resp.Header)

	// Convert DTOs to domain entities
	orgs := make([]*organization.Organization, len(orgDTOs))
	for i, orgDTO := range orgDTOs {
		orgs[i] = orgDTO.toDomain()
	}

	return orgs, pagination, nil
}

// Create creates a new organization
func (r *OrganizationRepositoryImpl) Create(ctx context.Context, org *organization.Organization) error {
	// Convert domain entity to DTO
	orgDTO := fromDomainOrganization(org)

	// Create request body
	body, err := json.Marshal(orgDTO)
	if err != nil {
		return fmt.Errorf("failed to marshal organization: %w", err)
	}

	req, err := createRequest(ctx, http.MethodPost, r.baseURL+"/organizations", body, r.authHeaders)
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	resp, err := r.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to execute request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		return handleErrorResponse(resp)
	}

	// Parse the response to get the created organization's ID
	var createdOrgDTO organizationDTO
	if err := json.NewDecoder(resp.Body).Decode(&createdOrgDTO); err != nil {
		return fmt.Errorf("failed to decode response: %w", err)
	}

	// Update the ID of the original organization object
	org.ID = createdOrgDTO.ID

	return nil
}

// Update updates an existing organization
func (r *OrganizationRepositoryImpl) Update(ctx context.Context, org *organization.Organization) error {
	// Convert domain entity to DTO
	orgDTO := fromDomainOrganization(org)

	// Create request body
	body, err := json.Marshal(orgDTO)
	if err != nil {
		return fmt.Errorf("failed to marshal organization: %w", err)
	}

	path := fmt.Sprintf("/organizations/%d", org.ID)
	req, err := createRequest(ctx, http.MethodPut, r.baseURL+path, body, r.authHeaders)
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	resp, err := r.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to execute request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return handleErrorResponse(resp)
	}

	return nil
}

// Delete deletes an organization by ID
func (r *OrganizationRepositoryImpl) Delete(ctx context.Context, id int) error {
	path := fmt.Sprintf("/organizations/%d", id)
	req, err := createRequest(ctx, http.MethodDelete, r.baseURL+path, nil, r.authHeaders)
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	resp, err := r.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to execute request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNoContent {
		return handleErrorResponse(resp)
	}

	return nil
}

// GetByName retrieves an organization by name within an account
func (r *OrganizationRepositoryImpl) GetByName(ctx context.Context, accountID int, name string) (*organization.Organization, error) {
	// Construct query parameters to filter by name and account
	query := url.Values{}
	query.Set("account_id", strconv.Itoa(accountID))
	query.Set("name", name)

	path := "/organizations?" + query.Encode()
	req, err := createRequest(ctx, http.MethodGet, r.baseURL+path, nil, r.authHeaders)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	resp, err := r.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to execute request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, handleErrorResponse(resp)
	}

	var orgDTOs []organizationDTO
	if err := json.NewDecoder(resp.Body).Decode(&orgDTOs); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	// Check if we found any organizations
	if len(orgDTOs) == 0 {
		return nil, fmt.Errorf("no organization found with name '%s' in account %d", name, accountID)
	}

	// Return the first match (assuming the API filters exactly by name)
	return orgDTOs[0].toDomain(), nil
}

// GetStatistics retrieves organization statistics
func (r *OrganizationRepositoryImpl) GetStatistics(ctx context.Context, id int) (map[string]interface{}, error) {
	path := fmt.Sprintf("/organizations/%d/statistics", id)
	req, err := createRequest(ctx, http.MethodGet, r.baseURL+path, nil, r.authHeaders)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	resp, err := r.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to execute request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, handleErrorResponse(resp)
	}

	var stats map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&stats); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return stats, nil
}

// UpdateTags updates the tags for an organization
func (r *OrganizationRepositoryImpl) UpdateTags(ctx context.Context, id int, tags []string) error {
	path := fmt.Sprintf("/organizations/%d/tags", id)

	// Create request body
	tagData := map[string][]string{"tags": tags}
	body, err := json.Marshal(tagData)
	if err != nil {
		return fmt.Errorf("failed to marshal tags: %w", err)
	}

	req, err := createRequest(ctx, http.MethodPut, r.baseURL+path, body, r.authHeaders)
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	resp, err := r.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to execute request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return handleErrorResponse(resp)
	}

	return nil
}

// organizationDTO is the data transfer object for an organization
type organizationDTO struct {
	ID          int                    `json:"id"`
	AccountID   int                    `json:"account_id"`
	Name        string                 `json:"name"`
	Description string                 `json:"description,omitempty"`
	Status      string                 `json:"status"`
	Address     addressDTO             `json:"address,omitempty"`
	ContactInfo contactInfoDTO         `json:"contact_info,omitempty"`
	CreatedAt   string                 `json:"created_at"`
	UpdatedAt   string                 `json:"updated_at"`
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

// toDomain converts a DTO to a domain entity
func (dto *organizationDTO) toDomain() *organization.Organization {
	org := &organization.Organization{
		ID:          dto.ID,
		AccountID:   dto.AccountID,
		Name:        dto.Name,
		Description: dto.Description,
		Status:      organization.OrganizationStatus(dto.Status),
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
		Settings:   dto.Settings,
		Tags:       dto.Tags,
		Industry:   dto.Industry,
		AgentCount: dto.AgentCount,
	}

	// Parse timestamps
	if dto.CreatedAt != "" {
		org.CreatedAt, _ = parseTime(dto.CreatedAt)
	}
	if dto.UpdatedAt != "" {
		org.UpdatedAt, _ = parseTime(dto.UpdatedAt)
	}

	return org
}

// fromDomainOrganization converts a domain entity to a DTO
func fromDomainOrganization(org *organization.Organization) *organizationDTO {
	return &organizationDTO{
		ID:          org.ID,
		AccountID:   org.AccountID,
		Name:        org.Name,
		Description: org.Description,
		Status:      string(org.Status),
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
		Settings:   org.Settings,
		Tags:       org.Tags,
		Industry:   org.Industry,
		AgentCount: org.AgentCount,
		CreatedAt:  org.CreatedAt.Format(time.RFC3339),
		UpdatedAt:  org.UpdatedAt.Format(time.RFC3339),
	}
}
