package main

import (
	"fmt"
	"log"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// User model
type User struct {
	gorm.Model
	PhoneNumber string `gorm:"type:varchar(255);uniqueIndex"`
	Password    string
	Balance     float64
}

// Transaction model
type Transaction struct {
	gorm.Model
	UserID      uint
	PhoneNumber string
	Amount      float64
	Type        string // "top-up" or "transfer"
}

var db *gorm.DB
var loggedInUserID uint

func main() {
	var err error
	var dsn = fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local", "root", "pasword", "localhost", 3306, "project")
	db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Error connecting to database: %v", err)
	}

	// Auto migrate the schema
	err = db.AutoMigrate(&User{}, &Transaction{})
	if err != nil {
		fmt.Println("Failed to perform AutoMigrate:", err)
		return
	}

	var input string
	for input != "0" {
		fmt.Println("Menu:")
		fmt.Println("1. Register")
		fmt.Println("2. Login")
		fmt.Println("3. Read Account")
		fmt.Println("4. Update Account")
		fmt.Println("5. Delete Account")
		fmt.Println("6. Top-up")
		fmt.Println("7. Transfer")
		fmt.Println("8. History Top-up")
		fmt.Println("9. History Transfer")
		fmt.Println("10. Melihat profil user lain")
		fmt.Println("0. Exit")
		fmt.Print("Input: ")
		fmt.Scanln(&input)

		switch input {
		case "1":
			register()
		case "2":
			login()
		case "3":
			read()
		case "4":
			update()
		case "5":
			delete()
		case "6":
			topUp()
		case "7":
			transfer()
		case "8":
			historyTopUp()
		case "9":
			historyTransfer()
		case "10":
			viewProfile()
		case "0":
			fmt.Println("Terimakasih telah bertransaksi")
		default:
			fmt.Println("Pilihan tidak valid")
		}
	}
}

func register() {
	var phoneNumber, password string
	fmt.Print("Masukkan nomor telepon: ")
	fmt.Scanln(&phoneNumber)
	fmt.Print("Masukkan password: ")
	fmt.Scanln(&password)

	user := User{PhoneNumber: phoneNumber, Password: password}
	result := db.Create(&user)
	if result.Error != nil {
		fmt.Println("Gagal mendaftar:", result.Error)
		return
	}
	fmt.Println("Berhasil mendaftar")
}

func login() {
	var phoneNumber, password string
	fmt.Print("Masukkan nomor telepon: ")
	fmt.Scanln(&phoneNumber)
	fmt.Print("Masukkan password: ")
	fmt.Scanln(&password)

	var user User
	result := db.Where("phone_number = ? AND password = ?", phoneNumber, password).First(&user)
	if result.Error != nil {
		fmt.Println("Gagal login:", result.Error)
		return
	}
	fmt.Println("Berhasil login")
	loggedInUserID = user.ID
}

func read() {
	if loggedInUserID == 0 {
		fmt.Println("Silakan login terlebih dahulu")
		return
	}

	var user User
	result := db.First(&user, loggedInUserID)
	if result.Error != nil {
		fmt.Println("Gagal melihat profil user:", result.Error)
		return
	}

	fmt.Printf("Nomor Telepon: %s, Saldo: %f\n", user.PhoneNumber, user.Balance)
}

func update() {
	if loggedInUserID == 0 {
		fmt.Println("Silakan login terlebih dahulu")
		return
	}

	var newPassword string
	fmt.Print("Masukkan password baru: ")
	fmt.Scanln(&newPassword)

	var user User
	result := db.First(&user, loggedInUserID)
	if result.Error != nil {
		fmt.Println("Gagal update:", result.Error)
		return
	}

	if loggedInUserID != user.ID {
		fmt.Println("Anda tidak memiliki izin untuk mengubah akun ini")
		return
	}

	user.Password = newPassword
	db.Save(&user)
	fmt.Println("Berhasil update password")
}

func delete() {
	if loggedInUserID == 0 {
		fmt.Println("Silakan login terlebih dahulu")
		return
	}

	var user User
	result := db.First(&user, loggedInUserID)
	if result.Error != nil {
		fmt.Println("Gagal menghapus akun:", result.Error)
		return
	}

	if loggedInUserID != user.ID {
		fmt.Println("Anda tidak memiliki izin untuk menghapus akun ini")
		return
	}

	db.Unscoped().Delete(&user)
	fmt.Println("Akun berhasil dihapus")
}

