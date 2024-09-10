package srm

import (
	"fmt"
	"log"
	"math/rand"
	"os"
	"os/exec"
	"path/filepath"
	"time"
)

func randomListenerPort() int {
	// Инициализация генератора случайных чисел
	rand.NewSource(time.Microsecond.Nanoseconds())

	// Первое число в диапазоне, кратное 4
	min := 5000 + 4 - (5000 % 4)
	// Последнее число в диапазоне, кратное 4
	max := 6000 - (6000 % 4)

	// Количество чисел, кратных 4, в диапазоне
	count := (max-min)/4 + 1

	// Выбор случайного индекса
	randomIndex := rand.Intn(count)

	// Возвращение случайного числа, кратного 4
	return min + randomIndex*4
}

func createUserHome(username, localStorage string) string {

	// Создаем путь к домашней директории пользователя
	userHome := filepath.Join(localStorage, username+"-home")

	// Создаем директорию
	if err := os.Mkdir(userHome, 0777); err != nil {
		if os.IsExist(err) {
			fmt.Println("Папка уже существует.")
			return userHome
		} else {
			log.Fatalf("Ошибка при создании папки: %v", err)
		}
	}

	// Формируем команду rsync
	source := filepath.Join(localStorage, "FILE_BASE") + "/"
	dest := userHome + "/"
	cmd := exec.Command("rsync", "-a", source, dest)

	// Выполняем rsync
	if err := cmd.Run(); err != nil {
		log.Fatalf("Ошибка при выполнении rsync: %v", err)
	}

	fmt.Println("Папки успешно созданы.")
	return userHome
}
