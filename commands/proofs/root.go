package proofs

import (
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	lc "github.com/tendermint/light-client"
	"github.com/tendermint/light-client/commands"
	"github.com/tendermint/light-client/proofs"
	"github.com/tendermint/tendermint/rpc/client"
)

// RootCmd represents the base command when called without any subcommands
var RootCmd = &cobra.Command{
	Use:   "proof",
	Short: "Get and store merkle proofs for blockchain data",
	Long: `Proofs allows you to validate data and merkle proofs.

These proofs tie the data to a checkpoint, which is managed by "seeds".
Here we can validate these proofs and import/export them to prove specific
data to other peers as needed.
`,
}

type ProofCommander struct {
	node client.Client
	lc.Prover
	ProverFunc func(client.Client) lc.Prover
	proofs.Presenters
}

// Init uses configuration info to create a network connection
// as well as initializing the prover
func (p *ProofCommander) Init() {
	endpoint := viper.GetString(commands.NodeFlag)
	p.node = client.NewHTTP(endpoint, "/websockets")
	p.Prover = p.ProverFunc(p.node)
}

func (p ProofCommander) RegisterGet(parent *cobra.Command) {
	// we add each subcommand here, so we can register the
	// ProofCommander in one swoop
	parent.AddCommand(p.MakeGetCmd(
		true,
		"get",
		"Get a proof from the tendermint node",
	))
}
func (p ProofCommander) RegisterList(parent *cobra.Command) {
	parent.AddCommand(p.MakeGetCmd(
		false,
		"list",
		"Get a list proof from the tendermint node",
	))
}

const (
	heightFlag = "height"
	appFlag    = "app"
	keyFlag    = "key"
)
