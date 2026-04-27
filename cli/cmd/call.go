package cmd

import (
	"fmt"
	"os"

	"github.com/ashczar77/mockingjay/internal/twilio"
	"github.com/spf13/cobra"
)

var (
	twilioSID    string
	twilioToken  string
	twilioFrom   string
	twilioTo     string
	twilioWebhook string
	callRecord   bool
)

var callCmd = &cobra.Command{
	Use:   "call",
	Short: "Make a real phone call via Twilio",
	Long:  `Initiate an outbound call through Twilio to test your voice AI agent.`,
	Run: func(cmd *cobra.Command, args []string) {
		if err := runCall(); err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
	},
}

func init() {
	callCmd.Flags().StringVar(&twilioSID, "account-sid", os.Getenv("TWILIO_ACCOUNT_SID"), "Twilio Account SID")
	callCmd.Flags().StringVar(&twilioToken, "auth-token", os.Getenv("TWILIO_AUTH_TOKEN"), "Twilio Auth Token")
	callCmd.Flags().StringVar(&twilioFrom, "from", os.Getenv("TWILIO_FROM_NUMBER"), "Twilio phone number to call from")
	callCmd.Flags().StringVar(&twilioTo, "to", "", "Phone number to call")
	callCmd.Flags().StringVar(&twilioWebhook, "webhook", "", "TwiML webhook URL for call instructions")
	callCmd.Flags().BoolVar(&callRecord, "record", false, "Record the call")
	callCmd.MarkFlagRequired("to")
	callCmd.MarkFlagRequired("webhook")
}

func runCall() error {
	if twilioSID == "" || twilioToken == "" || twilioFrom == "" {
		return fmt.Errorf("Twilio credentials required: --account-sid, --auth-token, --from (or set TWILIO_ACCOUNT_SID, TWILIO_AUTH_TOKEN, TWILIO_FROM_NUMBER)")
	}

	fmt.Println("🐦 MockingJay - Real Call Test")
	fmt.Println()
	fmt.Printf("📞 Calling %s from %s...\n", twilioTo, twilioFrom)
	if callRecord {
		fmt.Println("🔴 Recording enabled")
	}
	fmt.Println()

	client := twilio.NewClient(twilio.Config{
		AccountSID: twilioSID,
		AuthToken:  twilioToken,
		FromNumber: twilioFrom,
		ToNumber:   twilioTo,
		WebhookURL: twilioWebhook,
	})

	result, err := client.MakeCall(callRecord)
	if err != nil {
		return err
	}

	fmt.Printf("📊 Call Result:\n")
	fmt.Printf("  SID:      %s\n", result.CallSID)
	fmt.Printf("  Status:   %s\n", result.Status)
	fmt.Printf("  Duration: %s\n", result.Duration)
	if result.RecordURL != "" {
		fmt.Printf("  Recording: %s\n", result.RecordURL)
	}

	if result.Status != "completed" {
		return fmt.Errorf("call did not complete successfully (status: %s)", result.Status)
	}

	fmt.Println("\n✨ Call completed successfully!")
	return nil
}
