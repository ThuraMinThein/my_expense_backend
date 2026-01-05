# My Expense Backend

A RESTful API backend for the "My Expense" mobile expense tracking application built with Go and Gin framework.

## Features

- **User Authentication** (Email/Password + Google OAuth)
- **Expense Management** (Create, Read, Delete)
- **Analytics Dashboard** (Daily, Weekly, Monthly reports)
- **JWT-based Security**
- **PostgreSQL Database with GORM**
- **Clean Architecture Pattern**

## API Endpoints

### Authentication
- `POST /auth/sign-up` - User registration
- `POST /auth/login` - User login
- `POST /auth/refresh` - Refresh JWT token
- `GET /auth/google` - Get Google OAuth URL
- `GET /auth/google/callback` - Google OAuth callback
- `POST /auth/google/token` - Google login with token
- `POST /auth/logout` - User logout (protected)

### Expenses
- `POST /expenses` - Create new expense (protected)
- `GET /expenses?from=YYYY-MM-DD&to=YYYY-MM-DD` - List expenses (protected, defaults to last 30 days)
- `DELETE /expenses/:id` - Delete expense (protected)

### Analytics
- `GET /analytics/daily?date=YYYY-MM-DD` - Daily usage statistics (protected)
- `GET /analytics/weekly?week=YYYY-WWW` - Weekly usage with daily breakdown (protected)
- `GET /analytics/monthly?month=YYYY-MM` - Monthly usage by category (protected)

## Architecture

The project follows clean architecture principles:

```
internal/app/
├── models/          # Data models and database schemas
├── repositories/    # Data access layer (database operations)
├── services/        # Business logic layer
├── handlers/        # HTTP request handlers (controllers)
├── routes/          # Route definitions and middleware
└── helper/          # Utility functions
```

## Requirements

- Go 1.24+
- PostgreSQL 12+
- Redis (for caching, optional)

## Getting Started

1. **Clone the repository**
   ```bash
   git clone <repository-url>
   cd my_expense_backend
   ```

2. **Install dependencies**
   ```bash
   go mod tidy
   ```

3. **Set up environment variables**
   ```bash
   cp .env.example .env
   # Edit .env with your database credentials and JWT secrets
   ```

4. **Run database migrations**
   ```bash
   # Make sure PostgreSQL is running and accessible
   go run cmd/migrate/main.go
   ```

5. **Start the server**
   ```bash
   go run cmd/server/main.go
   ```

The API will be available at `http://localhost:8080`

## API Usage Examples

### Create Expense
```bash
curl -X POST http://localhost:8080/expenses \
  -H "Authorization: Bearer <your-jwt-token>" \
  -H "Content-Type: application/json" \
  -d '{
    "amount": 5000,
    "category": "Food",
    "note": "Lunch",
    "expense_date": "2026-01-05"
  }'
```

### List Expenses
```bash
curl -X GET "http://localhost:8080/expenses?from=2026-01-01&to=2026-01-31" \
  -H "Authorization: Bearer <your-jwt-token>"
```

### Get Daily Analytics
```bash
curl -X GET "http://localhost:8080/analytics/daily?date=2026-01-05" \
  -H "Authorization: Bearer <your-jwt-token>"
```

### Get Monthly Analytics
```bash
curl -X GET "http://localhost:8080/analytics/monthly?month=2026-01" \
  -H "Authorization: Bearer <your-jwt-token>"
```

## Security

- JWT-based authentication with refresh tokens
- Password hashing with bcrypt
- User-scoped data access (users can only access their own data)
- Input validation and sanitization
- CORS configuration for cross-origin requests

## Testing

```bash
# Run all tests
go test ./...

# Run tests with coverage
go test -cover ./...

# Run specific test
go test ./internal/app/services/expense_service_test.go
```

### Environment Variables
- `DB_HOST`: PostgreSQL host
- `DB_PORT`: PostgreSQL port
- `DB_USER`: PostgreSQL username
- `DB_PASSWORD`: PostgreSQL password
- `DB_NAME`: Database name
- `JWT_SECRET`: JWT signing secret
- `GOOGLE_CLIENT_ID`: Google OAuth client ID
- `GOOGLE_CLIENT_SECRET`: Google OAuth client secret

## Contributing

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add some amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## Support

For support and questions, please open an issue in the GitHub repository.