export interface Product {
  referenceID: string
  name: string
  description: string
  image: string
  payment: PaymentInfo
  stock: StockInfo
  price: PriceInfo
  logistics: LogisticsInfo
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

export interface LogisticsInfo {
  dimensions: {
    height: number
    length: number
    width: number
  }
  availabilityTime: AvailabilityTime[]
  costs: CostInfo[]
}

export interface AvailabilityTime {
  timeRange: {
    from: string
    to: string
  }
  daysOfWeek: string[]
}

export interface CostInfo {
  condition: string
  type: string
  value: number
  timeRange: {
    from: string
    to: string
  }
}

export interface CreateProductRequest {
  referenceID: string
  name: string
  description: string
  image: string
  payment: PaymentInfo
  stock: StockInfo
  price: PriceInfo
  logistics: LogisticsInfo
}
