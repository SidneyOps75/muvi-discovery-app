// Main JavaScript functionality for Muvi Discovery App

// Trailer functionality
function openTrailerModal(videoKey, videoTitle) {
    console.log('Opening trailer modal with key:', videoKey, 'title:', videoTitle);
    
    const modal = document.getElementById('trailerModal');
    const iframe = document.getElementById('trailerFrame');
    const title = document.getElementById('trailerTitle');
    
    console.log('Modal elements found:', {
        modal: !!modal,
        iframe: !!iframe,
        title: !!title
    });
    
    if (modal && iframe && title) {
        // Set the YouTube embed URL
        const embedUrl = `https://www.youtube.com/embed/${videoKey}?autoplay=1&rel=0`;
        console.log('Setting iframe src to:', embedUrl);
        
        iframe.src = embedUrl;
        title.textContent = videoTitle || 'Trailer';
        
        // Show the modal
        modal.style.display = 'block';
        document.body.style.overflow = 'hidden'; // Prevent background scrolling
        
        // Close modal when clicking outside
        modal.onclick = function(event) {
            if (event.target === modal) {
                closeTrailerModal();
            }
        };
        
        // Close modal with Escape key
        document.addEventListener('keydown', handleTrailerEscape);
        
        console.log('Trailer modal opened successfully');
    } else {
        console.error('Could not find required modal elements');
    }
}

function closeTrailerModal() {
    console.log('Closing trailer modal');
    
    const modal = document.getElementById('trailerModal');
    const iframe = document.getElementById('trailerFrame');
    
    if (modal && iframe) {
        // Hide the modal
        modal.style.display = 'none';
        document.body.style.overflow = 'auto'; // Restore scrolling
        
        // Stop the video by clearing the src
        iframe.src = '';
        
        // Remove escape key listener
        document.removeEventListener('keydown', handleTrailerEscape);
        
        console.log('Trailer modal closed successfully');
    } else {
        console.error('Could not find modal elements to close');
    }
}

function handleTrailerEscape(event) {
    if (event.key === 'Escape') {
        closeTrailerModal();
    }
}

// Watchlist functionality
function addToWatchlist(id, type, title, posterPath, releaseDate, voteAverage, buttonElement) {
    const item = {
        id: id,
        type: type,
        title: title,
        poster_path: posterPath,
        release_date: releaseDate,
        vote_average: voteAverage
    };

    fetch('/api/watchlist', {
        method: 'POST',
        headers: {
            'Content-Type': 'application/json',
        },
        body: JSON.stringify(item)
    })
    .then(response => response.json())
    .then(data => {
        if (data.status === 'success') {
            showNotification('Added to watchlist!', 'success');
            // Update button if provided
            if (buttonElement) {
                buttonElement.textContent = 'Remove from Watchlist';
                buttonElement.className = 'btn btn-secondary';
                buttonElement.onclick = () => removeFromWatchlist(id, type, buttonElement);
            }
            
            // Update watchlist count
            updateWatchlistCount();
        }
    })
    .catch(error => {
        console.error('Error:', error);
        showNotification('Failed to add to watchlist', 'error');
    });
}

function removeFromWatchlist(id, type, buttonElement) {
    fetch(`/api/watchlist/${id}?type=${type}`, {
        method: 'DELETE'
    })
    .then(response => response.json())
    .then(data => {
        if (data.status === 'success') {
            showNotification('Removed from watchlist!', 'success');
            
            // If on watchlist page, remove the item
            if (window.location.pathname === '/watchlist') {
                location.reload();
            } else if (buttonElement) {
                // Update button
                buttonElement.textContent = 'Add to Watchlist';
                buttonElement.className = 'btn btn-primary';
                buttonElement.onclick = () => addToWatchlist(id, type, '', '', '', 0, buttonElement);
            }
            
            // Update watchlist count
            updateWatchlistCount();
        }
    })
    .catch(error => {
        console.error('Error:', error);
        showNotification('Failed to remove from watchlist', 'error');
    });
}

function toggleWatched(id, type) {
    fetch(`/api/watchlist/${id}/toggle?type=${type}`, {
        method: 'PUT'
    })
    .then(response => response.json())
    .then(data => {
        if (data.status === 'success') {
            showNotification('Updated watch status!', 'success');
            // Reload the page to reflect changes
            location.reload();
        }
    })
    .catch(error => {
        console.error('Error:', error);
        showNotification('Failed to update watch status', 'error');
    });
}

