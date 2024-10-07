package main

import (
	"encoding/base64"
	"errors"
	"fmt"
	"regexp"
	"strconv"
	"sync"
	"time"
)

// Struct untuk merepresentasikan item menu
type MenuItem struct {
	Name     string
	Price    float64
	Quantity int
}

// Interface untuk item yang dapat dipesan
type Orderable interface {
	GetPrice() float64
	GetName() string
}

// Method untuk MenuItem agar memenuhi interface Orderable
func (m MenuItem) GetPrice() float64 {
	return m.Price
}

func (m MenuItem) GetName() string {
	return m.Name
}

// Method untuk menambah kuantitas pesanan (menggunakan pointer)
func (m *MenuItem) AddQuantity(qty int) {
	m.Quantity += qty
}

// Fungsi untuk menambah item ke menu
func AddMenuItem(name string, price float64, menu map[string]*MenuItem) {
	menu[name] = &MenuItem{name, price, 0}
}

// Fungsi untuk memproses pesanan (menggunakan goroutine)
func processOrder(orderChan chan string, wg *sync.WaitGroup) {
	defer wg.Done()
	for order := range orderChan {
		fmt.Printf("Processing order: %s\n", order)
		time.Sleep(2 * time.Second) // simulasi waktu pemrosesan
	}
}

// Fungsi untuk menangani panic dan recover
func handleError() {
	if r := recover(); r != nil {
		fmt.Println("Recovered from error:", r)
	}
}

// Validasi input harga menggunakan regular expression
func validatePrice(price string) error {
	regex := regexp.MustCompile(`^\d+(\.\d{1,2})?$`)
	if !regex.MatchString(price) {
		return errors.New("invalid price format")
	}
	return nil
}

// Fungsi utama
func main() {
	defer fmt.Println("Program selesai")
	defer handleError()

	menu := make(map[string]*MenuItem)
	var wg sync.WaitGroup

	// Input nama dan harga item
	var name, priceStr string
	fmt.Print("Masukkan nama item: ")
	fmt.Scan(&name)
	fmt.Print("Masukkan harga item: ")
	fmt.Scan(&priceStr)

	// Validasi harga
	err := validatePrice(priceStr)
	if err != nil {
		panic("Harga tidak valid")
	}

	// Konversi harga ke float
	price, err := strconv.ParseFloat(priceStr, 64)
	if err != nil {
		panic(err)
	}

	// Tambahkan item ke menu
	AddMenuItem(name, price, menu)

	// Simulasi input pesanan
	orderChan := make(chan string, 3)
	wg.Add(1)
	go processOrder(orderChan, &wg)

	// Kirim pesanan melalui channel
	for i := 1; i <= 3; i++ {
		orderChan <- fmt.Sprintf("Pesanan %d: %s", i, name)
	}

	// Menutup channel dan menunggu goroutine selesai
	close(orderChan)
	wg.Wait()

	// Encode detail pesanan menggunakan base64
	encoded := base64.StdEncoding.EncodeToString([]byte("Order detail: " + name))
	fmt.Println("Encoded order detail:", encoded)
}
