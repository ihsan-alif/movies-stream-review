# Movies Stream Review

A full-stack application for stream movies (YouTube trailer) and sharing reviews. Built with React/Vite, Go and MongoDB.

## Features

- User authentication and registration
- Write, view movies and reviews
- Personalized movie recommendations

## Tech Stack

### Frontend
- **React** with Vite
- **Axios** for API communication
- **React Router** for navigation
- Context API for state management

### Backend
- **Go** with gin-gonic framework
- **JWT** for authentication
- **RESTful API** design

## Setup & Installation

### Frontend
```bash
cd Client/movies-stream-review-client
npm install
npm run dev
```

### Backend
```bash
cd Server/MoviesStreamReview
go mod download
go run main.go
```

## API Endpoints

- `POST /users/register` - User registration
- `POST /users/login` - User login
- `GET /movies` - Get all movies
- `GET /movies/:imdb_id` - Get specific movie (protected)
- `PATCH /movies/updatereview/:imdb_id` - Add review (protected)
- `GET /movies/recommended` - Get recommendations (protected)

## Authentication

The application uses JWT tokens for authentication. Protected endpoints require a valid token in the Authorization header.
