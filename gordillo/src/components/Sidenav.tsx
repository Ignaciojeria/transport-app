import React from 'react';
import { ChevronDown, ChevronRight, Users, Wrench, Eye, ShoppingCart, Package } from 'lucide-react';

interface SidenavProps {
  activeSection: string;
  onSectionChange: (section: string) => void;
}

const Sidenav: React.FC<SidenavProps> = ({ activeSection, onSectionChange }) => {
  const menuItems = [
    {
      id: 'recursos-humanos',
      label: 'Recursos Humanos',
      icon: Users,
      hasSubmenu: true
    },
    {
      id: 'produccion',
      label: 'Producción',
      icon: Wrench,
      hasSubmenu: true
    },
    {
      id: 'modulo-vigia',
      label: 'Módulo Vigía',
      icon: Eye,
      hasSubmenu: true
    },
    {
      id: 'adquisiciones',
      label: 'Adquisiciones',
      icon: ShoppingCart,
      hasSubmenu: true
    },
    {
      id: 'bodega',
      label: 'Bodega',
      icon: Package,
      hasSubmenu: true,
      submenu: [
        'Recepcionar Item O/C',
        'Crear Retiro',
        'Busqueda y Edición de Retiros',
        'Indice Confirmar Retiros',
        'Gestión de Productos'
      ]
    }
  ];

  return (
    <div className="w-64 bg-gray-900 text-white h-screen fixed left-0 top-16 z-10 shadow-xl">
      <div className="p-6 border-b border-gray-700">
        <h2 className="text-xl font-bold text-white">PANEL</h2>
      </div>
      
      <nav className="px-4 py-4">
        {menuItems.map((item) => (
          <div key={item.id} className="mb-2">
            <button
              onClick={() => onSectionChange(item.id)}
              className={`w-full flex items-center justify-between p-3 rounded-lg transition-all duration-200 ${
                activeSection === item.id 
                  ? 'bg-blue-600 text-white shadow-md' 
                  : 'text-gray-300 hover:bg-gray-800 hover:text-white'
              }`}
            >
              <div className="flex items-center">
                <item.icon className="w-5 h-5 mr-3" />
                <span className="text-sm font-medium">{item.label}</span>
              </div>
              {item.hasSubmenu && (
                activeSection === item.id ? 
                  <ChevronDown className="w-4 h-4" /> : 
                  <ChevronRight className="w-4 h-4" />
              )}
            </button>
            
            {activeSection === item.id && item.submenu && (
              <div className="ml-8 mt-2 space-y-1 bg-gray-800 rounded-lg p-2">
                {item.submenu.map((subItem, index) => (
                  <button
                    key={index}
                    onClick={() => onSectionChange(`${item.id}-${index}`)}
                    className="block w-full text-left p-2 text-sm text-gray-400 hover:text-white hover:bg-gray-700 rounded transition-colors duration-200"
                  >
                    • {subItem}
                  </button>
                ))}
              </div>
            )}
          </div>
        ))}
      </nav>
    </div>
  );
};

export default Sidenav;
