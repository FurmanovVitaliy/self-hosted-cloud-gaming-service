package srm

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"os/exec"
)

type XServer struct {
	ScreenNumber string `json:"screen"`
	Card         string `json:"card"`
	Port         string `json:"port_name"`
	Connector    int    `json:"connector_id"`
	Plane        int    `json:"plane_id"`
	Used         bool
}

func initXservers() {
	scriptPath := "/home/vitalii/Desktop/new-enable-all-screens.sh"

	// Установка прав на выполнение скрипта (необходимо запустить с правами администратора)
	if err := exec.Command("chmod", "+x", scriptPath).Run(); err != nil {
		fmt.Printf("Failed to set execute permission: %s\n", err)
		return
	}

	// Запуск скрипта для включения всех экранов
	cmd := exec.Command("sh", scriptPath)
	output, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Printf("Error running script: %s, Output: %s\n", err, output)
	} else {
		fmt.Printf("Script output: %s\n", output)
	}
}

func jsonToXserver() []xServer {
	jsonFile, err := os.Open("/home/vitalii/Desktop/display-info.json")
	if err != nil {
		fmt.Println("Error opening JSON file:", err)
		return nil
	}
	defer jsonFile.Close()

	byteValue, err := io.ReadAll(jsonFile)
	if err != nil {
		fmt.Println("Error reading JSON file:", err)
		return nil
	}

	// We initialize our array of XServers
	var xservers []xServer

	// we unmarshal our byteArray which contains our JSON file's content into 'xservers'
	json.Unmarshal(byteValue, &xservers)

	return xservers
}
