package display

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/FurmanovVitaliy/pixel-cloud/pkg/errors"
)

type XServerRepository interface {
	Populate([]XServer) error
	GetAll() ([]XServer, error)
}

type XServerService struct {
	xServerRepository XServerRepository
}

// custom errors
var (
	ErrInvalidScriptPath = errors.New(418, "DS", "00001", "invalid script path")
	ErrInvalidJsonFile   = errors.New(418, "DS", "00002", "invalid json file")
)

func NewXServerService(xServerRepository XServerRepository) *XServerService {
	return &XServerService{xServerRepository: xServerRepository}
}

func (s *XServerService) GetAll() (xServers []XServer) {
	xServers, _ = s.xServerRepository.GetAll()
	return
}

func (s *XServerService) PopulateViaLocalScript(scriptPath, resultJsonFile string) (err error) {

	if scriptPath == "" || filepath.Ext(scriptPath) != ".sh" {
		return ErrInvalidScriptPath
	}

	if resultJsonFile == "" || filepath.Ext(resultJsonFile) != ".json" {
		return ErrInvalidJsonFile
	}

	if err = exec.Command("chmod", "+x", scriptPath).Run(); err != nil {
		return fmt.Errorf("failed to set execute permission: %w", err)
	}

	cmd := exec.Command("sh", scriptPath)
	output, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Printf("Error running script: %v, Output: %s\n", err, output)
		return fmt.Errorf("error running script: %w, Output: %s", err, output)
	}
	fmt.Printf("Script output: \n%s\n", output)

	res, err := os.Open(resultJsonFile)
	if err != nil {
		return fmt.Errorf("error opening JSON file: %w", err)
	}
	defer res.Close()

	byteValue, err := io.ReadAll(res)
	if err != nil {
		return fmt.Errorf("error reading JSON file: %w", err)
	}

	var xservers []XServer
	if err = json.Unmarshal(byteValue, &xservers); err != nil {
		return fmt.Errorf("error unmarshaling JSON: %w", err)
	}

	if err = s.xServerRepository.Populate(xservers); err != nil {
		return fmt.Errorf("error populating repository: %w", err)
	}

	return nil
}
