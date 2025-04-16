import { FaPiggyBank, FaWallet, FaUniversity } from 'react-icons/fa'

const accountIcons = {
  checking: <FaWallet className="text-blue-500" size={24} />,
  savings: <FaPiggyBank className="text-green-500" size={24} />,
  business: <FaUniversity className="text-purple-500" size={24} />
}

export default function AccountCard({ account }) {
  return (
    <div className="bg-white rounded-lg shadow overflow-hidden hover:shadow-md transition-shadow">
      <div className="p-6">
        <div className="flex justify-between items-start">
          <div>
            <h3 className="text-lg font-semibold text-gray-800">{account.name}</h3>
            <p className="text-gray-500 text-sm">{account.accountNumber}</p>
          </div>
          {accountIcons[account.type] || accountIcons.checking}
        </div>
        
        <div className="mt-6">
          <p className="text-gray-500 text-sm">Current Balance</p>
          <p className="text-2xl font-bold">${account.balance.toFixed(2)}</p>
        </div>
        
        <div className="mt-4 pt-4 border-t border-gray-100 flex justify-between">
          <span className="text-sm text-gray-500">APY: {account.apy || '0.5'}%</span>
          <button className="text-sm text-indigo-600 hover:text-indigo-800">
            View Details
          </button>
        </div>
      </div>
    </div>
  )
}