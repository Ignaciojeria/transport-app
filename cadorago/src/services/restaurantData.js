/**
 * Servicio de datos del restaurante Cadorago
 */
export const restaurantData = {
  "id": "f8679ed4-1b19-49ff-9d97-0b51967a86bd",
  "coverImage":"https://storage.googleapis.com/micartapro-menus/menus/01KCW67YKSV455GBVDT88S4072/gallery/portadav2.webp",
  "footerImage":"https://storage.googleapis.com/micartapro-menus/menus/01KCW67YKSV455GBVDT88S4072/gallery/logov2.webp",
  "businessInfo": {
    "businessName": "cadorago",
    "whatsapp": "+56957857558",
    "businessHours": [
      "Lunes a Martes: 9h a 16h",
      "Miércoles, Jueves, Sábado y Domingo: hasta las 20h",
      "Viernes: Cerrado"
    ]
  },
  "menu": [
    {
      "title": "Menú",
      "items": [
        {
          "title": "Pollo a la plancha",
          "sides": [
            { "name": "Con puré", "price": 3000 },
            { "name": "Con arroz", "price": 3000 }
          ]
        },
        {
          "title": "Completo italiano",
          "price": 2500
        },
        {
          "title": "Hamburguesa",
          "sides": [
            { "name": "Sola", "price": 4100 }
          ]
        },
        {
          "title": "chacareros",
          "price": 7000
        }
      ]
    },
    {
      "title": "Postres",
      "items": [
        {
          "title": "mote con huesillo",
          "price": 4000
        },
        {
          "title": "leche asada",
          "price": 3000
        }
      ]
    }
  ]
};

/**
 * Obtiene los datos del restaurante
 * @returns {Object} Datos del restaurante
 */
export function getRestaurantData() {
  return restaurantData;
}

