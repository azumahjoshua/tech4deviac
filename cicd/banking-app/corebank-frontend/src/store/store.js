import { configureStore } from '@reduxjs/toolkit'
import accountReducer from './slices/accountSlice'
import transactionReducer from './slices/transactionSlice'

export default configureStore({
  reducer: {
    accounts: accountReducer,
    transactions: transactionReducer,
    auth: (state = { user: { name: 'John Doe' } }) => state
  },
  middleware: (getDefaultMiddleware) =>
    getDefaultMiddleware({
      serializableCheck: {
        // Ignore these field paths in all actions
        ignoredActionPaths: ['payload.date'],
        // Ignore these paths in the state
        ignoredPaths: ['transactions.items.date']
      }
    })
})