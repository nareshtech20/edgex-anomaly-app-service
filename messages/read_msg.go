package messages

import (
	"fmt"
	"os"
	"net/http"
	"io/ioutil"
	"encoding/json"
	//"bytes"

	log "github.com/sirupsen/logrus"
	mqtt "github.com/eclipse/paho.mqtt.golang"

)


func makeMessageHandler() mqtt.MessageHandler {
	
	return func(client mqtt.Client, msg mqtt.Message) {

		fmt.Printf("Received message: %s from topic: %s\n", msg.Payload(), msg.Topic())

		log.Info("New message..."+string(msg.Payload()))

		device := "insulin-injector"
		command := "WriteBoolValue"
		settings := make(map[string]string)
		settings["Bool"] = "true"
		settings["EnableRandomization_Bool"] = "false"

		jsonData, err := json.Marshal(settings)
		if err != nil {
			log.Error("Json Marshal...")
		}
		res, err := sendCommand(device, command, "GET", string(msg.Payload()))
		if err != nil {
			log.Error("Error - send live data to anomaly service...")
		}
		log.Debug("sendCommand.."+res)
		log.Debug("jsonData..."+string(jsonData))

	}
}

func sendCommand(deviceName string, commandName string, method string, value string) (string, error) {

	log.Info("Sending anomaly update..."+value)
	url := "http://10.239.80.228:8081/edge/save?deviceId=34&value="+value+"&sensor=glucose"

	req, err := http.NewRequest(method, url, nil)
	if err != nil {
		return "", err
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	return string(respBody), nil
}


func Subscribe() {

	opts := mqtt.NewClientOptions().AddBroker("tcp://edgex-mqtt-broker:1883")
	//opts.SetDefaultPublishHandler(messageHandler)
	opts.SetDefaultPublishHandler(makeMessageHandler())
	client := mqtt.NewClient(opts)
	token := client.Connect()
	token.Wait()

	if token.Error() != nil {
		fmt.Println(token.Error())
		os.Exit(1)
	}

	token = client.Subscribe("glucose", 0, nil)
	token.Wait()
	fmt.Printf("Successfully subscribed to topic: %s\n", "glucose")

	select {} // block forever
    	// Start a goroutine to keep the application running until interrupted.
    	//go func() {
        //	select {}
    	//}()	
}
