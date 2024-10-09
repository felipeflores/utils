package utils

import "github.com/felipeflores/utils/config"

// func LoadJson[C any](configPath string) (*C, error) {

// }
// teste
func LoadYaml[C any](configPath string) (*C, error) {
	return config.LoadConfig[C](configPath)
}
