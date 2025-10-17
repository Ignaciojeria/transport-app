import React, { useState } from 'react'
import { type CreateProductRequest, type Attribute, type Component, type Attachment } from '../types/product'
import { Card, CardContent } from './ui/Card'
import { Plus, Save, X, Image, Upload, Video, Camera, Package, Truck, DollarSign, Layers, Star, Tag, Download, MessageSquare, ChevronDown } from 'lucide-react'
import MDEditor from '@uiw/react-md-editor'
import ReactMarkdown from 'react-markdown'
import './MarkdownEditor.css'

interface CreateProductProps {
  onSave: (product: CreateProductRequest) => void
  onCancel: () => void
}

const CreateProduct: React.FC<CreateProductProps> = ({ onSave, onCancel }) => {
  const [formData, setFormData] = useState<CreateProductRequest>({
    referenceID: '',
    name: '',
    descriptionMarkdown: '',
    status: {
      isAvailable: true,
      isFeatured: false,
      allowReviews: true
    },
    attachments: [],
    properties: {
      sku: '',
      brand: '',
      barcode: ''
    },
    purchaseConditions: {
      fixed: { minUnits: 0, maxUnits: 0, multiplesOf: 0 },
      weight: { minWeight: 0, maxWeight: 0, multiplesOf: 0 },
      volume: { minVolume: 0, maxVolume: 0, multiplesOf: 0 }
    },
    attributes: [],
    categories: [],
    digitalBundle: {
      hasDigitalContent: false,
      type: 'downloadable',
      title: '',
    description: '',
      access: {
        method: 'link',
        url: '',
        expiresInDays: 30
      }
    },
    welcomeMessageMarkdown: '',
    media: {
      videos: [],
      gallery: []
    },
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
    cost: {
      fixedCost: 0,
      weight: { unitSize: 1, costPerUnit: 0 },
      volume: { unitSize: 1, constPerUnit: 0 }
    },
    components: [],
    logistics: {
      dimensions: { height: 0, length: 0, width: 0 },
      weight: 0,
      availabilityTime: [{
          timeRange: { from: '09:00', to: '22:00' },
          daysOfWeek: ['mon', 'tue', 'wed', 'thu', 'fri']
      }],
      deliveryFees: [{
          condition: 'prime',
          type: 'fixed',
          value: 0,
          timeRange: { from: '09:00', to: '18:00' }
      }]
    }
  })

  const [isSubmitting, setIsSubmitting] = useState(false)
  const [showCategoryModal, setShowCategoryModal] = useState(false)
  const [selectedCategory, setSelectedCategory] = useState('')
  const [newCategoryName, setNewCategoryName] = useState('')
  const [newCategoryParent, setNewCategoryParent] = useState('')

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

  // Funciones para manejar arrays dinámicos
  const addAttachment = () => {
    const newAttachment: Attachment = {
      name: '',
      description: '',
      url: '',
      type: 'pdf',
      sizeKb: 0
    }
    setFormData(prev => ({
      ...prev,
      attachments: [...prev.attachments, newAttachment]
    }))
      }

  const removeAttachment = (index: number) => {
    setFormData(prev => ({
      ...prev,
      attachments: prev.attachments.filter((_, i) => i !== index)
    }))
  }

  const addAttribute = () => {
    const newAttribute: Attribute = { name: '', value: '' }
      setFormData(prev => ({
        ...prev,
      attributes: [...prev.attributes, newAttribute]
    }))
  }

  const removeAttribute = (index: number) => {
    setFormData(prev => ({
      ...prev,
      attributes: prev.attributes.filter((_, i) => i !== index)
    }))
  }

  // Función para generar ID jerárquico automáticamente
  const generateCategoryId = (name: string, parentId: string | null) => {
    const cleanName = name.toLowerCase().replace(/[^a-z0-9]/g, '-').replace(/-+/g, '-').replace(/^-|-$/g, '')
    return parentId ? `${parentId}/${cleanName}` : cleanName
  }

  // Función para obtener la ruta completa de una categoría
  const getCategoryPath = (categoryId: string): string => {
    const category = formData.categories.find(cat => cat.id === categoryId)
    if (!category) return ''
    
    if (!category.parent) {
      return category.name
    }
    
    const parentPath: string = getCategoryPath(category.parent)
    return parentPath ? `${parentPath} / ${category.name}` : category.name
  }

  const addCategory = () => {
    setShowCategoryModal(true)
    setNewCategoryName('')
    setNewCategoryParent('')
  }

  const handleAddCategory = () => {
    if (!newCategoryName.trim()) return

    const parentId = newCategoryParent || null
    const newId = generateCategoryId(newCategoryName, parentId)
    
    const newCategory = {
      id: newId,
      name: newCategoryName.trim(),
      parent: parentId
    }

    setFormData(prev => ({
      ...prev,
      categories: [...prev.categories, newCategory]
    }))

    setShowCategoryModal(false)
    setNewCategoryName('')
    setNewCategoryParent('')
  }

  const removeCategory = (index: number) => {
      setFormData(prev => ({
        ...prev,
      categories: prev.categories.filter((_, i) => i !== index)
    }))
  }

  const addComponent = () => {
    const newComponent: Component = {
      type: 'base',
      name: '',
      required: true,
      stock: {}
    }
    setFormData(prev => ({
      ...prev,
      components: [...prev.components, newComponent]
    }))
  }

  const removeComponent = (index: number) => {
    setFormData(prev => ({
      ...prev,
      components: prev.components.filter((_, i) => i !== index)
    }))
  }

  const addAvailabilityTime = () => {
    setFormData(prev => ({
      ...prev,
      logistics: {
        ...prev.logistics,
        availabilityTime: [...prev.logistics.availabilityTime, {
          timeRange: { from: '09:00', to: '22:00' },
          daysOfWeek: ['mon', 'tue', 'wed', 'thu', 'fri']
        }]
      }
    }))
  }

  const removeAvailabilityTime = (index: number) => {
    setFormData(prev => ({
      ...prev,
      logistics: {
        ...prev.logistics,
        availabilityTime: prev.logistics.availabilityTime.filter((_, i) => i !== index)
      }
    }))
  }

  const addDeliveryFee = () => {
    setFormData(prev => ({
      ...prev,
      logistics: {
        ...prev.logistics,
        deliveryFees: [...prev.logistics.deliveryFees, {
          condition: '',
          type: 'fixed',
          value: 0,
          timeRange: { from: '09:00', to: '18:00' }
        }]
      }
    }))
  }

  const removeDeliveryFee = (index: number) => {
    setFormData(prev => ({
      ...prev,
      logistics: {
        ...prev.logistics,
        deliveryFees: prev.logistics.deliveryFees.filter((_, i) => i !== index)
      }
    }))
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
    <div className="max-w-7xl mx-auto p-6">
      <div className="mb-6">
        <h1 className="text-2xl font-bold text-gray-900">Crear Nuevo Producto</h1>
        <p className="text-gray-600">Completa la información del producto siguiendo el contrato especificado</p>
      </div>

          <form onSubmit={handleSubmit} className="space-y-6">
            {/* Información Básica */}
        <Card>
          <CardContent className="p-6">
            <h3 className="text-lg font-semibold text-gray-800 mb-4 flex items-center">
              <Package className="w-5 h-5 mr-2" />
              Información Básica
            </h3>
            <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
              <div>
                <label className="block text-sm font-medium text-gray-700 mb-1">
                  ID de Referencia
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
          </CardContent>
        </Card>


        {/* descriptionMarkdown */}
        <Card>
          <CardContent className="p-6">
            <h3 className="text-lg font-semibold text-gray-800 mb-4 flex items-center">
              <MessageSquare className="w-5 h-5 mr-2" />
                Descripción
            </h3>
            <div className="border border-gray-300 rounded-lg overflow-hidden">
              <MDEditor
                value={formData.descriptionMarkdown}
                onChange={(value) => handleInputChange('descriptionMarkdown', value || '')}
                preview="edit"
                hideToolbar={false}
                data-color-mode="light"
                height={300}
              />
            </div>
            {formData.descriptionMarkdown && (
              <div className="mt-2 p-3 bg-gray-50 rounded-lg border">
                <h5 className="text-sm font-medium text-gray-600 mb-2">Vista previa:</h5>
                <div className="prose prose-sm max-w-none">
                  <ReactMarkdown>{formData.descriptionMarkdown}</ReactMarkdown>
                </div>
              </div>
            )}
          </CardContent>
        </Card>

        {/* status */}
        <Card>
          <CardContent className="p-6">
            <h3 className="text-lg font-semibold text-gray-800 mb-4 flex items-center">
              <Star className="w-5 h-5 mr-2" />
              Estado del Producto
            </h3>
            <div className="grid grid-cols-1 md:grid-cols-3 gap-4">
              <div className="flex items-center space-x-3">
                <input
                  type="checkbox"
                  checked={formData.status.isAvailable}
                  onChange={(e) => handleInputChange('status.isAvailable', e.target.checked)}
                  className="w-4 h-4 text-blue-600 bg-gray-100 border-gray-300 rounded focus:ring-blue-500"
                />
                <label className="text-sm font-medium text-gray-700">
                  Disponible
              </label>
              </div>
              <div className="flex items-center space-x-3">
              <input
                  type="checkbox"
                  checked={formData.status.isFeatured}
                  onChange={(e) => handleInputChange('status.isFeatured', e.target.checked)}
                  className="w-4 h-4 text-blue-600 bg-gray-100 border-gray-300 rounded focus:ring-blue-500"
                />
                <label className="text-sm font-medium text-gray-700">
                  Destacado
                </label>
            </div>
              <div className="flex items-center space-x-3">
                <input
                  type="checkbox"
                  checked={formData.status.allowReviews}
                  onChange={(e) => handleInputChange('status.allowReviews', e.target.checked)}
                  className="w-4 h-4 text-blue-600 bg-gray-100 border-gray-300 rounded focus:ring-blue-500"
                />
                <label className="text-sm font-medium text-gray-700">
                  Permitir Reseñas
                </label>
              </div>
            </div>
          </CardContent>
        </Card>

        {/* attachments */}
        <Card>
          <CardContent className="p-6">
              <div className="flex items-center justify-between mb-4">
              <h3 className="text-lg font-semibold text-gray-800 flex items-center">
                <Upload className="w-5 h-5 mr-2" />
                Adjuntos
              </h3>
                <button
                  type="button"
                onClick={addAttachment}
                  className="px-3 py-1 bg-blue-600 text-white text-sm rounded-lg hover:bg-blue-700 transition-colors flex items-center"
                >
                  <Plus className="w-4 h-4 mr-1" />
                Agregar Adjunto
                </button>
              </div>
            <div className="space-y-3">
              {formData.attachments.map((attachment, index) => (
                <div key={index} className="grid grid-cols-1 md:grid-cols-3 gap-3 p-3 bg-gray-50 rounded-lg">
                  <input
                    type="text"
                    value={attachment.name}
                    onChange={(e) => handleArrayChange('attachments', index, 'name', e.target.value)}
                    placeholder="Nombre del archivo"
                    className="px-3 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent"
                  />
                  <input
                    type="text"
                    value={attachment.description}
                    onChange={(e) => handleArrayChange('attachments', index, 'description', e.target.value)}
                    placeholder="Descripción"
                    className="px-3 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent"
                  />
                  <div className="flex items-center gap-2">
                    <input
                      type="file"
                      onChange={(e) => {
                        const file = e.target.files?.[0]
                        if (file) {
                          // Aquí podrías manejar la subida del archivo
                          // Por ahora solo actualizamos el nombre del archivo
                          handleArrayChange('attachments', index, 'name', file.name)
                          handleArrayChange('attachments', index, 'type', file.type.split('/')[1] || 'pdf')
                          handleArrayChange('attachments', index, 'sizeKb', Math.round(file.size / 1024))
                        }
                      }}
                      className="flex-1 px-3 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent file:mr-4 file:py-2 file:px-4 file:rounded-lg file:border-0 file:text-sm file:font-semibold file:bg-blue-50 file:text-blue-700 hover:file:bg-blue-100"
                    />
                    <button
                      type="button"
                      onClick={() => removeAttachment(index)}
                      className="p-2 hover:bg-red-100 rounded transition-colors"
                      title="Eliminar adjunto"
                    >
                      <X className="w-4 h-4 text-red-500" />
                    </button>
                  </div>
                </div>
              ))}
            </div>
          </CardContent>
        </Card>

        {/* properties */}
        <Card>
          <CardContent className="p-6">
            <h3 className="text-lg font-semibold text-gray-800 mb-4 flex items-center">
              <Tag className="w-5 h-5 mr-2" />
              Propiedades
            </h3>
            <div className="grid grid-cols-1 md:grid-cols-3 gap-4">
              <div>
                <label className="block text-sm font-medium text-gray-700 mb-1">
                  SKU
                      </label>
                <input
                  type="text"
                  value={formData.properties.sku}
                  onChange={(e) => handleInputChange('properties.sku', e.target.value)}
                  className="w-full px-3 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent"
                        />
                      </div>
              <div>
                <label className="block text-sm font-medium text-gray-700 mb-1">
                  Marca
                </label>
                <input
                  type="text"
                  value={formData.properties.brand}
                  onChange={(e) => handleInputChange('properties.brand', e.target.value)}
                  className="w-full px-3 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent"
                />
              </div>
              <div>
                <label className="block text-sm font-medium text-gray-700 mb-1">
                  Código de Barras
                </label>
                <input
                  type="text"
                  value={formData.properties.barcode}
                  onChange={(e) => handleInputChange('properties.barcode', e.target.value)}
                  className="w-full px-3 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent"
                />
              </div>
            </div>
          </CardContent>
        </Card>

        {/* purchaseConditions */}
        <Card>
          <CardContent className="p-6">
            <h3 className="text-lg font-semibold text-gray-800 mb-4 flex items-center">
              <Package className="w-5 h-5 mr-2" />
              Condiciones de Compra
            </h3>
            <div className="space-y-6">
              {/* fixed */}
              <div>
                <h4 className="text-md font-medium text-gray-700 mb-3">Unidades Fijas</h4>
                <div className="grid grid-cols-1 md:grid-cols-3 gap-4">
                  <div>
                    <label className="block text-sm font-medium text-gray-600 mb-1">Mínimo de Unidades</label>
                    <input
                      type="number"
                      value={formData.purchaseConditions.fixed.minUnits || 0}
                      onChange={(e) => handleInputChange('purchaseConditions.fixed.minUnits', parseInt(e.target.value) || 0)}
                      className="w-full px-3 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent"
                      min="0"
                    />
                          </div>
                  <div>
                    <label className="block text-sm font-medium text-gray-600 mb-1">Máximo de Unidades</label>
                    <input
                      type="number"
                      value={formData.purchaseConditions.fixed.maxUnits || 0}
                      onChange={(e) => handleInputChange('purchaseConditions.fixed.maxUnits', parseInt(e.target.value) || 0)}
                      className="w-full px-3 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent"
                      min="0"
                    />
                        </div>
                  <div>
                    <label className="block text-sm font-medium text-gray-600 mb-1">Múltiplos de</label>
                    <input
                      type="number"
                      value={formData.purchaseConditions.fixed.multiplesOf || 0}
                      onChange={(e) => handleInputChange('purchaseConditions.fixed.multiplesOf', parseInt(e.target.value) || 0)}
                      className="w-full px-3 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent"
                      min="0"
                    />
                  </div>
                </div>
              </div>

              {/* weight */}
              <div>
                <h4 className="text-md font-medium text-gray-700 mb-3">Peso</h4>
                <div className="grid grid-cols-1 md:grid-cols-4 gap-4">
                  <div>
                    <label className="block text-sm font-medium text-gray-600 mb-1">Peso Mínimo</label>
                    <input
                      type="number"
                      value={formData.purchaseConditions.weight.minWeight || 0}
                      onChange={(e) => handleInputChange('purchaseConditions.weight.minWeight', parseInt(e.target.value) || 0)}
                      className="w-full px-3 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent"
                      min="0"
                    />
                        </div>
                  <div>
                    <label className="block text-sm font-medium text-gray-600 mb-1">Peso Máximo</label>
                    <input
                      type="number"
                      value={formData.purchaseConditions.weight.maxWeight || 0}
                      onChange={(e) => handleInputChange('purchaseConditions.weight.maxWeight', parseInt(e.target.value) || 0)}
                      className="w-full px-3 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent"
                      min="0"
                    />
                  </div>
                  <div>
                    <label className="block text-sm font-medium text-gray-600 mb-1">Múltiplos de</label>
                    <input
                      type="number"
                      value={formData.purchaseConditions.weight.multiplesOf || 0}
                      onChange={(e) => handleInputChange('purchaseConditions.weight.multiplesOf', parseInt(e.target.value) || 0)}
                      className="w-full px-3 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent"
                      min="0"
                    />
                  </div>
                  <div>
                    <label className="block text-sm font-medium text-gray-600 mb-1">Notas</label>
                    <input
                      type="text"
                      value={formData.purchaseConditions.weight.notes || ''}
                      onChange={(e) => handleInputChange('purchaseConditions.weight.notes', e.target.value)}
                      className="w-full px-3 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent"
                    />
                  </div>
                      </div>
                    </div>

              {/* volume */}
              <div>
                <h4 className="text-md font-medium text-gray-700 mb-3">Volumen</h4>
                <div className="grid grid-cols-1 md:grid-cols-3 gap-4">
                  <div>
                    <label className="block text-sm font-medium text-gray-600 mb-1">Volumen Mínimo</label>
                    <input
                      type="number"
                      value={formData.purchaseConditions.volume.minVolume || 0}
                      onChange={(e) => handleInputChange('purchaseConditions.volume.minVolume', parseInt(e.target.value) || 0)}
                        className="w-full px-3 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent"
                      min="0"
                    />
                    </div>
                  <div>
                    <label className="block text-sm font-medium text-gray-600 mb-1">Volumen Máximo</label>
                    <input
                      type="number"
                      value={formData.purchaseConditions.volume.maxVolume || 0}
                      onChange={(e) => handleInputChange('purchaseConditions.volume.maxVolume', parseInt(e.target.value) || 0)}
                      className="w-full px-3 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent"
                      min="0"
                    />
                  </div>
                  <div>
                    <label className="block text-sm font-medium text-gray-600 mb-1">Múltiplos de</label>
                    <input
                      type="number"
                      value={formData.purchaseConditions.volume.multiplesOf || 0}
                      onChange={(e) => handleInputChange('purchaseConditions.volume.multiplesOf', parseInt(e.target.value) || 0)}
                      className="w-full px-3 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent"
                      min="0"
                    />
                  </div>
                </div>
              </div>
            </div>
          </CardContent>
        </Card>

        {/* attributes */}
        <Card>
          <CardContent className="p-6">
            <div className="flex items-center justify-between mb-4">
              <h3 className="text-lg font-semibold text-gray-800 flex items-center">
                <Tag className="w-5 h-5 mr-2" />
                Atributos
              </h3>
                          <button
                            type="button"
                onClick={addAttribute}
                className="px-3 py-1 bg-blue-600 text-white text-sm rounded-lg hover:bg-blue-700 transition-colors flex items-center"
                          >
                <Plus className="w-4 h-4 mr-1" />
                Agregar Atributo
                          </button>
                        </div>
            <div className="space-y-3">
              {formData.attributes.map((attr, index) => (
                <div key={index} className="flex items-center gap-3 p-3 bg-gray-50 rounded-lg">
                              <input
                                type="text"
                    value={attr.name}
                    onChange={(e) => handleArrayChange('attributes', index, 'name', e.target.value)}
                    placeholder="Nombre del atributo"
                    className="flex-1 px-3 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent"
                  />
                  <input
                    type="text"
                    value={attr.value}
                    onChange={(e) => handleArrayChange('attributes', index, 'value', e.target.value)}
                    placeholder="Valor del atributo"
                                className="flex-1 px-3 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent"
                              />
                              <button
                                type="button"
                    onClick={() => removeAttribute(index)}
                    className="p-2 hover:bg-red-100 rounded transition-colors"
                    title="Eliminar atributo"
                  >
                    <X className="w-4 h-4 text-red-500" />
                  </button>
                </div>
              ))}
            </div>
          </CardContent>
        </Card>

        {/* categories */}
        <Card>
          <CardContent className="p-6">
            <h3 className="text-lg font-semibold text-gray-800 mb-4 flex items-center">
              <Layers className="w-5 h-5 mr-2" />
              Categorías
            </h3>
            
            {/* Selector de categoría */}
            <div className="mb-4">
              <label className="block text-sm font-medium text-gray-700 mb-2">Categorías</label>
              <div className="relative">
                <input
                  type="text"
                  value={selectedCategory}
                  onChange={(e) => setSelectedCategory(e.target.value)}
                  placeholder="Selecciona una categoría"
                  className="w-full px-3 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-green-500 focus:border-green-500 pr-10"
                />
                <ChevronDown className="absolute right-3 top-1/2 transform -translate-y-1/2 w-4 h-4 text-gray-400" />
              </div>
            </div>

            {/* Botón para agregar categoría */}
            <button
              type="button"
              onClick={addCategory}
              className="w-full px-4 py-2 border border-gray-300 rounded-lg text-gray-700 hover:bg-gray-50 transition-colors flex items-center justify-center"
            >
              <Plus className="w-4 h-4 mr-2" />
              Agregar Categoría
            </button>

            {/* Lista de categorías existentes */}
            {formData.categories.length > 0 && (
              <div className="mt-4">
                <h4 className="text-sm font-medium text-gray-700 mb-2">Categorías existentes:</h4>
                <div className="space-y-2">
                  {formData.categories.map((cat, index) => (
                    <div key={index} className="flex items-center justify-between p-2 bg-gray-50 rounded-lg">
                      <div className="flex-1">
                        <div className="text-sm font-medium text-gray-800">{cat.name}</div>
                        <div className="text-xs text-gray-500">ID: {cat.id}</div>
                        {cat.parent && (
                          <div className="text-xs text-gray-400">Padre: {getCategoryPath(cat.parent)}</div>
                        )}
                      </div>
                      <button
                        type="button"
                        onClick={() => removeCategory(index)}
                                className="p-1 hover:bg-red-100 rounded transition-colors"
                        title="Eliminar categoría"
                              >
                        <X className="w-4 h-4 text-red-500" />
                              </button>
                            </div>
                          ))}
                        </div>
                      </div>
                    )}
          </CardContent>
        </Card>

        {/* digitalBundle */}
        <Card>
          <CardContent className="p-6">
            <h3 className="text-lg font-semibold text-gray-800 mb-4 flex items-center">
              <Download className="w-5 h-5 mr-2" />
              Bundle Digital
            </h3>
            <div className="space-y-4">
              <div className="flex items-center space-x-3">
                <input
                  type="checkbox"
                  checked={formData.digitalBundle.hasDigitalContent}
                  onChange={(e) => handleInputChange('digitalBundle.hasDigitalContent', e.target.checked)}
                  className="w-4 h-4 text-blue-600 bg-gray-100 border-gray-300 rounded focus:ring-blue-500"
                />
                <label className="text-sm font-medium text-gray-700">
                  Incluir Contenido Digital
                        </label>
              </div>

              {formData.digitalBundle.hasDigitalContent && (
                <div className="space-y-4">
                  <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
                    <div>
                      <label className="block text-sm font-medium text-gray-700 mb-1">Tipo</label>
                      <select
                        value={formData.digitalBundle.type}
                        onChange={(e) => handleInputChange('digitalBundle.type', e.target.value)}
                        className="w-full px-3 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent"
                      >
                        <option value="downloadable">Descargable</option>
                        <option value="streaming">Streaming</option>
                        <option value="access">Acceso</option>
                      </select>
                    </div>
                    <div>
                      <label className="block text-sm font-medium text-gray-700 mb-1">Método de Acceso</label>
                      <select
                        value={formData.digitalBundle.access.method}
                        onChange={(e) => handleInputChange('digitalBundle.access.method', e.target.value)}
                        className="w-full px-3 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent"
                      >
                        <option value="link">Enlace</option>
                        <option value="email">Email</option>
                        <option value="sms">SMS</option>
                      </select>
                    </div>
                  </div>

                  <div>
                    <label className="block text-sm font-medium text-gray-700 mb-1">Título</label>
                    <input
                      type="text"
                      value={formData.digitalBundle.title}
                      onChange={(e) => handleInputChange('digitalBundle.title', e.target.value)}
                      className="w-full px-3 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent"
                    />
                  </div>

                  <div>
                    <label className="block text-sm font-medium text-gray-700 mb-1">Descripción</label>
                    <textarea
                      value={formData.digitalBundle.description}
                      onChange={(e) => handleInputChange('digitalBundle.description', e.target.value)}
                      className="w-full px-3 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent"
                      rows={3}
                    />
                  </div>

                  <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
                    <div>
                      <label className="block text-sm font-medium text-gray-700 mb-1">URL de Acceso</label>
                      <input
                        type="url"
                        value={formData.digitalBundle.access.url}
                        onChange={(e) => handleInputChange('digitalBundle.access.url', e.target.value)}
                        className="w-full px-3 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent"
                      />
                    </div>
                    <div>
                      <label className="block text-sm font-medium text-gray-700 mb-1">Expira en (días)</label>
                      <input
                        type="number"
                        value={formData.digitalBundle.access.expiresInDays}
                        onChange={(e) => handleInputChange('digitalBundle.access.expiresInDays', parseInt(e.target.value) || 30)}
                        className="w-full px-3 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent"
                        min="1"
                      />
                    </div>
                  </div>
                </div>
              )}
            </div>
          </CardContent>
        </Card>

        {/* welcomeMessageMarkdown */}
        <Card>
          <CardContent className="p-6">
            <h3 className="text-lg font-semibold text-gray-800 mb-4 flex items-center">
              <MessageSquare className="w-5 h-5 mr-2" />
              Mensaje de Bienvenida
            </h3>
                        <div className="border border-gray-300 rounded-lg overflow-hidden">
                          <MDEditor
                value={formData.welcomeMessageMarkdown}
                onChange={(value) => handleInputChange('welcomeMessageMarkdown', value || '')}
                            preview="edit"
                            hideToolbar={false}
                            data-color-mode="light"
                height={200}
                          />
                        </div>
            {formData.welcomeMessageMarkdown && (
                          <div className="mt-2 p-3 bg-gray-50 rounded-lg border">
                <h5 className="text-sm font-medium text-gray-600 mb-2">Vista previa:</h5>
                            <div className="prose prose-sm max-w-none">
                  <ReactMarkdown>{formData.welcomeMessageMarkdown}</ReactMarkdown>
                            </div>
                          </div>
                        )}
          </CardContent>
        </Card>

        {/* media */}
        <Card>
          <CardContent className="p-6">
            <h3 className="text-lg font-semibold text-gray-800 mb-4 flex items-center">
              <Camera className="w-5 h-5 mr-2" />
              Medios
            </h3>
            <div className="space-y-6">
              {/* videos */}
              <div>
                  <div className="flex items-center justify-between mb-4">
                    <h4 className="text-md font-medium text-gray-700">Videos</h4>
                  <button
                    type="button"
                    onClick={() => {
                      setFormData(prev => ({
                        ...prev,
                        media: {
                          ...prev.media,
                          videos: [...prev.media.videos, { title: '', platform: 'YouTube', url: '', thumbnail: '' }]
                        }
                      }))
                    }}
                    className="px-3 py-1 bg-blue-600 text-white text-sm rounded-lg hover:bg-blue-700 transition-colors flex items-center"
                  >
                    <Video className="w-4 h-4 mr-1" />
                    Agregar Video
                  </button>
                </div>
                <div className="space-y-3">
                  {formData.media.videos.map((video, index) => (
                    <div key={index} className="grid grid-cols-1 md:grid-cols-2 gap-3 p-3 bg-gray-50 rounded-lg">
                      <input
                        type="text"
                        value={video.title}
                        onChange={(e) => handleArrayChange('media.videos', index, 'title', e.target.value)}
                        placeholder="Título del video"
                        className="px-3 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent"
                      />
                      <input
                        type="text"
                        value={video.platform}
                        onChange={(e) => handleArrayChange('media.videos', index, 'platform', e.target.value)}
                        placeholder="Plataforma (YouTube, Vimeo, etc.)"
                        className="px-3 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent"
                      />
                      <input
                        type="url"
                        value={video.url}
                        onChange={(e) => handleArrayChange('media.videos', index, 'url', e.target.value)}
                        placeholder="URL del video"
                        className="px-3 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent"
                      />
                      <div className="flex items-center gap-2">
                        <input
                          type="url"
                          value={video.thumbnail}
                          onChange={(e) => handleArrayChange('media.videos', index, 'thumbnail', e.target.value)}
                          placeholder="URL del thumbnail"
                          className="flex-1 px-3 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent"
                        />
                        <button
                          type="button"
                          onClick={() => {
                            setFormData(prev => ({
                              ...prev,
                              media: {
                                ...prev.media,
                                videos: prev.media.videos.filter((_, i) => i !== index)
                              }
                            }))
                          }}
                          className="p-2 hover:bg-red-100 rounded transition-colors"
                          title="Eliminar video"
                        >
                          <X className="w-4 h-4 text-red-500" />
                        </button>
                    </div>
                  </div>
                ))}
              </div>
            </div>

              {/* gallery */}
              <div>
                  <div className="flex items-center justify-between mb-4">
                    <h4 className="text-md font-medium text-gray-700">Galería de Imágenes</h4>
                  <button
                    type="button"
                    onClick={() => {
                      setFormData(prev => ({
                        ...prev,
                        media: {
                          ...prev.media,
                          gallery: [...prev.media.gallery, '']
                        }
                      }))
                    }}
                    className="px-3 py-1 bg-blue-600 text-white text-sm rounded-lg hover:bg-blue-700 transition-colors flex items-center"
                  >
                    <Image className="w-4 h-4 mr-1" />
                    Agregar Imagen
                  </button>
                </div>
                <div className="space-y-3">
                  {formData.media.gallery.map((image, index) => (
                    <div key={index} className="flex items-center gap-3 p-3 bg-gray-50 rounded-lg">
                      <input
                        type="url"
                        value={image}
                        onChange={(e) => {
                          const newGallery = [...formData.media.gallery]
                          newGallery[index] = e.target.value
                          setFormData(prev => ({
                            ...prev,
                            media: { ...prev.media, gallery: newGallery }
                          }))
                        }}
                        placeholder="URL de imagen"
                        className="flex-1 px-3 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent"
                      />
                      <button
                        type="button"
                        onClick={() => {
                          setFormData(prev => ({
                            ...prev,
                            media: {
                              ...prev.media,
                              gallery: prev.media.gallery.filter((_, i) => i !== index)
                            }
                          }))
                        }}
                        className="p-2 hover:bg-red-100 rounded transition-colors"
                        title="Eliminar imagen"
                      >
                        <X className="w-4 h-4 text-red-500" />
                      </button>
                    </div>
                  ))}
                </div>
              </div>
            </div>
          </CardContent>
        </Card>

        {/* payment */}
        <Card>
          <CardContent className="p-6">
            <h3 className="text-lg font-semibold text-gray-800 mb-4 flex items-center">
              <DollarSign className="w-5 h-5 mr-2" />
              Pago
            </h3>
              <div className="grid grid-cols-1 md:grid-cols-3 gap-4">
                <div>
                <label className="block text-sm font-medium text-gray-700 mb-1">Moneda</label>
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
                <label className="block text-sm font-medium text-gray-700 mb-1">Proveedor</label>
                  <input
                    type="text"
                    value={formData.payment.provider}
                    onChange={(e) => handleInputChange('payment.provider', e.target.value)}
                    className="w-full px-3 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent"
                  />
                </div>
                <div>
                <label className="block text-sm font-medium text-gray-700 mb-1">Métodos de Pago</label>
                  <div className="space-y-2">
                  {[
                    { value: 'credit_card', label: 'Tarjeta de Crédito' },
                    { value: 'debit_card', label: 'Tarjeta de Débito' },
                    { value: 'transfer', label: 'Transferencia' }
                  ].map(method => (
                    <label key={method.value} className="flex items-center">
                        <input
                          type="checkbox"
                        checked={formData.payment.methods.includes(method.value)}
                          onChange={(e) => {
                            const methods = e.target.checked
                            ? [...formData.payment.methods, method.value]
                            : formData.payment.methods.filter(m => m !== method.value)
                            handleInputChange('payment.methods', methods)
                          }}
                          className="mr-2"
                        />
                      <span className="text-sm">{method.label}</span>
                      </label>
                    ))}
                  </div>
                </div>
              </div>
          </CardContent>
        </Card>

        {/* stock */}
        <Card>
          <CardContent className="p-6">
            <h3 className="text-lg font-semibold text-gray-800 mb-4 flex items-center">
              <Package className="w-5 h-5 mr-2" />
              Stock
            </h3>
            <div className="space-y-6">
              {/* fixed */}
                <div>
                <h4 className="text-md font-medium text-gray-700 mb-3">Unidades Fijas</h4>
                <div>
                  <label className="block text-sm font-medium text-gray-600 mb-1">Unidades Disponibles</label>
                  <input
                    type="number"
                    value={formData.stock.fixed.availableUnits}
                    onChange={(e) => handleInputChange('stock.fixed.availableUnits', parseInt(e.target.value) || 0)}
                    className="w-full px-3 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent"
                    min="0"
                  />
                </div>
              </div>

              {/* weight */}
                <div>
                <h4 className="text-md font-medium text-gray-700 mb-3">Peso</h4>
                <div>
                  <label className="block text-sm font-medium text-gray-600 mb-1">Peso Disponible</label>
                  <input
                    type="number"
                    value={formData.stock.weight.availableWeight}
                    onChange={(e) => handleInputChange('stock.weight.availableWeight', parseInt(e.target.value) || 0)}
                    className="w-full px-3 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent"
                    min="0"
                  />
                </div>
              </div>

              {/* volume */}
                <div>
                <h4 className="text-md font-medium text-gray-700 mb-3">Volumen</h4>
                <div>
                  <label className="block text-sm font-medium text-gray-600 mb-1">Volumen Disponible</label>
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
          </CardContent>
        </Card>

        {/* price */}
        <Card>
          <CardContent className="p-6">
            <h3 className="text-lg font-semibold text-gray-800 mb-4 flex items-center">
              <DollarSign className="w-5 h-5 mr-2" />
              Precio
            </h3>
            <div className="space-y-6">
                <div>
                <label className="block text-sm font-medium text-gray-700 mb-1">Precio Fijo</label>
                  <input
                    type="number"
                    value={formData.price.fixedPrice}
                    onChange={(e) => handleInputChange('price.fixedPrice', parseInt(e.target.value) || 0)}
                    className="w-full px-3 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent"
                    min="0"
                  />
                </div>

              {/* weight */}
                <div>
                <h4 className="text-md font-medium text-gray-700 mb-3">Precio por Peso</h4>
                <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
                  <div>
                    <label className="block text-sm font-medium text-gray-600 mb-1">Tamaño de Unidad</label>
                    <input
                      type="number"
                      value={formData.price.weight.unitSize}
                      onChange={(e) => handleInputChange('price.weight.unitSize', parseInt(e.target.value) || 1)}
                      className="w-full px-3 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent"
                      min="1"
                    />
                  </div>
                  <div>
                    <label className="block text-sm font-medium text-gray-600 mb-1">Precio por Unidad</label>
                  <input
                    type="number"
                    value={formData.price.weight.pricePerUnit}
                    onChange={(e) => handleInputChange('price.weight.pricePerUnit', parseInt(e.target.value) || 0)}
                    className="w-full px-3 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent"
                    min="0"
                  />
                </div>
                </div>
              </div>

              {/* volume */}
                <div>
                <h4 className="text-md font-medium text-gray-700 mb-3">Precio por Volumen</h4>
                <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
                  <div>
                    <label className="block text-sm font-medium text-gray-600 mb-1">Tamaño de Unidad</label>
                    <input
                      type="number"
                      value={formData.price.volume.unitSize}
                      onChange={(e) => handleInputChange('price.volume.unitSize', parseInt(e.target.value) || 1)}
                      className="w-full px-3 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent"
                      min="1"
                    />
                  </div>
                  <div>
                    <label className="block text-sm font-medium text-gray-600 mb-1">Precio por Unidad</label>
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
            </div>
          </CardContent>
        </Card>

        {/* cost */}
        <Card>
          <CardContent className="p-6">
            <h3 className="text-lg font-semibold text-gray-800 mb-4 flex items-center">
              <DollarSign className="w-5 h-5 mr-2" />
              Costo
            </h3>
            <div className="space-y-6">
              <div>
                <label className="block text-sm font-medium text-gray-700 mb-1">Costo Fijo</label>
                <input
                  type="number"
                  value={formData.cost.fixedCost}
                  onChange={(e) => handleInputChange('cost.fixedCost', parseInt(e.target.value) || 0)}
                  className="w-full px-3 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent"
                  min="0"
                />
              </div>

              {/* weight */}
              <div>
                <h4 className="text-md font-medium text-gray-700 mb-3">Costo por Peso</h4>
                <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
                  <div>
                    <label className="block text-sm font-medium text-gray-600 mb-1">Tamaño de Unidad</label>
                    <input
                      type="number"
                      value={formData.cost.weight.unitSize}
                      onChange={(e) => handleInputChange('cost.weight.unitSize', parseInt(e.target.value) || 1)}
                      className="w-full px-3 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent"
                      min="1"
                    />
                  </div>
                  <div>
                    <label className="block text-sm font-medium text-gray-600 mb-1">Costo por Unidad</label>
                    <input
                      type="number"
                      value={formData.cost.weight.costPerUnit}
                      onChange={(e) => handleInputChange('cost.weight.costPerUnit', parseInt(e.target.value) || 0)}
                      className="w-full px-3 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent"
                      min="0"
                    />
                  </div>
                </div>
              </div>

              {/* volume */}
              <div>
                <h4 className="text-md font-medium text-gray-700 mb-3">Costo por Volumen</h4>
                <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
                  <div>
                    <label className="block text-sm font-medium text-gray-600 mb-1">Tamaño de Unidad</label>
                    <input
                      type="number"
                      value={formData.cost.volume.unitSize}
                      onChange={(e) => handleInputChange('cost.volume.unitSize', parseInt(e.target.value) || 1)}
                      className="w-full px-3 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent"
                      min="1"
                    />
                  </div>
                  <div>
                    <label className="block text-sm font-medium text-gray-600 mb-1">Costo por Unidad</label>
                    <input
                      type="number"
                      value={formData.cost.volume.constPerUnit}
                      onChange={(e) => handleInputChange('cost.volume.constPerUnit', parseInt(e.target.value) || 0)}
                      className="w-full px-3 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent"
                      min="0"
                    />
                  </div>
                </div>
              </div>
            </div>
          </CardContent>
        </Card>

        {/* components */}
        <Card>
          <CardContent className="p-6">
            <div className="flex items-center justify-between mb-4">
              <h3 className="text-lg font-semibold text-gray-800 flex items-center">
                <Layers className="w-5 h-5 mr-2" />
                Componentes
              </h3>
              <button
                type="button"
                onClick={addComponent}
                className="px-3 py-1 bg-blue-600 text-white text-sm rounded-lg hover:bg-blue-700 transition-colors flex items-center"
              >
                <Plus className="w-4 h-4 mr-1" />
                Agregar Componente
              </button>
            </div>
            
            <div className="space-y-4">
              {formData.components.map((component, index) => (
                <div key={index} className="p-4 bg-gray-50 rounded-lg border">
                  <div className="grid grid-cols-1 md:grid-cols-2 gap-4 mb-4">
                    <div>
                      <label className="block text-sm font-medium text-gray-600 mb-1">Tipo</label>
                      <select
                        value={component.type}
                        onChange={(e) => handleArrayChange('components', index, 'type', e.target.value)}
                        className="w-full px-3 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent"
                      >
                        <option value="base">Base</option>
                        <option value="addon">Adicional</option>
                      </select>
                    </div>
                    <div>
                      <label className="block text-sm font-medium text-gray-600 mb-1">Nombre</label>
                      <input
                        type="text"
                        value={component.name}
                        onChange={(e) => handleArrayChange('components', index, 'name', e.target.value)}
                        className="w-full px-3 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent"
                      />
                    </div>
                    <div>
                      <label className="block text-sm font-medium text-gray-600 mb-1">Cantidad</label>
                      <input
                        type="text"
                        value={component.quantity || ''}
                        onChange={(e) => handleArrayChange('components', index, 'quantity', e.target.value)}
                        className="w-full px-3 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent"
                      />
                    </div>
                    <div>
                      <label className="block text-sm font-medium text-gray-600 mb-1">Descripción</label>
                      <input
                        type="text"
                        value={component.description || ''}
                        onChange={(e) => handleArrayChange('components', index, 'description', e.target.value)}
                        className="w-full px-3 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent"
                      />
                    </div>
                  </div>
                  
              <div className="grid grid-cols-1 md:grid-cols-3 gap-4 mb-4">
                    <div className="flex items-center space-x-2">
                      <input
                        type="checkbox"
                        checked={component.required}
                        onChange={(e) => handleArrayChange('components', index, 'required', e.target.checked)}
                        className="w-4 h-4 text-blue-600 bg-gray-100 border-gray-300 rounded focus:ring-blue-500"
                      />
                      <label className="text-sm font-medium text-gray-700">Requerido</label>
                    </div>
                    {component.type === 'addon' && (
                      <>
                <div>
                          <label className="block text-sm font-medium text-gray-600 mb-1">Precio</label>
                          <input
                            type="number"
                            value={component.price || 0}
                            onChange={(e) => handleArrayChange('components', index, 'price', parseInt(e.target.value) || 0)}
                            className="w-full px-3 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent"
                            min="0"
                          />
                        </div>
                        <div>
                          <label className="block text-sm font-medium text-gray-600 mb-1">Costo por Unidad</label>
                          <input
                            type="number"
                            value={component.cost?.unitCost || 0}
                            onChange={(e) => handleArrayChange('components', index, 'cost', { unitCost: parseInt(e.target.value) || 0 })}
                            className="w-full px-3 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent"
                            min="0"
                          />
                        </div>
                      </>
                    )}
                  </div>
                  
                  <div className="flex justify-end">
                    <button
                      type="button"
                      onClick={() => removeComponent(index)}
                      className="p-2 hover:bg-red-100 rounded transition-colors"
                      title="Eliminar componente"
                    >
                      <X className="w-4 h-4 text-red-500" />
                    </button>
                  </div>
                </div>
              ))}
            </div>
          </CardContent>
        </Card>

        {/* logistics */}
        <Card>
          <CardContent className="p-6">
            <h3 className="text-lg font-semibold text-gray-800 mb-4 flex items-center">
              <Truck className="w-5 h-5 mr-2" />
              Logística
            </h3>
            <div className="space-y-6">
              {/* dimensions */}
              <div>
                <h4 className="text-md font-medium text-gray-700 mb-3">Dimensiones</h4>
                <div className="grid grid-cols-1 md:grid-cols-3 gap-4">
                  <div>
                    <label className="block text-sm font-medium text-gray-600 mb-1">Alto</label>
                  <input
                    type="number"
                    value={formData.logistics.dimensions.height}
                    onChange={(e) => handleInputChange('logistics.dimensions.height', parseInt(e.target.value) || 0)}
                    className="w-full px-3 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent"
                    min="0"
                  />
                </div>
                <div>
                    <label className="block text-sm font-medium text-gray-600 mb-1">Largo</label>
                  <input
                    type="number"
                    value={formData.logistics.dimensions.length}
                    onChange={(e) => handleInputChange('logistics.dimensions.length', parseInt(e.target.value) || 0)}
                    className="w-full px-3 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent"
                    min="0"
                  />
                </div>
                <div>
                    <label className="block text-sm font-medium text-gray-600 mb-1">Ancho</label>
                  <input
                    type="number"
                    value={formData.logistics.dimensions.width}
                    onChange={(e) => handleInputChange('logistics.dimensions.width', parseInt(e.target.value) || 0)}
                    className="w-full px-3 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent"
                    min="0"
                  />
                  </div>
                </div>
              </div>

              <div>
                <label className="block text-sm font-medium text-gray-700 mb-1">Peso</label>
                <input
                  type="number"
                  value={formData.logistics.weight}
                  onChange={(e) => handleInputChange('logistics.weight', parseInt(e.target.value) || 0)}
                  className="w-full px-3 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent"
                  min="0"
                />
              </div>

              {/* availabilityTime */}
              <div>
                <div className="flex items-center justify-between mb-4">
                  <h4 className="text-md font-medium text-gray-700">Horarios de Disponibilidad</h4>
                  <button
                    type="button"
                    onClick={addAvailabilityTime}
                    className="px-3 py-1 bg-blue-600 text-white text-sm rounded-lg hover:bg-blue-700 transition-colors flex items-center"
                  >
                    <Plus className="w-4 h-4 mr-1" />
                    Agregar Horario
                  </button>
                </div>
                <div className="space-y-3">
                {formData.logistics.availabilityTime.map((time, index) => (
                    <div key={index} className="p-3 bg-gray-50 rounded-lg">
                      <div className="grid grid-cols-1 md:grid-cols-3 gap-4 mb-3">
                        <div>
                          <label className="block text-sm font-medium text-gray-600 mb-1">Desde</label>
                          <input
                            type="time"
                            value={time.timeRange.from}
                            onChange={(e) => handleArrayChange('logistics.availabilityTime', index, 'timeRange', { ...time.timeRange, from: e.target.value })}
                            className="w-full px-3 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent"
                          />
                        </div>
                        <div>
                          <label className="block text-sm font-medium text-gray-600 mb-1">Hasta</label>
                          <input
                            type="time"
                            value={time.timeRange.to}
                            onChange={(e) => handleArrayChange('logistics.availabilityTime', index, 'timeRange', { ...time.timeRange, to: e.target.value })}
                            className="w-full px-3 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent"
                          />
                        </div>
                        <div className="flex items-center justify-end">
                          <button
                            type="button"
                            onClick={() => removeAvailabilityTime(index)}
                            className="p-2 hover:bg-red-100 rounded transition-colors"
                            title="Eliminar horario"
                          >
                            <X className="w-4 h-4 text-red-500" />
                          </button>
                        </div>
                    </div>
                      <div>
                        <label className="block text-sm font-medium text-gray-600 mb-2">Días de la Semana</label>
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
              </div>

              {/* deliveryFees */}
              <div>
                <div className="flex items-center justify-between mb-4">
                  <h4 className="text-md font-medium text-gray-700">Tarifas de Envío</h4>
                  <button
                    type="button"
                    onClick={addDeliveryFee}
                    className="px-3 py-1 bg-blue-600 text-white text-sm rounded-lg hover:bg-blue-700 transition-colors flex items-center"
                  >
                    <Plus className="w-4 h-4 mr-1" />
                    Agregar Tarifa
                  </button>
                </div>
                <div className="space-y-3">
                  {formData.logistics.deliveryFees.map((fee, index) => (
                    <div key={index} className="p-3 bg-gray-50 rounded-lg">
                      <div className="grid grid-cols-1 md:grid-cols-4 gap-4 mb-3">
                        <div>
                          <label className="block text-sm font-medium text-gray-600 mb-1">Condición</label>
                          <input
                            type="text"
                            value={fee.condition}
                            onChange={(e) => handleArrayChange('logistics.deliveryFees', index, 'condition', e.target.value)}
                            className="w-full px-3 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent"
                          />
                        </div>
                        <div>
                          <label className="block text-sm font-medium text-gray-600 mb-1">Tipo</label>
                          <select
                            value={fee.type}
                            onChange={(e) => handleArrayChange('logistics.deliveryFees', index, 'type', e.target.value)}
                            className="w-full px-3 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent"
                          >
                            <option value="fixed">Fijo</option>
                            <option value="percentage">Porcentaje</option>
                          </select>
                        </div>
                        <div>
                          <label className="block text-sm font-medium text-gray-600 mb-1">Valor</label>
                          <input
                            type="number"
                            value={fee.value}
                            onChange={(e) => handleArrayChange('logistics.deliveryFees', index, 'value', parseInt(e.target.value) || 0)}
                            className="w-full px-3 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent"
                            min="0"
                          />
                        </div>
                        <div className="flex items-center justify-end">
                          <button
                            type="button"
                            onClick={() => removeDeliveryFee(index)}
                            className="p-2 hover:bg-red-100 rounded transition-colors"
                            title="Eliminar tarifa"
                          >
                            <X className="w-4 h-4 text-red-500" />
                          </button>
                        </div>
                      </div>
                      <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
                        <div>
                          <label className="block text-sm font-medium text-gray-600 mb-1">Desde</label>
                          <input
                            type="time"
                            value={fee.timeRange.from}
                            onChange={(e) => handleArrayChange('logistics.deliveryFees', index, 'timeRange', { ...fee.timeRange, from: e.target.value })}
                            className="w-full px-3 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent"
                          />
                        </div>
                        <div>
                          <label className="block text-sm font-medium text-gray-600 mb-1">Hasta</label>
                          <input
                            type="time"
                            value={fee.timeRange.to}
                            onChange={(e) => handleArrayChange('logistics.deliveryFees', index, 'timeRange', { ...fee.timeRange, to: e.target.value })}
                            className="w-full px-3 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent"
                          />
                        </div>
                      </div>
                  </div>
                ))}
              </div>
            </div>
            </div>
          </CardContent>
        </Card>

        {/* Botones de Acción */}
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

      {/* Modal para agregar categoría */}
      {showCategoryModal && (
        <div className="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center z-50">
          <div className="bg-white rounded-lg p-6 w-full max-w-md mx-4">
            <div className="flex items-center justify-between mb-4">
              <h3 className="text-lg font-semibold text-gray-800">Agregar Categoría</h3>
              <button
                type="button"
                onClick={() => setShowCategoryModal(false)}
                className="text-gray-400 hover:text-gray-600"
              >
                <X className="w-5 h-5" />
              </button>
            </div>

            <div className="space-y-4">
              <div>
                <label className="block text-sm font-medium text-gray-700 mb-1">Nombre</label>
                <input
                  type="text"
                  value={newCategoryName}
                  onChange={(e) => setNewCategoryName(e.target.value)}
                  placeholder="Nueva categoría"
                  className="w-full px-3 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-green-500 focus:border-green-500"
                />
              </div>

              <div>
                <label className="block text-sm font-medium text-gray-700 mb-1">Categoría Madre</label>
                <div className="relative">
                  <select
                    value={newCategoryParent}
                    onChange={(e) => setNewCategoryParent(e.target.value)}
                    className="w-full px-3 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-green-500 focus:border-green-500 appearance-none pr-10"
                  >
                    <option value="">Selecciona una categoría</option>
                    {formData.categories.map((cat) => (
                      <option key={cat.id} value={cat.id}>
                        {getCategoryPath(cat.id)}
                      </option>
                    ))}
                  </select>
                  <ChevronDown className="absolute right-3 top-1/2 transform -translate-y-1/2 w-4 h-4 text-gray-400 pointer-events-none" />
                </div>
              </div>

              {/* Vista previa del ID que se generará */}
              {newCategoryName && (
                <div className="p-3 bg-gray-50 rounded-lg">
                  <div className="text-sm text-gray-600">ID que se generará:</div>
                  <div className="text-sm font-mono text-gray-800">
                    {generateCategoryId(newCategoryName, newCategoryParent || null)}
                  </div>
                </div>
              )}
            </div>

            <div className="flex justify-end gap-3 mt-6">
              <button
                type="button"
                onClick={() => setShowCategoryModal(false)}
                className="px-4 py-2 border border-gray-300 text-gray-700 rounded-lg hover:bg-gray-50 transition-colors"
              >
                Cancelar
              </button>
              <button
                type="button"
                onClick={handleAddCategory}
                disabled={!newCategoryName.trim()}
                className="px-4 py-2 bg-green-600 text-white rounded-lg hover:bg-green-700 transition-colors disabled:opacity-50 disabled:cursor-not-allowed"
              >
                Agregar Categoría
              </button>
            </div>
          </div>
        </div>
      )}
    </div>
  )
}

export default CreateProduct