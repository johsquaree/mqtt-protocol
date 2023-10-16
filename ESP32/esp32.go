package mqtt

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
	// "mqtt://broker_adresi:port_numarası"
	// "mqtt://broker_adresi:879687"
	opts.AddBroker(fmt.Sprintf("mqtt://%s:%d", broker, port))

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

	// İleri, geri, sağ ve sol komutlarını temsil eden MQTT konuları oluşturuyoruz
	forwardTopic := "/move/forward"
	backwardTopic := "/move/backward"
	clockWiseTopic := "/control/clockwise"
	counterClockWiseTopic := "/control/counterclockwise"

	// İleri komutuna abone oluyoruz
	subscribe(client, forwardTopic, 2) // 2 adet ileri komutu

	// Geri komutuna abone oluyoruz
	subscribe(client, backwardTopic, 4) // 4 adet geri komutu

	// Sağa dönme komutuna abone oluyoruz
	subscribe(client, clockWiseTopic, 3) // 3 adet sağa dönme komutu

	// Sola dönme komutuna abone oluyoruz
	subscribe(client, counterClockWiseTopic, 1) // 1 adet sola dönme komutu

	// Sonsuz bir döngüde programı çalışır halde tutmak için kullanılır
	select {}
}

func subscribe(client mqtt.Client, topic string, count int) {
	for i := 0; i < count; i++ {
		token := client.Subscribe(fmt.Sprintf("%s/%d", topic, i), 0, messageHandler)
		token.Wait()
		fmt.Printf("Abone olundu: %s/%d\n", topic, i)
	}
}
