import { useState } from 'react'
import { useSelector } from 'react-redux'
import TransactionModal from '../components/TransactionModal'
import AccountModal from '../components/AccountModal'
import BalanceChart from '../components/BalanceChart'
import TransactionCard from '../components/TransactionCard'

export default function Dashboard() {
  const [showTransactionModal, setShowTransactionModal] = useState(false)
  const [showAccountModal, setShowAccountModal] = useState(false)
  const [transactionType, setTransactionType] = useState('')
  
  const { accounts, transactions } = useSelector((state) => ({
    accounts: state.accounts.items,
    transactions: state.transactions.recent
  }))

  const totalBalance = accounts.reduce((sum, account) => sum + account.balance, 0)

  const handleTransactionClick = (type) => {
    setTransactionType(type)
    setShowTransactionModal(true)
  }

  return (
    <div className="py-4">
      <h2 className="text-2xl font-bold text-gray-800">Dashboard</h2>
      
      <div className="mt-6 grid grid-cols-1 md:grid-cols-2 gap-6">
        {/* Balance Summary Card */}
        <div className="bg-white rounded-lg shadow p-6">
          <h3 className="text-lg font-medium text-gray-700">Total Balance</h3>
          <p className="text-4xl font-bold my-4">${totalBalance.toFixed(2)}</p>
          <BalanceChart accounts={accounts} />
        </div>
        
        {/* Quick Actions Card */}
        <div className="bg-white rounded-lg shadow p-6">
          <h3 className="text-lg font-medium text-gray-700">Quick Actions</h3>
          <div className="mt-4 space-y-3">
            <button
              onClick={() => handleTransactionClick('deposit')}
              className="w-full bg-green-600 hover:bg-green-700 text-white py-2 px-4 rounded"
            >
              New Deposit
            </button>
            <button
              onClick={() => handleTransactionClick('withdrawal')}
              className="w-full bg-blue-600 hover:bg-blue-700 text-white py-2 px-4 rounded"
            >
              New Withdrawal
            </button>
            <button
              onClick={() => setShowAccountModal(true)}
              className="w-full bg-indigo-600 hover:bg-indigo-700 text-white py-2 px-4 rounded"
            >
              New Account
            </button>
          </div>
        </div>
        
        {/* Recent Transactions Card */}
        <div className="md:col-span-2 bg-white rounded-lg shadow p-6">
          <h3 className="text-lg font-medium text-gray-700">Recent Transactions</h3>
          <div className="transaction-list mt-4 space-y-3">
            {transactions.length > 0 ? (
              transactions.map(transaction => (
                <TransactionCard key={transaction.id} transaction={transaction} />
              ))
            ) : (
              <p className="text-gray-500">No recent transactions</p>
            )}
          </div>
        </div>
      </div>
      
      {/* Modals */}
      <TransactionModal
        show={showTransactionModal}
        onHide={() => setShowTransactionModal(false)}
        type={transactionType}
      />
      <AccountModal
        show={showAccountModal}
        onHide={() => setShowAccountModal(false)}
      />
    </div>
  )
}