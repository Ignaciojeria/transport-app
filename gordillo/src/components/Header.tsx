import React from 'react';
import { Menu, X, User, ChevronDown } from 'lucide-react';

interface HeaderProps {
  onMenuToggle: () => void;
  isMenuOpen: boolean;
}

const Header: React.FC<HeaderProps> = ({ onMenuToggle, isMenuOpen }) => {
  return (
    <header className="bg-white shadow-lg border-b border-gray-200 h-16 flex items-center justify-between px-6 fixed top-0 left-0 right-0 z-20">
      <div className="flex items-center space-x-4">
        {/* Logo Grupo Gordillo */}
        <div className="flex items-center space-x-3">
          <div className="w-10 h-10 bg-gordillo-yellow flex items-center justify-center rounded-md shadow-sm">
            <div className="text-xs font-bold text-gray-800 text-center leading-tight">
              <div>Grupo</div>
              <div>Gordillo</div>
            </div>
          </div>
          <span className="text-xl font-bold text-gray-800">SisPro ERP</span>
        </div>

        {/* Hamburger Menu */}
        <button
          onClick={onMenuToggle}
          className="p-2 rounded-md hover:bg-gray-100 transition-colors duration-200"
        >
          {isMenuOpen ? <X className="w-5 h-5 text-gray-600" /> : <Menu className="w-5 h-5 text-gray-600" />}
        </button>
      </div>

      <div className="flex items-center space-x-6">
        {/* Version */}
        <div className="flex items-center space-x-2 text-gray-500">
          <span className="text-sm font-medium">v1.8.1</span>
        </div>

        {/* User Profile */}
        <div className="flex items-center space-x-3 bg-gray-50 px-3 py-2 rounded-lg hover:bg-gray-100 transition-colors cursor-pointer">
          <User className="w-5 h-5 text-gray-600" />
          <span className="text-sm font-medium text-gray-800">Alexander Maverick Gutierrez</span>
          <ChevronDown className="w-4 h-4 text-gray-600" />
        </div>
      </div>
    </header>
  );
};

export default Header;
