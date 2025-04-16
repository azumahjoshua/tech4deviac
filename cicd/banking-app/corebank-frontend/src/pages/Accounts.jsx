import { useState } from 'react'
import { useSelector } from 'react-redux'
import AccountModal from '../components/AccountModal'
import AccountCard from '../components/AccountCard'

export default function Accounts() {
  const [showAccountModal, setShowAccountModal] = useState(false)
  const { items: accounts } = useSelector((state) => state.accounts)

  return (
    <div className="py-4">
      <div className="flex justify-between items-center mb-6">
        <h2 className="text-2xl font-bold text-gray-800">Accounts</h2>
        <button
          onClick={() => setShowAccountModal(true)}
          className="bg-indigo-600 hover:bg-indigo-700 text-white py-2 px-4 rounded"
        >
          New Account
        </button>
      </div>

      <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
        {accounts.length > 0 ? (
          accounts.map((account) => (
            <AccountCard key={account.id} account={account} />
          ))
        ) : (
          <div className="col-span-full text-center py-12">
            <p className="text-gray-500 text-lg">No accounts found</p>
            <button
              onClick={() => setShowAccountModal(true)}
              className="mt-4 bg-indigo-600 hover:bg-indigo-700 text-white py-2 px-4 rounded"
            >
              Create Your First Account
            </button>
          </div>
        )}
      </div>

      <AccountModal
        show={showAccountModal}
        onHide={() => setShowAccountModal(false)}
      />
    </div>
  )
}