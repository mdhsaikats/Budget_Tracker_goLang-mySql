package main

func main() {
	connection := db.Connect()
	defer connection.Close()

}
