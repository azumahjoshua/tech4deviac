import { useState } from 'react'
import { useDispatch } from 'react-redux'
import { addAccount } from '../store/slices/accountSlice'
import axios from 'axios'

export default function AccountModal({ show, onHide }) {
  const dispatch = useDispatch()
  const [formData, setFormData] = useState({
    name: '',
    email: '',
    type: 'checking',
    initialBalance: ''
  })
  const [isLoading, setIsLoading] = useState(false)
  const [error, setError] = useState(null)

  const handleSubmit = async (e) => {
    e.preventDefault()
    setIsLoading(true)
    setError(null)

    // Create account data
    const newAccount = {
      owner: formData.name,
      email: formData.email,
      type: formData.type,
      balance: parseFloat(formData.initialBalance) || 0,
    }

    try {
      // Use the proxy path (/api) instead of direct URL
      const response = await axios.post('/accounts', newAccount)

      if (response.status === 201) { // 201 Created is more appropriate
        dispatch(addAccount(response.data))
        onHide()
        resetForm()
      }
    } catch (error) {
      console.error('Error creating account:', error)
      setError(error.response?.data?.message || 'Error creating account. Please try again.')
    } finally {
      setIsLoading(false)
    }
  }

  const resetForm = () => {
    setFormData({
      name: '',
      email: '',
      type: 'checking',
      initialBalance: ''
    })
  }

  const handleChange = (e) => {
    const { name, value } = e.target
    setFormData(prev => ({
      ...prev,
      [name]: value
    }))
  }

  if (!show) return null

  return (
    <div className="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center p-4 z-50">
      <div className="bg-white rounded-lg shadow-xl w-full max-w-md">
        <div className="flex justify-between items-center border-b p-4">
          <h3 className="text-lg font-semibold">New Account</h3>
          <button 
            onClick={onHide} 
            className="text-gray-500 hover:text-gray-700"
            disabled={isLoading}
          >
            &times;
          </button>
        </div>

        <form onSubmit={handleSubmit} className="p-4 space-y-4">
          {error && (
            <div className="bg-red-50 text-red-600 p-3 rounded-md text-sm">
              {error}
            </div>
          )}

          <div>
            <label className="block text-sm font-medium text-gray-700 mb-1">
              Account Name
            </label>
            <input
              type="text"
              name="name"
              value={formData.name}
              onChange={handleChange}
              className="w-full border border-gray-300 rounded-md p-2 focus:ring-2 focus:ring-indigo-500 focus:border-indigo-500"
              required
              disabled={isLoading}
            />
          </div>

          <div>
            <label className="block text-sm font-medium text-gray-700 mb-1">
              Email
            </label>
            <input
              type="email"
              name="email"
              value={formData.email}
              onChange={handleChange}
              className="w-full border border-gray-300 rounded-md p-2 focus:ring-2 focus:ring-indigo-500 focus:border-indigo-500"
              required
              disabled={isLoading}
            />
          </div>

          <div>
            <label className="block text-sm font-medium text-gray-700 mb-1">
              Account Type
            </label>
            <select
              name="type"
              value={formData.type}
              onChange={handleChange}
              className="w-full border border-gray-300 rounded-md p-2 focus:ring-2 focus:ring-indigo-500 focus:border-indigo-500"
              disabled={isLoading}
            >
              <option value="checking">Checking</option>
              <option value="savings">Savings</option>
              <option value="business">Business</option>
            </select>
          </div>

          <div>
            <label className="block text-sm font-medium text-gray-700 mb-1">
              Initial Balance
            </label>
            <div className="relative">
              <span className="absolute inset-y-0 left-0 flex items-center pl-3 text-gray-500">$</span>
              <input
                type="number"
                name="initialBalance"
                value={formData.initialBalance}
                onChange={handleChange}
                min="0"
                step="0.01"
                className="w-full border border-gray-300 rounded-md p-2 pl-8 focus:ring-2 focus:ring-indigo-500 focus:border-indigo-500"
                required
                disabled={isLoading}
              />
            </div>
          </div>

          <div className="flex justify-end space-x-3 pt-4">
            <button
              type="button"
              onClick={onHide}
              className="px-4 py-2 border border-gray-300 rounded-md text-gray-700 hover:bg-gray-50 disabled:opacity-50"
              disabled={isLoading}
            >
              Cancel
            </button>
            <button
              type="submit"
              className="px-4 py-2 bg-indigo-600 text-white rounded-md hover:bg-indigo-700 disabled:opacity-50 disabled:cursor-not-allowed"
              disabled={isLoading}
            >
              {isLoading ? 'Creating...' : 'Create Account'}
            </button>
          </div>
        </form>
      </div>
    </div>
  )
}


