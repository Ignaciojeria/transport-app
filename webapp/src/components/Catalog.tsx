import React, { useState } from 'react'
import { type ElectricTenantData } from '../db/collections'
import { type CreateProductRequest, type Product } from '../types/product'
import CreateProduct from './CreateProduct'
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from './ui/Card'
import { 
  Plus, 
  Package, 
  Edit, 
  Trash2, 
  Eye,
  Search,
  Filter
} from 'lucide-react'

interface CatalogProps {
  tenant: ElectricTenantData
}

const Catalog: React.FC<CatalogProps> = ({ tenant }) => {
  const [showCreateForm, setShowCreateForm] = useState(false)
  const [products, setProducts] = useState<Product[]>([])
  const [searchTerm, setSearchTerm] = useState('')

  const handleCreateProduct = async (productData: CreateProductRequest) => {
    try {
      // Aquí implementarías la llamada a la API para crear el producto
      console.log('Creando producto:', productData)
      
      // Simular creación exitosa
      const newProduct: Product = {
        ...productData,
        // Agregar ID generado por el servidor
        referenceID: productData.referenceID || `PROD_${Date.now()}`
      }
      
      setProducts(prev => [...prev, newProduct])
      setShowCreateForm(false)
      
      // Mostrar mensaje de éxito
      alert('Producto creado exitosamente')
    } catch (error) {
      console.error('Error al crear producto:', error)
      alert('Error al crear el producto')
    }
  }

  const handleDeleteProduct = (productId: string) => {
    if (confirm('¿Estás seguro de que quieres eliminar este producto?')) {
      setProducts(prev => prev.filter(p => p.referenceID !== productId))
    }
  }

  const filteredProducts = products.filter(product =>
    product.name.toLowerCase().includes(searchTerm.toLowerCase()) ||
    product.descriptionMarkdown.toLowerCase().includes(searchTerm.toLowerCase()) ||
    product.referenceID.toLowerCase().includes(searchTerm.toLowerCase())
  )

  return (
    <div className="p-6">
      {/* Header */}
      <div className="mb-6">
        <h1 className="text-2xl font-bold text-gray-800 mb-2">Catálogo de Productos</h1>
        <p className="text-gray-600">
          Gestiona los productos de {tenant.name}
        </p>
      </div>

      {/* Controles */}
      <div className="mb-6 flex flex-col sm:flex-row gap-4 justify-between">
        <div className="flex-1 max-w-md">
          <div className="relative">
            <Search className="absolute left-3 top-1/2 transform -translate-y-1/2 text-gray-400 w-4 h-4" />
            <input
              type="text"
              placeholder="Buscar productos..."
              value={searchTerm}
              onChange={(e) => setSearchTerm(e.target.value)}
              className="w-full pl-10 pr-4 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent"
            />
          </div>
        </div>
        
        <div className="flex gap-2">
          <button className="px-4 py-2 border border-gray-300 text-gray-700 rounded-lg hover:bg-gray-50 transition-colors flex items-center">
            <Filter className="w-4 h-4 mr-2" />
            Filtros
          </button>
          <button
            onClick={() => setShowCreateForm(true)}
            className="px-4 py-2 bg-blue-600 text-white rounded-lg hover:bg-blue-700 transition-colors flex items-center"
          >
            <Plus className="w-4 h-4 mr-2" />
            Nuevo Producto
          </button>
        </div>
      </div>

      {/* Formulario de creación */}
      {showCreateForm && (
        <div className="mb-6">
          <CreateProduct
            onSave={handleCreateProduct}
            onCancel={() => setShowCreateForm(false)}
          />
        </div>
      )}

      {/* Lista de productos */}
      {filteredProducts.length === 0 ? (
        <Card>
          <CardContent className="text-center py-12">
            <Package className="w-16 h-16 text-gray-300 mx-auto mb-4" />
            <h3 className="text-lg font-semibold text-gray-600 mb-2">
              {products.length === 0 ? 'No hay productos' : 'No se encontraron productos'}
            </h3>
            <p className="text-gray-500 mb-4">
              {products.length === 0 
                ? 'Comienza creando tu primer producto'
                : 'Intenta con otros términos de búsqueda'
              }
            </p>
            {products.length === 0 && (
              <button
                onClick={() => setShowCreateForm(true)}
                className="px-6 py-2 bg-blue-600 text-white rounded-lg hover:bg-blue-700 transition-colors flex items-center mx-auto"
              >
                <Plus className="w-4 h-4 mr-2" />
                Crear Primer Producto
              </button>
            )}
          </CardContent>
        </Card>
      ) : (
        <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
          {filteredProducts.map((product) => (
            <Card key={product.referenceID} className="hover:shadow-lg transition-shadow">
              <CardHeader>
                <div className="flex items-start justify-between">
                  <div className="flex-1">
                    <CardTitle className="text-lg">{product.name}</CardTitle>
                    <CardDescription className="mt-1">
                      ID: {product.referenceID}
                    </CardDescription>
                  </div>
                  <div className="flex gap-1">
                    <button className="p-1 hover:bg-gray-100 rounded">
                      <Eye className="w-4 h-4 text-gray-500" />
                    </button>
                    <button className="p-1 hover:bg-gray-100 rounded">
                      <Edit className="w-4 h-4 text-gray-500" />
                    </button>
                    <button 
                      onClick={() => handleDeleteProduct(product.referenceID)}
                      className="p-1 hover:bg-red-100 rounded"
                    >
                      <Trash2 className="w-4 h-4 text-red-500" />
                    </button>
                  </div>
                </div>
              </CardHeader>
              <CardContent>
                <div className="space-y-3">
                  <p className="text-sm text-gray-600 line-clamp-2">
                    {product.descriptionMarkdown}
                  </p>
                  
                  {product.media.gallery.length > 0 && (
                    <div className="aspect-video bg-gray-100 rounded-lg overflow-hidden">
                      <img 
                        src={product.media.gallery[0]} 
                        alt={product.name}
                        className="w-full h-full object-cover"
                        onError={(e) => {
                          e.currentTarget.style.display = 'none'
                        }}
                      />
                    </div>
                  )}
                  
                  <div className="flex justify-between items-center pt-2 border-t">
                    <div>
                      <span className="text-lg font-semibold text-gray-800">
                        ${product.price.fixedPrice.toLocaleString()}
                      </span>
                      <span className="text-sm text-gray-500 ml-1">
                        {product.payment.currency}
                      </span>
                    </div>
                    <div className="text-right">
                      <div className="text-sm text-gray-600">
                        Stock: {product.stock.fixed.availableUnits}
                      </div>
                      {product.stock.weight.availableWeight > 0 && (
                        <div className="text-xs text-gray-500">
                          Peso: {product.stock.weight.availableWeight}kg
                        </div>
                      )}
                    </div>
                  </div>
                  
                  <div className="flex flex-wrap gap-1">
                    {product.payment.methods.map((method) => (
                      <span 
                        key={method}
                        className="px-2 py-1 bg-blue-100 text-blue-800 text-xs rounded-full"
                      >
                        {method.replace('_', ' ')}
                      </span>
                    ))}
                  </div>
                </div>
              </CardContent>
            </Card>
          ))}
        </div>
      )}

      {/* Estadísticas */}
      {products.length > 0 && (
        <div className="mt-8 p-4 bg-gray-50 rounded-lg">
          <h3 className="font-semibold text-gray-700 mb-2">Estadísticas del Catálogo</h3>
          <div className="grid grid-cols-2 md:grid-cols-4 gap-4 text-sm">
            <div>
              <span className="text-gray-600">Total productos:</span>
              <span className="font-semibold ml-1">{products.length}</span>
            </div>
            <div>
              <span className="text-gray-600">Valor promedio:</span>
              <span className="font-semibold ml-1">
                ${products.length > 0 ? Math.round(products.reduce((sum, p) => sum + p.price.fixedPrice, 0) / products.length).toLocaleString() : 0}
              </span>
            </div>
            <div>
              <span className="text-gray-600">Stock total:</span>
              <span className="font-semibold ml-1">
                {products.reduce((sum, p) => sum + p.stock.fixed.availableUnits, 0)}
              </span>
            </div>
            <div>
              <span className="text-gray-600">Con imagen:</span>
              <span className="font-semibold ml-1">
                {products.filter(p => p.media.gallery.length > 0).length}
              </span>
            </div>
          </div>
        </div>
      )}
    </div>
  )
}

export default Catalog
