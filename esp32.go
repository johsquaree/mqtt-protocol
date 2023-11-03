package main

import (
	"fmt"
	"github.com/eclipse/paho.mqtt.golang"
)

var MQTTClient mqtt.Client // Global MQTT client variable

func connectHandler(client mqtt.Client) {
    fmt.Println("Connected") // Connection handler function
}

func connectLostHandler(client mqtt.Client, err error) {
    fmt.Printf("Connection lost: %v\n", err) // Connection lost handler function
}

func messageHandler(client mqtt.Client, message mqtt.Message) {
    fmt.Printf("Received message on topic: %s\n", message.Topic())
    fmt.Printf("Message: %s\n", message.Payload())
}

func Publish(message string) {
    token := MQTTClient.Publish("/movements", 0, false, message) // Publish a message to the "/movements" topic
    token.Wait()
    if token.Error() != nil {
        fmt.Println(token.Error())
    }
}

func InitializeMQTTClient() {
    var broker = "broker_url"
    var port = 1883
    opts := mqtt.NewClientOptions()
    opts.AddBroker(fmt.Sprintf("tcp://%s:%d", broker, port))
    opts.SetClientID("go_mqtt_client")
    opts.SetUsername("admin")
    opts.SetPassword("instar")
    opts.SetDefaultPublishHandler(messagePubHandler)
    opts.OnConnect = connectHandler
    opts.OnConnectionLost = connectLostHandler
    MQTTClient = mqtt.NewClient(opts)

    if token := MQTTClient.Connect(); token.Wait() && token.Error() != nil {
        panic(token.Error())
    }
}

func main() {
    // Initialize the MQTT client
    InitializeMQTTClient()

    // Subscribe to the MQTT client
    token := MQTTClient.Subscribe("/movement", 0, func(client mqtt.Client, message mqtt.Message) {
        fmt.Printf("Received message on topic %s: %s\n", message.Topic(), message.Payload())
        // Handle the movement command here
    })
    token.Wait()

    // Publish a movement command
    Publish("forward")
}
