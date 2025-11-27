/**
 * Servicio de datos del restaurante Cadorago
 */
export const restaurantData = {
  "contacto": {
    "whatssap": "+56944930403"
  },
  "horariosDeAtencion": [
    "lunes a viernes de 9 am a 4 de la tarde"
  ],
  "carta": [
    {
      "titulo": "Colaciones Diarias",
      "items": [
        {
          "titulo": "Pescado frito",
          "descripcion": "",
          "acompanamientos": [
            { "id": "pescado_porcion_sola", "nombre": "Porción sola", "precio": 2800 },
            { "id": "pescado_arroz", "nombre": "Arroz", "precio": 4200 },
            { "id": "pescado_pure", "nombre": "Puré", "precio": 4200 },
            { "id": "pescado_papas_fritas", "nombre": "Papas fritas", "precio": 4200 },
            { "id": "pescado_arroz_papas", "nombre": "Arroz + papas fritas", "precio": 5400 },
            { "id": "pescado_ensalada_chica", "nombre": "Ensalada chica", "precio": 3800 },
            { "id": "pescado_ensalada_grande", "nombre": "Ensalada grande", "precio": 5400 }
          ]
        },
        {
          "titulo": "Lomito de cerdo al horno",
          "descripcion": "",
          "acompanamientos": [
            { "id": "lomito_porcion_sola", "nombre": "Porción sola", "precio": 2800 },
            { "id": "lomito_arroz", "nombre": "Arroz", "precio": 4200 },
            { "id": "lomito_pure", "nombre": "Puré", "precio": 4200 },
            { "id": "lomito_papas_fritas", "nombre": "Papas fritas", "precio": 4200 },
            { "id": "lomito_arroz_papas", "nombre": "Arroz + papas fritas", "precio": 5400 },
            { "id": "lomito_ensalada_chica", "nombre": "Ensalada chica", "precio": 3800 },
            { "id": "lomito_ensalada_grande", "nombre": "Ensalada grande", "precio": 5400 }
          ]
        },
        {
          "titulo": "Pollo a la plancha",
          "descripcion": "",
          "acompanamientos": [
            { "id": "pollo_porcion_sola", "nombre": "Porción sola", "precio": 2800 },
            { "id": "pollo_arroz", "nombre": "Arroz", "precio": 4200 },
            { "id": "pollo_pure", "nombre": "Puré", "precio": 4200 },
            { "id": "pollo_papas_fritas", "nombre": "Papas fritas", "precio": 4200 },
            { "id": "pollo_arroz_papas", "nombre": "Arroz + papas fritas", "precio": 5400 },
            { "id": "pollo_ensalada_chica", "nombre": "Ensalada chica", "precio": 3800 },
            { "id": "pollo_ensalada_grande", "nombre": "Ensalada grande", "precio": 4200 }
          ]
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

