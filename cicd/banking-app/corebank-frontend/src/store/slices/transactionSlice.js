import { createSlice } from '@reduxjs/toolkit'

const initialState = {
  items: [
    {
      id: '1',
      accountId: '1',
      type: 'deposit',
      amount: 1000,
      description: 'Initial deposit',
      date: '2025-03-01T10:30:00Z'
    },
    {
      id: '2',
      accountId: '1',
      type: 'withdrawal',
      amount: 200,
      description: 'ATM withdrawal',
      date: '2025-03-05T14:15:00Z'
    }
  ],
  recent: []
}

export const transactionSlice = createSlice({
  name: 'transactions',
  initialState,
  reducers: {
    addTransaction: (state, action) => {
      state.items.unshift(action.payload)
      state.recent = state.items.slice(0, 5)
    },
    setRecentTransactions: (state) => {
      state.recent = state.items.slice(0, 5)
    }
  }
})

export const { addTransaction, setRecentTransactions } = transactionSlice.actions
export default transactionSlice.reducer