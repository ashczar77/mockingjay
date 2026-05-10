package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/ashczar77/mockingjay/internal/audio"
	"github.com/ashczar77/mockingjay/internal/classifier"
	"github.com/ashczar77/mockingjay/internal/reporter"
	"github.com/ashczar77/mockingjay/internal/twilio"
	"github.com/spf13/cobra"
)

var (
	calltestTo      string
	calltestWebhook string
	calltestAPIURL  string
	calltestExpect  string
)

var calltestCmd = &cobra.Command{
	Use:   "calltest",
	Short: "Call → transcribe → validate → report in one command",
	Long:  `Makes a real phone call, transcribes the recording, validates the transcript, and reports results to the dashboard.`,
	Run: func(cmd *cobra.Command, args []string) {
		if err := runCallTest(); err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
	},
}

func init() {
	calltestCmd.Flags().StringVar(&calltestTo, "to", "", "Phone number to call")
	calltestCmd.Flags().StringVar(&calltestWebhook, "webhook", "", "TwiML webhook URL for call instructions")
	calltestCmd.Flags().StringVar(&calltestAPIURL, "api-url", "", "Backend URL to report results (e.g. http://localhost:8080)")
	calltestCmd.Flags().StringVar(&calltestExpect, "expect", "", "Expected phrase in transcript (validates call quality)")
	calltestCmd.MarkFlagRequired("to")
	calltestCmd.MarkFlagRequired("webhook")
}

func runCallTest() error {
	sid := os.Getenv("TWILIO_ACCOUNT_SID")
	token := os.Getenv("TWILIO_AUTH_TOKEN")
	from := os.Getenv("TWILIO_FROM_NUMBER")
	dgKey := os.Getenv("DEEPGRAM_API_KEY")

	if sid == "" || token == "" || from == "" {
		return fmt.Errorf("TWILIO_ACCOUNT_SID, TWILIO_AUTH_TOKEN, TWILIO_FROM_NUMBER must be set")
	}
	if dgKey == "" {
		return fmt.Errorf("DEEPGRAM_API_KEY must be set")
	}

	fmt.Println("🐦 MockingJay - Call Quality Test")
	fmt.Println()

	// Step 1: Make the call
	fmt.Printf("📞 Step 1/3: Calling %s...\n", calltestTo)
	client := twilio.NewClient(twilio.Config{
		AccountSID: sid,
		AuthToken:  token,
		FromNumber: from,
		ToNumber:   calltestTo,
		WebhookURL: calltestWebhook,
	})
	result, err := client.MakeCall(true)
	if err != nil {
		return fmt.Errorf("call failed: %w", err)
	}
	fmt.Printf("   ✓ Call completed (SID: %s, Duration: %s)\n\n", result.CallSID, result.Duration)

	// Step 2: Download and transcribe the recording
	fmt.Println("🎙️  Step 2/3: Transcribing recording...")
	recordingURL, err := twilio.GetRecordingURL(sid, token, result.CallSID)
	if err != nil {
		return fmt.Errorf("failed to fetch recording: %w", err)
	}

	recordingPath := filepath.Join(os.TempDir(), fmt.Sprintf("mockingjay-%s.wav", result.CallSID))
	if err := audio.SaveAudio(recordingURL, recordingPath); err != nil {
		return fmt.Errorf("failed to download recording: %w", err)
	}
	defer os.Remove(recordingPath)

	t := audio.NewDeepgramTranscriber(dgKey)
	transcript, err := t.Transcribe(recordingPath)
	if err != nil {
		return fmt.Errorf("transcription failed: %w", err)
	}
	fmt.Printf("   ✓ Transcript: %q\n", transcript.Text)
	fmt.Printf("   ✓ Confidence: %.1f%%, Duration: %.1fs\n\n", transcript.Confidence*100, transcript.Duration)

	// Step 3: Validate using LLM classifier if available, else string match
	fmt.Println("✅ Step 3/3: Validating...")
	passed := true
	validationMsg := "no validation criteria specified"

	if calltestExpect != "" {
		c := classifier.New()
		if c != nil {
			result2, err := c.ClassifyTranscript(transcript.Text, calltestExpect)
			if err == nil {
				passed = result2.Confidence >= 0.7
				if passed {
					validationMsg = fmt.Sprintf("transcript matches intent %q (confidence: %.0f%% — %s)",
						calltestExpect, result2.Confidence*100, result2.Reasoning)
				} else {
					validationMsg = fmt.Sprintf("transcript does not match intent %q (confidence: %.0f%% — %s)",
						calltestExpect, result2.Confidence*100, result2.Reasoning)
				}
			} else {
				passed = strings.Contains(strings.ToLower(transcript.Text), strings.ToLower(calltestExpect))
				validationMsg = fmt.Sprintf("string match for %q: %v", calltestExpect, passed)
			}
		} else {
			// No LLM configured, fall back to string match
			passed = strings.Contains(strings.ToLower(transcript.Text), strings.ToLower(calltestExpect))
			if passed {
				validationMsg = fmt.Sprintf("transcript contains expected phrase %q", calltestExpect)
			} else {
				validationMsg = fmt.Sprintf("transcript missing expected phrase %q", calltestExpect)
			}
		}
	}

	if passed {
		fmt.Printf("   ✓ PASS — %s\n", validationMsg)
	} else {
		fmt.Printf("   ✗ FAIL — %s\n", validationMsg)
	}

	// Report to backend
	if calltestAPIURL != "" {
		r := reporter.NewClient(calltestAPIURL)

		// Save transcription
		r.ReportTranscription(reporter.Transcription{
			CallSID:   result.CallSID,
			AudioPath: recordingPath,
			Text:      transcript.Text,
			Confidence: transcript.Confidence,
			Duration:  transcript.Duration,
		})

		// Save test result
		errMsg := ""
		if !passed {
			errMsg = validationMsg
		}
		r.ReportRaw(reporter.RawResult{
			Scenario:  fmt.Sprintf("calltest:%s", calltestTo),
			Passed:    passed,
			LatencyMs: result.Duration.Milliseconds(),
			Error:     errMsg,
		})
		fmt.Printf("\n📤 Results saved to dashboard (%s)\n", calltestAPIURL)
	}

	fmt.Println()
	if !passed {
		return fmt.Errorf("call quality validation failed: %s", validationMsg)
	}
	fmt.Printf("✨ Call quality test passed! (completed at %s)\n", time.Now().Format("15:04:05"))
	return nil
}
