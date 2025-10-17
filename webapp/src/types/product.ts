// Interfaces del contrato anterior (mantenidas para compatibilidad)
export interface PromptItem {
  questionMarkdown: string
  type: 'text' | 'select' | 'multiselect' | 'number' | 'boolean'
  options?: string[]
  required?: boolean
  placeholderMarkdown?: string
}

// Nuevas interfaces basadas en el contrato actualizado
export interface ProductStatus {
  isAvailable: boolean
  isFeatured: boolean
  allowReviews: boolean
}

export interface Attachment {
  name: string
  description: string
  url: string
  type: string
  sizeKb: number
}

export interface ProductProperties {
  sku: string
  brand: string
  barcode: string
}

export interface PurchaseCondition {
  minUnits?: number
  maxUnits?: number
  multiplesOf?: number
  minWeight?: number
  maxWeight?: number
  minVolume?: number
  maxVolume?: number
  notes?: string
}

export interface PurchaseConditions {
  fixed: PurchaseCondition
  weight: PurchaseCondition
  volume: PurchaseCondition
}

export interface Attribute {
  name: string
  value: string
}

export interface Category {
  id: string
  name: string
  parent: string | null
}

export interface DigitalBundleAccess {
  method: string
  url: string
  expiresInDays: number
}

export interface DigitalBundle {
  hasDigitalContent: boolean
  type: string
  title: string
  description: string
  access: DigitalBundleAccess
}

export interface Video {
  title: string
  platform: string
  url: string
  thumbnail: string
}

export interface Media {
  videos: Video[]
  gallery: string[]
}

export interface PaymentInfo {
  currency: string
  methods: string[]
  provider: string
}

export interface StockInfo {
  fixed: {
    availableUnits: number
  }
  weight: {
    availableWeight: number
  }
  volume: {
    availableVolume: number
  }
}

export interface PriceInfo {
  fixedPrice: number
  weight: {
    unitSize: number
    pricePerUnit: number
  }
  volume: {
    unitSize: number
    pricePerUnit: number
  }
}

export interface CostInfo {
  fixedCost: number
  weight: {
    unitSize: number
    costPerUnit: number
  }
  volume: {
    unitSize: number
    constPerUnit: number
  }
}

export interface ComponentStock {
  fixed?: {
    availableUnits: number
  }
  weight?: {
    availableWeight: number
  }
  volume?: {
    availableVolume: number
  }
}

export interface ComponentCost {
  unitCost: number
}

export interface Component {
  type: 'base' | 'addon'
  name: string
  quantity?: string
  description?: string
  required: boolean
  price?: number
  cost?: ComponentCost
  stock: ComponentStock
}

export interface DeliveryFee {
  condition: string
  type: string
  value: number
  timeRange: {
    from: string
    to: string
  }
}

export interface LogisticsInfo {
  dimensions: {
    height: number
    length: number
    width: number
  }
  weight: number
  availabilityTime: AvailabilityTime[]
  deliveryFees: DeliveryFee[]
}

export interface AvailabilityTime {
  timeRange: {
    from: string
    to: string
  }
  daysOfWeek: string[]
}

// Interfaz principal del producto actualizada
export interface Product {
  referenceID: string
  name: string
  descriptionMarkdown: string
  status: ProductStatus
  attachments: Attachment[]
  properties: ProductProperties
  purchaseConditions: PurchaseConditions
  attributes: Attribute[]
  categories: Category[]
  digitalBundle: DigitalBundle
  welcomeMessageMarkdown: string
  media: Media
  payment: PaymentInfo
  stock: StockInfo
  price: PriceInfo
  cost: CostInfo
  components: Component[]
  logistics: LogisticsInfo
}

export interface CreateProductRequest extends Product {}
