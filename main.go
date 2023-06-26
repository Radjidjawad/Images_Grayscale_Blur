package main

import (
	"fmt"
	"image"
	"sync"
	"time"

	"github.com/disintegration/imaging"
)

func main() {

	src, err := imaging.Open("input.jpg")
	if err != nil {
		fmt.Printf("Erreur lors du chargement de l'image : %v\n", err)
		return
	}

	wg := sync.WaitGroup{}

	startTime := time.Now()
	for i := 0; i < 4; i++ {
		wg.Add(1)
		go processImageWithWaitGroup(src, &wg)
	}
	wg.Wait()
	duration := time.Since(startTime)
	fmt.Printf("Durée avec WaitGroup : %s\n", duration)

	startTime = time.Now()
	done := make(chan bool)
	for i := 0; i < 4; i++ {
		go processImageWithChannel(src, done)
	}
	for i := 0; i < 4; i++ {
		<-done
	}
	duration = time.Since(startTime)
	fmt.Printf("Durée avec Channel : %s\n", duration)
}

func processImageWithWaitGroup(src image.Image, wg *sync.WaitGroup) {
	defer wg.Done()

	grayscale := imaging.Grayscale(src)

	blurred := imaging.Blur(grayscale, 5.0)

	err := imaging.Save(blurred, fmt.Sprintf("output_wg_%d.jpg", time.Now().UnixNano()))
	if err != nil {
		fmt.Printf("Erreur lors de la sauvegarde de l'image : %v\n", err)
	}
}

func processImageWithChannel(src image.Image, done chan<- bool) {

	grayscale := imaging.Grayscale(src)

	blurred := imaging.Blur(grayscale, 5.0)

	err := imaging.Save(blurred, fmt.Sprintf("output_ch_%d.jpg", time.Now().UnixNano()))
	if err != nil {
		fmt.Printf("Erreur lors de la sauvegarde de l'image : %v\n", err)
	}

	done <- true
}
