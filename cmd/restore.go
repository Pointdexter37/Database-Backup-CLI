package cmd

import (
	"fmt"
	

	"github.com/spf13/cobra"
)

var restoreCmd = &cobra.Command{
	Use: "restore",
	Short: "Restore a database from a backup file",
	Run : func(cmd *cobra.Command, args []string){
		fmt.Println("Restore logic initiated")
	},
}
func init(){
	rootCmd.AddCommand(restoreCmd)
}




