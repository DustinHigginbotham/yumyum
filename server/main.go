package main

func main() {
	config := loadConfig()
	server := NewServer(config)

	// server starts by default on port 8666
	server.Start()
}
