import { createSlice } from '@reduxjs/toolkit'

const initialState = {
  items: [
    { id: '1', name: 'Checking Account', balance: 2500.75 },
    { id: '2', name: 'Savings Account', balance: 15000.50 }
  ]
}

export const accountSlice = createSlice({
  name: 'accounts',
  initialState,
  reducers: {
    addAccount: (state, action) => {
      state.items.push(action.payload)
    },
    updateBalance: (state, action) => {
      const { accountId, amount } = action.payload
      const account = state.items.find(acc => acc.id === accountId)
      if (account) {
        account.balance += amount
      }
    }
  }
})

export const { addAccount, updateBalance } = accountSlice.actions
export default accountSlice.reducer