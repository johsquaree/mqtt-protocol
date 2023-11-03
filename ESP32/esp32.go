package main

import (
	"fmt"
	"github.com/eclipse/paho.mqtt.golang"
)

var MQTTClient mqtt.Client // Global MQTT istemci değişkeni

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

func Publish(message string) {
    token := MQTTClient.Publish("/movements", 0, false, message)
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
    // MQTT istemcisini başlat
    InitializeMQTTClient()

    // MQTT istemcisine abone ol
    token := MQTTClient.Subscribe("/movement", 0, func(client mqtt.Client, message mqtt.Message) {
        fmt.Printf("Received message on topic %s: %s\n", message.Topic(), message.Payload())
        // Hareket komutunu burada işleyin
    })
    token.Wait()

    // Hareket komutunu yayınla
    Publish("forward")

}