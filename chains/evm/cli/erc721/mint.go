package erc721

import (
	"fmt"
	"github.com/ChainSafe/chainbridge-core/chains/evm/calls/contracts/erc721"
	"github.com/ChainSafe/chainbridge-core/chains/evm/calls/evmtransaction"
	"github.com/ChainSafe/chainbridge-core/chains/evm/cli/initialize"
	"github.com/ChainSafe/chainbridge-core/util"
	"math/big"

	"github.com/ChainSafe/chainbridge-core/chains/evm/calls/transactor"
	"github.com/ChainSafe/chainbridge-core/chains/evm/cli/flags"
	"github.com/ChainSafe/chainbridge-core/chains/evm/cli/logger"
	"github.com/ethereum/go-ethereum/common"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
)

var mintCmd = &cobra.Command{
	Use:   "mint",
	Short: "Mint an ERC721 token",
	Long:  "The mint subcommand mints a token on an ERC721 mintable contract",
	PreRun: func(cmd *cobra.Command, args []string) {
		logger.LoggerMetadata(cmd.Name(), cmd.Flags())
	},
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		return util.CallPersistentPreRun(cmd, args)
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		c, err := initialize.InitializeClient(url, senderKeyPair)
		if err != nil {
			return err
		}
		t, err := initialize.InitializeTransactor(gasPrice, evmtransaction.NewTransaction, c)
		if err != nil {
			return err
		}
		return MintCmd(cmd, args, erc721.NewErc721Contract(c, erc721Addr, t))
	},
	Args: func(cmd *cobra.Command, args []string) error {
		err := ValidateMintFlags(cmd, args)
		if err != nil {
			return err
		}

		err = ProcessMintFlags(cmd, args)
		return err
	},
}

func init() {
	BindMintFlags(mintCmd)
}

func BindMintFlags(cmd *cobra.Command) {
	cmd.Flags().StringVar(&Erc721Address, "contract", "", "ERC721 contract address")
	cmd.Flags().StringVar(&DstAddress, "recipient", "", "Recipient address")
	cmd.Flags().StringVar(&TokenId, "token", "", "ERC721 token ID")
	cmd.Flags().StringVar(&Metadata, "metadata", "", "ERC721 token metadata")
	flags.MarkFlagsAsRequired(cmd, "contract", "recipient", "token", "metadata")
}

func ValidateMintFlags(cmd *cobra.Command, args []string) error {
	if !common.IsHexAddress(Erc721Address) {
		return fmt.Errorf("invalid ERC721 contract address %s", Erc721Address)
	}
	if !common.IsHexAddress(DstAddress) {
		return fmt.Errorf("invalid recipient address %s", DstAddress)
	}
	return nil
}

func ProcessMintFlags(cmd *cobra.Command, args []string) error {
	erc721Addr = common.HexToAddress(Erc721Address)

	if !common.IsHexAddress(DstAddress) {
		dstAddress = senderKeyPair.CommonAddress()
	} else {
		dstAddress = common.HexToAddress(DstAddress)
	}

	var ok bool
	if tokenId, ok = big.NewInt(0).SetString(TokenId, 10); !ok {
		return fmt.Errorf("invalid token id value")
	}

	return err
}

func MintCmd(cmd *cobra.Command, args []string, erc721Contract *erc721.ERC721Contract) error {
	_, err = erc721Contract.Mint(
		tokenId, Metadata, dstAddress, transactor.TransactOptions{GasLimit: gasLimit},
	)
	if err != nil {
		return err
	}

	log.Info().Msgf("%v token minted", tokenId)
	return err
}