import React, { useState } from 'react';
import Header from './components/Header';
import Sidenav from './components/Sidenav';
import ProductForm from './components/ProductForm';
import RegistroHHForm from './components/RegistroHHForm';

interface Product {
  id: string;
  nombreProducto: string;
  descripcion: string;
  precio: number;
  fechaCreacion: Date;
}

interface RegistroHH {
  id: string;
  fechaEjecucion: string;
  orderDeTrabajo: string;
  actividad: string;
  trabajadores: {
    trabajador: string;
    horasNormales: number;
    horasExtras: number;
  }[];
  fechaCreacion: Date;
}

const App: React.FC = () => {
  const [isMenuOpen, setIsMenuOpen] = useState(true);
  const [activeSection, setActiveSection] = useState('bodega');
  const [products, setProducts] = useState<Product[]>([]);
  const [registrosHH, setRegistrosHH] = useState<RegistroHH[]>([]);

  const handleProductSubmit = (productData: { nombreProducto: string; descripcion: string; precio: number }) => {
    const newProduct: Product = {
      id: Date.now().toString(),
      ...productData,
      fechaCreacion: new Date()
    };
    setProducts(prev => [...prev, newProduct]);
    alert('Producto creado exitosamente!');
  };

  const handleRegistroHHSubmit = (registroData: { fechaEjecucion: string; orderDeTrabajo: string; actividad: string; trabajadores: { trabajador: string; horasNormales: number; horasExtras: number; }[] }) => {
    const newRegistro: RegistroHH = {
      id: Date.now().toString(),
      ...registroData,
      fechaCreacion: new Date()
    };
    setRegistrosHH(prev => [...prev, newRegistro]);
    alert(`Registro HH creado exitosamente para ${registroData.trabajadores.length} trabajador(es)!`);
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

    if (activeSection === 'produccion-2') { // Registro Manual HH por OT
      return (
        <div className="space-y-8">
          <RegistroHHForm onSubmit={handleRegistroHHSubmit} />
          
          {registrosHH.length > 0 && (
            <div className="bg-white rounded-xl shadow-xl p-6 border border-gray-100">
              <h2 className="text-2xl font-bold text-gray-800 mb-6">Registros HH Creados</h2>
              <div className="space-y-6">
                {registrosHH.map((registro) => (
                  <div key={registro.id} className="border border-gray-200 rounded-lg p-4">
                    <div className="grid grid-cols-1 md:grid-cols-4 gap-4 mb-4">
                      <div>
                        <span className="text-sm font-semibold text-gray-600">Fecha Ejecución:</span>
                        <p className="text-sm text-gray-900">{new Date(registro.fechaEjecucion).toLocaleDateString()}</p>
                      </div>
                      <div>
                        <span className="text-sm font-semibold text-gray-600">Orden Trabajo:</span>
                        <p className="text-sm text-gray-900">{registro.orderDeTrabajo || 'N/A'}</p>
                      </div>
                      <div>
                        <span className="text-sm font-semibold text-gray-600">Actividad:</span>
                        <p className="text-sm text-gray-900">{registro.actividad}</p>
                      </div>
                      <div>
                        <span className="text-sm font-semibold text-gray-600">Fecha Creación:</span>
                        <p className="text-sm text-gray-900">{registro.fechaCreacion.toLocaleDateString()}</p>
                      </div>
                    </div>
                    
                    <div className="overflow-x-auto">
                      <table className="min-w-full divide-y divide-gray-200">
                        <thead className="bg-gray-50">
                          <tr>
                            <th className="px-4 py-3 text-left text-xs font-semibold text-gray-600 uppercase tracking-wider">
                              Trabajador
                            </th>
                            <th className="px-4 py-3 text-left text-xs font-semibold text-gray-600 uppercase tracking-wider">
                              H. Normales
                            </th>
                            <th className="px-4 py-3 text-left text-xs font-semibold text-gray-600 uppercase tracking-wider">
                              H. Extras
                            </th>
                            <th className="px-4 py-3 text-left text-xs font-semibold text-gray-600 uppercase tracking-wider">
                              Total Horas
                            </th>
                          </tr>
                        </thead>
                        <tbody className="bg-white divide-y divide-gray-200">
                          {registro.trabajadores.map((trabajador, index) => (
                            <tr key={index} className="hover:bg-gray-50 transition-colors">
                              <td className="px-4 py-3 whitespace-nowrap text-sm font-medium text-gray-900">
                                {trabajador.trabajador}
                              </td>
                              <td className="px-4 py-3 whitespace-nowrap text-sm font-semibold text-blue-600">
                                {trabajador.horasNormales}h
                              </td>
                              <td className="px-4 py-3 whitespace-nowrap text-sm font-semibold text-orange-600">
                                {trabajador.horasExtras}h
                              </td>
                              <td className="px-4 py-3 whitespace-nowrap text-sm font-bold text-green-600">
                                {trabajador.horasNormales + trabajador.horasExtras}h
                              </td>
                            </tr>
                          ))}
                        </tbody>
                        <tfoot className="bg-gray-100">
                          <tr>
                            <td className="px-4 py-3 text-sm font-bold text-gray-900">TOTALES:</td>
                            <td className="px-4 py-3 text-sm font-bold text-blue-600">
                              {registro.trabajadores.reduce((sum, t) => sum + t.horasNormales, 0)}h
                            </td>
                            <td className="px-4 py-3 text-sm font-bold text-orange-600">
                              {registro.trabajadores.reduce((sum, t) => sum + t.horasExtras, 0)}h
                            </td>
                            <td className="px-4 py-3 text-sm font-bold text-green-600">
                              {registro.trabajadores.reduce((sum, t) => sum + t.horasNormales + t.horasExtras, 0)}h
                            </td>
                          </tr>
                        </tfoot>
                      </table>
                    </div>
                  </div>
                ))}
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
