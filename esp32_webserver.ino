//MucitPilot 2021/////
//ESP-32 Web Server Örneği
#include <WiFi.h>
#include <AsyncTCP.h>//https://github.com/khoih-prog/AsyncHTTPRequest_Generic/archive/master.zip
#include <ESPAsyncWebServer.h>//https://github.com/me-no-dev/ESPAsyncWebServer/archive/master.zip
#include <Adafruit_Sensor.h>//https://github.com/adafruit/Adafruit_Sensor
#include <DHT.h>//https://github.com/adafruit/DHT-sensor-library

// Ağ adı ve şifrenizi girin
const char* ssid = "****";
const char* password = "*****";

// 80 portunu dinleyen  Bir AsyncWebServer nesnesi yaratıyoruz
AsyncWebServer server(80);

//  /events adresinde bir EventSource nesnesi yaratıyoruz
AsyncEventSource events("/events");

// Zamanlama ile ilgili değişkenler
unsigned long sonZaman = 0;  
unsigned long beklemeSuresi = 5000; //5sn


//sıcaklık ve nem değişkenini tanımlıyoruz
float temperature;
float humidity;

//sabit ip için gerekli bilgiler///
// Sabit IP adresini girin
IPAddress local_IP(192, 168, 1, 180);
// Gateway IP adresini girin
IPAddress gateway(192, 168, 1, 1);
//subnetmask adresini girin
IPAddress subnet(255, 255, 0, 0);


//////////////DHT SENSÖR AYAR BLOĞU////////////////////
// DHT sensör GPIO14 pinine bağlı
#define DHTPIN 14  
// DHT tipini seçin
#define DHTTYPE    DHT11     // DHT 11
//#define DHTTYPE    DHT22     // DHT 22 (AM2302)
//#define DHTTYPE    DHT21     // DHT 21 (AM2301)
DHT dht(DHTPIN, DHTTYPE); //bir dht nesnesi oluşturuyoruz
//////////////////////////////////////////////////////////

///sensörden veri okuma fonksiyonu
void sensorOku(){
humidity = dht.readHumidity();
 delay(100);
  //veri okunamadıysa
  if (isnan(humidity)){
    Serial.println("DHT sensör nem verisi okunamadı!!!");
    humidity = 0.0;
  }
temperature = dht.readTemperature();
 delay(100);
  //veri okunamadıysa
  if (isnan(temperature)){
    Serial.println("DHT sensör sıcaklık verisi okunamadı!!!");
    temperature = 0.0;
  }
  
}

// Wifi bağlantı fonksiyonu
void initWiFi() {
    WiFi.mode(WIFI_STA);
    WiFi.begin(ssid, password);
    Serial.print("Kablosuz Ağa Bağlanıyor..");
    while (WiFi.status() != WL_CONNECTED) {
        Serial.print('.');
        delay(1000);
    }
    Serial.println(WiFi.localIP());//aldığı ip adresini yazdırıyoruz
}

String processor(const String& var){
  sensorOku();
  //Serial.println(var);
  if(var == "TEMPERATURE"){
    return String(temperature);
  }
  else if(var == "HUMIDITY"){
    return String(humidity);
  }

  return String();
}

//sunucu sayfası için gerekli HTML kodları içeren fonksiyon
const char index_html[] PROGMEM = R"rawliteral(
<!DOCTYPE HTML><html>
<head>
  <title>ESP32 Web Server</title>
  <meta name="viewport" content="width=device-width, initial-scale=1">
  <link rel="stylesheet" href="https://use.fontawesome.com/releases/v5.7.2/css/all.css" integrity="sha384-fnmOCqbTlWIlj8LyTjo7mOUStjsKC4pOpQbqyi7RrhN7udi9RwhKkMHpvLbHG9Sr" crossorigin="anonymous">
  <link rel="icon" href="data:,">
  <style>
    html {font-family: Arial; display: inline-block; text-align: center;}
    p { font-size: 1.2rem;}
    body {  margin: 0;}
    .topnav { overflow: hidden; background-color: #50B8B4; color: white; font-size: 1rem; }
    .content { padding: 20px; }
    .card { background-color: white; box-shadow: 2px 2px 12px 1px rgba(140,140,140,.5); }
    .cards { max-width: 800px; margin: 0 auto; display: grid; grid-gap: 2rem; grid-template-columns: repeat(auto-fit, minmax(200px, 1fr)); }
    .reading { font-size: 1.4rem; }
  </style>
</head>
<body>
  <div class="topnav">
    <h1>ESP-32 WEB SERVER</h1>
  </div>
  <div class="content">
    <div class="cards">
      <div class="card">
        <p><i class="fas fa-thermometer-half" style="color:#059e8a;"></i> SICAKLIK</p><p><span class="reading"><span id="temp">%TEMPERATURE%</span> &deg;C</span></p>
      </div>
      <div class="card">
        <p><i class="fas fa-tint" style="color:#00add6;"></i> NEM</p><p><span class="reading"><span id="hum">%HUMIDITY%</span> &percnt;</span></p>
      </div>
    </div>
  </div>
<script>
if (!!window.EventSource) {
 var source = new EventSource('/events');
 
 source.addEventListener('open', function(e) {
  console.log("Events Connected");
 }, false);
 source.addEventListener('error', function(e) {
  if (e.target.readyState != EventSource.OPEN) {
    console.log("Events Disconnected");
  }
 }, false);
 
 source.addEventListener('message', function(e) {
  console.log("message", e.data);
 }, false);
 
 source.addEventListener('temperature', function(e) {
  console.log("temperature", e.data);
  document.getElementById("temp").innerHTML = e.data;
 }, false);
 
 source.addEventListener('humidity', function(e) {
  console.log("humidity", e.data);
  document.getElementById("hum").innerHTML = e.data;
 }, false);
 
}
</script>
</body>
</html>)rawliteral";

void setup() {
  Serial.begin(115200);
  // Statik ip yapılandırması
  if (!WiFi.config(local_IP, gateway, subnet)) {
     Serial.println("Statik IP ayarlanamadı");
  }
  initWiFi();



  // Web Server'a gelen istekler için gerekli kısım
  server.on("/", HTTP_GET, [](AsyncWebServerRequest *request){
    request->send_P(200, "text/html", index_html, processor);
  });

  // Web server Event'ler için
  events.onConnect([](AsyncEventSourceClient *client){
    if(client->lastId()){
      Serial.printf("İstemci Yeniden Baglandi! Aldigi son mesaj ID'si: %u\n", client->lastId());
    }
   
    client->send("Merhaba!", NULL, millis(), 5000);
  });
  server.addHandler(&events);
  server.begin();
}

void loop() {
  if ((millis() - sonZaman) > beklemeSuresi) {//sensörden belirlenmiş zaman aralığında veri okuyoruz
    sensorOku();
    Serial.printf("Sıcaklık = %.2f ºC \n", temperature);
    Serial.printf("Nem = %.2f \n", humidity);
    Serial.println();

    // Sensörden okunan bilgileri EVENT olarak sunucuya gönder
    events.send("ping",NULL,millis());
    events.send(String(temperature).c_str(),"temperature",millis());
    events.send(String(humidity).c_str(),"humidity",millis());
   
    
    sonZaman = millis();
  }
}
