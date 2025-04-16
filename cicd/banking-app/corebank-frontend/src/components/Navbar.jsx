import { Link } from 'react-router-dom'
import { useSelector } from 'react-redux'
import { FaUser } from 'react-icons/fa'

export default function Navbar() {
  const { user } = useSelector((state) => state.auth)
  
  return (
    <nav className="bg-primary text-white shadow-lg">
      <div className="container max-w-7xl px-4 py-3">
        <div className="flex justify-between items-center">
          <Link to="/" className="text-xl font-bold text-[oklch(var(--card))]">
            CoreBank
          </Link>
          
          <div className="hidden md:flex space-x-6">
            <Link 
              to="/" 
              className="text-[oklch(var(--card)/0.9)] hover:text-[oklch(var(--card))] hover:bg-primary-dark px-3 py-2 rounded-lg transition-colors"
            >
              Dashboard
            </Link>
            <Link 
              to="/accounts" 
              className="text-[oklch(var(--card)/0.9)] hover:text-[oklch(var(--card))] hover:bg-primary-dark px-3 py-2 rounded-lg transition-colors"
            >
              Accounts
            </Link>
            <Link 
              to="/transactions" 
              className="text-[oklch(var(--card)/0.9)] hover:text-[oklch(var(--card))] hover:bg-primary-dark px-3 py-2 rounded-lg transition-colors"
            >
              Transactions
            </Link>
            
          </div>
          
          <div className="flex items-center space-x-2">
            <FaUser className="text-sm text-[oklch(var(--card)/0.8)]" />
            <span className="text-[oklch(var(--card)/0.9)]">
              Welcome, {user?.name || 'User'}
            </span>
          </div>
        </div>
      </div>
    </nav>
  )
}
