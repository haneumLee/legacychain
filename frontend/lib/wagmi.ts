import { createConfig, http } from 'wagmi'
import { localhost } from 'wagmi/chains'

// Custom chain configuration for Besu
export const besuLocal = {
  id: 1337,
  name: 'Besu Local',
  nativeCurrency: {
    decimals: 18,
    name: 'Ether',
    symbol: 'ETH',
  },
  rpcUrls: {
    default: {
      http: [process.env.NEXT_PUBLIC_RPC_URL || 'http://localhost:8545'],
    },
    public: {
      http: [process.env.NEXT_PUBLIC_RPC_URL || 'http://localhost:8545'],
    },
  },
  blockExplorers: {
    default: {
      name: 'Besu Explorer',
      url: 'http://localhost:8545',
    },
  },
  testnet: true,
}

export const config = createConfig({
  chains: [besuLocal as typeof localhost],
  transports: {
    [besuLocal.id]: http(),
  },
})