function updateWatchlistCount() {
    // This would typically fetch the current count from the server
    // For now, we'll just reload the page to update the count
    setTimeout(() => {
        location.reload();
    }, 1000);
}

// Notification system
function showNotification(message, type = 'info') {
    // Remove existing notifications
    const existingNotifications = document.querySelectorAll('.notification');
    existingNotifications.forEach(notification => notification.remove());

    // Create notification element
    const notification = document.createElement('div');
    notification.className = `notification notification-${type}`;
    notification.textContent = message;

    // Style the notification
    notification.style.cssText = `
        position: fixed;
        top: 20px;
        right: 20px;
        padding: 1rem 1.5rem;
        border-radius: 0.5rem;
        color: white;
        font-weight: 500;
        z-index: 1000;
        animation: slideIn 0.3s ease-out;
    `;

    // Set background color based on type
    switch (type) {
        case 'success':
            notification.style.backgroundColor = '#10b981';
            break;
        case 'error':
            notification.style.backgroundColor = '#ef4444';
            break;
        case 'warning':
            notification.style.backgroundColor = '#f59e0b';
            break;
        default:
            notification.style.backgroundColor = '#3b82f6';
    }

    // Add to page
    document.body.appendChild(notification);

    // Remove after 3 seconds
    setTimeout(() => {
        notification.style.animation = 'slideOut 0.3s ease-in';
        setTimeout(() => {
            notification.remove();
        }, 300);
    }, 3000);
}

// Add CSS animations for notifications
const style = document.createElement('style');
style.textContent = `
    @keyframes slideIn {
        from {
            transform: translateX(100%);
            opacity: 0;
        }
        to {
            transform: translateX(0);
            opacity: 1;
        }
    }
    
    @keyframes slideOut {
        from {
            transform: translateX(0);
            opacity: 1;
        }
        to {
            transform: translateX(100%);
            opacity: 0;
        }
    }
`;
document.head.appendChild(style);

// Search functionality
function initializeSearch() {
    const searchForm = document.querySelector('.search-form');
    const searchInput = document.querySelector('.search-input');
    
    if (searchForm && searchInput) {
        // Add debounced search suggestions (optional enhancement)
        let searchTimeout;
        
        searchInput.addEventListener('input', function() {
            clearTimeout(searchTimeout);
            searchTimeout = setTimeout(() => {
                // Could implement search suggestions here
            }, 300);
        });
    }
}

// Mobile navigation toggle
function initializeMobileNav() {
    const navToggle = document.querySelector('.nav-toggle');
    const navMenu = document.querySelector('.nav-menu');
    
    if (navToggle && navMenu) {
        navToggle.addEventListener('click', function() {
            navMenu.classList.toggle('active');
            navToggle.classList.toggle('active');
        });
    }
}

// Lazy loading for images
function initializeLazyLoading() {
    const images = document.querySelectorAll('img[data-src]');
    
    if ('IntersectionObserver' in window) {
        const imageObserver = new IntersectionObserver((entries, observer) => {
            entries.forEach(entry => {
                if (entry.isIntersecting) {
                    const img = entry.target;
                    img.src = img.dataset.src;
                    img.classList.remove('lazy');
                    imageObserver.unobserve(img);
                }
            });
        });
        
        images.forEach(img => imageObserver.observe(img));
    } else {
        // Fallback for browsers without IntersectionObserver
        images.forEach(img => {
            img.src = img.dataset.src;
        });
    }
}

// Infinite scroll for content pages
function initializeInfiniteScroll() {
    const loadMoreButton = document.querySelector('.load-more');
    
    if (loadMoreButton) {
        let loading = false;
        
        window.addEventListener('scroll', () => {
            if (loading) return;
            
            const { scrollTop, scrollHeight, clientHeight } = document.documentElement;
            
            if (scrollTop + clientHeight >= scrollHeight - 1000) {
                loading = true;
                loadMoreContent();
            }
        });
    }
}

function loadMoreContent() {
    // This would be implemented based on the current page
    // For now, it's a placeholder
    console.log('Loading more content...');
}

