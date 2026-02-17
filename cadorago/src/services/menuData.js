/**
 * Servicio para manejar el nuevo contrato de menú
 * Convierte el nuevo formato a un formato compatible con los componentes existentes
 */

/**
 * Obtiene el costo de un item o side basado en su pricing (costPerUnit)
 * @param {Object} pricing - Objeto de pricing con mode, unit, costPerUnit, baseUnit
 * @param {number} quantity - Cantidad (opcional, por defecto 1)
 * @returns {number} Costo total
 */
export function getCostFromPricing(pricing, quantity = 1) {
  if (!pricing || pricing.costPerUnit == null) return 0;
  const costPerUnit = pricing.costPerUnit || 0;
  if (pricing.mode === 'UNIT') {
    return costPerUnit * (quantity ?? 1);
  }
  const qty = quantity ?? (pricing.baseUnit || 1);
  const baseUnit = pricing.baseUnit || 1;
  return (qty / baseUnit) * costPerUnit;
}

/**
 * Obtiene el precio de un item o side basado en su pricing
 * @param {Object} pricing - Objeto de pricing con mode, unit, pricePerUnit, baseUnit
 * @param {number} quantity - Cantidad (opcional, por defecto usa baseUnit)
 * @returns {number} Precio a mostrar
 */
export function getPriceFromPricing(pricing, quantity = null) {
  if (!pricing) return 0;
  
  // Para modo UNIT: precio = pricePerUnit * cantidad (cantidad por defecto 1)
  if (pricing.mode === 'UNIT') {
    const qty = quantity !== null && quantity !== undefined ? quantity : 1;
    return (pricing.pricePerUnit || 0) * qty;
  }
  
  // Para otros modos (WEIGHT, VOLUME, etc.), calcular según cantidad
  // Si no se proporciona cantidad, usar baseUnit como referencia
  const qty = quantity !== null ? quantity : (pricing.baseUnit || 1);
  const baseUnit = pricing.baseUnit || 1;
  
  // Fórmula: precio = (cantidad / baseUnit) * pricePerUnit
  return (qty / baseUnit) * (pricing.pricePerUnit || 0);
}

/**
 * Obtiene los límites recomendados para un slider basado en el modo de pricing
 * @param {Object} pricing - Objeto de pricing
 * @returns {{min: number, max: number, step: number}}
 */
export function getPricingLimits(pricing) {
  if (!pricing) {
    return { min: 0, max: 1000, step: 1 };
  }
  
  const baseUnit = pricing.baseUnit || 1;
  
  switch (pricing.mode) {
    case 'WEIGHT':
      if (pricing.unit === 'GRAM') {
        return {
          min: baseUnit,
          max: 5000, // 5kg máximo
          step: 10 // Incrementos de 10g
        };
      } else if (pricing.unit === 'KILOGRAM') {
        return {
          min: baseUnit,
          max: 10, // 10kg máximo
          step: 0.1 // Incrementos de 100g
        };
      }
      break;
      
    case 'VOLUME':
      if (pricing.unit === 'MILLILITER') {
        return {
          min: baseUnit,
          max: 5000, // 5L máximo
          step: 50 // Incrementos de 50ml
        };
      } else if (pricing.unit === 'LITER') {
        return {
          min: baseUnit,
          max: 10, // 10L máximo
          step: 0.1 // Incrementos de 100ml
        };
      }
      break;
      
    case 'LENGTH':
      return {
        min: baseUnit,
        max: 100, // 100m máximo
        step: 0.1 // Incrementos de 10cm
      };
      
    case 'AREA':
      return {
        min: baseUnit,
        max: 100, // 100m² máximo
        step: 0.1 // Incrementos de 0.1m²
      };
  }
  
  // Default
  return {
    min: baseUnit,
    max: 1000,
    step: 1
  };
}

/**
 * Obtiene los bloques de description que tienen selectables (preferencias sin precio).
 * @param {Object} item - Item del menú
 * @returns {Array<{id: string, base: string, languages: object, selectables: object}>}
 */
export function getDescriptionSelectablesForItem(item) {
  const desc = item?.description || [];
  return desc.filter(d => d.selectables && d.selectables.options?.length > 0);
}

/**
 * Convierte el nuevo formato de menú al formato esperado por los componentes
 * @param {Object} menuData - Datos del menú en el nuevo formato
 * @returns {Object} Datos en formato compatible
 */
