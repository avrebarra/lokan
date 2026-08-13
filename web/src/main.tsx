import React from 'react'
import ReactDOM from 'react-dom/client'
import App from './App'
import './tokens.css'
import './index.css'

// resolve theme: stored preference, else light (dark is opt-in via the toggle)
const stored = localStorage.getItem('lokan-theme')
document.documentElement.dataset.theme = stored ?? 'light'

ReactDOM.createRoot(document.getElementById('root')!).render(
  <React.StrictMode>
    <App />
  </React.StrictMode>,
)
