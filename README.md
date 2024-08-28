Golang library for MessageBroker:

Example of use:

go get github.com/MeDaLL1ST/MBLibrary-GO


mb := Init("host:port", "api_key", "http")
	defer mb.Close()
	err := mb.Add("qwe", "----123------")
	if err != nil {
		fmt.Println(err)
	}
	err = mb.Subscribe("somekey")
	if err != nil {
		fmt.Println(err)
	}
	message, _ := mb.Read()

	keys, err := mb.List()
	if err != nil {
		log.Fatal(err)
	}

	for _, key := range keys {
		fmt.Println(key)
	}
	fmt.Println(mb.Info("qwe1"))
	log.Printf("Got message by WebSocket: %s", message)
