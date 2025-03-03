package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"sync"
	"time"

	"github.com/gorilla/mux"
)

// User modeli
type User struct {
	ID        int       `json:"id"`
	Username  string    `json:"username"`
	Email     string    `json:"email"`
	FullName  string    `json:"fullName"`
	CreatedAt time.Time `json:"createdAt"`
	IsActive  bool      `json:"isActive"`
}

// UserRepository veritabanı işlemleri için arayüz
type UserRepository interface {
	GetAll() ([]User, error)
	GetByID(id int) (*User, error)
	Create(user User) (*User, error)
	Update(user User) (*User, error)
	Delete(id int) error
}

// InMemoryUserRepository bellek içi veri depolama
type InMemoryUserRepository struct {
	users  []User
	nextID int
	mu     sync.RWMutex
}

// NewInMemoryUserRepository yeni bir repository oluşturur
func NewInMemoryUserRepository() *InMemoryUserRepository {
	return &InMemoryUserRepository{
		users: []User{
			{ID: 1, Username: "ahmet123", Email: "ahmet@ornek.com", FullName: "Ahmet Yılmaz", CreatedAt: time.Now().Add(-45 * 24 * time.Hour), IsActive: true},
			{ID: 2, Username: "mehmetk", Email: "mehmet@ornek.com", FullName: "Mehmet Kaya", CreatedAt: time.Now().Add(-30 * 24 * time.Hour), IsActive: true},
			{ID: 3, Username: "ayses", Email: "ayse@ornek.com", FullName: "Ayşe Şahin", CreatedAt: time.Now().Add(-10 * 24 * time.Hour), IsActive: true},
			{ID: 4, Username: "fatmad", Email: "fatma@ornek.com", FullName: "Fatma Demir", CreatedAt: time.Now().Add(-5 * 24 * time.Hour), IsActive: true},
			{ID: 5, Username: "mustafar", Email: "mustafa@ornek.com", FullName: "Mustafa Reis", CreatedAt: time.Now().Add(-2 * 24 * time.Hour), IsActive: true},
		},
		nextID: 6,
	}
}

// GetAll tüm kullanıcıları döndürür
func (r *InMemoryUserRepository) GetAll() ([]User, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return r.users, nil
}

// GetByID ID'ye göre kullanıcı döndürür
func (r *InMemoryUserRepository) GetByID(id int) (*User, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	for _, user := range r.users {
		if user.ID == id {
			return &user, nil
		}
	}
	return nil, fmt.Errorf("kullanıcı bulunamadı (ID: %d)", id)
}

// Create yeni bir kullanıcı oluşturur
func (r *InMemoryUserRepository) Create(user User) (*User, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	user.ID = r.nextID
	r.nextID++
	user.CreatedAt = time.Now()

	r.users = append(r.users, user)
	return &user, nil
}

// Update kullanıcı bilgilerini günceller
func (r *InMemoryUserRepository) Update(user User) (*User, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	for i, u := range r.users {
		if u.ID == user.ID {
			// CreatedAt değerini korumak
			user.CreatedAt = r.users[i].CreatedAt
			r.users[i] = user
			return &r.users[i], nil
		}
	}
	return nil, fmt.Errorf("kullanıcı bulunamadı (ID: %d)", user.ID)
}

// Delete kullanıcıyı siler
func (r *InMemoryUserRepository) Delete(id int) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	for i, user := range r.users {
		if user.ID == id {
			r.users = append(r.users[:i], r.users[i+1:]...)
			return nil
		}
	}
	return fmt.Errorf("kullanıcı bulunamadı (ID: %d)", id)
}

// UserHandler HTTP işleyicileri
type UserHandler struct {
	repo UserRepository
}

// NewUserHandler yeni bir handler oluşturur
func NewUserHandler(repo UserRepository) *UserHandler {
	return &UserHandler{repo: repo}
}

// GetAllUsers tüm kullanıcıları listeler
func (h *UserHandler) GetAllUsers(w http.ResponseWriter, r *http.Request) {
	users, err := h.repo.GetAll()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(users)
}

// GetUserByID ID'ye göre kullanıcı döndürür
func (h *UserHandler) GetUserByID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "Geçersiz kullanıcı ID", http.StatusBadRequest)
		return
	}

	user, err := h.repo.GetByID(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)
}

// CreateUser yeni kullanıcı oluşturur
func (h *UserHandler) CreateUser(w http.ResponseWriter, r *http.Request) {
	var user User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	newUser, err := h.repo.Create(user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(newUser)
}

// UpdateUser kullanıcı bilgilerini günceller
func (h *UserHandler) UpdateUser(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "Geçersiz kullanıcı ID", http.StatusBadRequest)
		return
	}

	var user User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	user.ID = id
	updatedUser, err := h.repo.Update(user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(updatedUser)
}

// DeleteUser kullanıcıyı siler
func (h *UserHandler) DeleteUser(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "Geçersiz kullanıcı ID", http.StatusBadRequest)
		return
	}

	if err := h.repo.Delete(id); err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// CORS ara yazılımı
func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}

// Logging ara yazılımı
func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		next.ServeHTTP(w, r)
		log.Printf("%s %s %s", r.Method, r.RequestURI, time.Since(start))
	})
}

func main() {
	// Repository oluştur
	repo := NewInMemoryUserRepository()

	// Handler oluştur
	handler := NewUserHandler(repo)

	// Router oluştur
	r := mux.NewRouter()

	// Ara yazılımları uygula
	r.Use(corsMiddleware)
	r.Use(loggingMiddleware)

	// Rotaları tanımla
	r.HandleFunc("/users", handler.GetAllUsers).Methods("GET")
	r.HandleFunc("/users/{id:[0-9]+}", handler.GetUserByID).Methods("GET")
	r.HandleFunc("/users", handler.CreateUser).Methods("POST")
	r.HandleFunc("/users/{id:[0-9]+}", handler.UpdateUser).Methods("PUT")
	r.HandleFunc("/users/{id:[0-9]+}", handler.DeleteUser).Methods("DELETE")
	r.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Sağlık durumu: Çalışıyor"))
	}).Methods("GET")

	// Sunucuyu başlat
	port := ":8080"
	log.Printf("Go Mikroservis %s portunda başlatılıyor", port)
	log.Fatal(http.ListenAndServe(port, r))
}