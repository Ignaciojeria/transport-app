// Generar UUID aleatorio para la demo
const generateUUID = () => {
  return 'xxxxxxxx-xxxx-4xxx-yxxx-xxxxxxxxxxxx'.replace(/[xy]/g, function(c) {
    const r = Math.random() * 16 | 0;
    const v = c == 'x' ? r : (r & 0x3 | 0x8);
    return v.toString(16);
  });
};

export const mockRouteData: any ={
    "documentID": "",
    "referenceID": generateUUID(),
    "createdAt": "2025-09-02T04:37:00Z",
    "planReferenceID": "5bfc89e2-8d23-4b1b-8bb7-b8a1b7b28195",
    "vehicle": {
        "plate": "vehicle_A",
        "startLocation": {
            "addressInfo": {
                "contact": {},
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
                "contact": {},
                "coordinates": {},
                "politicalArea": {}
            },
            "nodeInfo": {}
        },
        "timeWindow": {},
        "capacity": {
            "volume": 1000,
            "weight": 1000,
            "insurance": 1000,
            "deliveryUnitsQuantity": 1000
        }
    },
    "geometry": {
        "encoding": "polyline",
        "type": "linestring",
        "value": "bkdkE|c`oLy@sCsIjAyARM@O?KAMAuB[u@MwHoAOCUEYEMIOKIIIKIMCICKAMAS?CAk@x@}HLe@Je@L_@HWHUVo@l@oA~@uBJSNWNYDEFGDCFCh@_A|CmFNs@Vu@FKBCFKJOLMVYl@w@XS\\YTUTS^a@RStA_B`AoA~@kA~KeN`NuPtE{FnA}ANQ@CX_@nB_CZ_@\\]\\]ZWZW\\U^U^U`@S\\O\\O\\M^M`@M^KtBg@pA[j@M^If@Kl@Kn@Kp@KjAUhImBjDw@dDu@`Ck@lFqAzDw@|A]f@Mh@Ob@M^Mb@Qb@Q\\Md@Ud@S`@UzAw@pBiAfCoArH_EpE}BtDmBnCyAxCwA~GkDhCoA~E{BbCkAvAo@bGwCPKVKv@a@VM\\O^STK|Au@pEwBjB{@|Aq@vAo@fCiAJErAo@`@SpB{@fAi@nAm@h@Wl@[fB}@xAs@xHmDnAk@`Ae@~KeF`Ac@d@S^Q^Qb@S`@UXQXQf@]d@]z@q@~AoA|AoAnAeAjGkFxDeDb@a@d@a@VWZYXYNSPURYR[R_@R_@Pa@Ri@Pk@J_@J_@H]P}@N{@dBcKnAoHhAaHX}AZ_Bj@sCl@uCn@cDjBeJRaAFWH_@Ng@Ng@Rk@Tg@N[N[P[PYR[RYT[nAyAfBuBfBuBlByBlBuBpJoKrB}Bv@}@t@_At@aAt@aAt@eA~@qA|@sA~@wA`AwAXc@Xe@Ve@Ta@Pc@Rc@To@Ro@Rq@nAkElAmEp@aCd@iBb@aBxCqK|C}K|A}FnCkKlCcKRq@L_@H[HUHQFOFQFOFMHQHSJSHSNWR]LSHOHMPWPWPWd@k@j@u@lMgP`@k@`@m@`@k@^o@Xi@Xi@Vk@Vk@^aAzAiEFSZ}@dMa_@p@iBr@iB`@eAt@iBt@iBlGqOZy@Z{@Z{@Pk@Pk@Nm@XiAXkA`FuSdAeE^_BfAyEp@_DViAHe@Lm@h@}Ct@wEVcBZiBf@gCrGoY|@}D`@qBtA}GP}@N{@Ny@Lw@Ju@Ju@PwANyAb@eEhAyLbAsJrAwLBYBQ^qD\\{DPaBRaCf@yGh@gH\\gEdAkMXqDV_Eh@uI`@eGv@oKFeADeAFeAHcA|@mMJ}ALqBLsBVqDD_ABe@@e@@o@@o@AiAEcACe@Cg@M_Bu@{Ji@oHCc@Co@Cg@Cm@KwBIsAy@sKMiBi@eIIs@Ek@Is@Kw@Iy@Gs@Eu@KoAWqEC]CSWeECc@Ee@Ii@CSCUG]K]I_@K]M]K]Ui@Sa@Wa@S_@U[QUWYUU[WWUYSa@We@Wm@_@yGoDwAu@yAy@_DgBgGkDkGoDaEaCeIsEeBaAqPoJi@YOIsAu@c@c@KMGEMKy@m@GGGGGEEIYQUOUQ]]USUWUYUYYa@Y_@OYQ[EGCGEKYg@Yi@}@}A[i@_@g@CECEIKGGY[GGOOOO_@[UQYU]WECKGGLc@z@CDGJbAx@DDLJLLf@d@TTFDLLFP?B@B?BBJ?D?D@J?J?Z?l@?j@?`@@~B?J@J@HBF?vACFAH?F?H?T?b@E?eBBC?CAC?CCQKCCCAC?E?C@OBG@E@?D?D?h@?NAlA?pA?D?J~A@H?fA?L?J@?K?I?{@?g@@K@K@I?K?O?[?U?I?w@?IAIAGCG?wABI@K?K?M?w@?[?K?EBEBEBCDAJCl@IZEB?LANAl@CDCDE@GBE?C@E@IYi@}@}A[i@_@g@CECEIKGGY[GGOOOO_@[UQYUGSAMAE?E@I@EDIBEFGHIJCPGDCDCBEBEBE@G@G@UA]?{@CyC?E?M?m@?K?ECiGA_DAkCAyDAkE?m@?k@?iBCeC?yBAcC?mC?o@AgA?mD?EAcD?y@?[?eA?kA?m@VILCMBWH?Q?U?y@RCf@GhBShBSHAbAKJCJCHEJEHGHIHIDKDIFODOVgADUDMBKFMFIHMHIHGJGHEHEHCVETEREZEv@MLC\\GRGHCJC^O`@O^M^K^Kv@UXIPELE^Kn@QfA[pBk@FCFCPENENCBABP@D`@hDF`@@L?L?JAHANAHAH?J?H@`@DzF@Z?J?J@d@B`DBfB@p@?h@?NBvB?h@BjC@jABbC?NBCDADCFAnAYrE}@NEFAFCDCFCh@a@FEHELGLCb@KFCBABABCBADC@A@A@?BAxAi@hJcD|CgANENE~B{@n@Ut@Wp@UnIyC\\MnAc@RGBBDBB?DABABA`Am@r@c@p@c@PKr@c@HGbBeAXSz@i@bAm@ROp@a@FEfAq@z@k@HGLG?PATFXAnC?b@CdE?R?XA|C?VApESFuFlB}@\\C@C@AB?BABA~BAhCANg@Af@@@O@iC@_C@C?C@CBABA|@]tFmBRG@qE?W@}C?Y?SBeE?c@@oCFY?M?S?O?Y?SE[@oC@mD?a@@aB@kF@U?sA?s@?{@@sADk@@k@?kB@cB?a@?q@?Y@gB?Q?IAKAIAI?U@u@?O?}F?y@?E?C?C?S\\KTGf@S\\MZMRKTKZQRI^O\\MZMh@Q\\K`@MVIXI^Id@K`@Il@OVIVIVKVKVKLGVMAs@AGAEAEAC?K?E?GJkK?O@Y@o@@mB@mBBmB@oB@oBBkB@oB@mBBqB@uB@YAXAtBCpBAlBAnBCjBAnBAnBClBAlBAlBAn@AX?NKjK?F?D?JQJE@WNm@ViAb@i@To@Zc@Re@TSHSHUHSHk@R_@LUFQFSDUFYFc@Je@Jo@NODODG@UHSHa@N]Lc@Re@Rs@\\[LSJUH[JgKhDq@TIBQFCLC\\CN?P?J?t@?B?`A?`@?X?J?l@?^?lC?~@AjA?@?nA?`@?V?`@?~B?L?r@SFiDlA_@LkA`@oKtDEBG@AQA[ScEAIAKS}EAMAIUyE?OKmB?KCUAM?CASAUA_@?SAS?OCcBAe@GaHAs@C_A?YA]AaA@S@Q?m@?g@A_BAM?MAMAKG{D@O?S?O?a@ASCYCeI?OCkB?q@Aa@CcCAgA?G?_@?K?EEsEAgA?YAoA@U?UBWBWFYFWHa@Ty@DQFOHUz@uBDOFQDODUBSHy@Bo@@O?Q?MAMCUCQ?MCW?CGe@OgBYaDCKIcAEk@SIQDC@SDQBo@BUBe@BYBU@M@OAG?C?C@CBABAD?HDz@@TDr@?FBRBJDHBBF@L?HA^G|Cg@FCDCFCFALCNzABFDBDBDBYaDCKIcAEk@IaA?EC_@CSGSAIQiBAiAEuAGoB?e@?e@@a@@a@IGGEICKEMCOCM?M?AX?V@VBZ@J@LGFCBABA@AD?BCt@S?wDQB{Aw@CcBG??"
    },
    "visits": [
        {
            "type": "delivery",
            "addressInfo": {
                "addressLine1": "Calle 1007, la-florida",
                "contact": {
                    "fullName": "Roberto Silva"
                },
                "coordinates": {
                    "latitude": -33.5226641,
                    "longitude": -70.5996466
                },
                "politicalArea": {}
            },
            "nodeInfo": {},
            "sequenceNumber": 1,
            "timeWindow": {},
            "orders": [
                {
                    "referenceID": "107LA-A",
                    "documentID": "",
                    "deliveryUnits": [
                        {
                            "documentID": "",
                            "lpn": "CODE-1A",
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
                },
                {
                    "referenceID": "107LA-B",
                    "documentID": "",
                    "deliveryUnits": [
                        {
                            "documentID": "",
                            "lpn": "CODE-1B",
                            "items": [
                                {
                                    "description": "comida enlatada",
                                    "quantity": 5
                                }
                            ],
                            "volume": 1.5,
                            "weight": 8,
                            "price": 150
                        }
                    ]
                },
                {
                    "referenceID": "107LA-C",
                    "documentID": "",
                    "deliveryUnits": [
                        {
                            "documentID": "",
                            "lpn": "CODE-1C",
                            "items": [
                                {
                                    "description": "medicamentos",
                                    "quantity": 2
                                }
                            ],
                            "volume": 0.5,
                            "weight": 3,
                            "price": 75
                        }
                    ]
                },
                {
                    "referenceID": "107LA-D",
                    "documentID": "",
                    "deliveryUnits": [
                        {
                            "documentID": "",
                            "lpn": "CODE-1D",
                            "items": [
                                {
                                    "description": "productos de limpieza",
                                    "quantity": 3
                                }
                            ],
                            "volume": 3,
                            "weight": 15,
                            "price": 200
                        },
                        {
                            "documentID": "",
                            "lpn": "CODE-1E",
                            "items": [
                                {
                                    "description": "ropa deportiva",
                                    "quantity": 2
                                }
                            ],
                            "volume": 1,
                            "weight": 5,
                            "price": 120
                        },
                        {
                            "documentID": "",
                            "lpn": "CODE-1F",
                            "items": [
                                {
                                    "description": "herramientas",
                                    "quantity": 1
                                }
                            ],
                            "volume": 2.5,
                            "weight": 8,
                            "price": 180
                        }
                    ]
                }
            ]
        },
        {
            "type": "delivery",
            "addressInfo": {
                "addressLine1": "Calle 1007, la-florida",
                "contact": {
                    "fullName": "María Pérez"
                },
                "coordinates": {
                    "latitude": -33.5226641,
                    "longitude": -70.5996466
                },
                "politicalArea": {}
            },
            "nodeInfo": {},
            "sequenceNumber": 2,
            "timeWindow": {},
            "orders": [
                {
                    "referenceID": "107LA-E",
                    "documentID": "",
                    "deliveryUnits": [
                        {
                            "documentID": "",
                            "lpn": "CODE-1G",
                            "items": [
                                {
                                    "description": "libros",
                                    "quantity": 4
                                }
                            ],
                            "volume": 1.2,
                            "weight": 6,
                            "price": 90
                        }
                    ]
                }
            ]
        },
        {
            "type": "delivery",
            "addressInfo": {
                "addressLine1": "Calle 1007, la-florida",
                "contact": {
                    "fullName": "Carlos Mendoza"
                },
                "coordinates": {
                    "latitude": -33.5226641,
                    "longitude": -70.5996466
                },
                "politicalArea": {}
            },
            "nodeInfo": {},
            "sequenceNumber": 3,
            "timeWindow": {},
            "orders": [
                {
                    "referenceID": "107LA-F",
                    "documentID": "",
                    "deliveryUnits": [
                        {
                            "documentID": "",
                            "lpn": "CODE-1H",
                            "items": [
                                {
                                    "description": "electrodomésticos",
                                    "quantity": 1
                                }
                            ],
                            "volume": 4,
                            "weight": 25,
                            "price": 350
                        }
                    ]
                }
            ]
        },
        {
            "type": "delivery",
            "addressInfo": {
                "addressLine1": "Calle 1016, Piso 17, la-florida",
                "contact": {
                    "fullName": "Roberto Silva"
                },
                "coordinates": {
                    "latitude": -33.5231166,
                    "longitude": -70.5830913
                },
                "politicalArea": {}
            },
            "nodeInfo": {},
            "sequenceNumber": 6,
            "timeWindow": {},
            "orders": [
                {
                    "referenceID": "116LA",
                    "documentID": "",
                    "deliveryUnits": [
                        {
                            "documentID": "",
                            "lpn": "CODE-2",
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
                "addressLine1": "Calle 1009, Piso 10, la-florida",
                "contact": {
                    "fullName": "Gabriela Torres"
                },
                "coordinates": {
                    "latitude": -33.5301395,
                    "longitude": -70.5828204
                },
                "politicalArea": {}
            },
            "nodeInfo": {},
            "sequenceNumber": 6,
            "timeWindow": {},
            "orders": [
                {
                    "referenceID": "109LA",
                    "documentID": "",
                    "deliveryUnits": [
                        {
                            "documentID": "",
                            "lpn": "CODE-3",
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
                "addressLine1": "Calle 1000, Piso 1, la-florida",
                "contact": {
                    "fullName": "Ignacio Jeria"
                },
                "coordinates": {
                    "latitude": -33.5304825,
                    "longitude": -70.5854977
                },
                "politicalArea": {}
            },
            "nodeInfo": {},
            "sequenceNumber": 6,
            "timeWindow": {},
            "orders": [
                {
                    "referenceID": "100LA",
                    "documentID": "",
                    "deliveryUnits": [
                        {
                            "documentID": "",
                            "lpn": "CODE-4",
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
                "addressLine1": "Calle 1006, Piso 7, la-florida",
                "contact": {
                    "fullName": "Roberto Silva"
                },
                "coordinates": {
                    "latitude": -33.5414441,
                    "longitude": -70.5872566
                },
                "politicalArea": {}
            },
            "nodeInfo": {},
            "sequenceNumber": 6,
            "timeWindow": {},
            "orders": [
                {
                    "referenceID": "106LA",
                    "documentID": "",
                    "deliveryUnits": [
                        {
                            "documentID": "",
                            "lpn": "CODE-5",
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
                "addressLine1": "Calle 1013, Piso 14, la-florida",
                "contact": {
                    "fullName": "Ana Rodriguez"
                },
                "coordinates": {
                    "latitude": -33.5466664,
                    "longitude": -70.5596647
                },
                "politicalArea": {}
            },
            "nodeInfo": {},
            "sequenceNumber": 6,
            "timeWindow": {},
            "orders": [
                {
                    "referenceID": "113LA",
                    "documentID": "",
                    "deliveryUnits": [
                        {
                            "documentID": "",
                            "lpn": "CODE-6",
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
                "addressLine1": "Calle 1015, Piso 16, la-florida",
                "contact": {
                    "fullName": "Lucia Herrera"
                },
                "coordinates": {
                    "latitude": -33.5384557,
                    "longitude": -70.5767166
                },
                "politicalArea": {}
            },
            "nodeInfo": {},
            "sequenceNumber": 7,
            "timeWindow": {},
            "orders": [
                {
                    "referenceID": "115LA",
                    "documentID": "",
                    "deliveryUnits": [
                        {
                            "documentID": "",
                            "lpn": "CODE-7",
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
                "addressLine1": "Calle 1018, Piso 19, la-florida",
                "contact": {
                    "fullName": "Fernando Castro"
                },
                "coordinates": {
                    "latitude": -33.5341326,
                    "longitude": -70.5560076
                },
                "politicalArea": {}
            },
            "nodeInfo": {},
            "sequenceNumber": 8,
            "timeWindow": {},
            "orders": [
                {
                    "referenceID": "118LA",
                    "documentID": "",
                    "deliveryUnits": [
                        {
                            "documentID": "",
                            "lpn": "CODE-8",
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
                "addressLine1": "Calle 1002, Piso 3, la-florida",
                "contact": {
                    "fullName": "Maria Perez"
                },
                "coordinates": {
                    "latitude": -33.5332085,
                    "longitude": -70.5516135
                },
                "politicalArea": {}
            },
            "nodeInfo": {},
            "sequenceNumber": 9,
            "timeWindow": {},
            "orders": [
                {
                    "referenceID": "102LA",
                    "documentID": "",
                    "deliveryUnits": [
                        {
                            "documentID": "",
                            "lpn": "CODE-9",
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