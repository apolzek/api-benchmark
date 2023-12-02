package main

import (
    "context"
    "encoding/json"
    "io/ioutil"
    "log"
    "net/http"
    "os"
    "strconv"

    "github.com/jackc/pgx/v4/pgxpool"
)

var pool *pgxpool.Pool

func main() {
    // no need for .env file, env is configured in docker-compose
    //err := godotenv.Load()
    //if err != nil {
    //    log.Fatal("Error loading .env file")
    //}

    connectDB()
    defer pool.Close()

    http.HandleFunc("/user", userHandler)
    port := os.Getenv("PORT")
    if port == "" {
        port = "3000"
    }

    log.Printf("Server running on http://localhost:%s\n", port)
    http.ListenAndServe(":"+port, nil)
}

func connectDB() {
    var err error
    dbURL := "postgres://" + os.Getenv("POSTGRES_USER") + ":" +
        os.Getenv("POSTGRES_PASSWORD") + "@" +
        os.Getenv("POSTGRES_HOST") + "/" +
        os.Getenv("POSTGRES_DATABASE") //+ "?sslmode=disable"

    config, err := pgxpool.ParseConfig(dbURL)
    if err != nil {
        log.Fatalf("Unable to parse pool config: %v\n", err)
    }

    config.MaxConns = 60

    pool, err = pgxpool.ConnectConfig(context.Background(), config)
    if err != nil {
        log.Fatalf("Unable to connect to database: %v\n", err)
    }
}

func userHandler(w http.ResponseWriter, r *http.Request) {
    if r.Method != http.MethodPost {
        http.Error(w, "Not Found", http.StatusNotFound)
        return
    }

    body, err := ioutil.ReadAll(r.Body)
    if err != nil {
        http.Error(w, "Bad Request", http.StatusBadRequest)
        return
    }

    var data struct {
        Email    string `json:"email"`
        Password string `json:"password"`
    }
    err = json.Unmarshal(body, &data)
    if err != nil {
        http.Error(w, "Bad Request", http.StatusBadRequest)
        return
    }

    userID, err := CreateUser(data.Email, data.Password)
    if err != nil {
        log.Printf("Error creating user: %v\n", err)
        http.Error(w, "Internal Server Error", http.StatusInternalServerError)
        return
    }

    w.Header().Set("Location", "/user/"+strconv.Itoa(userID))
    w.WriteHeader(http.StatusCreated)
    json.NewEncoder(w).Encode(map[string]string{"message": "User created"})
}

func CreateUser(email, password string) (int, error) {
    var id int
    err := pool.QueryRow(context.Background(),
        "INSERT INTO users(email, password) VALUES($1, $2) RETURNING id",
        email, password).Scan(&id)
    return id, err
}

