package cmd

import (
	"fmt"
	

	"github.com/spf13/cobra"
)

var listCmd = &cobra.Command{
	Use: "list",
	Short: "List all available backups",
	Run: func(cmd *cobra.Command, args []string){
		fmt.Println("Fetching backup list ...")
	},
}

func init(){
	rootCmd.AddCommand(listCmd)
}