// import { useState } from 'react'
// import { useDispatch } from 'react-redux'
// import { addAccount } from '../store/slices/accountSlice'
// import axios from 'axios'

// export default function AccountModal({ show, onHide }) {
//   const dispatch = useDispatch()
//   const [formData, setFormData] = useState({
//     name: '',
//     email: '',
//     type: 'checking',
//     initialBalance: ''
//   })

//   const handleSubmit = async (e) => {
//     e.preventDefault()

//     // Create account data
//     const newAccount = {
//       owner: formData.name,
//       email: formData.email,
//       balance: parseFloat(formData.initialBalance),
//     }

//     try {
//       // Send a POST request to create the account
//       const response = await axios.post('http://localhost:8080/accounts', newAccount)

//       if (response.status === 200) {
//         // Dispatch the newly created account to Redux
//         dispatch(addAccount(response.data))

//         // Close the modal
//         onHide()

//         // Reset form data
//         setFormData({
//           name: '',
//           email: '',
//           type: 'checking',
//           initialBalance: ''
//         })
//       }
//     } catch (error) {
//       console.error('Error creating account:', error)
//       alert('Error creating account. Please try again.')
//     }
//   }

//   const handleChange = (e) => {
//     const { name, value } = e.target
//     setFormData((prev) => ({
//       ...prev,
//       [name]: value
//     }))
//   }

//   return (
//     show && (
//       <div className="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center p-4 z-50">
//         <div className="bg-white rounded-lg shadow-xl w-full max-w-md">
//           <div className="flex justify-between items-center border-b p-4">
//             <h3 className="text-lg font-semibold">New Account</h3>
//             <button onClick={onHide} className="text-gray-500 hover:text-gray-700">
//               &times;
//             </button>
//           </div>

//           <form onSubmit={handleSubmit} className="p-4 space-y-4">
//             <div>
//               <label className="block text-sm font-medium text-gray-700 mb-1">Account Name</label>
//               <input
//                 type="text"
//                 name="name"
//                 value={formData.name}
//                 onChange={handleChange}
//                 className="w-full border border-gray-300 rounded-md p-2"
//                 required
//               />
//             </div>

//             <div>
//               <label className="block text-sm font-medium text-gray-700 mb-1">Email</label>
//               <input
//                 type="email"
//                 name="email"
//                 value={formData.email}
//                 onChange={handleChange}
//                 className="w-full border border-gray-300 rounded-md p-2"
//                 required
//               />
//             </div>

//             <div>
//               <label className="block text-sm font-medium text-gray-700 mb-1">Initial Balance</label>
//               <div className="relative">
//                 <span className="absolute inset-y-0 left-0 flex items-center pl-3 text-gray-500">$</span>
//                 <input
//                   type="number"
//                   name="initialBalance"
//                   value={formData.initialBalance}
//                   onChange={handleChange}
//                   min="0"
//                   step="0.01"
//                   className="w-full border border-gray-300 rounded-md p-2 pl-8"
//                   required
//                 />
//               </div>
//             </div>

//             <div className="flex justify-end space-x-3 pt-4">
//               <button
//                 type="button"
//                 onClick={onHide}
//                 className="px-4 py-2 border border-gray-300 rounded-md text-gray-700 hover:bg-gray-50"
//               >
//                 Cancel
//               </button>
//               <button
//                 type="submit"
//                 className="px-4 py-2 bg-indigo-600 text-white rounded-md hover:bg-indigo-700"
//               >
//                 Create Account
//               </button>
//             </div>
//           </form>
//         </div>
//       </div>
//     )
//   )
// }
