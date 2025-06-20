// Package repository provides repository implementations for various domain entities.
package repository

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
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
func (r *OrganizationRepositoryImpl) Get(ctx context.Context, id string) (*organization.Organization, error) {
	path := fmt.Sprintf("/organizations/%s", id)
	req, err := createRequest(ctx, http.MethodGet, r.baseURL+path, nil, r.authHeaders)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	resp, err := r.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to execute request: %w", err)
	}
	defer func() {
		if err := resp.Body.Close(); err != nil {
			fmt.Fprintf(os.Stderr, "error closing response body: %v\n", err)
		}
	}()

	if resp.StatusCode != http.StatusOK {
		return nil, handleErrorResponse(resp)
	}

	var orgDTO organizationDTO
	if err := json.NewDecoder(resp.Body).Decode(&orgDTO); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return mapDTOToOrganization(&orgDTO)
}

// List retrieves multiple organizations based on filters
func (r *OrganizationRepositoryImpl) List(ctx context.Context, filters map[string]interface{}) ([]*organization.Organization, *repository.Pagination, error) {
	query := buildQueryParams(filters)
	path := "/organizations"
	if len(query) > 0 {
		path += "?" + query.Encode()
	}
	req, err := createRequest(ctx, http.MethodGet, r.baseURL+path, nil, r.authHeaders)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create request: %w", err)
	}
	resp, err := r.httpClient.Do(req)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to execute request: %w", err)
	}
	defer func() {
		if err := resp.Body.Close(); err != nil {
			fmt.Fprintf(os.Stderr, "error closing response body: %v\n", err)
		}
	}()
	if resp.StatusCode != http.StatusOK {
		return nil, nil, handleErrorResponse(resp)
	}
	var orgDTOs []organizationDTO
	if err := json.NewDecoder(resp.Body).Decode(&orgDTOs); err != nil {
		return nil, nil, fmt.Errorf("failed to decode response: %w", err)
	}
	// Extract pagination information from headers
	pagination := extractPagination(resp.Header)
	// Convert DTOs to domain entities
	orgs := make([]*organization.Organization, len(orgDTOs))
	for i, dto := range orgDTOs {
		org, err := mapDTOToOrganization(&dto)
		if err != nil {
			return nil, nil, fmt.Errorf("failed to convert organization DTO to domain model at index %d: %w", i, err)
		}
		orgs[i] = org
	}
	return orgs, &pagination, nil
}

// Create creates a new organization
func (r *OrganizationRepositoryImpl) Create(ctx context.Context, org *organization.Organization) (*organization.Organization, error) {
	dto := mapOrganizationToDTO(org)

	// Create request body
	body, err := json.Marshal(dto)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal organization: %w", err)
	}

	req, err := createRequest(ctx, http.MethodPost, r.baseURL+"/organizations", body, r.authHeaders)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	resp, err := r.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to execute request: %w", err)
	}
	defer func() {
		if err := resp.Body.Close(); err != nil {
			fmt.Fprintf(os.Stderr, "error closing response body: %v\n", err)
		}
	}()

	if resp.StatusCode != http.StatusCreated {
		return nil, handleErrorResponse(resp)
	}

	// Parse the response to get the created organization
	var createdDTO organizationDTO
	if err := json.NewDecoder(resp.Body).Decode(&createdDTO); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return mapDTOToOrganization(&createdDTO)
}

// Update updates an existing organization
func (r *OrganizationRepositoryImpl) Update(ctx context.Context, org *organization.Organization) (*organization.Organization, error) {
	dto := mapOrganizationToDTO(org)

	// Create request body
	body, err := json.Marshal(dto)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal organization: %w", err)
	}

	path := fmt.Sprintf("/organizations/%s", org.ID)
	req, err := createRequest(ctx, http.MethodPut, r.baseURL+path, body, r.authHeaders)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	resp, err := r.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to execute request: %w", err)
	}
	defer func() {
		if err := resp.Body.Close(); err != nil {
			fmt.Fprintf(os.Stderr, "error closing response body: %v\n", err)
		}
	}()

	if resp.StatusCode != http.StatusOK {
		return nil, handleErrorResponse(resp)
	}

	// Parse the response to get the updated organization
	var updatedDTO organizationDTO
	if err := json.NewDecoder(resp.Body).Decode(&updatedDTO); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return mapDTOToOrganization(&updatedDTO)
}

// Delete deletes an organization by ID
func (r *OrganizationRepositoryImpl) Delete(ctx context.Context, id string) error {
	path := fmt.Sprintf("/organizations/%s", id)
	req, err := createRequest(ctx, http.MethodDelete, r.baseURL+path, nil, r.authHeaders)
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	resp, err := r.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to execute request: %w", err)
	}
	defer func() {
		if err := resp.Body.Close(); err != nil {
			fmt.Fprintf(os.Stderr, "error closing response body: %v\n", err)
		}
	}()

	if resp.StatusCode != http.StatusNoContent {
		return handleErrorResponse(resp)
	}

	return nil
}

// GetUsers retrieves users associated with an organization
func (r *OrganizationRepositoryImpl) GetUsers(ctx context.Context, orgID string) ([]*organization.User, *repository.Pagination, error) {
	path := fmt.Sprintf("/organizations/%s/users", orgID)
	req, err := createRequest(ctx, http.MethodGet, r.baseURL+path, nil, r.authHeaders)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create request: %w", err)
	}

	resp, err := r.httpClient.Do(req)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to execute request: %w", err)
	}
	defer func() {
		if err := resp.Body.Close(); err != nil {
			fmt.Fprintf(os.Stderr, "error closing response body: %v\n", err)
		}
	}()

	if resp.StatusCode != http.StatusOK {
		return nil, nil, handleErrorResponse(resp)
	}

	var userDTOs []userDTO
	if err := json.NewDecoder(resp.Body).Decode(&userDTOs); err != nil {
		return nil, nil, fmt.Errorf("failed to decode response: %w", err)
	}

	// Extract pagination information from headers
	pagination := extractPagination(resp.Header)

	// Convert DTOs to domain entities
	users := make([]*organization.User, len(userDTOs))
	for i, userDTO := range userDTOs {
		users[i] = userDTO.toDomain()
	}

	return users, &pagination, nil
}

