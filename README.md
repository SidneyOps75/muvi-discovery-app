# Muvi Discovery App

[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)
[![Go Version](https://img.shields.io/badge/Go-1.19+-blue.svg)](https://golang.org)
[![Build Status](https://img.shields.io/badge/build-passing-brightgreen.svg)]()

A comprehensive entertainment discovery platform built with Go, HTML, and CSS. Discover movies and TV shows, manage your watchlist, and explore trending content with data from TMDB and OMDB APIs.

##  Features

### Core Features
- **Search & Discovery**: Real-time search for movies and TV shows
- **Detailed Information**: Comprehensive details including cast, ratings, plot, and more
- **Personal Watchlist**: Add/remove titles, mark as watched, and manage your collection
- **Trending Content**: Explore popular and trending movies and TV shows
- **Genre Filtering**: Browse content by categories and apply advanced filters
- **Multi-Source Ratings**: Integrated ratings from IMDB, Rotten Tomatoes, and TMDB
- **Responsive Design**: Optimized for both mobile and desktop viewing

### Technical Features
- **Go Backend**: Fast and efficient server built with Go
- **Clean Architecture**: Well-organized code structure with separation of concerns
- **File-based Storage**: Simple JSON file storage for watchlist data
- **API Integration**: Robust integration with TMDB and OMDB APIs
- **Error Handling**: Comprehensive error handling and user feedback
- **Rate Limiting**: Built-in rate limiting for API calls
- **Responsive UI**: Mobile-first responsive design

##  Tech Stack

### Backend
- **Go 1.19+** - Main programming language
- **Gorilla Mux** - HTTP router and URL matcher
- **Standard Library** - HTTP server, JSON handling, file I/O

### Frontend
- **HTML5** - Semantic markup
- **CSS3** - Modern styling with Flexbox and Grid
- **Vanilla JavaScript** - Interactive functionality
- **Google Fonts** - Inter font family

### APIs
- **TMDB API** - Movie/TV data, images, and trending content
- **OMDB API** - Additional ratings and plot information

##  Project Structure

```
muvi-discovery-app/
├── cmd/
│   └── main.go                 # Application entry point
├── internal/
│   ├── handlers/
│   │   └── handlers.go         # HTTP handlers
│   ├── models/
│   │   └── movie.go           # Data models and types
│   └── services/
│       ├── tmdb.go            # TMDB API service
│       ├── omdb.go            # OMDB API service
│       └── watchlist.go       # Watchlist management
├── web/
│   ├── static/
│   │   ├── css/
│   │   │   └── style.css      # Main stylesheet
│   │   ├── js/
│   │   │   └── main.js        # JavaScript functionality
│   │   └── images/
│   │       └── placeholder.jpg # Placeholder image
│   └── templates/
│       ├── base.html          # Base template
│       ├── home.html          # Homepage
│       ├── movies.html        # Movies listing
│       ├── movie_details.html # Movie details
│       ├── tv_shows.html      # TV shows listing
│       ├── tv_details.html    # TV show details
│       ├── search.html        # Search page
│       ├── discover.html      # Discovery page
│       └── watchlist.html     # Watchlist page
├── data/
│   └── watchlist.json         # User watchlist data
├── configs/                   # Configuration files
├── .env.example              # Environment variables example
├── .env                      # Environment variables
├── go.mod                    # Go module file
├── go.sum                    # Go dependencies
└── README.md                 # This file
```

##  Getting Started

### Prerequisites
- Go 1.19 or higher
- TMDB API key
- OMDB API key

### Installation

1. **Clone the repository**
   ```bash
   git clone https://github.com/SidneyOps75/muvi-discovery-app.git
   cd muvi-discovery-app
   ```

2. **Install dependencies**
   ```bash
   go mod tidy
   ```

3. **Set up environment variables**
   ```bash
   cp .env.example .env
   ```
   
   Edit `.env` and add your API keys:
   ```env
   TMDB_API_KEY=your_tmdb_api_key_here
   OMDB_API_KEY=your_omdb_api_key_here
   PORT=8080
   ```

4. **Create data directory**
   ```bash
   mkdir -p data
   ```

5. **Run the application**
   ```bash
   go run cmd/main.go
   ```

6. **Open your browser**
   Navigate to `http://localhost:3000`

### Getting API Keys

#### TMDB API Key
1. Visit [TMDB](https://www.themoviedb.org/)
2. Create an account and verify your email
3. Go to Settings > API
4. Request an API key (choose "Developer" option)
5. Fill out the application form
6. Copy your API key to the `.env` file

#### OMDB API Key
1. Visit [OMDB API](http://www.omdbapi.com/)
2. Click "API Key" tab
3. Choose a plan (free tier available)
4. Register and verify your email
5. Copy your API key to the `.env` file

## 📱 Usage

### Navigation
- **Home**: Trending movies and TV shows
- **Movies**: Browse and filter movies by category
- **TV Shows**: Browse and filter TV series
- **Search**: Search across all content with filters
- **Discover**: Advanced filtering and discovery tools
- **Watchlist**: Manage your saved content

### Features Guide

#### Search
- Use the search bar in the navigation
- Filter by movies or TV shows
- Browse results with pagination

#### Watchlist Management
- Click "Add to Watchlist" on any movie/show detail page
- Mark items as watched from the watchlist page
- Remove items you're no longer interested in
- Filter watchlist by watched/unwatched status

#### Discovery
- Use the Discover page for advanced filtering
- Filter by genre, year, rating, and more
- Sort results by popularity, rating, or release date

##  Development

### Building for Production
```bash
go build -o muvi-discovery-app cmd/main.go
```

### Running with Custom Port
```bash
PORT=3000 go run cmd/main.go
```

### Code Structure

#### Handlers
- Handle HTTP requests and responses
- Render templates with data
- Manage API endpoints

#### Services
- **TMDB Service**: Handles all TMDB API interactions
- **OMDB Service**: Handles OMDB API interactions
- **Watchlist Service**: Manages user watchlist data

#### Models
- Define data structures for movies, TV shows, and API responses
- Type-safe data handling

### Adding New Features

1. **Add new routes** in `cmd/main.go`
2. **Create handlers** in `internal/handlers/handlers.go`
3. **Add templates** in `web/templates/`
4. **Update CSS** in `web/static/css/style.css`
5. **Add JavaScript** in `web/static/js/main.js`

##  Deployment

### Docker Deployment
Create a `Dockerfile`:
```dockerfile
FROM golang:1.19-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -o main cmd/main.go

FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=builder /app/main .
COPY --from=builder /app/web ./web
COPY --from=builder /app/.env .
EXPOSE 8080
CMD ["./main"]
```

Build and run:
```bash
docker build -t muvi-discovery-app .
docker run -p 8080:8080 muvi-discovery-app
```

### Traditional Deployment
1. Build the application: `go build -o muvi-discovery-app cmd/main.go`
2. Copy the binary, `web/` directory, and `.env` file to your server
3. Run the binary: `./muvi-discovery-app`

### Environment Variables for Production
```env
TMDB_API_KEY=your_production_tmdb_key
OMDB_API_KEY=your_production_omdb_key
PORT=8080
```

##  Contributing

We welcome contributions to the Muvi Discovery App! By contributing, you agree that your contributions will be licensed under the same MIT License that covers the project.

### Development Workflow
1. Fork the repository
2. Create a feature branch: `git checkout -b feature/your-feature-name`
3. Make your changes
4. Test thoroughly
5. Commit with descriptive messages
6. Push to your fork
7. Create a pull request

### Code Style
- Follow Go conventions and best practices
- Use `gofmt` for code formatting
- Write clear, descriptive comments
- Keep functions focused and small


### Contribution Guidelines
- Ensure your code follows the existing style and patterns
- Add tests for new functionality
- Update documentation as needed
- Be respectful and constructive in discussions
- All contributions must be your own work or properly attributed

##  Performance

### Optimizations
- **Efficient API Calls**: Built-in rate limiting and caching
- **Responsive Images**: Optimized image loading
- **Minimal JavaScript**: Vanilla JS for better performance
- **Clean CSS**: Organized and efficient stylesheets

### Monitoring
- Server logs for debugging
- Error handling with user-friendly messages
- Performance metrics in browser console

##  Security

- **Environment Variables**: Secure API key management
- **Input Validation**: Server-side validation for all inputs
- **Error Handling**: Graceful error handling without exposing internals
- **HTTPS Ready**: Designed to work with HTTPS in production

##  License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

### What does this mean?

The MIT License is a permissive license that allows you to:

- ✅ **Use** the software for any purpose
- ✅ **Copy** and distribute the software
- ✅ **Modify** the software
- ✅ **Merge** it with other projects
- ✅ **Publish** and distribute modified versions
- ✅ **Sublicense** the software
- ✅ **Sell** copies of the software

**Requirements:**
- Include the original copyright notice and license text in any copy of the software
- The software is provided "as is" without warranty

For more information about the MIT License, visit [choosealicense.com/licenses/mit](https://choosealicense.com/licenses/mit/).

##  Acknowledgments

- **TMDB**: For providing comprehensive movie and TV data
- **OMDB**: For additional ratings and plot information
- **Go Community**: For the excellent standard library and ecosystem
- **Gorilla Mux**: For the robust HTTP router

## 📞 Support

If you encounter any issues or have questions:

1. Check the [Issues](../../issues) page for existing problems
2. Create a new issue with detailed information
3. Include steps to reproduce any bugs
4. Provide your environment details (OS, Go version, etc.)

---

**Happy coding! 🎬✨**