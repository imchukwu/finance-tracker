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
		listCmd := flag.NewFlagSet("list", flag.ExitOnError)
		category := listCmd.String("category", "", "Filter by category")
		transType := listCmd.String("type", "", "Filter by type (income/expense)")
		month := listCmd.Int("month", 0, "Filter by month (1-12)")
		year := listCmd.Int("year", 0, "Filter by year")
		
		listCmd.Parse(os.Args[2:])
		
		transactions, err := transactionStorage.LoadTransactions()
		if err != nil {
			log.Fatal("Error loading transactions:", err)
		}
		
		// Apply filters
		filtered := filterTransactions(transactions, *category, *transType, *month, *year)
		
		for _, t := range filtered {
			fmt.Println(t.Display())
		}
	
	case "edit":
		editCmd := flag.NewFlagSet("edit", flag.ExitOnError)
		id := editCmd.String("id", "", "Transaction ID to edit")
		amount := editCmd.Float64("amount", 0, "New amount (0 keeps current)")
		category := editCmd.String("category", "", "New category (empty keeps current)")
		notes := editCmd.String("notes", "", "New notes (empty keeps current)")
		transType := editCmd.String("type", "", "New type (empty keeps current)")
		
		editCmd.Parse(os.Args[2:])
		
		if *id == "" {
			editCmd.PrintDefaults()
			os.Exit(1)
		}
		
		editTransaction(*id, *amount, *category, *notes, *transType)

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

	case "report":
		reportCmd := flag.NewFlagSet("report", flag.ExitOnError)
		month := reportCmd.Int("month", 0, "Month to report (1-12)")
		year := reportCmd.Int("year", 0, "Year to report")
		
		reportCmd.Parse(os.Args[2:])
		
		generateReport(*month, *year)

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

func editTransaction(id string, amount float64, category, notes, transType string) {
    // Get existing transaction
    t, err := transactionStorage.GetTransactionByID(id)
    if err != nil {
        log.Fatal("Error finding transaction:", err)
    }

    // Apply changes (only fields that were provided)
    if amount > 0 {
        t.Amount = amount
    }
    if category != "" {
        t.Category = category
    }
    if notes != "" {
        t.Notes = notes
    }
    if transType != "" {
        t.Type = transType
    }

    // Validate and save
    if err := t.Validate(); err != nil {
        log.Fatal("Invalid transaction data:", err)
    }

    // Delete old and save new
    if err := transactionStorage.DeleteTransaction(id); err != nil {
        log.Fatal("Error removing old transaction:", err)
    }
    if err := transactionStorage.SaveTransaction(t); err != nil {
        log.Fatal("Error saving updated transaction:", err)
    }

    fmt.Println("Successfully updated transaction:", t.Display())
}

func generateID() string {
	// Simple ID generation - we'll improve this later
	return fmt.Sprintf("%d", time.Now().UnixNano())
}

func filterTransactions(transactions []*models.Transaction, category, transType string, month, year int) []*models.Transaction {
    var filtered []*models.Transaction
    
    for _, t := range transactions {
        match := true
        
        if category != "" && t.Category != category {
            match = false
        }
        if transType != "" && t.Type != transType {
            match = false
        }
        if month > 0 && t.Date.Month() != time.Month(month) {
            match = false
        }
        if year > 0 && t.Date.Year() != year {
            match = false
        }
        
        if match {
            filtered = append(filtered, t)
        }
    }
    
    return filtered
}

func generateReport(month, year int) {
    transactions, err := transactionStorage.LoadTransactions()
    if err != nil {
        log.Fatal("Error loading transactions:", err)
    }
    
    // Filter if needed
    var filtered []*models.Transaction
    if month > 0 || year > 0 {
        filtered = filterTransactions(transactions, "", "", month, year)
    } else {
        filtered = transactions
    }
    
    // Calculate totals
    var incomeTotal, expenseTotal float64
    categoryTotals := make(map[string]float64)
    
    for _, t := range filtered {
        if t.Type == "income" {
            incomeTotal += t.Amount
        } else {
            expenseTotal += t.Amount
        }
        categoryTotals[t.Category] += t.Amount
    }
    
    // Print report
    fmt.Printf("\n=== Financial Report ===\n")
    if month > 0 {
        fmt.Printf("Period: %s %d\n", time.Month(month), year)
    }
    fmt.Printf("Total Income: $%.2f\n", incomeTotal)
    fmt.Printf("Total Expenses: $%.2f\n", expenseTotal)
    fmt.Printf("Net Balance: $%.2f\n\n", incomeTotal-expenseTotal)
    
    fmt.Println("By Category:")
    for category, total := range categoryTotals {
        fmt.Printf("- %s: $%.2f\n", category, total)
    }
    fmt.Println()
}