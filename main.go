package main

import (
    "bufio"
    "context"
    "flag"
    "fmt"
    "os"
    "strings"

    "github.com/Azure/azure-sdk-for-go/profiles/latest/storage/mgmt/storage"
    "github.com/Azure/go-autorest/autorest/azure/auth"
)

const (
    colorRed   = "\033[31m"
    colorReset = "\033[0m"
)

func main() {
    var (
        subscriptionID string
        list           bool
    )

    flag.StringVar(&subscriptionID, "subscription", "", "Azure Subscription ID (optional)")
    flag.BoolVar(&list, "list", false, "List Azure storage accounts")
    flag.Parse()

    if subscriptionID == "" {
        fmt.Printf(colorRed+"Subscription ID was not provided."+colorReset+"\nDo you want to enter it now? (Y/n): ")
        reader := bufio.NewReader(os.Stdin)
        response, _ := reader.ReadString('\n')

        if strings.TrimSpace(strings.ToLower(response)) == "y" {
            fmt.Print("Enter Subscription ID: ")
            subscriptionID, _ = reader.ReadString('\n')
            subscriptionID = strings.TrimSpace(subscriptionID)
        }

        if subscriptionID == "" {
            fmt.Println(colorRed + "No Subscription ID provided. Exiting." + colorReset)
            os.Exit(1)
        }
    }

    if list {
        listStorageAccounts(subscriptionID)
    } else {
        fmt.Println("No action specified.")
        flag.PrintDefaults()
    }
}

func listStorageAccounts(subscriptionID string) {
    authorizer, err := auth.NewAuthorizerFromCLI()
    if err != nil {
        fmt.Printf(colorRed+"Authentication failed: %v\n"+colorReset, err)
        os.Exit(1)
    }

    storageAccountsClient := storage.NewAccountsClient(subscriptionID)
    storageAccountsClient.Authorizer = authorizer

    ctx := context.Background()
    result, err := storageAccountsClient.List(ctx)
    if err != nil {
        fmt.Printf(colorRed+"Failed to list storage accounts: %v\n"+colorReset, err)
        os.Exit(1)
    }

    for _, account := range result.Values() {
        fmt.Printf("Account Name: %s\n", *account.Name)
    }
}
