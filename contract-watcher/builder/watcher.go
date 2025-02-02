package builder

import (
	"time"

	solana_go "github.com/gagliardetto/solana-go"
	"github.com/wormhole-foundation/wormhole-explorer/contract-watcher/config"
	"github.com/wormhole-foundation/wormhole-explorer/contract-watcher/internal/ankr"
	"github.com/wormhole-foundation/wormhole-explorer/contract-watcher/internal/aptos"
	"github.com/wormhole-foundation/wormhole-explorer/contract-watcher/internal/evm"
	"github.com/wormhole-foundation/wormhole-explorer/contract-watcher/internal/metrics"
	"github.com/wormhole-foundation/wormhole-explorer/contract-watcher/internal/solana"
	"github.com/wormhole-foundation/wormhole-explorer/contract-watcher/internal/terra"
	"github.com/wormhole-foundation/wormhole-explorer/contract-watcher/storage"
	"github.com/wormhole-foundation/wormhole-explorer/contract-watcher/watcher"
	"go.uber.org/ratelimit"
	"go.uber.org/zap"
)

func CreateEVMWatcher(rateLimit int, chainURL string, wb config.WatcherBlockchainAddresses, repo *storage.Repository,
	metrics metrics.Metrics, logger *zap.Logger) watcher.ContractWatcher {
	evmLimiter := ratelimit.New(rateLimit, ratelimit.Per(time.Second))
	ankrClient := ankr.NewAnkrSDK(chainURL, evmLimiter, metrics)
	params := watcher.EVMParams{ChainID: wb.ChainID, Blockchain: wb.Name, SizeBlocks: wb.SizeBlocks,
		WaitSeconds: wb.WaitSeconds, InitialBlock: wb.InitialBlock, MethodsByAddress: wb.MethodsByAddress}
	return watcher.NewEVMWatcher(ankrClient, repo, params, metrics, logger)
}

func CreateSolanaWatcher(rateLimit int, chainURL string, wb config.WatcherBlockchain, logger *zap.Logger, repo *storage.Repository, metrics metrics.Metrics) watcher.ContractWatcher {
	contractAddress, err := solana_go.PublicKeyFromBase58(wb.Address)
	if err != nil {
		logger.Fatal("failed to parse solana contract address", zap.Error(err))
	}
	solanaLimiter := ratelimit.New(rateLimit, ratelimit.Per(time.Second))
	solanaClient := solana.NewSolanaSDK(chainURL, solanaLimiter, metrics, solana.WithRetries(3, 10*time.Second))
	params := watcher.SolanaParams{Blockchain: wb.Name, ContractAddress: contractAddress,
		SizeBlocks: wb.SizeBlocks, WaitSeconds: wb.WaitSeconds, InitialBlock: wb.InitialBlock}
	return watcher.NewSolanaWatcher(solanaClient, repo, params, metrics, logger)
}

func CreateTerraWatcher(rateLimit int, chainURL string, wb config.WatcherBlockchain, logger *zap.Logger, repo *storage.Repository, metrics metrics.Metrics) watcher.ContractWatcher {
	terraLimiter := ratelimit.New(rateLimit, ratelimit.Per(time.Second))
	terraClient := terra.NewTerraSDK(chainURL, terraLimiter, metrics)
	params := watcher.TerraParams{ChainID: wb.ChainID, Blockchain: wb.Name,
		ContractAddress: wb.Address, WaitSeconds: wb.WaitSeconds, InitialBlock: wb.InitialBlock}
	return watcher.NewTerraWatcher(terraClient, params, repo, metrics, logger)
}

func CreateAptosWatcher(rateLimit int, chainURL string, wb config.WatcherBlockchain, logger *zap.Logger, repo *storage.Repository, metrics metrics.Metrics) watcher.ContractWatcher {
	aptosLimiter := ratelimit.New(rateLimit, ratelimit.Per(time.Second))
	aptosClient := aptos.NewAptosSDK(chainURL, aptosLimiter, metrics)
	params := watcher.AptosParams{
		Blockchain:      wb.Name,
		ContractAddress: wb.Address,
		SizeBlocks:      wb.SizeBlocks,
		WaitSeconds:     wb.WaitSeconds,
		InitialBlock:    wb.InitialBlock}
	return watcher.NewAptosWatcher(aptosClient, params, repo, metrics, logger)
}

