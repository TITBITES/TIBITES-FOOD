package httpapi

// Transport DTOs aligned with OpenAPI schemas (field casing exact)

type CategoryDTO struct {
	ID          string  `json:"id"`
	Name        string  `json:"name"`
	Description *string `json:"description,omitempty"`
	SortOrder   int     `json:"sortOrder"`
	IsActive    bool    `json:"isActive"`
}

type MenuItemDTO struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	CategoryID  string `json:"categoryId"`
	PriceCents  int    `json:"priceCents"`
	Currency    string `json:"currency"`
	IsAvailable bool   `json:"isAvailable"`
}

type PaginationDTO struct {
	Page     int `json:"page"`
	PageSize int `json:"pageSize"`
	Total    int `json:"total"`
}

type CategoriesResponse struct {
	Categories []CategoryDTO `json:"categories"`
}

type MenuItemsResponse struct {
	Items      []MenuItemDTO `json:"items"`
	Pagination PaginationDTO `json:"pagination"`
}

// Orders DTOs (HTTP transport only)

type OrderItemDTO struct {
	ID             string `json:"id"`
	OrderID        string `json:"orderId"`
	MenuItemID     string `json:"menuItemId"`
	NameSnapshot   string `json:"nameSnapshot"`
	Quantity       int    `json:"quantity"`
	UnitPriceCents int    `json:"unitPriceCents"`
}

type OrderDTO struct {
	ID                  string   `json:"id"`
	CustomerID          *string  `json:"customerId"`
	Status              string   `json:"status"`
	TotalCents          int      `json:"totalCents"`
	Currency            string   `json:"currency"`
	DeliveryAddressID   *string  `json:"deliveryAddressId"`
	DeliveryInstructions *string `json:"deliveryInstructions"`
	TipCents            int      `json:"tipCents"`
	DeliveryEtaMinutes  *int     `json:"deliveryEtaMinutes"`
	DeliveredAt         *string  `json:"deliveredAt"`
	CourierName         *string  `json:"courierName"`
	DeliveryTrackingURL *string  `json:"deliveryTrackingUrl"`
	CreatedAt           string   `json:"createdAt"`
	UpdatedAt           string   `json:"updatedAt"`
	Items               []OrderItemDTO `json:"items"`
}

type OrdersResponse struct {
	Orders     []OrderDTO    `json:"orders"`
	Pagination PaginationDTO `json:"pagination"`
}

// OrderResponse mirrors Order schema for single-order endpoints.
// We keep a separate type name to follow task requirements while matching OpenAPI shape.
type OrderResponse = OrderDTO

// Payments DTOs (HTTP transport only)

type PaymentIntentDTO struct {
	ID                string  `json:"id"`
	OrderID           string  `json:"orderId"`
	Provider          string  `json:"provider"`
	ProviderIntentID  string  `json:"providerIntentId"`
	Status            string  `json:"status"`
	ClientSecret      *string `json:"clientSecret,omitempty"`
	AuthorizationURL  *string `json:"authorizationUrl,omitempty"`
	AmountCents       int     `json:"amountCents"`
	Currency          string  `json:"currency"`
	CreatedAt         string  `json:"createdAt"`
	UpdatedAt         string  `json:"updatedAt"`
}

type PaymentWebhookEventDTO struct {
	// Minimal placeholder for webhook receipt; opaque per provider
	Event   string                 `json:"event"`
	Data    map[string]interface{} `json:"data"`
}

// PaymentIntentResponse mirrors PaymentIntent schema for single endpoints.
type PaymentIntentResponse = PaymentIntentDTO