func topUp() {
	var phoneNumber string
	var amount float64
	fmt.Print("Masukkan nomor telepon: ")
	fmt.Scanln(&phoneNumber)
	fmt.Print("Masukkan jumlah top-up: ")
	fmt.Scanln(&amount)

	var user User
	result := db.Where("phone_number = ?", phoneNumber).First(&user)
	if result.Error != nil {
		fmt.Println("Gagal top-up:", result.Error)
		return
	}

	user.Balance += amount
	db.Save(&user)

	transaction := Transaction{UserID: user.ID, PhoneNumber: phoneNumber, Amount: amount, Type: "top-up"}
	db.Create(&transaction)
	fmt.Println("Top-up berhasil")
}

func transfer() {
	var senderPhoneNumber, receiverPhoneNumber string
	var amount float64
	fmt.Print("Masukkan nomor telepon pengirim: ")
	fmt.Scanln(&senderPhoneNumber)
	fmt.Print("Masukkan nomor telepon penerima: ")
	fmt.Scanln(&receiverPhoneNumber)
	fmt.Print("Masukkan jumlah transfer: ")
	fmt.Scanln(&amount)

	var sender, receiver User
	result := db.Where("phone_number = ?", senderPhoneNumber).First(&sender)
	if result.Error != nil {
		fmt.Println("Gagal transfer:", result.Error)
		return
	}

	result = db.Where("phone_number = ?", receiverPhoneNumber).First(&receiver)
	if result.Error != nil {
		fmt.Println("Gagal transfer:", result.Error)
		return
	}

	if sender.Balance < amount {
		fmt.Println("Saldo pengirim tidak mencukupi")
		return
	}

	sender.Balance -= amount
	receiver.Balance += amount
	db.Save(&sender)
	db.Save(&receiver)

	senderTransaction := Transaction{UserID: sender.ID, PhoneNumber: senderPhoneNumber, Amount: -amount, Type: "transfer"}
	receiverTransaction := Transaction{UserID: receiver.ID, PhoneNumber: receiverPhoneNumber, Amount: amount, Type: "transfer"}
	db.Create(&senderTransaction)
	db.Create(&receiverTransaction)
	fmt.Println("Transfer berhasil")
}

func historyTopUp() {
	var phoneNumber string
	fmt.Print("Masukkan nomor telepon: ")
	fmt.Scanln(&phoneNumber)

	var user User
	result := db.Where("phone_number = ?", phoneNumber).First(&user)
	if result.Error != nil {
		fmt.Println("Gagal menampilkan history top-up:", result.Error)
		return
	}

	var transactions []Transaction
	db.Where("user_id = ? AND type = ?", user.ID, "top-up").Find(&transactions)

	fmt.Println("History Top-up:")
	for _, transaction := range transactions {
		fmt.Printf("ID: %d, Nomor Telepon: %s, Jumlah: %f\n", transaction.ID, transaction.PhoneNumber, transaction.Amount)
	}
}

func historyTransfer() {
	var phoneNumber string
	fmt.Print("Masukkan nomor telepon: ")
	fmt.Scanln(&phoneNumber)

	var user User
	result := db.Where("phone_number = ?", phoneNumber).First(&user)
	if result.Error != nil {
		fmt.Println("Gagal menampilkan history transfer:", result.Error)
		return
	}

	var transactions []Transaction
	db.Where("user_id = ? AND type = ?", user.ID, "transfer").Find(&transactions)

	fmt.Println("History Transfer:")
	for _, transaction := range transactions {
		fmt.Printf("ID: %d, Nomor Telepon: %s, Jumlah: %f\n", transaction.ID, transaction.PhoneNumber, transaction.Amount)
	}
}

func viewProfile() {
	var phoneNumber string
	fmt.Print("Masukkan nomor telepon user yang ingin dilihat profilnya: ")
	fmt.Scanln(&phoneNumber)

	var user User
	result := db.Where("phone_number = ?", phoneNumber).First(&user)
	if result.Error != nil {
		fmt.Println("Gagal melihat profil user:", result.Error)
		return
	}

	fmt.Printf("Nomor Telepon: %s, Saldo: %f\n", user.PhoneNumber, user.Balance)
}
