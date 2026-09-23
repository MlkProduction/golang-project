package core_logger

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)


/*
Не очень понимаю, что делает код, но попробую
1. "Value возвращает значение, связанное с этим контекстом для ключа, или nil, если с ключом не связано никакое значение. 
При последовательных вызовах Value с одним и тем же ключом возвращается один и тот же результат."

ctx.Value("log") — достаём из context значение по ключу "log".
.(*Logger) — пытаемся утверждать:
«это значение имеет тип *Logger».

2. panic происходит, если по ключу "log" не удалось получить значение типа *Logger.
3. Возвращается сам логгер, точнее указатель *Logger, который лежал в context.
 
*/
func FromContext(ctx context.Context) *Logger {
	log, ok := ctx.Value("log").(*Logger)
	if !ok {
		panic("no log in context")
	}
	return log
}

type Logger struct {
	*zap.Logger // сам логгер от zap; благодаря embedding его методы доступны через наш Logger

	file *os.File // а file здесь нужен обёрточному Logger, чтобы иметь ссылку на открытый файл и корректно освободить ресурс.
}

func NewLogger(config Config) (*Logger, error) {
	// zapLvl = создали "контейнер" для уровня
	// Проверяем уровень в понятный для zap
	// если уровень не понятен, unmarshal log level
	zapLvl := zap.NewAtomicLevel()                                    
	if err := zapLvl.UnmarshalText([]byte(config.Level)); err != nil { 
		return nil, fmt.Errorf("unmarshal log level: %w", err)
	}

	// Создаём папку для логов, если её ещё нет.
	// Что будет, если папка уже существует? * os.MkdirAll не ругнётся, если директория уже существует.
	/*  logs нет → создаст;
		logs уже есть → продолжит без ошибки;
		по пути что-то не так → вернёт ошибку.
		*/
	if err := os.MkdirAll(config.Folder, 0755); err != nil {
		return nil, fmt.Errorf("mkdir log folder: %w", err)
	}

	// Формируем имя файла на основе текущего времени.
	timestamp := time.Now().UTC().Format("2006-01-02T15-04-05.000000")

	// Собираем полный путь:
	// config.Folder + имя файла.
	logFilePath := filepath.Join(
		config.Folder,
		fmt.Sprintf("%s.log", timestamp),
	)

	// Открываем/создаём файл.
	// Обрати внимание на флаги O_CREATE и O_WRONLY:
	// как они влияют на то, что произойдёт с файлом?
	logFile, err := os.OpenFile(logFilePath, os.O_CREATE|os.O_WRONLY, 0644) //«Создай файл, если его нет, и открой его для записи».
	if err != nil {
		return nil, fmt.Errorf("log file not open %w", err)
	}

	// Настраиваем формат, в котором zap будет превращать
	// информацию о логе в текст.
	zapConfig := zap.NewDevelopmentEncoderConfig()

	// Отдельно говорим, как отображать время в готовой записи лога.
	zapConfig.EncodeTime = zapcore.TimeEncoderOfLayout(
		"2006-01-02T15:04:05.000000",
	)

	// Encoder отвечает примерно за:
	// "как будет выглядеть готовая строка лога?"
	zapEncoder := zapcore.NewConsoleEncoder(zapConfig)

	// Создаём два направления, куда будут уходить ОДНИ И ТЕ ЖЕ логи:
	// 1. os.Stdout → консоль
	// 2. logFile  → файл
	//
	// NewTee объединяет эти направления.
	core := zapcore.NewTee(
		zapcore.NewCore(zapEncoder, zapcore.AddSync(os.Stdout), zapLvl),
		zapcore.NewCore(zapEncoder, zapcore.AddSync(logFile), zapLvl),
	)

	// Создаём настоящий zap.Logger на основе настроенного core.
	// AddCaller добавляет информацию о месте в коде,
	// откуда был вызван логгер.
	zapLogger := zap.New(core, zap.AddCaller())

	// Заворачиваем zap.Logger в наш Logger
	// и сохраняем открытый файл, чтобы потом его можно было закрыть.
	return &Logger{
		Logger: zapLogger,
		file:   logFile,
	}, nil
}

func (l *Logger) With(field ...zap.Field) *Logger {
	// Создаём новый Logger с дополнительными полями,
	// например: request_id, user_id и т.д.
	//
	// При этом сам файл НЕ открываем заново —
	// передаём тот же указатель на уже открытый файл.
	return &Logger{
		Logger: l.Logger.With(field...),
		file:   l.file,
	}
}

func (l *Logger) Close() {
	// Закрываем файл, который был открыт в NewLogger.
	// Зачем вообще его закрывать?
	if err := l.file.Close(); err != nil {
		fmt.Println("failed to close logger")
	}
}
