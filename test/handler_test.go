package test

import (
	"api-otto/internal/domain"
	"api-otto/internal/handler"
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/julienschmidt/httprouter"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// Mock Services
type MockBrandService struct {
	mock.Mock
}

// Update implements domain.BrandService.
func (m *MockBrandService) Update(brand *domain.Brand) error {
	panic("unimplemented")
}

type MockVoucherService struct {
	mock.Mock
}

// Update implements domain.VoucherService.
func (m *MockVoucherService) Update(voucher *domain.Voucher) error {
	panic("unimplemented")
}

type MockTransactionService struct {
	mock.Mock
}

// GetCustomerTransactions implements domain.TransactionService.
func (m *MockTransactionService) GetCustomerTransactions(customerID int64) ([]domain.Transaction, error) {
	panic("unimplemented")
}

// Brand Service Mock Methods
func (m *MockBrandService) Create(brand *domain.Brand) error {
	args := m.Called(brand)
	return args.Error(0)
}

func (m *MockBrandService) GetByID(id int64) (*domain.Brand, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Brand), args.Error(1)
}

func (m *MockBrandService) GetAll() ([]domain.Brand, error) {
	args := m.Called()
	return args.Get(0).([]domain.Brand), args.Error(1)
}

func (m *MockBrandService) Delete(id int64) error {
	args := m.Called(id)
	return args.Error(0)
}

func (m *MockBrandService) List() ([]domain.Brand, error) {
	args := m.Called()
	return args.Get(0).([]domain.Brand), args.Error(1)
}

// Voucher Service Mock Methods
func (m *MockVoucherService) Create(voucher *domain.Voucher) error {
	args := m.Called(voucher)
	return args.Error(0)
}

func (m *MockVoucherService) GetByID(id int64) (*domain.Voucher, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Voucher), args.Error(1)
}

func (m *MockVoucherService) GetByBrandID(brandID int64) ([]domain.Voucher, error) {
	args := m.Called(brandID)
	return args.Get(0).([]domain.Voucher), args.Error(1)
}

func (m *MockVoucherService) Delete(id int64) error {
	args := m.Called(id)
	return args.Error(0)
}

func (m *MockVoucherService) List() ([]domain.Voucher, error) {
	args := m.Called()
	return args.Get(0).([]domain.Voucher), args.Error(1)
}

// Transaction Service Mock Methods
func (m *MockTransactionService) CreateRedemption(transaction *domain.Transaction) error {
	args := m.Called(transaction)
	return args.Error(0)
}

func (m *MockTransactionService) GetTransactionByID(id int64) (*domain.Transaction, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Transaction), args.Error(1)
}

