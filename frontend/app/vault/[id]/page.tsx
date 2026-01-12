'use client'

import { use, useState } from 'react'
import { useAccount, useReadContract, useWriteContract, useWaitForTransactionReceipt } from 'wagmi'
import { INDIVIDUAL_VAULT_ABI } from '@/lib/contracts'
import Link from 'next/link'
import { WalletConnect } from '@/components/WalletConnect'
import { formatEther } from 'viem'

interface PageProps {
  params: Promise<{
    id: string
  }>
}

export default function VaultDetailPage({ params }: PageProps) {
  const { id } = use(params)
  const vaultAddress = id as `0x${string}`
  const { address, isConnected } = useAccount()

  // Read vault data
  const { data: owner } = useReadContract({
    address: vaultAddress,
    abi: INDIVIDUAL_VAULT_ABI,
    functionName: 'owner',
  })

  const { data: heartbeatInterval } = useReadContract({
    address: vaultAddress,
    abi: INDIVIDUAL_VAULT_ABI,
    functionName: 'heartbeatInterval',
  })

  const { data: maxInactivityPeriod } = useReadContract({
    address: vaultAddress,
    abi: INDIVIDUAL_VAULT_ABI,
    functionName: 'maxInactivityPeriod',
  })

  const { data: lastHeartbeat } = useReadContract({
    address: vaultAddress,
    abi: INDIVIDUAL_VAULT_ABI,
    functionName: 'lastHeartbeat',
  })

  const { data: isPaused } = useReadContract({
    address: vaultAddress,
    abi: INDIVIDUAL_VAULT_ABI,
    functionName: 'paused',
  })

  // Heartbeat commit
  const [commitHash, setCommitHash] = useState('')
  const { writeContract: commitHeartbeat, data: commitHash_tx, isPending: isCommitting } = useWriteContract()
  const { isLoading: isCommitConfirming, isSuccess: isCommitSuccess } = useWaitForTransactionReceipt({
    hash: commitHash_tx,
  })

  // Heartbeat reveal
  const [revealSecret, setRevealSecret] = useState('')
  const { writeContract: revealHeartbeat, data: revealHash, isPending: isRevealing } = useWriteContract()
  const { isLoading: isRevealConfirming, isSuccess: isRevealSuccess } = useWaitForTransactionReceipt({
    hash: revealHash,
  })

  const handleCommit = (e: React.FormEvent) => {
    e.preventDefault()
    if (!commitHash) return

    commitHeartbeat({
      address: vaultAddress,
      abi: INDIVIDUAL_VAULT_ABI,
      functionName: 'commitHeartbeat',
      args: [commitHash as `0x${string}`],
    })
  }

  const handleReveal = (e: React.FormEvent) => {
    e.preventDefault()
    if (!revealSecret) return

    revealHeartbeat({
      address: vaultAddress,
      abi: INDIVIDUAL_VAULT_ABI,
      functionName: 'revealHeartbeat',
      args: [revealSecret as `0x${string}`],
    })
  }

  const formatAddress = (addr: string) => {
    return `${addr.slice(0, 6)}...${addr.slice(-4)}`
  }

  const formatTimestamp = (timestamp: bigint | undefined) => {
    if (!timestamp) return 'N/A'
    const date = new Date(Number(timestamp) * 1000)
    return date.toLocaleString()
  }

  const formatDuration = (seconds: bigint | undefined) => {
    if (!seconds) return 'N/A'
    const days = Number(seconds) / (24 * 60 * 60)
    return `${days.toFixed(0)} days`
  }

  const getTimeUntilExpiry = () => {
    if (!lastHeartbeat || !maxInactivityPeriod) return null
    const expiryTime = Number(lastHeartbeat) + Number(maxInactivityPeriod)
    const now = Math.floor(Date.now() / 1000)
    const remaining = expiryTime - now

    if (remaining <= 0) return 'Expired'
    
    const days = Math.floor(remaining / (24 * 60 * 60))
    const hours = Math.floor((remaining % (24 * 60 * 60)) / 3600)
    
    if (days > 0) return `${days} days ${hours} hours`
    return `${hours} hours`
  }

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
              Please connect your wallet to view vault details
            </p>
            <WalletConnect />
          </div>
        </main>
      </div>
    )
  }

  const isOwner = owner && address && owner.toLowerCase() === address.toLowerCase()

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
          <div className="mb-6">
            <Link href="/vault" className="text-blue-600 hover:text-blue-700 font-semibold">
              ← Back to Vaults
            </Link>
          </div>

          <div className="mb-8">
            <h1 className="text-4xl font-bold text-gray-900 mb-2">Vault Details</h1>
            <p className="text-gray-600 font-mono">{vaultAddress}</p>
          </div>

          <div className="grid lg:grid-cols-3 gap-6">
            {/* Vault Info */}
            <div className="lg:col-span-2 space-y-6">
              {/* Status Card */}
              <div className="bg-white rounded-xl shadow-lg p-6">
                <h2 className="text-2xl font-semibold mb-4">Status</h2>
                <div className="grid md:grid-cols-2 gap-6">
                  <div>
                    <p className="text-sm text-gray-600 mb-1">Owner</p>
                    <p className="font-mono text-sm">{owner ? formatAddress(owner) : 'Loading...'}</p>
                    {isOwner && <span className="text-xs text-green-600 font-semibold">(You)</span>}
                  </div>
                  <div>
                    <p className="text-sm text-gray-600 mb-1">State</p>
                    <span className={`inline-block px-3 py-1 rounded-full text-sm font-semibold ${
                      isPaused ? 'bg-red-100 text-red-700' : 'bg-green-100 text-green-700'
                    }`}>
                      {isPaused ? 'Paused' : 'Active'}
                    </span>
                  </div>
                  <div>
                    <p className="text-sm text-gray-600 mb-1">Last Heartbeat</p>
                    <p className="font-semibold">{formatTimestamp(lastHeartbeat)}</p>
                  </div>
                  <div>
                    <p className="text-sm text-gray-600 mb-1">Time Until Expiry</p>
                    <p className="font-semibold text-orange-600">{getTimeUntilExpiry() || 'N/A'}</p>
                  </div>
                </div>
              </div>

              {/* Configuration Card */}
              <div className="bg-white rounded-xl shadow-lg p-6">
                <h2 className="text-2xl font-semibold mb-4">Configuration</h2>
                <div className="grid md:grid-cols-2 gap-6">
                  <div>
                    <p className="text-sm text-gray-600 mb-1">Heartbeat Interval</p>
                    <p className="font-semibold text-lg">{formatDuration(heartbeatInterval)}</p>
                  </div>
                  <div>
                    <p className="text-sm text-gray-600 mb-1">Max Inactivity Period</p>
                    <p className="font-semibold text-lg">{formatDuration(maxInactivityPeriod)}</p>
                  </div>
                </div>
              </div>

              {/* Heartbeat Management (Owner Only) */}
              {isOwner && (
                <div className="bg-white rounded-xl shadow-lg p-6">
                  <h2 className="text-2xl font-semibold mb-4">Heartbeat Management</h2>
                  
                  {/* Commit Phase */}
                  <div className="mb-6">
                    <h3 className="text-lg font-semibold mb-3">Step 1: Commit</h3>
                    <form onSubmit={handleCommit} className="space-y-4">
                      <div>
                        <label className="block text-sm font-medium text-gray-700 mb-2">
                          Commit Hash (keccak256 of secret)
                        </label>
                        <input
                          type="text"
                          placeholder="0x..."
                          value={commitHash}
                          onChange={(e) => setCommitHash(e.target.value)}
                          className="w-full px-4 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent font-mono text-sm"
                        />
                        <p className="text-xs text-gray-500 mt-1">
                          Generate with: keccak256(abi.encodePacked(secret))
                        </p>
                      </div>
                      <button
                        type="submit"
                        disabled={isCommitting || isCommitConfirming || !commitHash}
                        className="px-6 py-3 bg-blue-600 text-white rounded-lg hover:bg-blue-700 transition-colors font-semibold disabled:bg-gray-400 disabled:cursor-not-allowed"
                      >
                        {isCommitting ? 'Confirming...' : isCommitConfirming ? 'Processing...' : 'Commit Heartbeat'}
                      </button>
                      {isCommitSuccess && (
                        <p className="text-sm text-green-600">Commit successful! Now reveal your secret.</p>
                      )}
                    </form>
                  </div>

                  {/* Reveal Phase */}
                  <div className="pt-6 border-t border-gray-200">
                    <h3 className="text-lg font-semibold mb-3">Step 2: Reveal</h3>
                    <form onSubmit={handleReveal} className="space-y-4">
                      <div>
                        <label className="block text-sm font-medium text-gray-700 mb-2">
                          Secret (hex format)
                        </label>
                        <input
                          type="text"
                          placeholder="0x..."
                          value={revealSecret}
                          onChange={(e) => setRevealSecret(e.target.value)}
                          className="w-full px-4 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent font-mono text-sm"
                        />
                        <p className="text-xs text-gray-500 mt-1">
                          Must match the secret used in commit phase
                        </p>
                      </div>
                      <button
                        type="submit"
                        disabled={isRevealing || isRevealConfirming || !revealSecret}
                        className="px-6 py-3 bg-green-600 text-white rounded-lg hover:bg-green-700 transition-colors font-semibold disabled:bg-gray-400 disabled:cursor-not-allowed"
                      >
                        {isRevealing ? 'Confirming...' : isRevealConfirming ? 'Processing...' : 'Reveal Heartbeat'}
                      </button>
                      {isRevealSuccess && (
                        <p className="text-sm text-green-600">Heartbeat revealed successfully!</p>
                      )}
                    </form>
                  </div>
                </div>
              )}
            </div>

            {/* Sidebar */}
            <div className="space-y-6">
              {/* Quick Actions */}
              {isOwner && (
                <div className="bg-white rounded-xl shadow-lg p-6">
                  <h2 className="text-xl font-semibold mb-4">Quick Actions</h2>
                  <div className="space-y-3">
                    <button className="w-full px-4 py-3 bg-blue-600 text-white rounded-lg hover:bg-blue-700 transition-colors font-semibold">
                      Add Funds
                    </button>
                    <button className="w-full px-4 py-3 bg-gray-600 text-white rounded-lg hover:bg-gray-700 transition-colors font-semibold">
                      Manage Heirs
                    </button>
                    <button className="w-full px-4 py-3 bg-orange-600 text-white rounded-lg hover:bg-orange-700 transition-colors font-semibold">
                      {isPaused ? 'Resume Vault' : 'Pause Vault'}
                    </button>
                  </div>
                </div>
              )}

              {/* Info */}
              <div className="bg-blue-50 border border-blue-200 rounded-xl p-6">
                <h3 className="font-semibold text-blue-900 mb-2">About Heartbeat</h3>
                <p className="text-sm text-blue-800">
                  The commit-reveal mechanism ensures your activity proof cannot be front-run. 
                  First commit a hash, then reveal the secret within the allowed timeframe.
                </p>
              </div>
            </div>
          </div>
        </div>
      </main>
    </div>
  )
}
