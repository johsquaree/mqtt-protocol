package main

import (
	"fmt"
	"github.com/eclipse/paho.mqtt.golang" 
	// MQTT istemcisini içe aktarıyoruz
)

// MQTT broker ile bağlantı kurulduğunda çalışacak işlev
func connectHandler(client mqtt.Client) {
	fmt.Println("Bağlandi")
}

// Bağlantı kaybedildiğinde çalışacak işlev
func connectLostHandler(client mqtt.Client, err error) {
	fmt.Printf("Bağlanti kaybedildi: %v\n", err)
}

// MQTT mesajlarını işlemek için kullanılacak işlev
func messageHandler(client mqtt.Client, message mqtt.Message) {
	fmt.Printf("Received message on topic: %s\n", message.Topic())
	fmt.Printf("Message: %s\n", message.Payload())
}

func main() {
	// MQTT broker adresi
	var broker = "broker_adresi"
	// MQTT broker port numarası
	var port = 879687

	// MQTT istemcisini yapılandırmak için seçenekleri oluşturuyoruz
	opts := mqtt.NewClientOptions()

	// MQTT broker adresi ve portunu seçeneklere ekliyoruz 
	// "tcp://broker_adresi:port_numarası"
	// "tcp://broker_adresi:879687"
	opts.AddBroker(fmt.Sprintf("tcp://%s:%d", broker, port))

	// MQTT istemci kimliği
	opts.SetClientID("go_mqtt_client")

	// MQTT broker için kullanıcı adı
	opts.SetUsername("fgsafasfas")

	// MQTT broker için şifre
	opts.SetPassword("fdsfdsfs")

	// Bağlantı olayları için işlevleri ayarlıyoruz
	opts.OnConnect = connectHandler
	opts.OnConnectionLost = connectLostHandler

	// MQTT istemcisini oluşturuyoruz
	client := mqtt.NewClient(opts)

	// MQTT broker ile bağlantı kuruyoruz
	if token := client.Connect(); token.Wait() && token.Error() != nil {
		panic(token.Error())
	}

	// "/forward" adlı MQTT konusuna abone oluyoruz
	token := client.Subscribe("/forward", 0, messageHandler)

	// Abonelik işleminin sonucunu yazdırıyoruz
	if token.Wait() && token.Error() != nil {
		// panic işlemi anında sonlandırır.
		panic(token.Error())
	}

	// Sonsuz bir döngüde programı çalışır halde tutmak için kullanılır
	select {}
}