// Brand Handler Tests
func TestBrandHandler_Create(t *testing.T) {
	tests := []struct {
		name           string
		requestBody    interface{}
		mockBehavior   func(service *MockBrandService)
		expectedStatus int
		expectedBody   string
	}{
		{
			name: "Success Create Brand",
			requestBody: domain.Brand{
				Name:        "Test Brand",
				Description: "Test Description",
			},
			mockBehavior: func(service *MockBrandService) {
				service.On("Create", mock.Anything).Return(nil)
			},
			expectedStatus: http.StatusCreated,
			expectedBody: `{
    "status":201,
    "message":"Brand created successfully",
    "data":{
        "id":0,
        "name":"Test Brand",
        "description":"Test Description",
        "created_at":"0001-01-01T00:00:00Z",
        "updated_at":"0001-01-01T00:00:00Z"
    }
}`,
		},
		{
			name: "Invalid Request - Empty Name",
			requestBody: domain.Brand{
				Description: "Test Description",
			},
			mockBehavior:   func(service *MockBrandService) {},
			expectedStatus: http.StatusBadRequest,
			expectedBody:   `{"status":400,"message":"Nama brand tidak boleh kosong"}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockService := new(MockBrandService)
			tt.mockBehavior(mockService)
			handler := handler.NewBrandHandler(mockService)

			body, _ := json.Marshal(tt.requestBody)
			req := httptest.NewRequest(http.MethodPost, "/brand", bytes.NewBuffer(body))
			rec := httptest.NewRecorder()

			handler.Create(rec, req, nil)

			assert.Equal(t, tt.expectedStatus, rec.Code)
			assert.JSONEq(t, tt.expectedBody, rec.Body.String())
			mockService.AssertExpectations(t)
		})
	}
}

func TestBrandHandler_GetByID(t *testing.T) {
	tests := []struct {
		name           string
		brandID        string
		mockBehavior   func(service *MockBrandService)
		expectedStatus int
		expectedBody   string
	}{
		{
			name:    "Success Get Brand",
			brandID: "1",
			mockBehavior: func(service *MockBrandService) {
				service.On("GetByID", int64(1)).Return(&domain.Brand{
					ID:          1,
					Name:        "Test Brand",
					Description: "Test Description",
					CreatedAt:   time.Date(2024, 3, 1, 0, 0, 0, 0, time.UTC),
					UpdatedAt:   time.Date(2024, 3, 1, 0, 0, 0, 0, time.UTC),
				}, nil)
			},
			expectedStatus: http.StatusOK,
			expectedBody:   `{"status":200,"message":"Success","data":{"id":1,"name":"Test Brand","description":"Test Description","created_at":"2024-03-01T00:00:00Z","updated_at":"2024-03-01T00:00:00Z"}}`,
		},
		{
			name:    "Brand Not Found",
			brandID: "999",
			mockBehavior: func(service *MockBrandService) {
				service.On("GetByID", int64(999)).Return(nil, nil)
			},
			expectedStatus: http.StatusNotFound,
			expectedBody:   `{"status":404,"message":"Brand not found"}`,
		},
		{
			name:           "Invalid ID Format",
			brandID:        "abc",
			mockBehavior:   func(service *MockBrandService) {},
			expectedStatus: http.StatusBadRequest,
			expectedBody:   `{"status":400,"message":"Invalid ID"}`,
		},
		{
			name:           "Empty ID",
			brandID:        "",
			mockBehavior:   func(service *MockBrandService) {},
			expectedStatus: http.StatusBadRequest,
			expectedBody:   `{"status":400,"message":"Invalid ID"}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockService := new(MockBrandService)
			tt.mockBehavior(mockService)
			handler := handler.NewBrandHandler(mockService)

			// Create request with params
			req := httptest.NewRequest(http.MethodGet, "/brand/"+tt.brandID, nil)
			rec := httptest.NewRecorder()
			params := httprouter.Params{httprouter.Param{Key: "id", Value: tt.brandID}}

			handler.GetByID(rec, req, params)

			assert.Equal(t, tt.expectedStatus, rec.Code)
			assert.JSONEq(t, tt.expectedBody, rec.Body.String())
			mockService.AssertExpectations(t)
		})
	}
}

