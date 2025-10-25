import React, { useState } from 'react';
import Header from './components/Header';
import Sidenav from './components/Sidenav';
import ProductForm from './components/ProductForm';

interface Product {
  id: string;
  nombreProducto: string;
  descripcion: string;
  precio: number;
  fechaCreacion: Date;
}

const App: React.FC = () => {
  const [isMenuOpen, setIsMenuOpen] = useState(true);
  const [activeSection, setActiveSection] = useState('bodega');
  const [products, setProducts] = useState<Product[]>([]);

  const handleProductSubmit = (productData: { nombreProducto: string; descripcion: string; precio: number }) => {
    const newProduct: Product = {
      id: Date.now().toString(),
      ...productData,
      fechaCreacion: new Date()
    };
    setProducts(prev => [...prev, newProduct]);
    alert('Producto creado exitosamente!');
  };

  const renderContent = () => {
    if (activeSection === 'bodega-4') { // Gestión de Productos
      return (
        <div className="space-y-8">
          <ProductForm onSubmit={handleProductSubmit} />
          
          {products.length > 0 && (
            <div className="bg-white rounded-xl shadow-xl p-6 border border-gray-100">
              <h2 className="text-2xl font-bold text-gray-800 mb-6">Productos Creados</h2>
              <div className="overflow-x-auto">
                <table className="min-w-full divide-y divide-gray-200">
                  <thead className="bg-gray-50">
                    <tr>
                      <th className="px-6 py-4 text-left text-xs font-semibold text-gray-600 uppercase tracking-wider">
                        Nombre
                      </th>
                      <th className="px-6 py-4 text-left text-xs font-semibold text-gray-600 uppercase tracking-wider">
                        Descripción
                      </th>
                      <th className="px-6 py-4 text-left text-xs font-semibold text-gray-600 uppercase tracking-wider">
                        Precio
                      </th>
                      <th className="px-6 py-4 text-left text-xs font-semibold text-gray-600 uppercase tracking-wider">
                        Fecha Creación
                      </th>
                    </tr>
                  </thead>
                  <tbody className="bg-white divide-y divide-gray-200">
                    {products.map((product) => (
                      <tr key={product.id} className="hover:bg-gray-50 transition-colors">
                        <td className="px-6 py-4 whitespace-nowrap text-sm font-medium text-gray-900">
                          {product.nombreProducto}
                        </td>
                        <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-500">
                          {product.descripcion}
                        </td>
                        <td className="px-6 py-4 whitespace-nowrap text-sm font-semibold text-green-600">
                          ${product.precio.toLocaleString()}
                        </td>
                        <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-500">
                          {product.fechaCreacion.toLocaleDateString()}
                        </td>
                      </tr>
                    ))}
                  </tbody>
                </table>
              </div>
            </div>
          )}
        </div>
      );
    }

    return (
      <div className="bg-white rounded-xl shadow-xl p-8 border border-gray-100">
        <h1 className="text-3xl font-bold text-gray-800 mb-4">
          Bienvenido a SisPro ERP
        </h1>
        <p className="text-gray-600 text-lg">
          Seleccione una opción del menú lateral para comenzar.
        </p>
      </div>
    );
  };

  return (
    <div className="min-h-screen bg-gray-50">
      <Header onMenuToggle={() => setIsMenuOpen(!isMenuOpen)} isMenuOpen={isMenuOpen} />
      
      <div className="flex pt-16">
        {isMenuOpen && <Sidenav activeSection={activeSection} onSectionChange={setActiveSection} />}
        
        <main className={`flex-1 p-8 transition-all duration-300 ${isMenuOpen ? 'ml-64' : 'ml-0'}`}>
          {renderContent()}
        </main>
      </div>

      {/* Footer */}
      <footer className="fixed bottom-0 left-0 right-0 bg-white border-t border-gray-200 py-3 px-6 shadow-lg">
        <p className="text-sm text-gray-500 text-center">
          © 2025 - Aplicación Desarrollada por ServProTech SpA
        </p>
      </footer>
    </div>
  );
};

export default App;
