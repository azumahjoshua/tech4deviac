import { useState } from 'react'
import { useSelector } from 'react-redux'

export default function TransactionModal({ show, onHide, type }) {
  const [amount, setAmount] = useState('')
  const [description, setDescription] = useState('')
  const [accountId, setAccountId] = useState('')
  
  const { accounts } = useSelector((state) => state.accounts)
  
  const handleSubmit = (e) => {
    e.preventDefault()
    // Handle transaction submission
    onHide()
  }

  return (
    show && (
      <div className="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center p-4 z-50">
        <div className="bg-white rounded-lg shadow-xl w-full max-w-md">
          <div className="flex justify-between items-center border-b p-4">
            <h3 className="text-lg font-semibold">
              New {type === 'deposit' ? 'Deposit' : 'Withdrawal'}
            </h3>
            <button onClick={onHide} className="text-gray-500 hover:text-gray-700">
              &times;
            </button>
          </div>
          
          <form onSubmit={handleSubmit} className="p-4 space-y-4">
            <div>
              <label className="block text-sm font-medium text-gray-700 mb-1">Account</label>
              <select
                value={accountId}
                onChange={(e) => setAccountId(e.target.value)}
                className="w-full border border-gray-300 rounded-md p-2"
                required
              >
                <option value="">Select an account</option>
                {accounts.map(account => (
                  <option key={account.id} value={account.id}>
                    {account.name} (${account.balance.toFixed(2)})
                  </option>
                ))}
              </select>
            </div>
            
            <div>
              <label className="block text-sm font-medium text-gray-700 mb-1">Amount</label>
              <div className="relative">
                <span className="absolute inset-y-0 left-0 flex items-center pl-3 text-gray-500">$</span>
                <input
                  type="number"
                  value={amount}
                  onChange={(e) => setAmount(e.target.value)}
                  min="0.01"
                  step="0.01"
                  className="w-full border border-gray-300 rounded-md p-2 pl-8"
                  required
                />
              </div>
            </div>
            
            <div>
              <label className="block text-sm font-medium text-gray-700 mb-1">Description</label>
              <input
                type="text"
                value={description}
                onChange={(e) => setDescription(e.target.value)}
                className="w-full border border-gray-300 rounded-md p-2"
                required
              />
            </div>
            
            <div className="flex justify-end space-x-3 pt-4">
              <button
                type="button"
                onClick={onHide}
                className="px-4 py-2 border border-gray-300 rounded-md text-gray-700 hover:bg-gray-50"
              >
                Cancel
              </button>
              <button
                type="submit"
                className="px-4 py-2 bg-blue-600 text-white rounded-md hover:bg-blue-700"
              >
                Submit
              </button>
            </div>
          </form>
        </div>
      </div>
    )
  )
}