// Generar UUID aleatorio para la demo
const generateUUID = () => {
  return 'xxxxxxxx-xxxx-4xxx-yxxx-xxxxxxxxxxxx'.replace(/[xy]/g, function(c) {
    const r = Math.random() * 16 | 0;
    const v = c == 'x' ? r : (r & 0x3 | 0x8);
    return v.toString(16);
  });
};

export const mockRouteData: any ={
    "id":1,
    "referenceID": "6dd3adec-aa25-4147-9222-d3111b57a6f0",
    "createdAt": "2025-09-13T05:13:41Z",
    "planReferenceID": "f8aa394a-a6dd-4709-8e12-99357650b790",
    "vehicle": {
        "plate": "vehicle_1",
        "startLocation": {
            "addressInfo": {
                "coordinates": {
                    "latitude": -33.4505803,
                    "longitude": -70.7857318
                },
                "politicalArea": {}
            },
            "nodeInfo": {}
        },
        "endLocation": {
            "addressInfo": {
                "coordinates": {},
                "politicalArea": {}
            },
            "nodeInfo": {}
        },
        "timeWindow": {},
        "capacity": {
            "volume": 10000,
            "weight": 10000,
            "insurance": 10000
        }
    },
    "geometry": {
        "encoding": "polyline",
        "type": "linestring",
        "value": "bkdkE|c`oLy@sCsIjAyARM@O?KAMAuB[u@MwHoAOCUEYEMIOKIIIKIMCICKAMAS?CAk@x@}HLe@Je@L_@HWHUVo@l@oA~@uBJSNWNYDEFGDCFCh@_A|CmFNs@Vu@FKBCFKJOLMVYl@w@XS\\YTUTS^a@RStA_B`AoA~@kA~KeN`NuPtE{FnA}ANQ@CX_@nB_CZ_@\\]\\]ZWZW\\U^U^U`@S\\O\\O\\M^M`@M^KtBg@pA[j@M^If@Kl@Kn@Kp@KjAUhImBjDw@dDu@`Ck@lFqAzDw@|A]f@Mh@Ob@M^Mb@Qb@Q\\Md@Ud@S`@UzAw@pBiAfCoArH_EpE}BtDmBnCyAxCwA~GkDhCoA~E{BbCkAvAo@bGwCPKVKv@a@VM\\O^STK|Au@pEwBjB{@|Aq@vAo@fCiAJErAo@`@SpB{@fAi@nAm@h@Wl@[fB}@xAs@xHmDnAk@`Ae@~KeF`Ac@d@S^Q^Qb@S`@UXQXQf@]d@]z@q@~AoA|AoAnAeAjGkFxDeDb@a@d@a@VWZYXYNSPURYR[R_@R_@Pa@Ri@Pk@J_@J_@H]P}@N{@dBcKnAoHhAaHX}AZ_Bj@sCl@uCn@cDjBeJRaAFWH_@Ng@Ng@Rk@Tg@N[N[P[PYR[RYT[nAyAfBuBfBuBlByBlBuBpJoKrB}Bv@}@t@_At@aAt@aAt@eA~@qA|@sA~@wA`AwAXc@Xe@Ve@Ta@Pc@Rc@To@Ro@Rq@nAkElAmEp@aCd@iBb@aBxCqK|C}K|A}FnCkKlCcKRq@L_@H[HUHQFOFQFOFMHQHSJSHSNWR]LSHOHMPWPWPWd@k@j@u@lMgP`@k@`@m@`@k@^o@Xi@Xi@Vk@Vk@^aAzAiEFSZ}@dMa_@p@iBr@iB`@eAt@iBt@iBlGqOZy@Z{@Z{@Pk@Pk@Nm@XiAXkA`FuSdAeE^_BfAyEp@_DViAHe@Lm@h@}Ct@wEVcBZiBf@gCrGoY|@}D`@qBtA}GP}@N{@Ny@Lw@Ju@Ju@PwANyAb@eEhAyLbAsJrAwLBYBQ^qD\\{DPaBRaCf@yGh@gH\\gEdAkMXqDV_Eh@uI`@eGv@oKFeADeAFeAHcA|@mMJ}ALqBLsBVqDD_ABe@@e@@o@@o@AiAEcACe@Cg@M_Bu@{Ji@oHCc@Co@Cg@Cm@KwBIsAy@sKMiBi@eIIs@Ek@Is@Kw@Iy@Gs@Eu@KoAWqEC]CSWeECc@Ee@Ii@CSCUG]K]I_@K]M]K]Mi@EQIWO_@IUKUMSOUi@q@OQKICECEEGCGCICIAIAGAQAGCCACCCCAWAOAOAQCOCMEKEOIsBgAk@[s@_@CKEMACAGAMAK?MFUDUX}AJe@Lq@DWBMVyAX_BP_AF[X{A@GTkAFa@Ji@DQPcADU`@yBNw@?EDQBW@O?E@mC@w@?gABuC@gB?g@@s@?a@@u@?MW?kGEWAI?mAA?M@gC@oC@MgAAK?UAUAa@Nc@Ni@RcA^KB@E?E?G?A?a@@M?_B?w@@kC@uC@oD@_C?s@@yB?O?S@S?MAKAIEWKg@CSCQAKCe@D}D@KHyF?QHuFNEFCb@OvBw@wBv@c@NGBODItF?PIxFAJE|DAVAL?FAF?TMDuBp@gFhBEBs@T{Af@{Af@mHzBuBp@GBODQF}@Z}@XsAb@}Bt@EBSFQDgA^KBaGlBG@MFwC~@yBr@c@LoAb@c@PYHe@NeElAI?K?I?MAGAEAOGWGUQYU]WECKGGLc@z@CDGJbAx@DDLJLLf@d@TTFDLLFP?B@B?BBJ?D?D@J?J?Z?l@?????????j@?`@@~B?J@J@HBF?vACFAH?F?H?T?b@O?{ABC?CAC?CCQKCCCAC?E?C@OBG@E@?D?D?h@?NAlA?pA?D?J~A@H?lA?R@?K?eA?g@@K@K@I?K?O?[?U?I?w@?IAIAGCG?wABI@K?K?M?w@?[?K?EBEBEBCDAJCl@IZEB?LANAl@CDCDE@GBE?C@E@IYi@}@}A[i@_@g@CECEIKGGY[GGOOOO_@[UQYU]WECKGGEQKUMMEKGWQQKk@a@}@m@aEwCIEKGWSEEOKUQ_Ao@yByAgAu@k@]EEMIUOIIiCkBy@o@w@o@aDgC_Au@m@g@WSi@i@YWWUyAuAWOMGQK]UQKQKOGQIWIUG]ESIUIWGYG[GWE[E]CeBQ}BSgQ}AgDSi@E??"
    },
    "visits": [
        {
            "type": "delivery",
            "addressInfo": {
                "addressLine1": "Calle 1005, Piso 3, La Florida",
                "coordinates": {
                    "latitude": -33.5374662,
                    "longitude": -70.5989191
                },
                "politicalArea": {}
            },
            "nodeInfo": {
                "referenceID": "1315d1ba-25e1-4a80-85fd-e6be38564caf"
            },
            "sequenceNumber": 1,
            "timeWindow": {},
            "orders": [
                {
                    "referenceID": "105LA",
                    "contact": {
                        "fullName": "Fernando Castro"
                    },
                    "deliveryUnits": [
                        {
                            "lpn": "LPN-123456",
                            "items": [
                                {
                                    "description": "pepsi zero 350ml",
                                    "quantity": 12
                                },
                                {
                                    "description": "fanta zero 350ml",
                                    "quantity": 12
                                }
                            ],
                            "volume": 2,
                            "weight": 12,
                            "price": 100
                        },
                        {
                            "lpn": "LPN-1234567",
                            "items": [
                                {
                                    "description": "caja contenedora bebidas",
                                    "quantity": 1
                                }
                            ],
                            "volume": 2,
                            "weight": 12,
                            "price": 100
                        }
                    ]
                }
            ]
        },
        {
            "type": "delivery",
            "addressInfo": {
                "addressLine1": "Calle 1004, Piso 2, La Florida",
                "coordinates": {
                    "latitude": -33.536736,
                    "longitude": -70.5887186
                },
                "politicalArea": {}
            },
            "nodeInfo": {
                "referenceID": "0c97431f-2008-49ab-9fdf-6e6d7748ab9d"
            },
            "sequenceNumber": 2,
            "timeWindow": {},
            "orders": [
                {
                    "referenceID": "104LA",
                    "contact": {
                        "fullName": "Lucia Herrera"
                    },
                    "deliveryUnits": [
                        {
                            "items": [
                                {
                                    "description": "bebida 350ml",
                                    "quantity": 13
                                }
                            ],
                            "volume": 2,
                            "weight": 12,
                            "price": 100
                        }
                    ]
                }
            ]
        },
        {
            "type": "delivery",
            "addressInfo": {
                "addressLine1": "Calle 1007, La Florida",
                "coordinates": {
                    "latitude": -33.5226641,
                    "longitude": -70.5996466
                },
                "politicalArea": {}
            },
            "nodeInfo": {
                "referenceID": "f006430d-c0e2-4f25-8f6f-95f71a9cd24c"
            },
            "sequenceNumber": 3,
            "timeWindow": {},
            "orders": [
                {
                    "referenceID": "100LA-C",
                    "contact": {
                        "fullName": "Roberto Silva"
                    },
                    "deliveryUnits": [
                        {
                            "items": [
                                {
                                    "description": "Medicamentos",
                                    "quantity": 2
                                }
                            ],
                            "volume": 1,
                            "weight": 3,
                            "price": 75
                        }
                    ]
                },
                {
                    "referenceID": "101LA",
                    "contact": {
                        "fullName": "María Pérez"
                    },
                    "deliveryUnits": [
                        {
                            "items": [
                                {
                                    "description": "Libros",
                                    "quantity": 4
                                }
                            ],
                            "volume": 1,
                            "weight": 6,
                            "price": 90
                        }
                    ]
                },
                {
                    "referenceID": "100LA-A",
                    "contact": {
                        "fullName": "Roberto Silva"
                    },
                    "deliveryUnits": [
                        {
                            "lpn": "LPN-123456",
                            "items": [
                                {
                                    "description": "Bebida 350ml",
                                    "quantity": 13
                                }
                            ],
                            "volume": 2,
                            "weight": 12,
                            "price": 100
                        }
                    ]
                },
                {
                    "referenceID": "102LA",
                    "contact": {
                        "fullName": "Carlos Mendoza"
                    },
                    "deliveryUnits": [
                        {
                            "items": [
                                {
                                    "description": "Electrodomésticos",
                                    "quantity": 1
                                }
                            ],
                            "volume": 4,
                            "weight": 25,
                            "price": 350
                        }
                    ]
                },
                {
                    "referenceID": "100LA-B",
                    "contact": {
                        "fullName": "Roberto Silva"
                    },
                    "deliveryUnits": [
                        {
                            "items": [
                                {
                                    "description": "Comida enlatada",
                                    "quantity": 5
                                }
                            ],
                            "volume": 1,
                            "weight": 8,
                            "price": 150
                        }
                    ]
                }
            ]
        },
        {
            "type": "delivery",
            "addressInfo": {
                "addressLine1": "Calle 1003, Piso 1, La Florida",
                "coordinates": {
                    "latitude": -33.5068584,
                    "longitude": -70.5895279
                },
                "politicalArea": {}
            },
            "nodeInfo": {
                "referenceID": "3205e27e-f062-4da6-a51b-7d2e0a70c295"
            },
            "sequenceNumber": 4,
            "timeWindow": {},
            "orders": [
                {
                    "referenceID": "103LA",
                    "contact": {
                        "fullName": "Ana Rodriguez"
                    },
                    "deliveryUnits": [
                        {
                            "items": [
                                {
                                    "description": "bebida 350ml",
                                    "quantity": 13
                                }
                            ],
                            "volume": 2,
                            "weight": 12,
                            "price": 100
                        }
                    ]
                }
            ]
        }
    ]
}

export const mockDeliveryStates = {
  // Todas las entregas están pendientes (en ruta) por defecto
  // No hay entregas completadas inicialmente
}