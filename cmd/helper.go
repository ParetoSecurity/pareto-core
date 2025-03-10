package cmd

import (
	"encoding/json"
	"net"
	"os"

	"github.com/ParetoSecurity/agent/claims"
	shared "github.com/ParetoSecurity/agent/shared"
	"github.com/caarlos0/log"
	"github.com/samber/lo"
	"github.com/spf13/cobra"
)

func runHelper() {
	// Get the socket from file descriptor 0
	file := os.NewFile(0, "socket")
	listener, err := net.FileListener(file)
	if err != nil {
		log.Error("Failed to create listener, not running in systemd context")
		os.Exit(1)
	}
	defer listener.Close()
	log.WithField("socket", shared.SocketPath).WithField("version", shared.Version).Info("Listening on socket")

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Debugf("Failed to accept connection: %v\n", err)
			continue
		}

		handleConnection(conn)
		break
	}
}

func handleConnection(conn net.Conn) {
	defer conn.Close()
	log.Info("Connection received")

	// Read input from connection
	decoder := json.NewDecoder(conn)
	var input map[string]string
	if err := decoder.Decode(&input); err != nil {
		log.Debugf("Failed to decode input: %v\n", err)
		return
	}
	uuid, ok := input["uuid"]
	if !ok {
		log.Debugf("UUID not found in input")
		return
	}
	log.Debugf("Received UUID: %s", uuid)

	status := map[string]bool{}
	for _, claim := range claims.All {
		for _, chk := range claim.Checks {
			if chk.IsRunnable() && chk.RequiresRoot() && uuid == chk.UUID() {
				log.Infof("Running check %s\n", chk.UUID())
				if chk.Run() != nil {
					log.Warnf("Failed to run check %s\n", chk.UUID())
					continue
				}
				log.Infof("Check %s completed\n", chk.UUID())
				status[chk.UUID()] = chk.Passed()
			}
		}
	}

	// Handle the request
	response, err := json.Marshal(status)
	if err != nil {
		log.Debugf("Failed to marshal response: %v\n", err)
		return
	}
	if _, err = conn.Write(response); err != nil {
		log.Debugf("Failed to write to connection: %v\n", err)
	}
}

var helperCmd = &cobra.Command{
	Use:   "helper [--socket]",
	Short: "A root helper",
	Long:  `A root helper that listens on a Unix domain socket and responds to authenticated requests.`,
	Run: func(cmd *cobra.Command, args []string) {

		socketFlag, _ := cmd.Flags().GetString("socket")
		if lo.IsNotEmpty(socketFlag) {
			shared.SocketPath = socketFlag
		}

		runHelper()
	},
}

func init() {
	rootCmd.AddCommand(helperCmd)
	helperCmd.Flags().Bool("socket", false, "socket path")
}
