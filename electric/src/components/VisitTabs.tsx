import { Play, CheckCircle, XCircle } from 'lucide-react'

interface VisitTabsProps {
  activeTab: 'en-ruta' | 'entregados' | 'no-entregados'
  onTabChange: (tab: 'en-ruta' | 'entregados' | 'no-entregados') => void
  inRouteUnits: number
  deliveredUnits: number
  notDeliveredUnits: number
}

export function VisitTabs({ 
  activeTab, 
  onTabChange, 
  inRouteUnits, 
  deliveredUnits, 
  notDeliveredUnits 
}: VisitTabsProps) {
  
  const getTabIcon = (tabId: 'en-ruta' | 'entregados' | 'no-entregados') => {
    switch (tabId) {
      case 'en-ruta':
        return <Play className="w-3 h-3" />
      case 'entregados':
        return <CheckCircle className="w-3 h-3" />
      case 'no-entregados':
        return <XCircle className="w-3 h-3" />
      default:
        return null
    }
  }

  const tabs = [
    {
      id: 'en-ruta' as const,
      label: 'En ruta',
      count: inRouteUnits,
      icon: getTabIcon('en-ruta')
    },
    {
      id: 'entregados' as const,
      label: 'Entregados',
      count: deliveredUnits,
      icon: getTabIcon('entregados')
    },
    {
      id: 'no-entregados' as const,
      label: 'No entregados',
      count: notDeliveredUnits,
      icon: getTabIcon('no-entregados')
    }
  ]

  return (
    <div className="sticky top-0 z-20 bg-white/80 backdrop-blur border-b">
      <div className="flex">
        {tabs.map((tab) => (
          <button
            key={tab.id}
            onClick={() => onTabChange(tab.id)}
            className={`flex-1 py-3 px-2 text-center text-xs font-medium transition-all duration-200 border-b-2 ${
              activeTab === tab.id
                ? 'bg-gradient-to-r from-blue-50 to-indigo-50 border-indigo-500 text-indigo-700'
                : 'bg-gray-50 border-transparent text-gray-500 hover:bg-gray-100'
            }`}
          >
            <div className="flex flex-col items-center space-y-1">
              <div className="flex items-center space-x-1">
                {tab.icon}
                <span className="truncate">{tab.label}</span>
              </div>
              <span className={`${
                activeTab === tab.id ? 'bg-indigo-200 text-indigo-800' : 'bg-gray-200 text-gray-600'
              } px-2 py-0.5 rounded-full text-xs font-bold`}>
                ({tab.count})
              </span>
            </div>
          </button>
        ))}
      </div>
    </div>
  )
}
