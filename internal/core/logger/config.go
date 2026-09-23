package core_logger

import (
	"fmt"

	"github.com/kelseyhightower/envconfig"
)

type Config struct {
	Level  string `envconfig:"LEVEL" required:"true"` // уровень лога
	Folder string `envconfig:"FOLDER" required:"true"` // куда будет записывать
}

func NewConfig() (Config, error) {                                // создаём функцию, возвращаем конфиг, ошибку
	var config Config									// переменная типа конфиг

	err := envconfig.Process("LOGGER", &config)         // префикс нужен, что бы взять из .env 
														// например LOGGER_LEVEl, LOGGER_FOLDER
	if err != nil { //									// проверяем на ошибку
		return Config{}, fmt.Errorf("process envconfig: %w", err)  // возвращаем пустой конфиг и ошибку
	}

	return config, nil	// ошибки нет, возвращаем заполненый конфиг 
}

func NewConfigMust() Config {
	config, err := NewConfig()

	if err != nil {

		err = fmt.Errorf("get logger config: %w", err)
		panic(err)
	}
	return config
}
