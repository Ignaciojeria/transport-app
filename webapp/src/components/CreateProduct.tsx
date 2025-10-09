import React, { useState } from 'react'
import { type CreateProductRequest } from '../types/product'
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from './ui/Card'
import { Plus, Save, X } from 'lucide-react'

interface CreateProductProps {
  onSave: (product: CreateProductRequest) => void
  onCancel: () => void
}

const CreateProduct: React.FC<CreateProductProps> = ({ onSave, onCancel }) => {
  const [formData, setFormData] = useState<CreateProductRequest>({
    referenceID: '',
    name: '',
    description: '',
    image: '',
    payment: {
      currency: 'CLP',
      methods: ['credit_card', 'debit_card', 'transfer'],
      provider: 'Transbank'
    },
    stock: {
      fixed: { availableUnits: 0 },
      weight: { availableWeight: 0 },
      volume: { availableVolume: 0 }
    },
    price: {
      fixedPrice: 0,
      weight: { unitSize: 1, pricePerUnit: 0 },
      volume: { unitSize: 1, pricePerUnit: 0 }
    },
    logistics: {
      dimensions: { height: 0, length: 0, width: 0 },
      availabilityTime: [
        {
          timeRange: { from: '09:00', to: '22:00' },
          daysOfWeek: ['mon', 'tue', 'wed', 'thu', 'fri']
        }
      ],
      costs: [
        {
          condition: 'prime',
          type: 'fixed',
          value: 0,
          timeRange: { from: '09:00', to: '18:00' }
        }
      ]
    }
  })

  const [isSubmitting, setIsSubmitting] = useState(false)

  const handleInputChange = (path: string, value: any) => {
    setFormData(prev => {
      const keys = path.split('.')
      const newData = { ...prev }
      let current = newData as any
      
      for (let i = 0; i < keys.length - 1; i++) {
        current = current[keys[i]]
      }
      
      current[keys[keys.length - 1]] = value
      return newData
    })
  }

  const handleArrayChange = (path: string, index: number, field: string, value: any) => {
    setFormData(prev => {
      const keys = path.split('.')
      const newData = { ...prev }
      let current = newData as any
      
      for (let i = 0; i < keys.length - 1; i++) {
        if (Array.isArray(current[keys[i]])) {
          current[keys[i]] = [...current[keys[i]]]
        } else {
          current[keys[i]] = { ...current[keys[i]] }
        }
        current = current[keys[i]]
      }
      
      current[index] = { ...current[index], [field]: value }
      return newData
    })
  }

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault()
    setIsSubmitting(true)
    
    try {
      await onSave(formData)
    } catch (error) {
      console.error('Error al crear producto:', error)
    } finally {
      setIsSubmitting(false)
    }
  }

  return (
    <div className="max-w-4xl mx-auto p-6">
      <Card>
        <CardHeader>
          <CardTitle className="flex items-center">
            <Plus className="w-5 h-5 mr-2" />
            Crear Nuevo Producto
          </CardTitle>
          <CardDescription>
            Completa la información del producto siguiendo el contrato especificado
          </CardDescription>
        </CardHeader>
        <CardContent>
          <form onSubmit={handleSubmit} className="space-y-6">
            {/* Información Básica */}
            <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
              <div>
                <label className="block text-sm font-medium text-gray-700 mb-1">
                  Reference ID
                </label>
                <input
                  type="text"
                  value={formData.referenceID}
                  onChange={(e) => handleInputChange('referenceID', e.target.value)}
                  className="w-full px-3 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent"
                  required
                />
              </div>
              <div>
                <label className="block text-sm font-medium text-gray-700 mb-1">
                  Nombre del Producto
                </label>
                <input
                  type="text"
                  value={formData.name}
                  onChange={(e) => handleInputChange('name', e.target.value)}
                  className="w-full px-3 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent"
                  required
                />
              </div>
            </div>

            <div>
              <label className="block text-sm font-medium text-gray-700 mb-1">
                Descripción
              </label>
              <textarea
                value={formData.description}
                onChange={(e) => handleInputChange('description', e.target.value)}
                className="w-full px-3 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent"
                rows={3}
                required
              />
            </div>

            <div>
              <label className="block text-sm font-medium text-gray-700 mb-1">
                URL de Imagen
              </label>
              <input
                type="url"
                value={formData.image}
                onChange={(e) => handleInputChange('image', e.target.value)}
                className="w-full px-3 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent"
                required
              />
            </div>

            {/* Información de Pago */}
            <div className="border-t pt-6">
              <h3 className="text-lg font-semibold text-gray-800 mb-4">Información de Pago</h3>
              <div className="grid grid-cols-1 md:grid-cols-3 gap-4">
                <div>
                  <label className="block text-sm font-medium text-gray-700 mb-1">
                    Moneda
                  </label>
                  <select
                    value={formData.payment.currency}
                    onChange={(e) => handleInputChange('payment.currency', e.target.value)}
                    className="w-full px-3 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent"
                  >
                    <option value="CLP">CLP</option>
                    <option value="USD">USD</option>
                    <option value="EUR">EUR</option>
                  </select>
                </div>
                <div>
                  <label className="block text-sm font-medium text-gray-700 mb-1">
                    Proveedor
                  </label>
                  <input
                    type="text"
                    value={formData.payment.provider}
                    onChange={(e) => handleInputChange('payment.provider', e.target.value)}
                    className="w-full px-3 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent"
                    required
                  />
                </div>
                <div>
                  <label className="block text-sm font-medium text-gray-700 mb-1">
                    Métodos de Pago
                  </label>
                  <div className="space-y-2">
                    {['credit_card', 'debit_card', 'transfer'].map(method => (
                      <label key={method} className="flex items-center">
                        <input
                          type="checkbox"
                          checked={formData.payment.methods.includes(method)}
                          onChange={(e) => {
                            const methods = e.target.checked
                              ? [...formData.payment.methods, method]
                              : formData.payment.methods.filter(m => m !== method)
                            handleInputChange('payment.methods', methods)
                          }}
                          className="mr-2"
                        />
                        <span className="text-sm capitalize">{method.replace('_', ' ')}</span>
                      </label>
                    ))}
                  </div>
                </div>
              </div>
            </div>

            {/* Stock */}
            <div className="border-t pt-6">
              <h3 className="text-lg font-semibold text-gray-800 mb-4">Stock</h3>
              <div className="grid grid-cols-1 md:grid-cols-3 gap-4">
                <div>
                  <label className="block text-sm font-medium text-gray-700 mb-1">
                    Unidades Disponibles
                  </label>
                  <input
                    type="number"
                    value={formData.stock.fixed.availableUnits}
                    onChange={(e) => handleInputChange('stock.fixed.availableUnits', parseInt(e.target.value) || 0)}
                    className="w-full px-3 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent"
                    min="0"
                  />
                </div>
                <div>
                  <label className="block text-sm font-medium text-gray-700 mb-1">
                    Peso Disponible (kg)
                  </label>
                  <input
                    type="number"
                    value={formData.stock.weight.availableWeight}
                    onChange={(e) => handleInputChange('stock.weight.availableWeight', parseInt(e.target.value) || 0)}
                    className="w-full px-3 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent"
                    min="0"
                  />
                </div>
                <div>
                  <label className="block text-sm font-medium text-gray-700 mb-1">
                    Volumen Disponible (m³)
                  </label>
                  <input
                    type="number"
                    value={formData.stock.volume.availableVolume}
                    onChange={(e) => handleInputChange('stock.volume.availableVolume', parseInt(e.target.value) || 0)}
                    className="w-full px-3 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent"
                    min="0"
                  />
                </div>
              </div>
            </div>

            {/* Precios */}
            <div className="border-t pt-6">
              <h3 className="text-lg font-semibold text-gray-800 mb-4">Precios</h3>
              <div className="grid grid-cols-1 md:grid-cols-3 gap-4">
                <div>
                  <label className="block text-sm font-medium text-gray-700 mb-1">
                    Precio Fijo
                  </label>
                  <input
                    type="number"
                    value={formData.price.fixedPrice}
                    onChange={(e) => handleInputChange('price.fixedPrice', parseInt(e.target.value) || 0)}
                    className="w-full px-3 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent"
                    min="0"
                  />
                </div>
                <div>
                  <label className="block text-sm font-medium text-gray-700 mb-1">
                    Precio por Peso
                  </label>
                  <input
                    type="number"
                    value={formData.price.weight.pricePerUnit}
                    onChange={(e) => handleInputChange('price.weight.pricePerUnit', parseInt(e.target.value) || 0)}
                    className="w-full px-3 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent"
                    min="0"
                  />
                </div>
                <div>
                  <label className="block text-sm font-medium text-gray-700 mb-1">
                    Precio por Volumen
                  </label>
                  <input
                    type="number"
                    value={formData.price.volume.pricePerUnit}
                    onChange={(e) => handleInputChange('price.volume.pricePerUnit', parseInt(e.target.value) || 0)}
                    className="w-full px-3 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent"
                    min="0"
                  />
                </div>
              </div>
            </div>

            {/* Logística */}
            <div className="border-t pt-6">
              <h3 className="text-lg font-semibold text-gray-800 mb-4">Logística</h3>
              <div className="grid grid-cols-1 md:grid-cols-3 gap-4 mb-4">
                <div>
                  <label className="block text-sm font-medium text-gray-700 mb-1">
                    Altura (cm)
                  </label>
                  <input
                    type="number"
                    value={formData.logistics.dimensions.height}
                    onChange={(e) => handleInputChange('logistics.dimensions.height', parseInt(e.target.value) || 0)}
                    className="w-full px-3 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent"
                    min="0"
                  />
                </div>
                <div>
                  <label className="block text-sm font-medium text-gray-700 mb-1">
                    Largo (cm)
                  </label>
                  <input
                    type="number"
                    value={formData.logistics.dimensions.length}
                    onChange={(e) => handleInputChange('logistics.dimensions.length', parseInt(e.target.value) || 0)}
                    className="w-full px-3 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent"
                    min="0"
                  />
                </div>
                <div>
                  <label className="block text-sm font-medium text-gray-700 mb-1">
                    Ancho (cm)
                  </label>
                  <input
                    type="number"
                    value={formData.logistics.dimensions.width}
                    onChange={(e) => handleInputChange('logistics.dimensions.width', parseInt(e.target.value) || 0)}
                    className="w-full px-3 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent"
                    min="0"
                  />
                </div>
              </div>

              {/* Horarios de Disponibilidad */}
              <div className="mb-4">
                <h4 className="text-md font-medium text-gray-700 mb-2">Horarios de Disponibilidad</h4>
                {formData.logistics.availabilityTime.map((time, index) => (
                  <div key={index} className="grid grid-cols-1 md:grid-cols-3 gap-4 mb-2 p-3 bg-gray-50 rounded-lg">
                    <div>
                      <label className="block text-sm font-medium text-gray-600 mb-1">
                        Desde
                      </label>
                      <input
                        type="time"
                        value={time.timeRange.from}
                        onChange={(e) => handleArrayChange('logistics.availabilityTime', index, 'timeRange', { ...time.timeRange, from: e.target.value })}
                        className="w-full px-3 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent"
                      />
                    </div>
                    <div>
                      <label className="block text-sm font-medium text-gray-600 mb-1">
                        Hasta
                      </label>
                      <input
                        type="time"
                        value={time.timeRange.to}
                        onChange={(e) => handleArrayChange('logistics.availabilityTime', index, 'timeRange', { ...time.timeRange, to: e.target.value })}
                        className="w-full px-3 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent"
                      />
                    </div>
                    <div>
                      <label className="block text-sm font-medium text-gray-600 mb-2">
                        Días de la Semana
                      </label>
                      <div className="grid grid-cols-2 gap-2">
                        {[
                          { value: 'mon', label: 'Lunes' },
                          { value: 'tue', label: 'Martes' },
                          { value: 'wed', label: 'Miércoles' },
                          { value: 'thu', label: 'Jueves' },
                          { value: 'fri', label: 'Viernes' },
                          { value: 'sat', label: 'Sábado' },
                          { value: 'sun', label: 'Domingo' }
                        ].map((day) => (
                          <label key={day.value} className="flex items-center">
                            <input
                              type="checkbox"
                              checked={time.daysOfWeek.includes(day.value)}
                              onChange={(e) => {
                                const selectedDays = e.target.checked
                                  ? [...time.daysOfWeek, day.value]
                                  : time.daysOfWeek.filter(d => d !== day.value)
                                
                                // Actualizar directamente el array en el estado
                                setFormData(prev => {
                                  const newData = { ...prev }
                                  const newAvailabilityTime = [...newData.logistics.availabilityTime]
                                  newAvailabilityTime[index] = {
                                    ...newAvailabilityTime[index],
                                    daysOfWeek: selectedDays
                                  }
                                  newData.logistics = {
                                    ...newData.logistics,
                                    availabilityTime: newAvailabilityTime
                                  }
                                  return newData
                                })
                              }}
                              className="mr-2 rounded border-gray-300 text-blue-600 focus:ring-blue-500"
                            />
                            <span className="text-sm text-gray-700">{day.label}</span>
                          </label>
                        ))}
                      </div>
                    </div>
                  </div>
                ))}
              </div>

              {/* Costos */}
              <div>
                <h4 className="text-md font-medium text-gray-700 mb-2">Costos de Logística</h4>
                {formData.logistics.costs.map((cost, index) => (
                  <div key={index} className="grid grid-cols-1 md:grid-cols-4 gap-4 mb-2 p-3 bg-gray-50 rounded-lg">
                    <div>
                      <label className="block text-sm font-medium text-gray-600 mb-1">
                        Condición
                      </label>
                      <input
                        type="text"
                        value={cost.condition}
                        onChange={(e) => handleArrayChange('logistics.costs', index, 'condition', e.target.value)}
                        className="w-full px-3 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent"
                      />
                    </div>
                    <div>
                      <label className="block text-sm font-medium text-gray-600 mb-1">
                        Tipo
                      </label>
                      <select
                        value={cost.type}
                        onChange={(e) => handleArrayChange('logistics.costs', index, 'type', e.target.value)}
                        className="w-full px-3 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent"
                      >
                        <option value="fixed">Fijo</option>
                        <option value="percentage">Porcentaje</option>
                      </select>
                    </div>
                    <div>
                      <label className="block text-sm font-medium text-gray-600 mb-1">
                        Valor
                      </label>
                      <input
                        type="number"
                        value={cost.value}
                        onChange={(e) => handleArrayChange('logistics.costs', index, 'value', parseInt(e.target.value) || 0)}
                        className="w-full px-3 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent"
                        min="0"
                      />
                    </div>
                    <div>
                      <label className="block text-sm font-medium text-gray-600 mb-1">
                        Horario
                      </label>
                      <div className="flex gap-1">
                        <input
                          type="time"
                          value={cost.timeRange.from}
                          onChange={(e) => handleArrayChange('logistics.costs', index, 'timeRange', { ...cost.timeRange, from: e.target.value })}
                          className="w-full px-2 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent text-xs"
                        />
                        <input
                          type="time"
                          value={cost.timeRange.to}
                          onChange={(e) => handleArrayChange('logistics.costs', index, 'timeRange', { ...cost.timeRange, to: e.target.value })}
                          className="w-full px-2 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent text-xs"
                        />
                      </div>
                    </div>
                  </div>
                ))}
              </div>
            </div>

            {/* Botones */}
            <div className="flex justify-end gap-4 pt-6 border-t">
              <button
                type="button"
                onClick={onCancel}
                className="px-6 py-2 border border-gray-300 text-gray-700 rounded-lg hover:bg-gray-50 transition-colors flex items-center"
              >
                <X className="w-4 h-4 mr-2" />
                Cancelar
              </button>
              <button
                type="submit"
                disabled={isSubmitting}
                className="px-6 py-2 bg-blue-600 text-white rounded-lg hover:bg-blue-700 transition-colors flex items-center disabled:opacity-50"
              >
                <Save className="w-4 h-4 mr-2" />
                {isSubmitting ? 'Guardando...' : 'Guardar Producto'}
              </button>
            </div>
          </form>
        </CardContent>
      </Card>
    </div>
  )
}

export default CreateProduct
