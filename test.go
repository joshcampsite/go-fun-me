package main

import (
	"fmt"
	"log" // Use log for fatal errors
	"os"  // Import os package to access environment variables
	"time"

	"github.com/joho/godotenv" // Import godotenv

	"github.com/reiver/go-atproto/com/atproto/repo"
	"github.com/reiver/go-atproto/com/atproto/server"
)

// No longer hardcode credentials here
// var handle string = "josh.campsite.social"
// var password string = "BRP6XjQwEFdjBWrcyZT74z1U"

func main() {
	// Load environment variables from .env file
	// It's good practice to load this early.
	err := godotenv.Load() // Loads .env from the current directory
	if err != nil {
		// Log a warning instead of fatal if .env is optional (e.g., vars could be set globally)
		// In this case, we likely depend on it, but a warning allows fallback if needed.
		log.Println("Warning: Could not load .env file:", err)
		// If .env is strictly required, uncomment the next line:
		// log.Fatal("Error loading .env file")
	}

	// Get credentials from environment variables
	handle := os.Getenv("ATPROTO_HANDLE")
	password := os.Getenv("ATPROTO_PASSWORD")

	// --- Essential Check ---
	// Verify that the variables were actually loaded
	if handle == "" {
		log.Fatal("Error: ATPROTO_HANDLE not found in environment variables or .env file.")
	}
	if password == "" {
		log.Fatal("Error: ATPROTO_PASSWORD not found in environment variables or .env file.")
	}
	// --- End Check ---

	fmt.Println("Attempting to log in as:", handle) // Use the loaded handle
	// Login
	var bearerToken string
	var userDID string // Store the DID after login
	{
		var dst server.CreateSessionResponse
		// Use the variables loaded from the environment
		err := server.CreateSession(&dst, handle, password)
		if err != nil {
			fmt.Println("Error logging in:", err) // Print the specific error
			return                                // Exit if login fails
		}
		bearerToken = dst.AccessJWT
		userDID = dst.DID // Store the user's DID from the response
		fmt.Println("Login successful!")
		fmt.Println("User DID:", userDID)
	}

	fmt.Println("Attempting to create post...")
	// Post
	var post map[string]any
	{
		when := time.Now().UTC().Format(time.RFC3339Nano)
		postText := "Test post via Go ATProto client at " + when
		post = map[string]any{
			"$type":     "app.bsky.feed.post",
			"text":      postText,
			"createdAt": when,
			// Optional: Add language if desired
			// "langs": []string{"en"},
		}
		fmt.Println("Post content:", postText)
	}

	var dst repo.CreateRecordResponse
	{
		// Use the user's DID obtained during login as the repo identifier
		var repoName string = userDID
		var collection string = "app.bsky.feed.post"
		fmt.Printf("Creating record in repo '%s', collection '%s'\n", repoName, collection)

		err := repo.CreateRecord(&dst, bearerToken, repoName, collection, post)
		if err != nil {
			fmt.Println("Error creating record:", err) // Print the specific error
			return                                     // Exit if record creation fails
		}
		fmt.Println("Post creation successful!")
		fmt.Println("Post URI:", dst.URI)
		fmt.Println("Post CID:", dst.CID)
	}

	fmt.Println("Finished.")
}
