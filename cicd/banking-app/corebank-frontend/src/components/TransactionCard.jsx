import { FaArrowUp, FaArrowDown } from 'react-icons/fa'

export default function TransactionCard({ transaction }) {
  const isDeposit = transaction.type === 'deposit'
  
  return (
    <div className="flex items-center justify-between p-4 border border-gray-200 rounded-lg hover:bg-gray-50">
      <div className="flex items-center space-x-4">
        <div className={`p-3 rounded-full ${isDeposit ? 'bg-green-100 text-green-600' : 'bg-red-100 text-red-600'}`}>
          {isDeposit ? <FaArrowDown /> : <FaArrowUp />}
        </div>
        <div>
          <p className="font-medium">{transaction.description}</p>
          <p className="text-sm text-gray-500">
            {new Date(transaction.date).toLocaleDateString()}
          </p>
        </div>
      </div>
      <div className={`font-bold ${isDeposit ? 'text-green-600' : 'text-red-600'}`}>
        {isDeposit ? '+' : '-'}${transaction.amount.toFixed(2)}
      </div>
    </div>
  )
}