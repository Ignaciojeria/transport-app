/**
 * Servicio de datos del restaurante Cadorago
 */
export const restaurantData = {
  contacto: {
    whatssap: "+56957857558"
  },
  horariosDeAtencion: [
    "lunes a viernes de 9 am a 4 de la tarde"
  ],
  carta: [
    {
      titulo: "menu de la casa",
      items: [
        {
          titulo: "Lentejas Caseras",
          descripción: "sabrosas lentejas con verduras y arroz",
          precio: 7000
        },
        {
          titulo: "Fideos con salsa blanca",
          descripción: "sabrosos fideos con salsa alfredo",
          precio: 7000
        }
      ]
    },
    {
      titulo: "Bebestibles",
      items: [
        {
          titulo: "CocaCola Zero",
          descripción: "",
          precio: 1000
        },
        {
          titulo: "Fideos con salsa blanca",
          descripción: "",
          precio: 7000
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