// Theme toggle functionality (optional enhancement)
function initializeThemeToggle() {
    const themeToggle = document.querySelector('.theme-toggle');
    
    if (themeToggle) {
        themeToggle.addEventListener('click', function() {
            document.body.classList.toggle('dark-theme');
            
            // Save preference to localStorage
            const isDark = document.body.classList.contains('dark-theme');
            localStorage.setItem('theme', isDark ? 'dark' : 'light');
        });
        
        // Load saved theme
        const savedTheme = localStorage.getItem('theme');
        if (savedTheme === 'dark') {
            document.body.classList.add('dark-theme');
        }
    }
}

// Keyboard shortcuts
function initializeKeyboardShortcuts() {
    document.addEventListener('keydown', function(e) {
        // Focus search with '/' key
        if (e.key === '/' && !e.target.matches('input, textarea')) {
            e.preventDefault();
            const searchInput = document.querySelector('.search-input');
            if (searchInput) {
                searchInput.focus();
            }
        }
        
        // Escape to clear search
        if (e.key === 'Escape') {
            const searchInput = document.querySelector('.search-input');
            if (searchInput && document.activeElement === searchInput) {
                searchInput.blur();
            }
        }
    });
}

// Form validation
function initializeFormValidation() {
    const forms = document.querySelectorAll('form');
    
    forms.forEach(form => {
        form.addEventListener('submit', function(e) {
            const requiredFields = form.querySelectorAll('[required]');
            let isValid = true;
            
            requiredFields.forEach(field => {
                if (!field.value.trim()) {
                    isValid = false;
                    field.classList.add('error');
                    
                    // Remove error class on input
                    field.addEventListener('input', function() {
                        this.classList.remove('error');
                    }, { once: true });
                }
            });
            
            if (!isValid) {
                e.preventDefault();
                showNotification('Please fill in all required fields', 'error');
            }
        });
    });
}

// Smooth scrolling for anchor links
function initializeSmoothScrolling() {
    const anchorLinks = document.querySelectorAll('a[href^="#"]');
    
    anchorLinks.forEach(link => {
        link.addEventListener('click', function(e) {
            e.preventDefault();
            
            const targetId = this.getAttribute('href').substring(1);
            const targetElement = document.getElementById(targetId);
            
            if (targetElement) {
                targetElement.scrollIntoView({
                    behavior: 'smooth',
                    block: 'start'
                });
            }
        });
    });
}

// Global error handler for JavaScript errors
window.addEventListener('error', function(e) {
    console.error('JavaScript Error:', e.error);
    showNotification('An error occurred. Please try again.', 'error');
});

// Initialize all functionality when DOM is loaded
document.addEventListener('DOMContentLoaded', function() {
    try {
        initializeSearch();
        initializeMobileNav();
        initializeLazyLoading();
        initializeInfiniteScroll();
        initializeThemeToggle();
        initializeKeyboardShortcuts();
        initializeFormValidation();
        initializeSmoothScrolling();
        
        // Add loading states to buttons
        const buttons = document.querySelectorAll('.btn');
        buttons.forEach(button => {
            button.addEventListener('click', function() {
                if (this.type === 'submit' || this.onclick) {
                    this.classList.add('loading');
                    setTimeout(() => {
                        this.classList.remove('loading');
                    }, 2000);
                }
            });
        });
    } catch (error) {
        console.error('Initialization error:', error);
        showNotification('Failed to initialize some features', 'warning');
    }
});

// Service Worker registration removed to prevent errors
// Uncomment and create sw.js file if PWA capabilities are needed
// if ('serviceWorker' in navigator) {
//     window.addEventListener('load', function() {
//         navigator.serviceWorker.register('/sw.js')
//             .then(function(registration) {
//                 console.log('ServiceWorker registration successful');
//             })
//             .catch(function(err) {
//                 console.log('ServiceWorker registration failed');
//             });
//     });
// }

// Error handling for images
document.addEventListener('error', function(e) {
    if (e.target.tagName === 'IMG') {
        e.target.src = '/static/images/placeholder.jpg';
    }
}, true);

// Performance monitoring
function logPerformance() {
    if ('performance' in window) {
        window.addEventListener('load', function() {
            setTimeout(() => {
                const perfData = performance.getEntriesByType('navigation')[0];
                console.log('Page load time:', perfData.loadEventEnd - perfData.loadEventStart, 'ms');
            }, 0);
        });
    }
}

logPerformance();