package main

import (
	"fmt"
	"github.com/eclipse/paho.mqtt.golang"
)

func connectHandler(client mqtt.Client) {
	fmt.Println("Bağlandı")
}

func connectLostHandler(client mqtt.Client, err error) {
	fmt.Printf("Bağlanti kaybedildi: %v\n", err)
}

func messageHandler(client mqtt.Client, message mqtt.Message) {
	fmt.Printf("Received message on topic: %s\n", message.Topic())
	fmt.Printf("Message: %s\n", message.Payload())
}

func subscribeMovements(client mqtt.Client, apiMovementTopic string) {
	movements := []string{"forward", "backward", "clockwise", "counterclockwise",}

	for _, movement := range movements {
		topic := fmt.Sprintf("%s/%s", apiMovementTopic, movement)
		token := client.Subscribe(topic, 0, messageHandler)
		token.Wait()
		fmt.Printf("Abone olundu: %s\n", topic)
	}
}

func main() {
	var broker = "broker_adresi"
	var port = 1883

	opts := mqtt.NewClientOptions()
	opts.AddBroker(fmt.Sprintf("mqtt://%s:%d", broker, port))
	opts.SetClientID("go_mqtt_client")
	opts.SetUsername("kullanici_adi")
	opts.SetPassword("sifre")
	opts.OnConnect = connectHandler
	opts.OnConnectionLost = connectLostHandler

	client := mqtt.NewClient(opts)

	if token := client.Connect(); token.Wait() && token.Error() != nil {
		panic(token.Error())
	}

	apiMovementTopic := "/api/movement"
	subscribeMovements(client, apiMovementTopic)

}
