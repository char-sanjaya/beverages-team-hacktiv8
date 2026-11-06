package helper

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

var reader = bufio.NewReader(os.Stdin)

// ReadString membaca input teks dan menolak input kosong
func ReadString(prompt string) string {
	for {
		fmt.Print(prompt)
		input, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println("[!] Terjadi kesalahan pembacaan input, silakan coba lagi.")
			continue
		}
		input = strings.TrimSpace(input)
		if input == "" {
			fmt.Println("[!] Input tidak boleh kosong!")
			continue
		}
		return input
	}
}

// ReadPositiveInt membaca bilangan bulat > 0 (menolak huruf, desimal, nol, dan minus)
func ReadPositiveInt(prompt string) int {
	for {
		fmt.Print(prompt)
		input, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println("[!] Gagal membaca input.")
			continue
		}
		input = strings.TrimSpace(input)

		val, err := strconv.Atoi(input)
		if err != nil {
			fmt.Println("[!] Input ditolak! Harus berupa angka bulat (tidak boleh huruf atau simbol).")
			continue
		}
		if val <= 0 {
			fmt.Println("[!] Input ditolak! Angka tidak boleh minus atau nol.")
			continue
		}
		return val
	}
}

// ReadPositiveFloat membaca nominal uang > 0 (menolak huruf, nol, dan angka minus)
func ReadPositiveFloat(prompt string) float64 {
	for {
		fmt.Print(prompt)
		input, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println("[!] Gagal membaca input.")
			continue
		}
		input = strings.TrimSpace(input)

		val, err := strconv.ParseFloat(input, 64)
		if err != nil {
			fmt.Println("[!] Input ditolak! Nominal harus berupa angka (tidak boleh huruf).")
			continue
		}
		if val <= 0 {
			fmt.Println("[!] Input ditolak! Nominal tidak boleh minus atau nol.")
			continue
		}
		return val
	}
}
