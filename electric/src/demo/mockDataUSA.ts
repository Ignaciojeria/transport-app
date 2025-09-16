import type { Route } from '../domain/route'

export const mockRouteUSA: Route = {
  id: 1,
  referenceID: "ROUTE-MIA-001",
  planReferenceID: "PLAN-MIAMI-001",
  createdAt: "2025-09-16T10:00:00Z",
  geometry: {
    type: "linestring",
    encoding: "polyline",
    value: "_p~iF~ps|U_ulLnnqC_mqNvxq@"
  },
  vehicle: {
    plate: "FL7A2BC",
    capacity: {
      weight: 1000,
      volume: 1000,
      insurance: 20000
    },
    skills: ["temperature_controlled", "fragile"],
    startLocation: {
      addressInfo: {
        addressLine1: "1200 Biscayne Blvd",
        addressLine2: "Downtown",
        coordinates: {
          latitude: 25.7617,
          longitude: -80.1918
        },
        politicalArea: {
          adminAreaLevel1: "Florida",
          adminAreaLevel2: "Miami-Dade County",
          adminAreaLevel3: "Miami",
          adminAreaLevel4: "",
          code: "US-FL"
        },
        zipCode: "33132"
      },
      nodeInfo: {
        referenceID: "NODE-DOWNTOWN-MIA"
      }
    },
    endLocation: {
      addressInfo: {
        addressLine1: "1200 Biscayne Blvd",
        addressLine2: "Downtown",
        coordinates: {
          latitude: 25.7617,
          longitude: -80.1918
        },
        politicalArea: {
          adminAreaLevel1: "Florida",
          adminAreaLevel2: "Miami-Dade County",
          adminAreaLevel3: "Miami",
          adminAreaLevel4: "",
          code: "US-FL"
        },
        zipCode: "33132"
      },
      nodeInfo: {
        referenceID: "NODE-DOWNTOWN-MIA"
      }
    },
    timeWindow: {
      start: "08:00:00",
      end: "18:00:00"
    }
  },
  visits: [
    // Visit 1: Multiple customers at the same address (South Beach)
    {
      sequenceNumber: 1,
      type: "delivery",
      unassignedReason: "",
      addressInfo: {
        addressLine1: "450 Lincoln Rd",
        addressLine2: "South Beach",
        coordinates: {
          latitude: 25.7907,
          longitude: -80.1349
        },
        politicalArea: {
          adminAreaLevel1: "Florida",
          adminAreaLevel2: "Miami-Dade County",
          adminAreaLevel3: "Miami Beach",
          adminAreaLevel4: "",
          code: "US-FL"
        },
        zipCode: "33139"
      },
      nodeInfo: {
        referenceID: "NODE-SOUTH-BEACH-001"
      },
      serviceTime: 10,
      timeWindow: {
        start: "09:00:00",
        end: "12:00:00"
      },
      orders: [
        // Customer 1: Michael Johnson (3 orders)
        {
          referenceID: "ORDER-US-001A",
          contact: {
            fullName: "Michael Johnson",
            phone: "+1 305 555-0123",
            email: "m.johnson@email.com",
            nationalID: "123-45-6789"
          },
          instructions: "Deliver to front desk",
          deliveryUnits: [
            {
              lpn: "US-LPN-001A",
              weight: 0.5,
              volume: 0.01,
              price: 1199.00,
              skills: ["fragile"],
              items: [
                {
                  sku: "APPL-IPH15-256-US",
                  description: "Apple iPhone 15 Pro 256GB",
                  quantity: 1
                }
              ],
              evidences: []
            }
          ]
        },
        {
          referenceID: "ORDER-US-001B",
          contact: {
            fullName: "Michael Johnson",
            phone: "+1 305 555-0123",
            email: "m.johnson@email.com",
            nationalID: "123-45-6789"
          },
          instructions: "Second order - same customer",
          deliveryUnits: [
            {
              lpn: "US-LPN-001B",
              weight: 0.2,
              volume: 0.005,
              price: 129.00,
              skills: ["fragile"],
              items: [
                {
                  sku: "CASE-APPL-LEATHER-US",
                  description: "Apple Leather Case with MagSafe",
                  quantity: 1
                },
                {
                  sku: "SCRN-PROT-APPL-US",
                  description: "Tempered Glass Screen Protector",
                  quantity: 1
                }
              ],
              evidences: []
            }
          ]
        },
        {
          referenceID: "ORDER-US-001C",
          contact: {
            fullName: "Michael Johnson",
            phone: "+1 305 555-0123",
            email: "m.johnson@email.com",
            nationalID: "123-45-6789"
          },
          instructions: "Third order - AirPods",
          deliveryUnits: [
            {
              lpn: "US-LPN-001C",
              weight: 0.1,
              volume: 0.002,
              price: 249.00,
              skills: [],
              items: [
                {
                  sku: "APPL-AIRPODS-PRO2-US",
                  description: "Apple AirPods Pro 2nd Generation",
                  quantity: 1
                }
              ],
              evidences: []
            }
          ]
        },
        // Customer 2: Sarah Williams
        {
          referenceID: "ORDER-US-002",
          contact: {
            fullName: "Sarah Williams",
            phone: "+1 305 555-0456",
            email: "sarah.williams@gmail.com",
            nationalID: "987-65-4321"
          },
          instructions: "Call before delivery",
          deliveryUnits: [
            {
              lpn: "US-LPN-002",
              weight: 1.5,
              volume: 0.008,
              price: 899.00,
              skills: ["fragile"],
              items: [
                {
                  sku: "APPL-IPAD-AIR-US",
                  description: "iPad Air 256GB Wi-Fi + Cellular",
                  quantity: 1
                }
              ],
              evidences: []
            }
          ]
        },
        // Customer 3: David Brown
        {
          referenceID: "ORDER-US-003",
          contact: {
            fullName: "David Brown",
            phone: "+1 305 555-0789",
            email: "david.brown@hotmail.com",
            nationalID: "456-78-9123"
          },
          instructions: "Concierge accepts deliveries",
          deliveryUnits: [
            {
              lpn: "US-LPN-003",
              weight: 4.2,
              volume: 0.015,
              price: 1999.00,
              skills: ["fragile"],
              items: [
                {
                  sku: "APPL-MBP-14-M3-US",
                  description: "MacBook Pro 14-inch M3 Pro 512GB",
                  quantity: 1
                }
              ],
              evidences: []
            }
          ]
        }
      ]
    },
    // Visits 2-10: Single customers at different Miami locations
    {
      sequenceNumber: 2,
      type: "delivery",
      unassignedReason: "",
      addressInfo: {
        addressLine1: "3401 NE 1st Ave",
        addressLine2: "Design District",
        coordinates: {
          latitude: 25.8067,
          longitude: -80.1918
        },
        politicalArea: {
          adminAreaLevel1: "Florida",
          adminAreaLevel2: "Miami-Dade County",
          adminAreaLevel3: "Miami",
          adminAreaLevel4: "",
          code: "US-FL"
        },
        zipCode: "33137"
      },
      nodeInfo: {
        referenceID: "NODE-DESIGN-DISTRICT-001"
      },
      serviceTime: 8,
      timeWindow: {
        start: "10:00:00",
        end: "14:00:00"
      },
      orders: [
        {
          referenceID: "ORDER-US-004",
          contact: {
            fullName: "Jennifer Martinez",
            phone: "+1 305 555-1234",
            email: "jennifer.martinez@outlook.com",
            nationalID: "789-12-3456"
          },
          instructions: "Blue building with glass facade",
          deliveryUnits: [
            {
              lpn: "US-LPN-004",
              weight: 0.3,
              volume: 0.003,
              price: 399.00,
              skills: [],
              items: [
                {
                  sku: "APPL-WATCH-S9-US",
                  description: "Apple Watch Series 9 GPS 45mm",
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
        addressLine1: "1111 SW 1st Ave",
        addressLine2: "Brickell",
        coordinates: {
          latitude: 25.7617,
          longitude: -80.1918
        },
        politicalArea: {
          adminAreaLevel1: "Florida",
          adminAreaLevel2: "Miami-Dade County",
          adminAreaLevel3: "Miami",
          adminAreaLevel4: "",
          code: "US-FL"
        },
        zipCode: "33130"
      },
      nodeInfo: {
        referenceID: "NODE-BRICKELL-001"
      },
      serviceTime: 6,
      timeWindow: {
        start: "11:00:00",
        end: "15:00:00"
      },
      orders: [
        {
          referenceID: "ORDER-US-005",
          contact: {
            fullName: "Robert Davis",
            phone: "+1 305 555-2345",
            email: "robert.davis@gmail.com",
            nationalID: "321-65-4987"
          },
          instructions: "Apt 1502 - Ring Robert",
          deliveryUnits: [
            {
              lpn: "US-LPN-005",
              weight: 8.5,
              volume: 0.02,
              price: 1299.00,
              skills: ["heavy"],
              items: [
                {
                  sku: "SONY-TV-55-4K-US",
                  description: "Sony 55' 4K OLED Smart TV",
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
        addressLine1: "2700 NW 7th St",
        addressLine2: "Little Havana",
        coordinates: {
          latitude: 25.7742,
          longitude: -80.2345
        },
        politicalArea: {
          adminAreaLevel1: "Florida",
          adminAreaLevel2: "Miami-Dade County",
          adminAreaLevel3: "Miami",
          adminAreaLevel4: "",
          code: "US-FL"
        },
        zipCode: "33135"
      },
      nodeInfo: {
        referenceID: "NODE-LITTLE-HAVANA-001"
      },
      serviceTime: 5,
      timeWindow: {
        start: "09:30:00",
        end: "13:30:00"
      },
      orders: [
        {
          referenceID: "ORDER-US-006",
          contact: {
            fullName: "Maria Garcia",
            phone: "+1 305 555-3456",
            email: "maria.garcia@yahoo.com",
            nationalID: "654-32-1098"
          },
          instructions: "Only deliver to Maria",
          deliveryUnits: [
            {
              lpn: "US-LPN-006",
              weight: 3.2,
              volume: 0.012,
              price: 599.00,
              skills: ["fragile"],
              items: [
                {
                  sku: "NINJ-COFFEE-MAKER-US",
                  description: "Ninja Coffee Maker with Frother",
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
        addressLine1: "19501 Biscayne Blvd",
        addressLine2: "Aventura",
        coordinates: {
          latitude: 25.9565,
          longitude: -80.1389
        },
        politicalArea: {
          adminAreaLevel1: "Florida",
          adminAreaLevel2: "Miami-Dade County",
          adminAreaLevel3: "Aventura",
          adminAreaLevel4: "",
          code: "US-FL"
        },
        zipCode: "33180"
      },
      nodeInfo: {
        referenceID: "NODE-AVENTURA-001"
      },
      serviceTime: 7,
      timeWindow: {
        start: "12:00:00",
        end: "16:00:00"
      },
      orders: [
        {
          referenceID: "ORDER-US-007",
          contact: {
            fullName: "Christopher Lee",
            phone: "+1 305 555-4567",
            email: "christopher.lee@gmail.com",
            nationalID: "987-32-1654"
          },
          instructions: "Beachfront condo - gate code 1234",
          deliveryUnits: [
            {
              lpn: "US-LPN-007",
              weight: 1.1,
              volume: 0.006,
              price: 449.00,
              skills: [],
              items: [
                {
                  sku: "JBL-XTREME3-BT-US",
                  description: "JBL Xtreme 3 Bluetooth Speaker",
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
        addressLine1: "9700 Collins Ave",
        addressLine2: "Bal Harbour",
        coordinates: {
          latitude: 25.8919,
          longitude: -80.1234
        },
        politicalArea: {
          adminAreaLevel1: "Florida",
          adminAreaLevel2: "Miami-Dade County",
          adminAreaLevel3: "Bal Harbour",
          adminAreaLevel4: "",
          code: "US-FL"
        },
        zipCode: "33154"
      },
      nodeInfo: {
        referenceID: "NODE-BAL-HARBOUR-001"
      },
      serviceTime: 6,
      timeWindow: {
        start: "13:00:00",
        end: "17:00:00"
      },
      orders: [
        {
          referenceID: "ORDER-US-008",
          contact: {
            fullName: "Amanda Wilson",
            phone: "+1 305 555-5678",
            email: "amanda.wilson@gmail.com",
            nationalID: "123-98-7456"
          },
          instructions: "Luxury condo - Valet parking",
          deliveryUnits: [
            {
              lpn: "US-LPN-008",
              weight: 2.8,
              volume: 0.01,
              price: 799.00,
              skills: ["fragile"],
              items: [
                {
                  sku: "DYSON-AIRWRAP-US",
                  description: "Dyson Airwrap Multi-Styler Complete",
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
        addressLine1: "1 Key Biscayne Blvd",
        addressLine2: "Key Biscayne",
        coordinates: {
          latitude: 25.6912,
          longitude: -80.1567
        },
        politicalArea: {
          adminAreaLevel1: "Florida",
          adminAreaLevel2: "Miami-Dade County",
          adminAreaLevel3: "Key Biscayne",
          adminAreaLevel4: "",
          code: "US-FL"
        },
        zipCode: "33149"
      },
      nodeInfo: {
        referenceID: "NODE-KEY-BISCAYNE-001"
      },
      serviceTime: 8,
      timeWindow: {
        start: "14:00:00",
        end: "18:00:00"
      },
      orders: [
        {
          referenceID: "ORDER-US-009",
          contact: {
            fullName: "James Thompson",
            phone: "+1 305 555-6789",
            email: "james.thompson@hotmail.com",
            nationalID: "456-12-3789"
          },
          instructions: "Island home with security gate",
          deliveryUnits: [
            {
              lpn: "US-LPN-009",
              weight: 0.2,
              volume: 0.004,
              price: 299.00,
              skills: [],
              items: [
                {
                  sku: "LOGI-MX-MASTER3-US",
                  description: "Logitech MX Master 3 Wireless Mouse",
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
        addressLine1: "15701 Biscayne Blvd",
        addressLine2: "North Miami Beach",
        coordinates: {
          latitude: 25.9234,
          longitude: -80.1456
        },
        politicalArea: {
          adminAreaLevel1: "Florida",
          adminAreaLevel2: "Miami-Dade County",
          adminAreaLevel3: "North Miami Beach",
          adminAreaLevel4: "",
          code: "US-FL"
        },
        zipCode: "33160"
      },
      nodeInfo: {
        referenceID: "NODE-N-MIAMI-BEACH-001"
      },
      serviceTime: 9,
      timeWindow: {
        start: "15:00:00",
        end: "18:00:00"
      },
      orders: [
        {
          referenceID: "ORDER-US-010",
          contact: {
            fullName: "Lisa Anderson",
            phone: "+1 305 555-7890",
            email: "lisa.anderson@gmail.com",
            nationalID: "789-45-6123"
          },
          instructions: "Ocean View Resort - Front Desk",
          deliveryUnits: [
            {
              lpn: "US-LPN-010",
              weight: 0.8,
              volume: 0.009,
              price: 549.00,
              skills: ["fragile"],
              items: [
                {
                  sku: "KINDLE-OASIS-32GB-US",
                  description: "Kindle Oasis 32GB with Cellular",
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
        addressLine1: "8888 SW 136th St",
        addressLine2: "Pinecrest",
        coordinates: {
          latitude: 25.6567,
          longitude: -80.3123
        },
        politicalArea: {
          adminAreaLevel1: "Florida",
          adminAreaLevel2: "Miami-Dade County",
          adminAreaLevel3: "Pinecrest",
          adminAreaLevel4: "",
          code: "US-FL"
        },
        zipCode: "33156"
      },
      nodeInfo: {
        referenceID: "NODE-PINECREST-001"
      },
      serviceTime: 5,
      timeWindow: {
        start: "16:00:00",
        end: "18:00:00"
      },
      orders: [
        {
          referenceID: "ORDER-US-011",
          contact: {
            fullName: "Kevin Rodriguez",
            phone: "+1 305 555-8901",
            email: "kevin.rodriguez@yahoo.com",
            nationalID: "321-78-9456"
          },
          instructions: "Suburban home - ring doorbell",
          deliveryUnits: [
            {
              lpn: "US-LPN-011",
              weight: 4.5,
              volume: 0.018,
              price: 899.00,
              skills: ["heavy"],
              items: [
                {
                  sku: "PANA-MICRO-1200W-US",
                  description: "Panasonic Microwave 1200W Stainless",
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
        addressLine1: "12000 SW 88th St",
        addressLine2: "Kendall",
        coordinates: {
          latitude: 25.6789,
          longitude: -80.3567
        },
        politicalArea: {
          adminAreaLevel1: "Florida",
          adminAreaLevel2: "Miami-Dade County",
          adminAreaLevel3: "Kendall",
          adminAreaLevel4: "",
          code: "US-FL"
        },
        zipCode: "33176"
      },
      nodeInfo: {
        referenceID: "NODE-KENDALL-001"
      },
      serviceTime: 10,
      timeWindow: {
        start: "16:30:00",
        end: "18:00:00"
      },
      orders: [
        {
          referenceID: "ORDER-US-012",
          contact: {
            fullName: "Michelle Taylor",
            phone: "+1 305 555-9012",
            email: "michelle.taylor@gmail.com",
            nationalID: "654-98-7321"
          },
          instructions: "End of residential street",
          deliveryUnits: [
            {
              lpn: "US-LPN-012",
              weight: 0.6,
              volume: 0.005,
              price: 349.00,
              skills: [],
              items: [
                {
                  sku: "ANKER-POWERBANK-20K-US",
                  description: "Anker PowerCore 20,000mAh USB-C",
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
