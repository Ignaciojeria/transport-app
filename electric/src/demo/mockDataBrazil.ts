import type { Route } from '../domain/route'

export const mockRouteBrazil: Route = {
  id: 1,
  referenceID: "ROTA-UBA-001",
  planReferenceID: "PLANO-UBATUBA-001",
  createdAt: "2025-09-16T10:00:00Z",
  geometry: {
    type: "linestring",
    encoding: "polyline",
    value: "_p~iF~ps|U_ulLnnqC_mqNvxq@"
  },
  vehicle: {
    plate: "ABC1D23",
    capacity: {
      weight: 1000,
      volume: 1000,
      insurance: 15000
    },
    skills: ["refrigerado", "fragil"],
    startLocation: {
      addressInfo: {
        addressLine1: "Av. Iperoig, 690",
        addressLine2: "Centro",
        coordinates: {
          latitude: -23.433889,
          longitude: -45.071111
        },
        politicalArea: {
          adminAreaLevel1: "São Paulo",
          adminAreaLevel2: "Ubatuba",
          adminAreaLevel3: "Centro",
          adminAreaLevel4: "",
          code: "BR-SP"
        },
        zipCode: "11680-000"
      },
      nodeInfo: {
        referenceID: "NODE-CENTRO-UBA"
      }
    },
    endLocation: {
      addressInfo: {
        addressLine1: "Av. Iperoig, 690",
        addressLine2: "Centro", 
        coordinates: {
          latitude: -23.433889,
          longitude: -45.071111
        },
        politicalArea: {
          adminAreaLevel1: "São Paulo",
          adminAreaLevel2: "Ubatuba",
          adminAreaLevel3: "Centro",
          adminAreaLevel4: "",
          code: "BR-SP"
        },
        zipCode: "11680-000"
      },
      nodeInfo: {
        referenceID: "NODE-CENTRO-UBA"
      }
    },
    timeWindow: {
      start: "08:00:00",
      end: "18:00:00"
    }
  },
  visits: [
    // Visita 1: Múltiples clientes en la misma dirección (Praia Grande)
    {
      sequenceNumber: 1,
      type: "delivery",
      unassignedReason: "",
      addressInfo: {
        addressLine1: "Rua Santa Rita, 156",
        addressLine2: "Praia Grande",
        coordinates: {
          latitude: -23.4342,
          longitude: -45.0834
        },
        politicalArea: {
          adminAreaLevel1: "São Paulo",
          adminAreaLevel2: "Ubatuba",
          adminAreaLevel3: "Praia Grande",
          adminAreaLevel4: "",
          code: "BR-SP"
        },
        zipCode: "11680-220"
      },
      nodeInfo: {
        referenceID: "NODE-PRAIA-GRANDE-001"
      },
      serviceTime: 10,
      timeWindow: {
        start: "09:00:00",
        end: "12:00:00"
      },
      orders: [
        // Cliente 1: João Silva (3 pedidos)
        {
          referenceID: "PEDIDO-BR-001A",
          contact: {
            fullName: "João Silva Santos",
            phone: "+55 12 99876-5432",
            email: "joao.santos@email.com.br",
            nationalID: "123.456.789-01"
          },
          instructions: "Entregar na recepção do prédio",
          deliveryUnits: [
            {
              lpn: "BR-LPN-001A",
              weight: 2.5,
              volume: 0.01,
              price: 899.90,
              skills: ["fragil"],
              items: [
                {
                  sku: "SAMS-A54-128-BR",
                  description: "Smartphone Samsung Galaxy A54 128GB",
                  quantity: 1
                }
              ],
              evidences: []
            }
          ]
        },
        {
          referenceID: "PEDIDO-BR-001B", 
          contact: {
            fullName: "João Silva Santos",
            phone: "+55 12 99876-5432",
            email: "joao.santos@email.com.br",
            nationalID: "123.456.789-01"
          },
          instructions: "Segundo pedido - mesmo cliente",
          deliveryUnits: [
            {
              lpn: "BR-LPN-001B",
              weight: 0.8,
              volume: 0.005,
              price: 299.90,
              skills: ["fragil"],
              items: [
                {
                  sku: "CASE-SAMS-BLU-BR",
                  description: "Capa para Samsung Galaxy Azul",
                  quantity: 1
                },
                {
                  sku: "PELIC-TEMP-SAMS-BR",
                  description: "Película de vidro temperado",
                  quantity: 1
                }
              ],
              evidences: []
            }
          ]
        },
        {
          referenceID: "PEDIDO-BR-001C",
          contact: {
            fullName: "João Silva Santos", 
            phone: "+55 12 99876-5432",
            email: "joao.santos@email.com.br",
            nationalID: "123.456.789-01"
          },
          instructions: "Terceiro pedido - fones de ouvido",
          deliveryUnits: [
            {
              lpn: "BR-LPN-001C",
              weight: 0.3,
              volume: 0.002,
              price: 189.90,
              skills: [],
              items: [
                {
                  sku: "FONE-BT-JBL-BR",
                  description: "Fone Bluetooth JBL Tune 510BT",
                  quantity: 1
                }
              ],
              evidences: []
            }
          ]
        },
        // Cliente 2: Maria Oliveira
        {
          referenceID: "PEDIDO-BR-002",
          contact: {
            fullName: "Maria Oliveira Costa",
            phone: "+55 12 98765-4321", 
            email: "maria.oliveira@gmail.com",
            nationalID: "987.654.321-02"
          },
          instructions: "Ligar antes de entregar",
          deliveryUnits: [
            {
              lpn: "BR-LPN-002",
              weight: 1.2,
              volume: 0.008,
              price: 599.90,
              skills: ["fragil"],
              items: [
                {
                  sku: "TAB-SAMS-A8-BR",
                  description: "Tablet Samsung Galaxy Tab A8 32GB",
                  quantity: 1
                }
              ],
              evidences: []
            }
          ]
        },
        // Cliente 3: Carlos Santos
        {
          referenceID: "PEDIDO-BR-003",
          contact: {
            fullName: "Carlos Santos Pereira",
            phone: "+55 12 97654-3210",
            email: "carlos.pereira@hotmail.com", 
            nationalID: "456.789.123-03"
          },
          instructions: "Porteiro recebe as entregas",
          deliveryUnits: [
            {
              lpn: "BR-LPN-003",
              weight: 3.8,
              volume: 0.015,
              price: 1299.90,
              skills: ["fragil"],
              items: [
                {
                  sku: "NOTE-DELL-I5-BR",
                  description: "Notebook Dell Inspiron 15 i5 8GB",
                  quantity: 1
                }
              ],
              evidences: []
            }
          ]
        }
      ]
    },
    // Visitas 2-10: Clientes únicos em diferentes locais de Ubatuba
    {
      sequenceNumber: 2,
      type: "delivery",
      unassignedReason: "", 
      addressInfo: {
        addressLine1: "Av. Leovigildo Dias Vieira, 1200",
        addressLine2: "Itaguá",
        coordinates: {
          latitude: -23.4410,
          longitude: -45.0691
        },
        politicalArea: {
          adminAreaLevel1: "São Paulo",
          adminAreaLevel2: "Ubatuba",
          adminAreaLevel3: "Itaguá",
          adminAreaLevel4: "",
          code: "BR-SP"
        },
        zipCode: "11680-450"
      },
      nodeInfo: {
        referenceID: "NODE-ITAGUA-001"
      },
      serviceTime: 8,
      timeWindow: {
        start: "10:00:00",
        end: "14:00:00"
      },
      orders: [
        {
          referenceID: "PEDIDO-BR-004",
          contact: {
            fullName: "Ana Paula Rodrigues",
            phone: "+55 12 96543-2109",
            email: "ana.rodrigues@outlook.com",
            nationalID: "789.123.456-04"
          },
          instructions: "Casa amarela com portão azul",
          deliveryUnits: [
            {
              lpn: "BR-LPN-004",
              weight: 0.5,
              volume: 0.003,
              price: 149.90,
              skills: [],
              items: [
                {
                  sku: "REL-FIT-XIAO-BR",
                  description: "Relógio Fitness Xiaomi Mi Band 7",
                  quantity: 1
                }
              ],
              evidences: []
            }
          ]
        }
      ]
    },
    {
      sequenceNumber: 3,
      type: "delivery",
      unassignedReason: "",
      addressInfo: {
        addressLine1: "Rua Prof. Thomaz Galhardo, 89",
        addressLine2: "Perequê-Açu", 
        coordinates: {
          latitude: -23.4521,
          longitude: -45.0563
        },
        politicalArea: {
          adminAreaLevel1: "São Paulo",
          adminAreaLevel2: "Ubatuba",
          adminAreaLevel3: "Perequê-Açu",
          adminAreaLevel4: "",
          code: "BR-SP"
        },
        zipCode: "11680-590"
      },
      nodeInfo: {
        referenceID: "NODE-PEREQUE-001"
      },
      serviceTime: 6,
      timeWindow: {
        start: "11:00:00", 
        end: "15:00:00"
      },
      orders: [
        {
          referenceID: "PEDIDO-BR-005",
          contact: {
            fullName: "Roberto Ferreira Lima",
            phone: "+55 12 95432-1098",
            email: "roberto.lima@gmail.com",
            nationalID: "321.654.987-05"
          },
          instructions: "Apto 302 - Interfone Roberto",
          deliveryUnits: [
            {
              lpn: "BR-LPN-005",
              weight: 4.2,
              volume: 0.02,
              price: 799.90,
              skills: ["pesado"],
              items: [
                {
                  sku: "TV-LG-32-SMART-BR",
                  description: "Smart TV LG 32' LED HD",
                  quantity: 1
                }
              ],
              evidences: []
            }
          ]
        }
      ]
    },
    {
      sequenceNumber: 4,
      type: "delivery",
      unassignedReason: "",
      addressInfo: {
        addressLine1: "Rua Guarani, 345",
        addressLine2: "Maranduba",
        coordinates: {
          latitude: -23.4234,
          longitude: -45.1012
        },
        politicalArea: {
          adminAreaLevel1: "São Paulo",
          adminAreaLevel2: "Ubatuba",
          adminAreaLevel3: "Maranduba", 
          adminAreaLevel4: "",
          code: "BR-SP"
        },
        zipCode: "11680-780"
      },
      nodeInfo: {
        referenceID: "NODE-MARANDUBA-001"
      },
      serviceTime: 5,
      timeWindow: {
        start: "09:30:00",
        end: "13:30:00"
      },
      orders: [
        {
          referenceID: "PEDIDO-BR-006",
          contact: {
            fullName: "Fernanda Alves Souza",
            phone: "+55 12 94321-0987",
            email: "fernanda.souza@yahoo.com.br",
            nationalID: "654.321.098-06"
          },
          instructions: "Entregar somente com a Fernanda",
          deliveryUnits: [
            {
              lpn: "BR-LPN-006",
              weight: 1.8,
              volume: 0.012,
              price: 459.90,
              skills: ["fragil"],
              items: [
                {
                  sku: "CAFE-NEST-DOLCE-BR",
                  description: "Cafeteira Nespresso Dolce Gusto Genio S",
                  quantity: 1
                }
              ],
              evidences: []
            }
          ]
        }
      ]
    },
    {
      sequenceNumber: 5,
      type: "delivery",
      unassignedReason: "",
      addressInfo: {
        addressLine1: "Av. Marginal, 1567",
        addressLine2: "Enseada",
        coordinates: {
          latitude: -23.4456,
          longitude: -45.0789
        },
        politicalArea: {
          adminAreaLevel1: "São Paulo",
          adminAreaLevel2: "Ubatuba",
          adminAreaLevel3: "Enseada",
          adminAreaLevel4: "",
          code: "BR-SP"
        },
        zipCode: "11680-120"
      },
      nodeInfo: {
        referenceID: "NODE-ENSEADA-001"
      },
      serviceTime: 7,
      timeWindow: {
        start: "12:00:00",
        end: "16:00:00"
      },
      orders: [
        {
          referenceID: "PEDIDO-BR-007",
          contact: {
            fullName: "Lucas Mendes Barbosa",
            phone: "+55 12 93210-9876",
            email: "lucas.barbosa@gmail.com",
            nationalID: "987.321.654-07"
          },
          instructions: "Casa de praia - portão verde",
          deliveryUnits: [
            {
              lpn: "BR-LPN-007", 
              weight: 0.9,
              volume: 0.006,
              price: 329.90,
              skills: [],
              items: [
                {
                  sku: "SPEAK-JBL-FLIP6-BR",
                  description: "Caixa de Som JBL Flip 6 Bluetooth",
                  quantity: 1
                }
              ],
              evidences: []
            }
          ]
        }
      ]
    },
    {
      sequenceNumber: 6,
      type: "delivery",
      unassignedReason: "",
      addressInfo: {
        addressLine1: "Rua das Palmeiras, 234",
        addressLine2: "Toninhas",
        coordinates: {
          latitude: -23.4890,
          longitude: -45.0678
        },
        politicalArea: {
          adminAreaLevel1: "São Paulo",
          adminAreaLevel2: "Ubatuba",
          adminAreaLevel3: "Toninhas",
          adminAreaLevel4: "",
          code: "BR-SP"
        },
        zipCode: "11680-310"
      },
      nodeInfo: {
        referenceID: "NODE-TONINHAS-001"
      },
      serviceTime: 6,
      timeWindow: {
        start: "13:00:00",
        end: "17:00:00"
      },
      orders: [
        {
          referenceID: "PEDIDO-BR-008",
          contact: {
            fullName: "Patrícia Gomes Silva",
            phone: "+55 12 92109-8765",
            email: "patricia.gomes@gmail.com",
            nationalID: "123.987.456-08"
          },
          instructions: "Condomínio Palmeiras - Torre B",
          deliveryUnits: [
            {
              lpn: "BR-LPN-008",
              weight: 2.1,
              volume: 0.01,
              price: 679.90,
              skills: ["fragil"],
              items: [
                {
                  sku: "AIR-FYER-PHIL-BR",
                  description: "Air Fryer Philips Walita Viva Collection",
                  quantity: 1
                }
              ],
              evidences: []
            }
          ]
        }
      ]
    },
    {
      sequenceNumber: 7,
      type: "delivery",
      unassignedReason: "",
      addressInfo: {
        addressLine1: "Rua Dr. Otávio Ribeiro, 567",
        addressLine2: "Saco da Ribeira",
        coordinates: {
          latitude: -23.5123,
          longitude: -45.0234
        },
        politicalArea: {
          adminAreaLevel1: "São Paulo",
          adminAreaLevel2: "Ubatuba",
          adminAreaLevel3: "Saco da Ribeira",
          adminAreaLevel4: "",
          code: "BR-SP"
        },
        zipCode: "11680-890"
      },
      nodeInfo: {
        referenceID: "NODE-RIBEIRA-001"
      },
      serviceTime: 8,
      timeWindow: {
        start: "14:00:00",
        end: "18:00:00"
      },
      orders: [
        {
          referenceID: "PEDIDO-BR-009",
          contact: {
            fullName: "Ricardo Oliveira Nunes",
            phone: "+55 12 91098-7654",
            email: "ricardo.nunes@hotmail.com",
            nationalID: "456.123.789-09"
          },
          instructions: "Casa com cerca de madeira",
          deliveryUnits: [
            {
              lpn: "BR-LPN-009",
              weight: 0.4,
              volume: 0.004,
              price: 189.90,
              skills: [],
              items: [
                {
                  sku: "MOUSE-LOG-MX3-BR",
                  description: "Mouse Logitech MX Master 3 Wireless",
                  quantity: 1
                }
              ],
              evidences: []
            }
          ]
        }
      ]
    },
    {
      sequenceNumber: 8,
      type: "delivery",
      unassignedReason: "",
      addressInfo: {
        addressLine1: "Av. Atlântica, 890",
        addressLine2: "Tenório",
        coordinates: {
          latitude: -23.4567,
          longitude: -45.0345
        },
        politicalArea: {
          adminAreaLevel1: "São Paulo",
          adminAreaLevel2: "Ubatuba",
          adminAreaLevel3: "Tenório",
          adminAreaLevel4: "",
          code: "BR-SP"
        },
        zipCode: "11680-540"
      },
      nodeInfo: {
        referenceID: "NODE-TENORIO-001"
      },
      serviceTime: 9,
      timeWindow: {
        start: "15:00:00",
        end: "18:00:00"
      },
      orders: [
        {
          referenceID: "PEDIDO-BR-010",
          contact: {
            fullName: "Juliana Costa Martins",
            phone: "+55 12 90987-6543",
            email: "juliana.martins@gmail.com",
            nationalID: "789.456.123-10"
          },
          instructions: "Pousada Maré Alta - Recepção",
          deliveryUnits: [
            {
              lpn: "BR-LPN-010",
              weight: 1.5,
              volume: 0.009,
              price: 399.90,
              skills: ["fragil"],
              items: [
                {
                  sku: "KINDLE-PAPER-11-BR",
                  description: "Kindle Paperwhite 11ª Geração 16GB",
                  quantity: 1
                }
              ],
              evidences: []
            }
          ]
        }
      ]
    },
    {
      sequenceNumber: 9,
      type: "delivery",
      unassignedReason: "",
      addressInfo: {
        addressLine1: "Rua das Gaivotas, 123",
        addressLine2: "Lazaro",
        coordinates: {
          latitude: -23.3456,
          longitude: -44.9789
        },
        politicalArea: {
          adminAreaLevel1: "São Paulo",
          adminAreaLevel2: "Ubatuba",
          adminAreaLevel3: "Lázaro",
          adminAreaLevel4: "",
          code: "BR-SP"
        },
        zipCode: "11680-670"
      },
      nodeInfo: {
        referenceID: "NODE-LAZARO-001"
      },
      serviceTime: 5,
      timeWindow: {
        start: "16:00:00",
        end: "18:00:00"
      },
      orders: [
        {
          referenceID: "PEDIDO-BR-011",
          contact: {
            fullName: "Marcos Vinícius da Silva",
            phone: "+55 12 89876-5432",
            email: "marcos.silva@yahoo.com.br",
            nationalID: "321.789.456-11"
          },
          instructions: "Sítio do Marcos - estrada de terra",
          deliveryUnits: [
            {
              lpn: "BR-LPN-011",
              weight: 3.0,
              volume: 0.018,
              price: 559.90,
              skills: ["pesado"],
              items: [
                {
                  sku: "MICRO-PANA-32L-BR",
                  description: "Micro-ondas Panasonic 32L Inox",
                  quantity: 1
                }
              ],
              evidences: []
            }
          ]
        }
      ]
    },
    {
      sequenceNumber: 10,
      type: "delivery",
      unassignedReason: "",
      addressInfo: {
        addressLine1: "Estrada da Fortaleza, 456",
        addressLine2: "Fortaleza",
        coordinates: {
          latitude: -23.5234,
          longitude: -45.1567
        },
        politicalArea: {
          adminAreaLevel1: "São Paulo",
          adminAreaLevel2: "Ubatuba",
          adminAreaLevel3: "Fortaleza",
          adminAreaLevel4: "",
          code: "BR-SP"
        },
        zipCode: "11680-930"
      },
      nodeInfo: {
        referenceID: "NODE-FORTALEZA-001"
      },
      serviceTime: 10,
      timeWindow: {
        start: "16:30:00",
        end: "18:00:00"
      },
      orders: [
        {
          referenceID: "PEDIDO-BR-012",
          contact: {
            fullName: "Camila Andrade Santos",
            phone: "+55 12 88765-4321",
            email: "camila.andrade@gmail.com",
            nationalID: "654.987.321-12"
          },
          instructions: "Casa da Camila - final da estrada",
          deliveryUnits: [
            {
              lpn: "BR-LPN-012",
              weight: 0.7,
              volume: 0.005,
              price: 249.90,
              skills: [],
              items: [
                {
                  sku: "POWER-BANK-20K-BR",
                  description: "Power Bank Xiaomi 20.000mAh USB-C",
                  quantity: 1
                }
              ],
              evidences: []
            }
          ]
        }
      ]
    }
  ]
}