func TestBrandHandler_GetAll(t *testing.T) {
	tests := []struct {
		name           string
		mockBehavior   func(service *MockBrandService)
		expectedStatus int
		expectedBody   string
	}{
		{
			name: "Success Get All Brands",
			mockBehavior: func(service *MockBrandService) {
				service.On("List").Return([]domain.Brand{
					{
						ID:          1,
						Name:        "Brand One",
						Description: "Description One",
						CreatedAt:   time.Date(2024, 3, 1, 0, 0, 0, 0, time.UTC),
						UpdatedAt:   time.Date(2024, 3, 1, 0, 0, 0, 0, time.UTC),
					},
					{
						ID:          2,
						Name:        "Brand Two",
						Description: "Description Two",
						CreatedAt:   time.Date(2024, 3, 1, 0, 0, 0, 0, time.UTC),
						UpdatedAt:   time.Date(2024, 3, 1, 0, 0, 0, 0, time.UTC),
					},
				}, nil)
			},
			expectedStatus: http.StatusOK,
			expectedBody: `{
				"status": 200,
				"message": "Success",
				"data": [
					{
						"id": 1,
						"name": "Brand One",
						"description": "Description One",
						"created_at": "2024-03-01T00:00:00Z",
						"updated_at": "2024-03-01T00:00:00Z"
					},
					{
						"id": 2,
						"name": "Brand Two",
						"description": "Description Two",
						"created_at": "2024-03-01T00:00:00Z",
						"updated_at": "2024-03-01T00:00:00Z"
					}
				]
			}`,
		},
		{
			name: "Empty Brand List",
			mockBehavior: func(service *MockBrandService) {
				service.On("List").Return([]domain.Brand{}, nil)
			},
			expectedStatus: http.StatusOK,
			expectedBody: `{
				"status": 200,
				"message": "Success",
				"data": []
			}`,
		},
		{
			name: "Internal Server Error",
			mockBehavior: func(service *MockBrandService) {
				service.On("List").Return([]domain.Brand{}, errors.New("database error"))
			},
			expectedStatus: http.StatusInternalServerError,
			expectedBody: `{
				"status": 500,
				"message": "database error"
			}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockService := new(MockBrandService)
			tt.mockBehavior(mockService)
			handler := handler.NewBrandHandler(mockService)

			req := httptest.NewRequest(http.MethodGet, "/brands", nil)
			rec := httptest.NewRecorder()

			handler.GetAll(rec, req, nil)

			assert.Equal(t, tt.expectedStatus, rec.Code)
			assert.JSONEq(t, tt.expectedBody, rec.Body.String())
			mockService.AssertExpectations(t)
		})
	}
}

// Voucher Handler Tests
func TestVoucherHandler_Create(t *testing.T) {
	validUntil := time.Date(2024, 12, 31, 23, 59, 59, 0, time.UTC)

	tests := []struct {
		name           string
		requestBody    interface{}
		mockBehavior   func(service *MockVoucherService)
		expectedStatus int
		expectedBody   string
	}{
		{
			name: "Success Create Voucher",
			requestBody: domain.Voucher{
				BrandID:    1,
				Code:       "VOUCHER123",
				Name:       "Test Voucher",
				Points:     50000,
				ValidUntil: validUntil,
			},
			mockBehavior: func(service *MockVoucherService) {
				service.On("Create", mock.Anything).Return(nil)
			},
			expectedStatus: http.StatusCreated,
			expectedBody:   `{"status":201,"message":"Voucher created successfully","data":{"id":0,"brand_id":1,"code":"VOUCHER123","name":"Test Voucher","description":"","points":50000,"valid_until":"2024-12-31T23:59:59Z","created_at":"0001-01-01T00:00:00Z","updated_at":"0001-01-01T00:00:00Z"}}`,
		},
		{
			name: "Invalid Request - Empty Name",
			requestBody: domain.Voucher{
				BrandID:    1,
				Code:       "VOUCHER123",
				Points:     50000,
				ValidUntil: validUntil,
				// Name dikosongkan untuk memicu error
			},
			mockBehavior: func(service *MockVoucherService) {
				service.On("Create", mock.Anything).Return(errors.New("Nama voucher tidak boleh kosong"))
			},
			expectedStatus: http.StatusInternalServerError,
			expectedBody:   `{"status":500,"message":"Nama voucher tidak boleh kosong"}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockService := new(MockVoucherService)
			tt.mockBehavior(mockService)
			handler := handler.NewVoucherHandler(mockService)

			body, _ := json.Marshal(tt.requestBody)
			req := httptest.NewRequest(http.MethodPost, "/voucher", bytes.NewBuffer(body))
			req.Header.Set("Content-Type", "application/json")
			rec := httptest.NewRecorder()

			handler.Create(rec, req, nil)

			assert.Equal(t, tt.expectedStatus, rec.Code)
			assert.JSONEq(t, tt.expectedBody, rec.Body.String())
			mockService.AssertExpectations(t)
		})
	}
}

