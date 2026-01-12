'use client'

import { useAccount, useReadContract } from 'wagmi'
import { VAULT_FACTORY_ADDRESS, VAULT_FACTORY_ABI, INDIVIDUAL_VAULT_ABI } from '@/lib/contracts'
import Link from 'next/link'
import { WalletConnect } from '@/components/WalletConnect'

export default function DashboardPage() {
  const { address, isConnected } = useAccount()

  const { data: vaults, isLoading } = useReadContract({
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
              Please connect your wallet to view your dashboard
            </p>
            <WalletConnect />
          </div>
        </main>
      </div>
    )
  }

  const vaultAddresses = (vaults as readonly `0x${string}`[]) || []

  return (
    <div className="min-h-screen bg-gradient-to-br from-blue-50 to-indigo-100">
      <header className="border-b bg-white/80 backdrop-blur-sm">
        <div className="container mx-auto px-4 py-4 flex justify-between items-center">
          <Link href="/" className="text-2xl font-bold text-blue-600">
            LegacyChain
          </Link>
          <div className="flex gap-4 items-center">
            <Link
              href="/vault"
              className="px-4 py-2 text-blue-600 hover:text-blue-700 font-semibold"
            >
              My Vaults
            </Link>
            <WalletConnect />
          </div>
        </div>
      </header>

      <main className="container mx-auto px-4 py-8">
        <div className="max-w-7xl mx-auto">
          <div className="mb-8">
            <h1 className="text-4xl font-bold text-gray-900 mb-2">Dashboard</h1>
            <p className="text-gray-600">Monitor your vaults and heartbeat activity</p>
          </div>

          {/* Statistics */}
          <div className="grid md:grid-cols-4 gap-6 mb-8">
            <div className="bg-white rounded-xl shadow-md p-6">
              <p className="text-sm text-gray-600 mb-1">Total Vaults</p>
              <p className="text-3xl font-bold text-blue-600">
                {isLoading ? '...' : vaultAddresses.length}
              </p>
            </div>
            <div className="bg-white rounded-xl shadow-md p-6">
              <p className="text-sm text-gray-600 mb-1">Active Vaults</p>
              <p className="text-3xl font-bold text-green-600">
                {isLoading ? '...' : vaultAddresses.length}
              </p>
            </div>
            <div className="bg-white rounded-xl shadow-md p-6">
              <p className="text-sm text-gray-600 mb-1">Pending Heartbeats</p>
              <p className="text-3xl font-bold text-orange-600">0</p>
            </div>
            <div className="bg-white rounded-xl shadow-md p-6">
              <p className="text-sm text-gray-600 mb-1">Alerts</p>
              <p className="text-3xl font-bold text-red-600">0</p>
            </div>
          </div>

          <div className="grid lg:grid-cols-3 gap-6">
            {/* Heartbeat Timeline */}
            <div className="lg:col-span-2 bg-white rounded-xl shadow-lg p-6">
              <h2 className="text-2xl font-semibold mb-6">Heartbeat Timeline</h2>
              
              {isLoading ? (
                <div className="text-center py-12">
                  <div className="inline-block animate-spin rounded-full h-12 w-12 border-b-2 border-blue-600"></div>
                </div>
              ) : vaultAddresses.length === 0 ? (
                <div className="text-center py-12">
                  <div className="text-6xl mb-4">📊</div>
                  <p className="text-gray-600 mb-4">No vaults to display</p>
                  <Link
                    href="/vault/create"
                    className="inline-block px-6 py-3 bg-blue-600 text-white rounded-lg hover:bg-blue-700 transition-colors"
                  >
                    Create Your First Vault
                  </Link>
                </div>
              ) : (
                <div className="space-y-4">
                  {vaultAddresses.map((vaultAddress) => (
                    <VaultTimeline key={vaultAddress} address={vaultAddress} />
                  ))}
                </div>
              )}
            </div>

            {/* Alerts & Quick Actions */}
            <div className="space-y-6">
              {/* Expiry Warnings */}
              <div className="bg-white rounded-xl shadow-lg p-6">
                <h2 className="text-xl font-semibold mb-4">Expiry Warnings</h2>
                <div className="space-y-3">
                  {vaultAddresses.length === 0 ? (
                    <p className="text-sm text-gray-500">No alerts</p>
                  ) : (
                    vaultAddresses.map((vaultAddress) => (
                      <VaultAlert key={vaultAddress} address={vaultAddress} />
                    ))
                  )}
                </div>
              </div>

              {/* Quick Heartbeat */}
              <div className="bg-blue-50 border border-blue-200 rounded-xl p-6">
                <h3 className="font-semibold text-blue-900 mb-3">Quick Heartbeat</h3>
                <p className="text-sm text-blue-800 mb-4">
                  Send heartbeat to all your vaults at once
                </p>
                <button className="w-full px-4 py-3 bg-blue-600 text-white rounded-lg hover:bg-blue-700 transition-colors font-semibold">
                  Send Heartbeat to All
                </button>
              </div>

              {/* Help */}
              <div className="bg-gray-50 border border-gray-200 rounded-xl p-6">
                <h3 className="font-semibold text-gray-900 mb-3">Need Help?</h3>
                <p className="text-sm text-gray-600 mb-4">
                  Learn more about managing your digital legacy
                </p>
                <Link
                  href="/docs"
                  className="text-blue-600 hover:text-blue-700 font-semibold text-sm"
                >
                  View Documentation →
                </Link>
              </div>
            </div>
          </div>
        </div>
      </main>
    </div>
  )
}

