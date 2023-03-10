package main

import (
	"os"
	"regexp"
	"strings"

	"gopkg.in/yaml.v3"
)

func ReadYml(filename string) (map[string]interface{}, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	decoder := yaml.NewDecoder(file)
	var config map[string]interface{}
	err = decoder.Decode(&config)
	if err != nil {
		return nil, err
	}
	return config, nil
}

func ProcessCommandOne(response string) ([]Container, error) {

	lines := strings.Split(response, "\n")

	lines = lines[1:]
	for i, line := range lines {
		lines[i] = regexp.MustCompile(" {2}").ReplaceAllString(line, ",")
		lines[i] = regexp.MustCompile(",+").ReplaceAllString(lines[i], ",")
	}

	containers := []Container{}
	container := Container{}
	for _, line := range lines {
		if line == "" {
			continue
		}
		split := strings.Split(line, ",")

		if len(split) <= 6 {
			container = Container{
				ID:      split[0],
				Created: split[3],
				Status:  split[4],
				Ports:   "Not known",
				Names:   split[5],
			}
		} else {
			container = Container{
				ID:      split[0],
				Created: split[3],
				Status:  split[4],
				Ports:   split[5],
				Names:   split[7],
			}
		}
		containers = append(containers, container)
	}

	return containers, nil
}

func ProcessCommandTwo(response string) ([]Images, error) {

	lines := strings.Split(response, "\n")
	lines = lines[1:]
	for i, line := range lines {
		lines[i] = regexp.MustCompile(" {2}").ReplaceAllString(line, ",")
		lines[i] = regexp.MustCompile(",+").ReplaceAllString(lines[i], ",")
	}

	images := []Images{}
	image := Images{}
	for _, line := range lines {
		if line == "" {
			continue
		}
		split := strings.Split(line, ",")
		image = Images{
			Repository: split[0],
			Tag:        split[1],
			ImageID:    split[2],
			Created:    split[3],
			Size:       split[4],
		}
		images = append(images, image)
	}

	return images, nil
}
