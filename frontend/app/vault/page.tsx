'use client'

import { useAccount, useReadContract } from 'wagmi'
import { VAULT_FACTORY_ADDRESS, VAULT_FACTORY_ABI } from '@/lib/contracts'
import Link from 'next/link'
import { WalletConnect } from '@/components/WalletConnect'

export default function VaultListPage() {
  const { address, isConnected } = useAccount()

  const { data: vaults, isLoading, error } = useReadContract({
    address: VAULT_FACTORY_ADDRESS,
    abi: VAULT_FACTORY_ABI,
    functionName: 'getUserVaults',
    args: address ? [address] : undefined,
    query: {
      enabled: !!address,
    },
  })

  if (!isConnected) {
    return (
      <div className="min-h-screen bg-gradient-to-br from-blue-50 to-indigo-100">
        <header className="border-b bg-white/80 backdrop-blur-sm">
          <div className="container mx-auto px-4 py-4 flex justify-between items-center">
            <Link href="/" className="text-2xl font-bold text-blue-600">
              LegacyChain
            </Link>
            <WalletConnect />
          </div>
        </header>

        <main className="container mx-auto px-4 py-16">
          <div className="max-w-2xl mx-auto text-center">
            <h2 className="text-3xl font-bold mb-4">Connect Your Wallet</h2>
            <p className="text-gray-600 mb-8">
              Please connect your wallet to view your vaults
            </p>
            <WalletConnect />
          </div>
        </main>
      </div>
    )
  }

  return (
    <div className="min-h-screen bg-gradient-to-br from-blue-50 to-indigo-100">
      <header className="border-b bg-white/80 backdrop-blur-sm">
        <div className="container mx-auto px-4 py-4 flex justify-between items-center">
          <Link href="/" className="text-2xl font-bold text-blue-600">
            LegacyChain
          </Link>
          <WalletConnect />
        </div>
      </header>

      <main className="container mx-auto px-4 py-8">
        <div className="max-w-6xl mx-auto">
          <div className="flex justify-between items-center mb-8">
            <div>
              <h1 className="text-4xl font-bold text-gray-900 mb-2">My Vaults</h1>
              <p className="text-gray-600">Manage your digital inheritance vaults</p>
            </div>
            <Link
              href="/vault/create"
              className="px-6 py-3 bg-blue-600 text-white rounded-lg hover:bg-blue-700 transition-colors font-semibold"
            >
              + Create New Vault
            </Link>
          </div>

          {isLoading && (
            <div className="text-center py-12">
              <div className="inline-block animate-spin rounded-full h-12 w-12 border-b-2 border-blue-600"></div>
              <p className="mt-4 text-gray-600">Loading your vaults...</p>
            </div>
          )}

          {error && (
            <div className="bg-red-50 border border-red-200 rounded-lg p-6 text-center">
              <p className="text-red-600">Error loading vaults: {error.message}</p>
            </div>
          )}

          {!isLoading && !error && vaults && (
            <>
              {(vaults as readonly `0x${string}`[]).length === 0 ? (
                <div className="bg-white rounded-xl shadow-lg p-12 text-center">
                  <div className="text-6xl mb-4">📦</div>
                  <h2 className="text-2xl font-semibold mb-2">No Vaults Yet</h2>
                  <p className="text-gray-600 mb-6">
                    Create your first vault to get started with digital inheritance
                  </p>
                  <Link
                    href="/vault/create"
                    className="inline-block px-8 py-4 bg-blue-600 text-white rounded-lg hover:bg-blue-700 transition-colors font-semibold"
                  >
                    Create Your First Vault
                  </Link>
                </div>
              ) : (
                <div className="grid md:grid-cols-2 lg:grid-cols-3 gap-6">
                  {(vaults as readonly `0x${string}`[]).map((vaultAddress, index) => (
                    <VaultCard key={vaultAddress} address={vaultAddress} index={index} />
                  ))}
                </div>
              )}
            </>
          )}
        </div>
      </main>
    </div>
  )
}

function VaultCard({ address, index }: { address: `0x${string}`; index: number }) {
  const formatAddress = (addr: string) => {
    return `${addr.slice(0, 6)}...${addr.slice(-4)}`
  }

  return (
    <Link
      href={`/vault/${address}`}
      className="block bg-white rounded-xl shadow-md hover:shadow-xl transition-shadow p-6"
    >
      <div className="flex items-start justify-between mb-4">
        <div>
          <h3 className="text-lg font-semibold text-gray-900">Vault #{index + 1}</h3>
          <p className="text-sm text-gray-500 font-mono mt-1">{formatAddress(address)}</p>
        </div>
        <div className="px-3 py-1 bg-green-100 text-green-700 rounded-full text-xs font-semibold">
          Active
        </div>
      </div>

      <div className="space-y-3 text-sm">
        <div className="flex justify-between">
          <span className="text-gray-600">Status:</span>
          <span className="font-semibold text-green-600">Operational</span>
        </div>
        <div className="flex justify-between">
          <span className="text-gray-600">Last Heartbeat:</span>
          <span className="font-semibold">-</span>
        </div>
        <div className="flex justify-between">
          <span className="text-gray-600">Heirs:</span>
          <span className="font-semibold">-</span>
        </div>
      </div>

      <div className="mt-6 pt-4 border-t border-gray-200">
        <button className="w-full text-center text-blue-600 font-semibold hover:text-blue-700">
          View Details →
        </button>
      </div>
    </Link>
  )
}
