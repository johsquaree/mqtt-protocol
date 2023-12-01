package main

import (
	"fmt"
	"net/http"
	"strings"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

// MQTTClient represents the MQTT client used for communication.
var MQTTClient mqtt.Client

// connectHandler is a callback function called upon successful MQTT connection.
func connectHandler(client mqtt.Client) {
	fmt.Println("Connected")
}

// connectLostHandler is a callback function called when the MQTT connection is lost.
func connectLostHandler(client mqtt.Client, err error) {
	fmt.Printf("Connection lost: %v\n", err)
}

// Publish sends an MQTT message with the specified payload to the "/movements" topic.
func Publish(message string) {
	token := MQTTClient.Publish("/movements", 0, false, message)
	token.Wait()
	if token.Error() != nil {
		fmt.Println(token.Error())
	}
}

// messagePubHandler is a callback function called when an MQTT message is published.
func messagePubHandler(client mqtt.Client, message mqtt.Message) {
	fmt.Printf("Published message on topic: %s\n", message.Topic())
	fmt.Printf("Message: %s\n", message.Payload())
}

// InitializeMQTTClient initializes and connects the MQTT client.
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

// messageHandler is a callback function called when an MQTT message is received on the "/movement" topic.
func messageHandler(client mqtt.Client, message mqtt.Message) {
	fmt.Printf("Received message on topic: %s\n", message.Topic())
	fmt.Printf("Message: %s\n", message.Payload())
	// Send the payload message to an HTTP channel
	payload := message.Payload()
	sendPayloadToHTTPChannel(payload)
}

// sendPayloadToHTTPChannel sends the payload to an HTTP channel using a POST request.
func sendPayloadToHTTPChannel(payload []byte) {
	url := "http://http_channel_url"
	payloadStr := string(payload)
	resp, err := http.Post(url, "application/json", strings.NewReader(payloadStr))
	if err != nil {
		fmt.Println("Error sending payload to HTTP channel:", err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		fmt.Printf("HTTP request failed with status code: %d\n", resp.StatusCode)
	}
}

func main() {
	InitializeMQTTClient()
	token := MQTTClient.Subscribe("/movement", byte(2), func(client mqtt.Client, message mqtt.Message) {
		fmt.Printf("Received message on topic %s: %s\n", message.Topic(), message.Payload())
	})
	token.Wait()

	Publish("forward")
}
