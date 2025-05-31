package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/imchukwu/finance-tracker/models"
	"github.com/imchukwu/finance-tracker/storage"
)

func main() {

	if len(os.Args) < 2 {
		printHelp()
		return
		// fmt.Println("Available command:")
		// fmt.Println(" add - Add a new transaction")
		// fmt.Println(" list - List all transaction")

	}

	switch os.Args[1] {
	case "add":
		addCmd := flag.NewFlagSet("add", flag.ExitOnError)
		amount := addCmd.Float64("amount", 0.0, "Transaction amount")
		category := addCmd.String("category", "", "Transaction category")
		notes := addCmd.String("notes", "", "Additional notes")
		transType := addCmd.String("type", "expense", "Transaction type (income/expense)")

		addCmd.Parse(os.Args[2:])

		if *amount == 0.0 || *category == "" {
			addCmd.PrintDefaults()
			os.Exit(1)
		}

		addTransaction(*amount, *category, *notes, *transType)

	case "list":
		transactions, err := transactionStorage.LoadTransactions()
		if err != nil {
			log.Fatal("Failed to load transactions:", err)
		}

		for _, t := range transactions {
			fmt.Println(t.Display())
		}
	
	case "delete":
		deleteCmd := flag.NewFlagSet("delete", flag.ExitOnError)
		id := deleteCmd.String("id", "", "Transaction ID to delete")
		deleteCmd.Parse(os.Args[2:])
		
		if *id == "" {
			deleteCmd.PrintDefaults()
			os.Exit(1)
		}
		
		if err := transactionStorage.DeleteTransaction(*id); err != nil {
			log.Fatal("Failed to delete transaction:", err)
		}
		fmt.Println("Successfully deleted transaction")

	default:
		fmt.Printf("Unknown command: %s\n", os.Args[1])
		printHelp()
	}
}

func printHelp() {
	fmt.Println(`Personal Finance Trancker
	Usage:
	   add - Add new transaction
	      flags:
		     -amount float    Transaction amount
			 -category string Transaction category
			 -notes strng     Additional notes (optional)
			 -type string     Transaction type (income/expense)	(default "expense")
	
	list - Lst all transactions`)
}

var (
	transactionStorage storage.Storage
)

func init() {
    var err error
    transactionStorage, err = storage.NewJSONStorage("transactions.json")
    if err != nil {
        log.Fatalf("Failed to initialize storage: %v", err)
    }
}

func addTransaction(amount float64, category, notes, transType string) {
	// Validate transaction type
	if transType != "income" && transType != "expense" {
		log.Fatal("Transaction type must be either 'income' or 'expense'")
	}

    t := models.NewTransaction(
        generateID(),
        time.Now(),
        amount,
        category,
        notes,
        transType,
    )

    if err := t.Validate(); err != nil {
        log.Fatalf("Invalid transaction: %v", err)
    }

    if err := transactionStorage.SaveTransaction(t); err != nil {
        log.Fatalf("Failed to save transaction: %v", err)
    }

    fmt.Println("Successfully added transaction:", t.Display())

}

func generateID() string {
	// Simple ID generation - we'll improve this later
	return fmt.Sprintf("%d", time.Now().UnixNano())
}
