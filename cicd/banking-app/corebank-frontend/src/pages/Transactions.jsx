import { useState } from 'react'
import { useSelector } from 'react-redux'
import TransactionModal from '../components/TransactionModal'
import TransactionCard from '../components/TransactionCard'

export default function Transactions() {
  const [showTransactionModal, setShowTransactionModal] = useState(false)
  const [transactionType, setTransactionType] = useState('')
  const { items: transactions } = useSelector((state) => state.transactions)

  const handleTransactionClick = (type) => {
    setTransactionType(type)
    setShowTransactionModal(true)
  }

  return (
    <div className="py-4">
      <div className="flex justify-between items-center mb-6">
        <h2 className="text-2xl font-bold text-gray-800">Transactions</h2>
        <div className="space-x-3">
          <button
            onClick={() => handleTransactionClick('deposit')}
            className="bg-green-600 hover:bg-green-700 text-white py-2 px-4 rounded"
          >
            New Deposit
          </button>
          <button
            onClick={() => handleTransactionClick('withdrawal')}
            className="bg-blue-600 hover:bg-blue-700 text-white py-2 px-4 rounded"
          >
            New Withdrawal
          </button>
        </div>
      </div>

      <div className="bg-white rounded-lg shadow overflow-hidden">
        <div className="p-6">
          <div className="mb-4 flex justify-between items-center">
            <h3 className="text-lg font-medium text-gray-700">Transaction History</h3>
            <div className="flex space-x-2">
              <select className="border border-gray-300 rounded-md px-3 py-1 text-sm">
                <option>All Accounts</option>
                <option>Checking</option>
                <option>Savings</option>
              </select>
              <select className="border border-gray-300 rounded-md px-3 py-1 text-sm">
                <option>Last 30 Days</option>
                <option>Last 90 Days</option>
                <option>This Year</option>
              </select>
            </div>
          </div>

          <div className="space-y-3">
            {transactions.length > 0 ? (
              transactions.map((transaction) => (
                <TransactionCard key={transaction.id} transaction={transaction} />
              ))
            ) : (
              <div className="text-center py-8">
                <p className="text-gray-500">No transactions found</p>
              </div>
            )}
          </div>
        </div>
      </div>

      <TransactionModal
        show={showTransactionModal}
        onHide={() => setShowTransactionModal(false)}
        type={transactionType}
      />
    </div>
  )
}