export function adaptMenuData(menuData) {
  if (!menuData) return null;
  
  // Adaptar el menú: preservar estructura multiidioma (no convertir a strings)
  const adaptedMenu = (menuData.menu || []).map(category => ({
    title: category.title ?? '', // Preservar estructura multiidioma o string tal cual viene
    items: (category.items || []).map(item => ({
      id: item.id ?? '',
      title: item.title ?? '', // Preservar estructura multiidioma o string tal cual viene
      description: item.description ?? [],
      foodAttributes: item.foodAttributes ?? [], // Preservar atributos alimentarios
      price: getPriceFromPricing(item.pricing),
      pricing: item.pricing, // Mantener el pricing original para uso futuro
      photoUrl: item.photoUrl ?? '', // Incluir photoUrl del item
      station: item.station ?? null, // KITCHEN | BAR para vista y pedido
      sides: (item.sides || []).map(side => ({
        id: side.id ?? '',
        name: side.name ?? '', // Preservar estructura multiidioma o string tal cual viene
        foodAttributes: side.foodAttributes ?? [], // Preservar atributos alimentarios
        price: getPriceFromPricing(side.pricing),
        pricing: side.pricing, // Mantener el pricing original
        photoUrl: side.photoUrl ?? '', // Incluir photoUrl del side
        station: side.station ?? null // KITCHEN | BAR
      }))
    }))
  }));
  
  // presentationStyle del contrato: HERO | MODERN. Por defecto HERO (cards horizontales, imagen a la derecha).
  const presentationStyle = (menuData.presentationStyle || 'HERO').toUpperCase();
  const normalized = presentationStyle === 'MODERN' ? 'MODERN' : 'HERO';

  // Preservar businessInfo y asegurar currency por defecto (CLP) si no viene del backend
  const businessInfo = menuData.businessInfo
    ? { ...menuData.businessInfo, currency: menuData.businessInfo.currency || 'CLP' }
    : { currency: 'CLP' };

  return {
    ...menuData,
    businessInfo,
    menu: adaptedMenu,
    presentationStyle: normalized
  };
}

/**
 * JSON por defecto para la ruta /test
 */
export const DEFAULT_TEST_MENU = {
  "id": "menu-123e4567-e89b-12d3-a456-426614174000",
  "presentationStyle": "HERO",
  "coverImage": "https://example.com/images/cover.jpg",
  "footerImage": "https://example.com/images/logo.png",
  "businessInfo": {
    "businessName": "La Pizzería del Centro",
    "whatsapp": "+56912345678",
    "currency": "CLP",
    "businessHours": [
      "Lunes a Viernes: 12:00 - 22:00",
      "Sábado y Domingo: 13:00 - 23:00"
    ]
  },
  "menu": [
    {
      "title": "Pizzas",
      "items": [
        {
          "title": "Pizza Margherita",
          "description": [{ "base": "Tomate, mozzarella y albahaca fresca", "languages": {} }],
          "pricing": {
            "mode": "UNIT",
            "unit": "EACH",
            "pricePerUnit": 8990,
            "baseUnit": 0
          },
          "sides": [
            {
              "name": "Tamaño Grande",
              "pricing": {
                "mode": "UNIT",
                "unit": "EACH",
                "pricePerUnit": 11990,
                "baseUnit": 0
              }
            },
            {
              "name": "Extra Queso",
              "pricing": {
                "mode": "UNIT",
                "unit": "EACH",
                "pricePerUnit": 2000,
                "baseUnit": 0
              }
            }
          ]
        },
        {
          "title": "Pizza Pepperoni",
          "description": [{ "base": "Pepperoni, queso mozzarella y orégano", "languages": {} }],
          "pricing": {
            "mode": "UNIT",
            "unit": "EACH",
            "pricePerUnit": 10990,
            "baseUnit": 0
          }
        }
      ]
    },
    {
      "title": "Bebidas",
      "items": [
        {
          "title": "Coca Cola",
          "description": [{ "base": "Refresco de cola 500ml", "languages": {} }],
          "pricing": {
            "mode": "UNIT",
            "unit": "EACH",
            "pricePerUnit": 1500,
            "baseUnit": 0
          },
          "sides": [
            {
              "name": "Tamaño 1.5L",
              "pricing": {
                "mode": "UNIT",
                "unit": "EACH",
                "pricePerUnit": 2500,
                "baseUnit": 0
              }
            }
          ]
        },
        {
          "title": "Agua Mineral",
          "pricing": {
            "mode": "UNIT",
            "unit": "EACH",
            "pricePerUnit": 800,
            "baseUnit": 0
          }
        }
      ]
    },
    {
      "title": "Frutas y Verduras",
      "items": [
        {
          "title": "Palta Hass",
          "description": [{ "base": "Palta Hass fresca, venta por kilos", "languages": {} }],
          "pricing": {
            "mode": "WEIGHT",
            "unit": "KILOGRAM",
            "pricePerUnit": 6000,
            "baseUnit": 1
          }
        }
      ]
    }
  ],
  "deliveryOptions": [
    {
      "type": "PICKUP",
      "requireTime": true,
      "timeRequestType": "WINDOW",
      "timeWindows": [
        {
          "start": "12:00",
          "end": "14:00"
        },
        {
          "start": "18:00",
          "end": "22:00"
        }
      ]
    }
  ]
};