func CreateOasisWatcher(rateLimit int, chainURL string, wb config.WatcherBlockchainAddresses, logger *zap.Logger, repo *storage.Repository, metrics metrics.Metrics) watcher.ContractWatcher {
	oasisLimiter := ratelimit.New(rateLimit, ratelimit.Per(time.Second))
	oasisClient := evm.NewEvmSDK(chainURL, oasisLimiter, metrics)
	params := watcher.EVMParams{
		ChainID:          wb.ChainID,
		Blockchain:       wb.Name,
		SizeBlocks:       wb.SizeBlocks,
		WaitSeconds:      wb.WaitSeconds,
		InitialBlock:     wb.InitialBlock,
		MethodsByAddress: wb.MethodsByAddress}
	return watcher.NewEvmStandarWatcher(oasisClient, params, repo, metrics, logger)
}

func CreateMoonbeamWatcher(rateLimit int, chainURL string, wb config.WatcherBlockchainAddresses, logger *zap.Logger, repo *storage.Repository, metrics metrics.Metrics) watcher.ContractWatcher {
	moonbeamLimiter := ratelimit.New(rateLimit, ratelimit.Per(time.Second))
	moonbeamClient := evm.NewEvmSDK(chainURL, moonbeamLimiter, metrics)
	params := watcher.EVMParams{
		ChainID:          wb.ChainID,
		Blockchain:       wb.Name,
		SizeBlocks:       wb.SizeBlocks,
		WaitSeconds:      wb.WaitSeconds,
		InitialBlock:     wb.InitialBlock,
		MethodsByAddress: wb.MethodsByAddress}
	return watcher.NewEvmStandarWatcher(moonbeamClient, params, repo, metrics, logger)
}

func CreateCeloWatcher(rateLimit int, chainURL string, wb config.WatcherBlockchainAddresses, logger *zap.Logger, repo *storage.Repository, metrics metrics.Metrics) watcher.ContractWatcher {
	celoLimiter := ratelimit.New(rateLimit, ratelimit.Per(time.Second))
	celoClient := evm.NewEvmSDK(chainURL, celoLimiter, metrics)
	params := watcher.EVMParams{
		ChainID:          wb.ChainID,
		Blockchain:       wb.Name,
		SizeBlocks:       wb.SizeBlocks,
		WaitSeconds:      wb.WaitSeconds,
		InitialBlock:     wb.InitialBlock,
		MethodsByAddress: wb.MethodsByAddress}
	return watcher.NewEvmStandarWatcher(celoClient, params, repo, metrics, logger)
}

func CreateOptimismWatcher(rateLimit int, chainURL string, wb config.WatcherBlockchainAddresses, logger *zap.Logger, repo *storage.Repository, metrics metrics.Metrics) watcher.ContractWatcher {
	optimismLimiter := ratelimit.New(rateLimit, ratelimit.Per(time.Second))
	optimismClient := evm.NewEvmSDK(chainURL, optimismLimiter, metrics)
	params := watcher.EVMParams{
		ChainID:          wb.ChainID,
		Blockchain:       wb.Name,
		SizeBlocks:       wb.SizeBlocks,
		WaitSeconds:      wb.WaitSeconds,
		InitialBlock:     wb.InitialBlock,
		MethodsByAddress: wb.MethodsByAddress}
	return watcher.NewEvmStandarWatcher(optimismClient, params, repo, metrics, logger)
}

func CreateArbitrumWatcher(rateLimit int, chainURL string, wb config.WatcherBlockchainAddresses, logger *zap.Logger, repo *storage.Repository, metrics metrics.Metrics) watcher.ContractWatcher {
	arbitrumLimiter := ratelimit.New(rateLimit, ratelimit.Per(time.Second))
	arbitrumClient := evm.NewEvmSDK(chainURL, arbitrumLimiter, metrics)
	params := watcher.EVMParams{
		ChainID:          wb.ChainID,
		Blockchain:       wb.Name,
		SizeBlocks:       wb.SizeBlocks,
		WaitSeconds:      wb.WaitSeconds,
		InitialBlock:     wb.InitialBlock,
		MethodsByAddress: wb.MethodsByAddress}
	return watcher.NewEvmStandarWatcher(arbitrumClient, params, repo, metrics, logger)
}
