package config

import walletutil "imsdk/internal/client/pkg/utils/wallet-util"

func GetServerWallet() *walletutil.Wallet {
	return serverWallet
}
