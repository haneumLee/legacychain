'use client'

import { useState } from 'react'
import { useAccount, useWriteContract, useWaitForTransactionReceipt } from 'wagmi'
import { parseEther } from 'viem'
import { VAULT_FACTORY_ADDRESS, VAULT_FACTORY_ABI } from '@/lib/contracts'
import Link from 'next/link'
import { WalletConnect } from '@/components/WalletConnect'

interface Heir {
  address: string
  weight: number
}

export default function CreateVaultPage() {
  const { address, isConnected } = useAccount()
  const { writeContract, data: hash, isPending } = useWriteContract()
  const { isLoading: isConfirming, isSuccess } = useWaitForTransactionReceipt({
    hash,
  })

  const [heirs, setHeirs] = useState<Heir[]>([{ address: '', weight: 100 }])
  const [heartbeatInterval, setHeartbeatInterval] = useState(30) // days
  const [maxInactivity, setMaxInactivity] = useState(90) // days

  const addHeir = () => {
    setHeirs([...heirs, { address: '', weight: 0 }])
  }

  const removeHeir = (index: number) => {
    setHeirs(heirs.filter((_, i) => i !== index))
  }

  const updateHeir = (index: number, field: keyof Heir, value: string | number) => {
    const newHeirs = [...heirs]
    newHeirs[index] = { ...newHeirs[index], [field]: value }
    setHeirs(newHeirs)
  }

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault()

    if (!isConnected || !address) {
      alert('Please connect your wallet first')
      return
    }

    // Validate heirs
    const validHeirs = heirs.filter(h => h.address && h.weight > 0)
    if (validHeirs.length === 0) {
      alert('Please add at least one heir with valid address and weight')
      return
    }

    const totalWeight = validHeirs.reduce((sum, h) => sum + h.weight, 0)
    if (totalWeight !== 100) {
      alert(`Total weight must be 100% (current: ${totalWeight}%)`)
      return
    }

    try {
      // Convert days to seconds
      const heartbeatIntervalSec = heartbeatInterval * 24 * 60 * 60
      const maxInactivitySec = maxInactivity * 24 * 60 * 60

      writeContract({
        address: VAULT_FACTORY_ADDRESS,
        abi: VAULT_FACTORY_ABI,
        functionName: 'createVault',
        args: [
          validHeirs.map(h => h.address as `0x${string}`),
          validHeirs.map(h => BigInt(h.weight)),
          BigInt(heartbeatIntervalSec),
          BigInt(maxInactivitySec),
        ],
      })
    } catch (error) {
      console.error('Error creating vault:', error)
      alert('Failed to create vault')
    }
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
              Please connect your wallet to create a vault
            </p>
            <WalletConnect />
          </div>
        </main>
      </div>
    )
  }

  if (isSuccess) {
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
          <div className="max-w-2xl mx-auto text-center space-y-6">
            <div className="text-6xl">✅</div>
            <h2 className="text-3xl font-bold text-green-600">Vault Created Successfully!</h2>
            <p className="text-gray-600">
              Transaction Hash: <span className="font-mono text-sm">{hash}</span>
            </p>
            <div className="flex gap-4 justify-center">
              <Link
                href="/vault"
                className="px-6 py-3 bg-blue-600 text-white rounded-lg hover:bg-blue-700 transition-colors"
              >
                View My Vaults
              </Link>
              <Link
                href="/"
                className="px-6 py-3 bg-white border-2 border-blue-600 text-blue-600 rounded-lg hover:bg-blue-50 transition-colors"
              >
                Back to Home
              </Link>
            </div>
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
        <div className="max-w-3xl mx-auto">
          <div className="mb-8">
            <h1 className="text-4xl font-bold text-gray-900 mb-2">Create New Vault</h1>
            <p className="text-gray-600">Set up your digital inheritance vault</p>
          </div>

          <form onSubmit={handleSubmit} className="space-y-8 bg-white rounded-xl shadow-lg p-8">
            {/* Heartbeat Configuration */}
            <section>
              <h2 className="text-2xl font-semibold mb-4">Heartbeat Configuration</h2>
              
              <div className="grid md:grid-cols-2 gap-6">
                <div>
                  <label className="block text-sm font-medium text-gray-700 mb-2">
                    Heartbeat Interval (days)
                  </label>
                  <input
                    type="number"
                    min="1"
                    max="365"
                    value={heartbeatInterval}
                    onChange={(e) => setHeartbeatInterval(Number(e.target.value))}
                    className="w-full px-4 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent"
                    required
                  />
                  <p className="text-sm text-gray-500 mt-1">
                    How often you need to check in
                  </p>
                </div>

                <div>
                  <label className="block text-sm font-medium text-gray-700 mb-2">
                    Max Inactivity (days)
                  </label>
                  <input
                    type="number"
                    min="1"
                    max="730"
                    value={maxInactivity}
                    onChange={(e) => setMaxInactivity(Number(e.target.value))}
                    className="w-full px-4 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent"
                    required
                  />
                  <p className="text-sm text-gray-500 mt-1">
                    Vault expires after this period
                  </p>
                </div>
              </div>
            </section>

            {/* Heirs Configuration */}
            <section>
              <div className="flex justify-between items-center mb-4">
                <h2 className="text-2xl font-semibold">Beneficiaries (Heirs)</h2>
                <button
                  type="button"
                  onClick={addHeir}
                  className="px-4 py-2 bg-blue-500 text-white rounded-lg hover:bg-blue-600 transition-colors"
                >
                  + Add Heir
                </button>
              </div>

              <div className="space-y-4">
                {heirs.map((heir, index) => (
                  <div key={index} className="flex gap-4 items-start p-4 border border-gray-200 rounded-lg">
                    <div className="flex-1">
                      <label className="block text-sm font-medium text-gray-700 mb-2">
                        Ethereum Address
                      </label>
                      <input
                        type="text"
                        placeholder="0x..."
                        value={heir.address}
                        onChange={(e) => updateHeir(index, 'address', e.target.value)}
                        className="w-full px-4 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent font-mono text-sm"
                        required
                      />
                    </div>

                    <div className="w-32">
                      <label className="block text-sm font-medium text-gray-700 mb-2">
                        Weight (%)
                      </label>
                      <input
                        type="number"
                        min="0"
                        max="100"
                        value={heir.weight}
                        onChange={(e) => updateHeir(index, 'weight', Number(e.target.value))}
                        className="w-full px-4 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent"
                        required
                      />
                    </div>

                    {heirs.length > 1 && (
                      <button
                        type="button"
                        onClick={() => removeHeir(index)}
                        className="mt-8 px-4 py-2 bg-red-500 text-white rounded-lg hover:bg-red-600 transition-colors"
                      >
                        Remove
                      </button>
                    )}
                  </div>
                ))}
              </div>

              <p className="text-sm text-gray-500 mt-2">
                Total weight: {heirs.reduce((sum, h) => sum + h.weight, 0)}% (must equal 100%)
              </p>
            </section>

            {/* Submit Button */}
            <div className="flex gap-4">
              <button
                type="submit"
                disabled={isPending || isConfirming}
                className="flex-1 px-8 py-4 bg-blue-600 text-white rounded-lg hover:bg-blue-700 transition-colors font-semibold text-lg disabled:bg-gray-400 disabled:cursor-not-allowed"
              >
                {isPending ? 'Confirming in Wallet...' : isConfirming ? 'Creating Vault...' : 'Create Vault'}
              </button>
              <Link
                href="/"
                className="px-8 py-4 bg-white border-2 border-gray-300 text-gray-700 rounded-lg hover:bg-gray-50 transition-colors font-semibold text-lg"
              >
                Cancel
              </Link>
            </div>
          </form>
        </div>
      </main>
    </div>
  )
}