func TestVoucherHandler_GetByID(t *testing.T) {
	validUntil := time.Date(2024, 12, 31, 23, 59, 59, 0, time.UTC)

	tests := []struct {
		name           string
		voucherID      string
		mockBehavior   func(service *MockVoucherService)
		expectedStatus int
		expectedBody   string
	}{
		{
			name:      "Success Get Voucher",
			voucherID: "1",
			mockBehavior: func(service *MockVoucherService) {
				service.On("GetByID", int64(1)).Return(&domain.Voucher{
					ID:          1,
					BrandID:     1,
					Code:        "VOUCHER123",
					Name:        "Test Voucher",
					Description: "Test Description",
					Points:      50000,
					ValidUntil:  validUntil,
					CreatedAt:   time.Date(2024, 3, 1, 0, 0, 0, 0, time.UTC),
					UpdatedAt:   time.Date(2024, 3, 1, 0, 0, 0, 0, time.UTC),
				}, nil)
			},
			expectedStatus: http.StatusOK,
			expectedBody: `{
				"status": 200,
				"message": "Success",
				"data": {
					"id": 1,
					"brand_id": 1,
					"code": "VOUCHER123",
					"name": "Test Voucher",
					"description": "Test Description",
					"points": 50000,
					"valid_until": "2024-12-31T23:59:59Z",
					"created_at": "2024-03-01T00:00:00Z",
					"updated_at": "2024-03-01T00:00:00Z"
				}
			}`,
		},
		{
			name:      "Voucher Not Found",
			voucherID: "999",
			mockBehavior: func(service *MockVoucherService) {
				service.On("GetByID", int64(999)).Return(nil, nil)
			},
			expectedStatus: http.StatusNotFound,
			expectedBody: `{
				"status": 404,
				"message": "Voucher not found"
			}`,
		},
		{
			name:           "Invalid ID Format",
			voucherID:      "abc",
			mockBehavior:   func(service *MockVoucherService) {},
			expectedStatus: http.StatusBadRequest,
			expectedBody: `{
				"status": 400,
				"message": "Invalid ID"
			}`,
		},
		{
			name:           "Empty ID",
			voucherID:      "",
			mockBehavior:   func(service *MockVoucherService) {},
			expectedStatus: http.StatusBadRequest,
			expectedBody: `{
				"status": 400,
				"message": "ID is required"
			}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockService := new(MockVoucherService)
			tt.mockBehavior(mockService)
			handler := handler.NewVoucherHandler(mockService)

			req := httptest.NewRequest(http.MethodGet, "/voucher/"+tt.voucherID, nil)
			rec := httptest.NewRecorder()
			params := httprouter.Params{httprouter.Param{Key: "id", Value: tt.voucherID}}

			handler.GetByID(rec, req, params)

			assert.Equal(t, tt.expectedStatus, rec.Code)
			assert.JSONEq(t, tt.expectedBody, rec.Body.String())
			mockService.AssertExpectations(t)
		})
	}
}

func TestVoucherHandler_GetByBrandID(t *testing.T) {
	validUntil := time.Date(2024, 12, 31, 23, 59, 59, 0, time.UTC)

	tests := []struct {
		name           string
		brandID        string
		mockBehavior   func(service *MockVoucherService)
		expectedStatus int
		expectedBody   string
	}{
		{
			name:    "Success Get Vouchers By Brand ID",
			brandID: "1",
			mockBehavior: func(service *MockVoucherService) {
				service.On("GetByBrandID", int64(1)).Return([]domain.Voucher{
					{
						ID:          1,
						BrandID:     1,
						Code:        "VOUCHER123",
						Name:        "Test Voucher 1",
						Description: "Test Description 1",
						Points:      50000,
						ValidUntil:  validUntil,
						CreatedAt:   time.Date(2024, 3, 1, 0, 0, 0, 0, time.UTC),
						UpdatedAt:   time.Date(2024, 3, 1, 0, 0, 0, 0, time.UTC),
					},
					{
						ID:          2,
						BrandID:     1,
						Code:        "VOUCHER456",
						Name:        "Test Voucher 2",
						Description: "Test Description 2",
						Points:      75000,
						ValidUntil:  validUntil,
						CreatedAt:   time.Date(2024, 3, 1, 0, 0, 0, 0, time.UTC),
						UpdatedAt:   time.Date(2024, 3, 1, 0, 0, 0, 0, time.UTC),
					},
				}, nil)
			},
			expectedStatus: http.StatusOK,
			expectedBody: `{
				"status": 200,
				"message": "Success",
				"data": [
					{
						"id": 1,
						"brand_id": 1,
						"code": "VOUCHER123",
						"name": "Test Voucher 1",
						"description": "Test Description 1",
						"points": 50000,
						"valid_until": "2024-12-31T23:59:59Z",
						"created_at": "2024-03-01T00:00:00Z",
						"updated_at": "2024-03-01T00:00:00Z"
					},
					{
						"id": 2,
						"brand_id": 1,
						"code": "VOUCHER456",
						"name": "Test Voucher 2",
						"description": "Test Description 2",
						"points": 75000,
						"valid_until": "2024-12-31T23:59:59Z",
						"created_at": "2024-03-01T00:00:00Z",
						"updated_at": "2024-03-01T00:00:00Z"
					}
				]
			}`,
		},
		{
			name:    "Brand Not Found",
			brandID: "999",
			mockBehavior: func(service *MockVoucherService) {
				service.On("GetByBrandID", int64(999)).Return([]domain.Voucher{}, errors.New("brand not found"))
			},
			expectedStatus: http.StatusNotFound,
			expectedBody: `{
				"status": 404,
				"message": "Brand not found"
			}`,
		},
		{
			name:    "No Vouchers Found",
			brandID: "2",
			mockBehavior: func(service *MockVoucherService) {
				service.On("GetByBrandID", int64(2)).Return([]domain.Voucher{}, nil)
			},
			expectedStatus: http.StatusNotFound,
			expectedBody: `{
				"status": 404,
				"message": "Brand tidak memiliki voucher"
			}`,
		},
		{
			name:           "Invalid Brand ID Format",
			brandID:        "abc",
			mockBehavior:   func(service *MockVoucherService) {},
			expectedStatus: http.StatusBadRequest,
			expectedBody: `{
				"status": 400,
				"message": "Invalid brand ID"
			}`,
		},
		{
			name:           "Empty Brand ID",
			brandID:        "",
			mockBehavior:   func(service *MockVoucherService) {},
			expectedStatus: http.StatusBadRequest,
			expectedBody: `{
				"status": 400,
				"message": "Brand ID is required"
			}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockService := new(MockVoucherService)
			tt.mockBehavior(mockService)
			handler := handler.NewVoucherHandler(mockService)

			req := httptest.NewRequest(http.MethodGet, "/brand/"+tt.brandID+"/vouchers", nil)
			rec := httptest.NewRecorder()
			params := httprouter.Params{httprouter.Param{Key: "id", Value: tt.brandID}}

			handler.GetByBrandID(rec, req, params)

			assert.Equal(t, tt.expectedStatus, rec.Code)
			assert.JSONEq(t, tt.expectedBody, rec.Body.String())
			mockService.AssertExpectations(t)
		})
	}
}