function VaultTimeline({ address }: { address: `0x${string}` }) {
  const { data: lastHeartbeat } = useReadContract({
    address,
    abi: INDIVIDUAL_VAULT_ABI,
    functionName: 'lastHeartbeat',
  })

  const { data: heartbeatInterval } = useReadContract({
    address,
    abi: INDIVIDUAL_VAULT_ABI,
    functionName: 'heartbeatInterval',
  })

  const formatAddress = (addr: string) => {
    return `${addr.slice(0, 6)}...${addr.slice(-4)}`
  }

  const formatTimestamp = (timestamp: bigint | undefined) => {
    if (!timestamp) return 'Never'
    const date = new Date(Number(timestamp) * 1000)
    return date.toLocaleString()
  }

  const getNextHeartbeat = () => {
    if (!lastHeartbeat || !heartbeatInterval) return 'N/A'
    const next = Number(lastHeartbeat) + Number(heartbeatInterval)
    const date = new Date(next * 1000)
    return date.toLocaleString()
  }

  return (
    <div className="border border-gray-200 rounded-lg p-4 hover:border-blue-300 transition-colors">
      <div className="flex justify-between items-start mb-3">
        <div>
          <Link
            href={`/vault/${address}`}
            className="font-mono text-sm text-blue-600 hover:text-blue-700"
          >
            {formatAddress(address)}
          </Link>
        </div>
        <span className="px-2 py-1 bg-green-100 text-green-700 rounded text-xs font-semibold">
          Active
        </span>
      </div>

      <div className="space-y-2 text-sm">
        <div className="flex justify-between">
          <span className="text-gray-600">Last Heartbeat:</span>
          <span className="font-semibold">{formatTimestamp(lastHeartbeat)}</span>
        </div>
        <div className="flex justify-between">
          <span className="text-gray-600">Next Due:</span>
          <span className="font-semibold">{getNextHeartbeat()}</span>
        </div>
      </div>

      <Link
        href={`/vault/${address}`}
        className="block mt-4 text-center px-4 py-2 bg-blue-50 text-blue-600 rounded hover:bg-blue-100 transition-colors text-sm font-semibold"
      >
        Send Heartbeat
      </Link>
    </div>
  )
}

function VaultAlert({ address }: { address: `0x${string}` }) {
  const { data: lastHeartbeat } = useReadContract({
    address,
    abi: INDIVIDUAL_VAULT_ABI,
    functionName: 'lastHeartbeat',
  })

  const { data: maxInactivityPeriod } = useReadContract({
    address,
    abi: INDIVIDUAL_VAULT_ABI,
    functionName: 'maxInactivityPeriod',
  })

  const formatAddress = (addr: string) => {
    return `${addr.slice(0, 6)}...${addr.slice(-4)}`
  }

  const getTimeRemaining = () => {
    if (!lastHeartbeat || !maxInactivityPeriod) return null
    
    const expiryTime = Number(lastHeartbeat) + Number(maxInactivityPeriod)
    const now = Math.floor(Date.now() / 1000)
    const remaining = expiryTime - now

    if (remaining <= 0) return { text: 'Expired', color: 'red', days: 0 }
    
    const days = Math.floor(remaining / (24 * 60 * 60))
    
    if (days <= 7) {
      return { text: `${days} days remaining`, color: 'red', days }
    } else if (days <= 30) {
      return { text: `${days} days remaining`, color: 'orange', days }
    }
    
    return null
  }

  const timeInfo = getTimeRemaining()

  if (!timeInfo) return null

  return (
    <Link
      href={`/vault/${address}`}
      className={`block p-3 rounded-lg border-2 ${
        timeInfo.color === 'red' 
          ? 'bg-red-50 border-red-200' 
          : 'bg-orange-50 border-orange-200'
      }`}
    >
      <div className="flex items-center justify-between mb-1">
        <span className="font-mono text-xs">{formatAddress(address)}</span>
        <span className={`text-xs font-semibold ${
          timeInfo.color === 'red' ? 'text-red-700' : 'text-orange-700'
        }`}>
          ⚠️ {timeInfo.text}
        </span>
      </div>
      <p className="text-xs text-gray-600">Click to send heartbeat</p>
    </Link>
  )
}
