import { useState } from 'react'
import reactLogo from './assets/react.svg'
import viteLogo from '/vite.svg'
import './App.css'

function App() {
  const [count, setCount] = useState(0)

  return (
    <div className="min-h-screen bg-gradient-to-br from-blue-50 to-indigo-100 flex flex-col items-center justify-center p-8">
      <div className="flex gap-8 mb-8">
        <a href="https://vite.dev" target="_blank" className="hover:scale-110 transition-transform">
          <img src={viteLogo} className="logo h-16 w-16" alt="Vite logo" />
        </a>
        <a href="https://react.dev" target="_blank" className="hover:scale-110 transition-transform">
          <img src={reactLogo} className="logo react h-16 w-16" alt="React logo" />
        </a>
      </div>
      
      <h1 className="text-4xl font-bold text-gray-800 mb-8">Vite + React + Tailwind</h1>
      
      <div className="bg-white rounded-lg shadow-lg p-8 max-w-md w-full text-center">
        <button 
          onClick={() => setCount((count) => count + 1)}
          className="bg-blue-500 hover:bg-blue-600 text-white font-semibold py-3 px-6 rounded-lg transition-colors duration-200 mb-4"
        >
          count is {count}
        </button>
        <p className="text-gray-600 mb-4">
          Edit <code className="bg-gray-100 px-2 py-1 rounded text-sm">src/App.tsx</code> and save to test HMR
        </p>
        <div className="bg-green-50 border border-green-200 rounded-lg p-4">
          <p className="text-green-700 font-medium">🎉 ¡Tailwind CSS está funcionando!</p>
          <p className="text-green-600 text-sm mt-1">Las clases de Tailwind se están aplicando correctamente.</p>
        </div>
      </div>
      
      <p className="text-gray-500 mt-8 text-center max-w-md">
        Click on the Vite and React logos to learn more
      </p>
    </div>
  )
}

export default App