// AddUser adds a user to an organization
func (r *OrganizationRepositoryImpl) AddUser(ctx context.Context, orgID string, user *organization.User) error {
	path := fmt.Sprintf("/organizations/%s/users", orgID)

	// Convert domain entity to DTO
	userDTO := fromDomainUser(user)

	// Create request body
	body, err := json.Marshal(userDTO)
	if err != nil {
		return fmt.Errorf("failed to marshal user: %w", err)
	}

	req, err := createRequest(ctx, http.MethodPost, r.baseURL+path, body, r.authHeaders)
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	resp, err := r.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to execute request: %w", err)
	}
	defer func() {
		if err := resp.Body.Close(); err != nil {
			fmt.Fprintf(os.Stderr, "error closing response body: %v\n", err)
		}
	}()

	if resp.StatusCode != http.StatusCreated {
		return handleErrorResponse(resp)
	}

	return nil
}

// RemoveUser removes a user from an organization
func (r *OrganizationRepositoryImpl) RemoveUser(ctx context.Context, orgID string, userID string) error {
	path := fmt.Sprintf("/organizations/%s/users/%s", orgID, userID)
	req, err := createRequest(ctx, http.MethodDelete, r.baseURL+path, nil, r.authHeaders)
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	resp, err := r.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to execute request: %w", err)
	}
	defer func() {
		if err := resp.Body.Close(); err != nil {
			fmt.Fprintf(os.Stderr, "error closing response body: %v\n", err)
		}
	}()

	if resp.StatusCode != http.StatusNoContent {
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

type userDTO struct {
	ID        string `json:"id"`
	Email     string `json:"email"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Role      string `json:"role"`
	Status    string `json:"status"`
}

// mapDTOToOrganization maps a DTO to a domain entity
func mapDTOToOrganization(dto *organizationDTO) (*organization.Organization, error) {
	// Parse timestamps
	createdAt, err := parseTime(dto.CreatedAt)
	if err != nil {
		return nil, fmt.Errorf("failed to parse created_at: %w", err)
	}

	updatedAt, err := parseTime(dto.UpdatedAt)
	if err != nil {
		return nil, fmt.Errorf("failed to parse updated_at: %w", err)
	}

	// Map address
	address := organization.Address{
		Street1: dto.Address.Street1,
		Street2: dto.Address.Street2,
		City:    dto.Address.City,
		State:   dto.Address.State,
		ZipCode: dto.Address.ZipCode,
		Country: dto.Address.Country,
	}

	// Map contact info
	contactInfo := organization.ContactInfo{
		Name:        dto.ContactInfo.Name,
		Email:       dto.ContactInfo.Email,
		PhoneNumber: dto.ContactInfo.PhoneNumber,
		Title:       dto.ContactInfo.Title,
	}

	// Create and return the domain entity
	org := &organization.Organization{
		ID:          strconv.Itoa(dto.ID), // Convert int ID to string for domain model
		AccountID:   dto.AccountID,
		Name:        dto.Name,
		Description: dto.Description,
		Status:      dto.Status,
		Address:     address,
		ContactInfo: contactInfo,
		CreatedAt:   createdAt,
		UpdatedAt:   updatedAt,
		Settings:    dto.Settings,
		Tags:        dto.Tags,
		Industry:    dto.Industry,
		AgentCount:  dto.AgentCount,
	}

	return org, nil
}

// mapOrganizationToDTO maps a domain entity to a DTO
func mapOrganizationToDTO(org *organization.Organization) *organizationDTO {
	// Convert string ID to int
	var idInt int
	if org.ID != "" {
		var err error
		idInt, err = strconv.Atoi(org.ID)
		if err != nil {
			// If conversion fails, use 0 or handle the error appropriately
			// For now, we'll log it and use 0
			fmt.Printf("Warning: Failed to convert organization ID '%s' to int: %v\n", org.ID, err)
		}
	}

	return &organizationDTO{
		ID:          idInt,
		AccountID:   org.AccountID,
		Name:        org.Name,
		Description: org.Description,
		Status:      org.Status,
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
		CreatedAt:  org.CreatedAt.Format(time.RFC3339),
		UpdatedAt:  org.UpdatedAt.Format(time.RFC3339),
		Settings:   org.Settings,
		Tags:       org.Tags,
		Industry:   org.Industry,
		AgentCount: org.AgentCount,
	}
}

// User DTO conversions

// toDomain converts a user DTO to a domain entity
func (dto *userDTO) toDomain() *organization.User {
	return &organization.User{
		ID:        dto.ID,
		Email:     dto.Email,
		FirstName: dto.FirstName,
		LastName:  dto.LastName,
		Role:      dto.Role,
		Status:    dto.Status,
	}
}

// fromDomainUser converts a domain user to a DTO
func fromDomainUser(user *organization.User) *userDTO {
	return &userDTO{
		ID:        user.ID,
		Email:     user.Email,
		FirstName: user.FirstName,
		LastName:  user.LastName,
		Role:      user.Role,
		Status:    user.Status,
	}
}

// Helper functions for this repository have been moved to http_utils.go
// Use the shared utility functions from that file instead of defining them here