func TestTransactionHandler_CreateRedemption(t *testing.T) {
	tests := []struct {
		name           string
		requestBody    interface{}
		mockBehavior   func(service *MockTransactionService)
		expectedStatus int
		expectedBody   string
	}{
		{
			name: "Success Create Transaction",
			requestBody: domain.Transaction{
				CustomerID: 1,
				TotalPoints: 75000,
				Items: []domain.TransactionItem{
					{
						VoucherID:  1,
						PointsUsed: 50000,
					},
					{
						VoucherID:  2,
						PointsUsed: 25000,
					},
				},
			},
			mockBehavior: func(service *MockTransactionService) {
				service.On("CreateRedemption", mock.AnythingOfType("*domain.Transaction")).Return(nil)
			},
			expectedStatus: http.StatusCreated,
			expectedBody: `{
				"status": 201,
				"message": "Redemption created successfully",
				"data": {
					"id": 0,
					"customer_id": 1,
					"total_points": 75000,
					"status": "",
					"items": [
						{
							"id": 0,
							"transaction_id": 0,
							"voucher_id": 1,
							"points_used": 50000,
							"created_at": "0001-01-01T00:00:00Z"
						},
						{
							"id": 0,
							"transaction_id": 0,
							"voucher_id": 2,
							"points_used": 25000,
							"created_at": "0001-01-01T00:00:00Z"
						}
					],
					"created_at": "0001-01-01T00:00:00Z",
					"updated_at": "0001-01-01T00:00:00Z"
				}
			}`,
		},
		{
			name: "Invalid Request - Empty Customer ID",
			requestBody: domain.Transaction{
				TotalPoints: 50000,
				Items: []domain.TransactionItem{
					{
						VoucherID:  1,
						PointsUsed: 50000,
					},
				},
			},
			mockBehavior: func(service *MockTransactionService) {
				service.On("CreateRedemption", mock.AnythingOfType("*domain.Transaction")).Return(errors.New("customer ID is required"))
			},
			expectedStatus: http.StatusInternalServerError,
			expectedBody: `{
				"status": 500,
				"message": "customer ID is required"
			}`,
		},
		{
			name: "Invalid Request - No Items",
			requestBody: domain.Transaction{
				CustomerID:  1,
				TotalPoints: 50000,
				Items:      []domain.TransactionItem{},
			},
			mockBehavior: func(service *MockTransactionService) {
				service.On("CreateRedemption", mock.AnythingOfType("*domain.Transaction")).Return(errors.New("transaction items are required"))
			},
			expectedStatus: http.StatusInternalServerError,
			expectedBody: `{
				"status": 500,
				"message": "transaction items are required"
			}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockService := new(MockTransactionService)
			tt.mockBehavior(mockService)
			handler := handler.NewTransactionHandler(mockService)

			body, _ := json.Marshal(tt.requestBody)
			req := httptest.NewRequest(http.MethodPost, "/transaction/redemption", bytes.NewBuffer(body))
			req.Header.Set("Content-Type", "application/json")
			rec := httptest.NewRecorder()

			handler.CreateRedemption(rec, req, nil)

			assert.Equal(t, tt.expectedStatus, rec.Code)
			assert.JSONEq(t, tt.expectedBody, rec.Body.String())
			mockService.AssertExpectations(t)
		})
	}
}

func TestTransactionHandler_GetByID(t *testing.T) {
	tests := []struct {
		name           string
		transactionID string
		mockBehavior   func(service *MockTransactionService)
		expectedStatus int
		expectedBody   string
	}{
		{
			name:         "Success Get Transaction",
			transactionID: "1",
			mockBehavior: func(service *MockTransactionService) {
				service.On("GetTransactionByID", int64(1)).Return(&domain.Transaction{
					ID:          1,
					CustomerID:  1,
					TotalPoints: 75000,
					Status:      "completed",
					Items: []domain.TransactionItem{
						{
							ID:            1,
							TransactionID: 1,
							VoucherID:     1,
							PointsUsed:    50000,
							CreatedAt:     time.Date(2024, 3, 1, 0, 0, 0, 0, time.UTC),
						},
						{
							ID:            2,
							TransactionID: 1,
							VoucherID:     2,
							PointsUsed:    25000,
							CreatedAt:     time.Date(2024, 3, 1, 0, 0, 0, 0, time.UTC),
						},
					},
					CreatedAt: time.Date(2024, 3, 1, 0, 0, 0, 0, time.UTC),
					UpdatedAt: time.Date(2024, 3, 1, 0, 0, 0, 0, time.UTC),
				}, nil)
			},
			expectedStatus: http.StatusOK,
			expectedBody: `{
				"status": 200,
				"message": "Success",
				"data": {
					"id": 1,
					"customer_id": 1,
					"total_points": 75000,
					"status": "completed",
					"items": [
						{
							"id": 1,
							"transaction_id": 1,
							"voucher_id": 1,
							"points_used": 50000,
							"created_at": "2024-03-01T00:00:00Z"
						},
						{
							"id": 2,
							"transaction_id": 1,
							"voucher_id": 2,
							"points_used": 25000,
							"created_at": "2024-03-01T00:00:00Z"
						}
					],
					"created_at": "2024-03-01T00:00:00Z",
					"updated_at": "2024-03-01T00:00:00Z"
				}
			}`,
		},
		{
			name:         "Transaction Not Found",
			transactionID: "999",
			mockBehavior: func(service *MockTransactionService) {
				service.On("GetTransactionByID", int64(999)).Return(nil, nil)
			},
			expectedStatus: http.StatusNotFound,
			expectedBody: `{
				"status": 404,
				"message": "Transaction not found"
			}`,
		},
		{
			name:         "Invalid ID Format",
			transactionID: "abc",
			mockBehavior: func(service *MockTransactionService) {},
			expectedStatus: http.StatusBadRequest,
			expectedBody: `{
				"status": 400,
				"message": "Invalid transaction ID"
			}`,
		},
		{
			name:         "Empty ID",
			transactionID: "",
			mockBehavior: func(service *MockTransactionService) {},
			expectedStatus: http.StatusBadRequest,
			expectedBody: `{
				"status": 400,
				"message": "Transaction ID is required"
			}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockService := new(MockTransactionService)
			tt.mockBehavior(mockService)
			handler := handler.NewTransactionHandler(mockService)

			req := httptest.NewRequest(http.MethodGet, "/transaction/"+tt.transactionID, nil)
			rec := httptest.NewRecorder()
			params := httprouter.Params{httprouter.Param{Key: "id", Value: tt.transactionID}}

			handler.GetTransactionByID(rec, req, params)

			assert.Equal(t, tt.expectedStatus, rec.Code)
			assert.JSONEq(t, tt.expectedBody, rec.Body.String())
			mockService.AssertExpectations(t)
		})
	}
}
