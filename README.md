Golang library for MessageBroker: https://github.com/MeDaLL1ST/MessageBroker

Example of use:

go get github.com/MeDaLL1ST/MBLibrary-GO


	import mbl "github.com/MeDaLL1ST/MBLibrary-GO"

 
	mb := mbl.InitMb("host:port", "api_key", "http")
	defer mb.Close()
	err := mb.Add("qwe", "----123------")
	if err != nil {
		fmt.Println(err)
	}
	err = mb.Subscribe("somekey")
	if err != nil {
		fmt.Println(err)
	}
 	go func() {
		fmt.Println(mb.ReadSync("qwe1", func() { fmt.Println(mb.ReadSync("qwe2")) }))
	}()
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

	gw := InitGw("host:port", "api_key", "http")
	gw.Add("qwe16", "data1", "topic3")
	ans, err := gw.ReadSync("qwe17", "topic3", func() { gw.Add("qwe17", "data2", "topic3") })
	if err != nil {
		fmt.Println("Error:", err)
	}
	fmt.Println(ans)
