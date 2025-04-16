import { Chart as ChartJS, ArcElement, Tooltip, Legend } from 'chart.js'
import { Doughnut } from 'react-chartjs-2'
import { FaCircle } from 'react-icons/fa'

ChartJS.register(ArcElement, Tooltip, Legend)

export default function BalanceChart({ accounts }) {
  // Prepare chart data
  const data = {
    labels: accounts.map(account => account.name),
    datasets: [
      {
        data: accounts.map(account => account.balance),
        backgroundColor: [
          '#3B82F6', // blue-500
          '#10B981', // emerald-500
          '#8B5CF6', // violet-500
          '#EC4899', // pink-500
          '#F59E0B'  // amber-500
        ],
        borderWidth: 0,
      },
    ],
  }

  const options = {
    cutout: '70%',
    plugins: {
      legend: {
        display: false,
      },
      tooltip: {
        callbacks: {
          label: (context) => {
            const label = context.label || ''
            const value = context.raw || 0
            return `${label}: $${value.toFixed(2)}`
          }
        }
      }
    },
    maintainAspectRatio: false
  }

  // Custom legend component
  const LegendItem = ({ account, color }) => (
    <div className="flex items-center mt-2">
      <FaCircle className="text-xs mr-2" style={{ color }} />
      <span className="text-sm text-gray-600">
        {account.name}: <span className="font-medium">${account.balance.toFixed(2)}</span>
      </span>
    </div>
  )

  return (
    <div className="mt-4">
      <div className="relative h-48">
        <Doughnut data={data} options={options} />
      </div>
      
      <div className="mt-4">
        {accounts.map((account, index) => (
          <LegendItem 
            key={account.id}
            account={account}
            color={data.datasets[0].backgroundColor[index]}
          />
        ))}
      </div>
    </div>
  )
}