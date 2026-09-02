# SocialNetwork 🚀

A minimalist, modern social media platform for nerdy people. Built with **Go** on the backend and **React + Vite** on the frontend.

## 🌟 Features

- **Authentication**: JWT-based login and registration system with email activation.
- **Home Feed**: View posts from yourself and users you follow.
- **Explore Feed**: Discover and follow new users whose content you might like.
- **Post Management**: Create new posts seamlessly with a responsive drawer interface.
- **Social Graph**: Follow and unfollow users instantly without page reloads.
- **Modern UI**: Clean, minimalist design with Dark Mode support built using `shadcn/ui` and Tailwind CSS.

---

## 🛠️ Tech Stack

### Backend
- **Go** (Golang)
- **Chi Router** for REST API routing
- **PostgreSQL** for relational database
- **Redis** for rate-limiting and caching
- **Swagger** for API documentation
- **Resend** for transactional emails

### Frontend
- **React** (Bootstrapped with Vite)
- **TypeScript**
- **Tailwind CSS**
- **shadcn/ui** for accessible and customizable UI components
- **SWR** for data fetching and state mutation

---

## 🚀 Getting Started

### Prerequisites
- Go 1.20+
- Node.js & npm
- PostgreSQL
- Redis
- Docker (optional, for running via docker-compose)

### Backend Setup

1. Configure your environment variables in `.envrc` or export them.
2. Ensure PostgreSQL and Redis are running.
3. Install dependencies and run the server (with `air` for hot-reload if configured):
   ```bash
   make run
   # OR
   go run cmd/api/*.go
   ```
4. API Documentation will be available at: `http://localhost:8080/swagger/index.html`

### Frontend Setup

1. Navigate to the frontend directory:
   ```bash
   cd web
   ```
2. Install dependencies:
   ```bash
   npm install
   ```
3. Start the Vite development server:
   ```bash
   npm run dev
   ```
4. Access the web app at `http://localhost:5173`.

---

## 🧪 Testing

### Testing Rate Limiter
You can test the backend rate limiter using `autocannon`:
```bash
npx autocannon -r 22 -d 1 -c 1 --renderStatusCodes http://localhost:8080/v1/health
```

---

## 📄 License
This project is open-source and available under the MIT License.
