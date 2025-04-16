import { BrowserRouter as Router, Routes, Route } from 'react-router-dom'
import { Provider } from 'react-redux'
import { useEffect } from 'react'
import { useDispatch } from 'react-redux'
import store from './store/store'
import Navbar from './components/Navbar'
import Dashboard from './pages/Dashboard'
import Accounts from './pages/Accounts'
import Transactions from './pages/Transactions'
import { setRecentTransactions } from './store/slices/transactionSlice'

// Component to initialize app data
function AppInitializer({ children }) {
  const dispatch = useDispatch()

  useEffect(() => {
    // Initialize recent transactions
    dispatch(setRecentTransactions())
    
    // Here you could add other initialization logic like:
    // - Fetching initial account data from API
    // - Checking authentication status
    // - Loading user preferences
  }, [dispatch])

  return children
}

function App() {
  return (
    <Provider store={store}>
      <Router>
        <AppInitializer>
          <div className="min-h-screen bg-gray-50">
            <Navbar />
            <main className="container mx-auto px-4 py-6">
              <Routes>
                <Route path="/" element={<Dashboard />} />
                <Route path="/accounts" element={<Accounts />} />
                <Route path="/transactions" element={<Transactions />} />
              </Routes>
            </main>
          </div>
        </AppInitializer>
      </Router>
    </Provider>
  )
}

export default App