# BoWatt Assignment - RAG Pipeline

A RAG pipeline that allows users to query documents using AI. Built with Go backend and Next.js frontend.

## Tech Stack

- **Backend**: Go 1.24, Gin framework
- **Frontend**: Next.js 15, React, TailwindCSS
- **AI Services**: Hugging Face API
- **Containerization**: Docker

## Project Setup

### Prerequisites

- Docker and Docker Compose
- Node.js 20.x or later
- Go 1.24 or later
- Hugging Face API key

### Getting Started

1. Clone the repository:
```bash
git clone https://github.com/enjoycod1ng/BoWatt.git
cd BoWatt
```

2. Set up environment variables:
```bash
cp .env.example .env
```

3. Add your Hugging Face API key to `.env`:
```properties
HUGGINGFACE_API_KEY=your_api_key
```

4. Build and run with Docker:
```bash
docker-compose up --build
```

The application will be available at:
- Frontend: http://localhost:3000
- Backend: http://localhost:8080

### Local Development

#### Backend
```bash
cd backend
go mod download
go run main.go
```

#### Frontend
```bash
cd frontend
npm install
npm run dev
```

## Architecture and Design Decisions

### Backend Architecture

```plaintext
backend/
├── api/
│   └── handlers/    # Request handlers
├── models/          # Data models
├── services/        # Business logic
└── utils/           # Helper functions
```

### Frontend Architecture

```plaintext
frontend/
├── app/
│   ├── components/  # Reusable UI components
│   └── services/    # API integration
└── public/          # Static assets
```

### Future Improvements (Not Implemented)

1. Document versioning
2. Advanced query filtering
3. Response streaming
4. Automated testing suite

## License

[MIT License](LICENSE)
