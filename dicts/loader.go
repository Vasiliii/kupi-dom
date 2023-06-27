package dicts

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"strconv"
)

func Load(dict *map[string]bool, filename string) error {
	file, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		(*dict)[scanner.Text()] = true
		fmt.Println(scanner.Text())
	}
	return nil
}

func LoadInt64(dict *map[int64]struct{}, filename string) error {
	file, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	var id int64
	for {
		_, err = fmt.Fscanf(file, "%d\n", &id)
		if err != nil {
			if err == io.EOF {
				break
			}
			return err
		}
		fmt.Println(id)
		(*dict)[id] = struct{}{}
	}
	return nil
}

func AddInt64(dict *map[int64]struct{}, id int64, filename string) error {
	Mu.Lock()
	defer Mu.Unlock()
	file, err := os.OpenFile(filename, os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = fmt.Fprintf(file, "%d\n", id)
	if err != nil {
		return err
	}
	(*dict)[id] = struct{}{}
	return nil
}

func RemoveInt64(dict *map[int64]struct{}, id int64, filename string) error {
	Mu.Lock()
	defer Mu.Unlock()

	// Удаляем идентификатор из словаря
	delete(*dict, id)

	// Читаем содержимое файла
	fileData, err := ioutil.ReadFile(filename)
	if err != nil {
		return err
	}

	// Преобразуем содержимое файла в список идентификаторов
	ids := make([]int64, 0)
	scanner := bufio.NewScanner(bytes.NewReader(fileData))
	for scanner.Scan() {
		line := scanner.Text()
		parsedID, err := strconv.ParseInt(line, 10, 64)
		if err != nil {
			continue
		}
		ids = append(ids, parsedID)
	}

	// Удаляем идентификатор из списка
	updatedIDs := make([]int64, 0)
	for _, existingID := range ids {
		if existingID != id {
			updatedIDs = append(updatedIDs, existingID)
		}
	}

	// Записываем обновленный список обратно в файл
	file, err := os.OpenFile(filename, os.O_WRONLY|os.O_TRUNC, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	for _, updatedID := range updatedIDs {
		_, err := fmt.Fprintf(file, "%d\n", updatedID)
		if err != nil {
			return err
		}
	}

	return nil
}
