import React, { useState } from 'react';

interface ProductFormData {
  nombreProducto: string;
  descripcion: string;
  precio: number;
}

interface ProductFormProps {
  onSubmit: (product: ProductFormData) => void;
}

const ProductForm: React.FC<ProductFormProps> = ({ onSubmit }) => {
  const [formData, setFormData] = useState<ProductFormData>({
    nombreProducto: '',
    descripcion: '',
    precio: 0
  });

  const handleInputChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    const { name, value } = e.target;
    setFormData(prev => ({
      ...prev,
      [name]: name === 'precio' ? parseFloat(value) || 0 : value
    }));
  };

  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault();
    if (formData.nombreProducto && formData.descripcion && formData.precio > 0) {
      onSubmit(formData);
      setFormData({
        nombreProducto: '',
        descripcion: '',
        precio: 0
      });
    }
  };

  return (
    <div className="bg-white rounded-xl shadow-xl p-8 max-w-2xl mx-auto border border-gray-100">
      <div className="mb-8">
        <h1 className="text-3xl font-bold text-gray-800 mb-2">Crear</h1>
        <p className="text-gray-600 text-lg">Producto</p>
      </div>

      <form onSubmit={handleSubmit} className="space-y-8">
        <div className="space-y-2">
          <label htmlFor="nombreProducto" className="block text-sm font-semibold text-gray-700">
            Nombre:
          </label>
          <input
            type="text"
            id="nombreProducto"
            name="nombreProducto"
            value={formData.nombreProducto}
            onChange={handleInputChange}
            className="w-full px-4 py-3 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500 focus:border-transparent transition-all duration-200 text-gray-700"
            placeholder="Ingrese el nombre del producto"
            required
          />
        </div>

        <div className="space-y-2">
          <label htmlFor="descripcion" className="block text-sm font-semibold text-gray-700">
            Descripción:
          </label>
          <input
            type="text"
            id="descripcion"
            name="descripcion"
            value={formData.descripcion}
            onChange={handleInputChange}
            className="w-full px-4 py-3 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500 focus:border-transparent transition-all duration-200 text-gray-700"
            placeholder="Ingrese la descripción del producto"
            required
          />
        </div>

        <div className="space-y-2">
          <label htmlFor="precio" className="block text-sm font-semibold text-gray-700">
            Precio:
          </label>
          <input
            type="number"
            id="precio"
            name="precio"
            value={formData.precio}
            onChange={handleInputChange}
            className="w-full px-4 py-3 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500 focus:border-transparent transition-all duration-200 text-gray-700"
            placeholder="Ingrese el precio del producto"
            min="0"
            step="0.01"
            required
          />
        </div>

        <div className="flex justify-center pt-6">
          <button
            type="submit"
            className="bg-blue-600 text-white px-12 py-4 rounded-lg font-semibold hover:bg-blue-700 transition-all duration-200 focus:outline-none focus:ring-2 focus:ring-blue-500 focus:ring-offset-2 shadow-lg hover:shadow-xl transform hover:-translate-y-0.5"
          >
            CREAR
          </button>
        </div>
      </form>
    </div>
  );
};

export default ProductForm;
