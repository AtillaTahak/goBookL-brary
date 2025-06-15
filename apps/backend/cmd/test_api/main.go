package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

const baseURL = "http://localhost:8080"

func main() {
	fmt.Println("🧪 Testing Book Library API...")

	time.Sleep(2 * time.Second)

	testHealthCheck()

	testUserRegistration()

	token := testUserLogin()

	testBookOperations(token)

	fmt.Println("✅ All tests completed!")
}

func testHealthCheck() {
	fmt.Println("\n📋 Testing health check...")
	resp, err := http.Get(baseURL + "/")
	if err != nil {
		fmt.Printf("❌ Health check failed: %v\n", err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode == 200 {
		fmt.Println("✅ Health check passed")
	} else {
		fmt.Printf("❌ Health check failed with status: %d\n", resp.StatusCode)
	}
}

func testUserRegistration() {
	fmt.Println("\n👤 Testing user registration...")

	user := map[string]string{
		"username": "testuser",
		"password": "testpass123",
		"email":    "test@example.com",
	}

	jsonData, _ := json.Marshal(user)
	resp, err := http.Post(baseURL+"/auth/register", "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		fmt.Printf("❌ Registration failed: %v\n", err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode == 201 {
		fmt.Println("✅ User registration passed")
	} else {
		fmt.Printf("❌ User registration failed with status: %d\n", resp.StatusCode)
	}
}

func testUserLogin() string {
	fmt.Println("\n🔐 Testing user login...")

	user := map[string]string{
		"username": "testuser",
		"password": "testpass123",
	}

	jsonData, _ := json.Marshal(user)
	resp, err := http.Post(baseURL+"/auth/login", "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		fmt.Printf("❌ Login failed: %v\n", err)
		return ""
	}
	defer resp.Body.Close()

	if resp.StatusCode == 200 {
		var result map[string]interface{}
		json.NewDecoder(resp.Body).Decode(&result)
		token := result["token"].(string)
		fmt.Println("✅ User login passed")
		return token
	} else {
		fmt.Printf("❌ User login failed with status: %d\n", resp.StatusCode)
		return ""
	}
}

func testBookOperations(token string) {
	fmt.Println("\n📚 Testing book operations...")

	// Test get books
	resp, err := http.Get(baseURL + "/books")
	if err != nil {
		fmt.Printf("❌ Get books failed: %v\n", err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode == 200 {
		fmt.Println("✅ Get books passed")
	} else {
		fmt.Printf("❌ Get books failed with status: %d\n", resp.StatusCode)
	}

	// Test create book (protected)
	if token != "" {
		book := map[string]interface{}{
			"title":  "Test Book",
			"author": "Test Author",
			"year":   2024,
			"genre":  "Test Genre",
		}

		jsonData, _ := json.Marshal(book)
		req, _ := http.NewRequest("POST", baseURL+"/books", bytes.NewBuffer(jsonData))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer "+token)

		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			fmt.Printf("❌ Create book failed: %v\n", err)
			return
		}
		defer resp.Body.Close()

		if resp.StatusCode == 201 {
			fmt.Println("✅ Create book passed")
		} else {
			fmt.Printf("❌ Create book failed with status: %d\n", resp.StatusCode)
		}
	}
}
