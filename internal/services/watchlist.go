package services

import (
	"encoding/json"
	"fmt"
	"os"
	"sync"
	"time"

	"muvi-discovery-app/internal/models"
)

// WatchlistService manages user watchlists using local file storage
type WatchlistService struct {
	mu        sync.RWMutex
	watchlist map[string]models.WatchlistItem // key is "type:id"
	filePath  string
}

func NewWatchlistService(filePath string) *WatchlistService {
	ws := &WatchlistService{
		watchlist: make(map[string]models.WatchlistItem),
		filePath:  filePath,
	}
	
	// Load existing data
	ws.loadFromFile()
	
	return ws
}

func (ws *WatchlistService) loadFromFile() {
	ws.mu.Lock()
	defer ws.mu.Unlock()

	data, err := os.ReadFile(ws.filePath)
	if err != nil {
		// File doesn't exist or can't be read, start with empty watchlist
		return
	}

	var items []models.WatchlistItem
	if err := json.Unmarshal(data, &items); err != nil {
		// Invalid JSON, start with empty watchlist
		return
	}

	// Convert slice to map for faster lookups
	for _, item := range items {
		key := fmt.Sprintf("%s:%d", item.Type, item.ID)
		ws.watchlist[key] = item
	}
}

func (ws *WatchlistService) saveToFile() error {
	// Convert map to slice
	items := make([]models.WatchlistItem, 0, len(ws.watchlist))
	for _, item := range ws.watchlist {
		items = append(items, item)
	}

	data, err := json.MarshalIndent(items, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(ws.filePath, data, 0644)
}

func (ws *WatchlistService) AddItem(item models.WatchlistItem) error {
	ws.mu.Lock()
	defer ws.mu.Unlock()

	key := fmt.Sprintf("%s:%d", item.Type, item.ID)
	
	// Check if item already exists
	if _, exists := ws.watchlist[key]; exists {
		return fmt.Errorf("item already in watchlist")
	}

	item.AddedAt = time.Now()
	item.Watched = false
	ws.watchlist[key] = item

	return ws.saveToFile()
}

func (ws *WatchlistService) RemoveItem(itemType string, id int) error {
	ws.mu.Lock()
	defer ws.mu.Unlock()

	key := fmt.Sprintf("%s:%d", itemType, id)
	
	if _, exists := ws.watchlist[key]; !exists {
		return fmt.Errorf("item not found in watchlist")
	}

	delete(ws.watchlist, key)
	return ws.saveToFile()
}

func (ws *WatchlistService) ToggleWatched(itemType string, id int) error {
	ws.mu.Lock()
	defer ws.mu.Unlock()

	key := fmt.Sprintf("%s:%d", itemType, id)
	
	item, exists := ws.watchlist[key]
	if !exists {
		return fmt.Errorf("item not found in watchlist")
	}

	item.Watched = !item.Watched
	if item.Watched {
		now := time.Now()
		item.WatchedAt = &now
	} else {
		item.WatchedAt = nil
	}

	ws.watchlist[key] = item
	return ws.saveToFile()
}

func (ws *WatchlistService) GetAllItems() []models.WatchlistItem {
	ws.mu.RLock()
	defer ws.mu.RUnlock()

	items := make([]models.WatchlistItem, 0, len(ws.watchlist))
	for _, item := range ws.watchlist {
		items = append(items, item)
	}

	return items
}

func (ws *WatchlistService) GetItem(itemType string, id int) (*models.WatchlistItem, bool) {
	ws.mu.RLock()
	defer ws.mu.RUnlock()

	key := fmt.Sprintf("%s:%d", itemType, id)
	item, exists := ws.watchlist[key]
	if !exists {
		return nil, false
	}

	return &item, true
}

func (ws *WatchlistService) IsInWatchlist(itemType string, id int) bool {
	ws.mu.RLock()
	defer ws.mu.RUnlock()

	key := fmt.Sprintf("%s:%d", itemType, id)
	_, exists := ws.watchlist[key]
	return exists
}

func (ws *WatchlistService) GetWatchedItems() []models.WatchlistItem {
	ws.mu.RLock()
	defer ws.mu.RUnlock()

	var watched []models.WatchlistItem
	for _, item := range ws.watchlist {
		if item.Watched {
			watched = append(watched, item)
		}
	}

	return watched
}

func (ws *WatchlistService) GetUnwatchedItems() []models.WatchlistItem {
	ws.mu.RLock()
	defer ws.mu.RUnlock()

	var unwatched []models.WatchlistItem
	for _, item := range ws.watchlist {
		if !item.Watched {
			unwatched = append(unwatched, item)
		}
	}

	return unwatched
}

func (ws *WatchlistService) GetItemCount() int {
	ws.mu.RLock()
	defer ws.mu.RUnlock()

	return len(ws.watchlist)
}