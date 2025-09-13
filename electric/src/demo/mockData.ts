// Generar UUID aleatorio para la demo
const generateUUID = () => {
  return 'xxxxxxxx-xxxx-4xxx-yxxx-xxxxxxxxxxxx'.replace(/[xy]/g, function(c) {
    const r = Math.random() * 16 | 0;
    const v = c == 'x' ? r : (r & 0x3 | 0x8);
    return v.toString(16);
  });
};

export const mockRouteData: any ={
    "id": 1,
    "referenceID": generateUUID(),
    "createdAt": "2025-09-02T04:37:00Z",
    "planReferenceID": "5bfc89e2-8d23-4b1b-8bb7-b8a1b7b28195",
    "vehicle": {
        "plate": "vehicle_A",
        "startLocation": {
            "addressInfo": {
                "addressLine1": "Punto de inicio",
                "addressLine2": "",
                "coordinates": {
                    "latitude": -33.4505803,
                    "longitude": -70.7857318
                },
                "politicalArea": {
                    "adminAreaLevel1": "",
                    "adminAreaLevel2": "",
                    "adminAreaLevel3": "",
                    "adminAreaLevel4": "",
                    "code": ""
                },
                "zipCode": ""
            },
            "nodeInfo": {
                "referenceID": "NODE-START"
            }
        },
        "endLocation": {
            "addressInfo": {
                "addressLine1": "Punto de fin",
                "addressLine2": "",
                "coordinates": {
                    "latitude": -33.4505803,
                    "longitude": -70.7857318
                },
                "politicalArea": {
                    "adminAreaLevel1": "",
                    "adminAreaLevel2": "",
                    "adminAreaLevel3": "",
                    "adminAreaLevel4": "",
                    "code": ""
                },
                "zipCode": ""
            },
            "nodeInfo": {
                "referenceID": "NODE-END"
            }
        },
        "timeWindow": {
            "start": "08:00",
            "end": "18:00"
        },
        "capacity": {
            "volume": 1000,
            "weight": 1000,
            "insurance": 1000
        },
        "skills": []
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
                "addressLine2": "",
                "coordinates": {
                    "latitude": -33.5226641,
                    "longitude": -70.5996466
                },
                "politicalArea": {
                    "adminAreaLevel1": "",
                    "adminAreaLevel2": "",
                    "adminAreaLevel3": "",
                    "adminAreaLevel4": "",
                    "code": ""
                },
                "zipCode": ""
            },
            "nodeInfo": {
                "referenceID": "NODE-001"
            },
            "sequenceNumber": 1,
            "serviceTime": 300,
            "timeWindow": {
                "start": "09:00",
                "end": "17:00"
            },
            "unassignedReason": "",
            "orders": [
                {
                    "contact": {
                        "email": "roberto.silva@email.com",
                        "fullName": "Roberto Silva",
                        "nationalID": "12345678-9",
                        "phone": "+56912345678"
                    },
                    "referenceID": "107LA-A",
                    "instructions": "Entregar en recepción del edificio",
                    "deliveryUnits": [
                        {
                            "lpn": "CODE-1A",
                            "items": [
                                {
                                    "description": "bebida 350ml",
                                    "quantity": 13,
                                    "sku": "BEB-350-001"
                                }
                            ],
                            "volume": 2,
                            "weight": 12,
                            "price": 100,
                            "skills": [],
                            "evidences": []
                        }
                    ]
                },
                {
                    "contact": {
                        "email": "roberto.silva@email.com",
                        "fullName": "Roberto Silva",
                        "nationalID": "12345678-9",
                        "phone": "+56912345678"
                    },
                    "referenceID": "107LA-B",
                    "instructions": "Entregar en recepción del edificio",
                    "deliveryUnits": [
                        {
                            "lpn": "CODE-1B",
                            "items": [
                                {
                                    "description": "comida enlatada",
                                    "quantity": 5,
                                    "sku": "COM-ENL-001"
                                }
                            ],
                            "volume": 1.5,
                            "weight": 8,
                            "price": 150,
                            "skills": [],
                            "evidences": []
                        }
                    ]
                },
                {
                    "contact": {
                        "email": "roberto.silva@email.com",
                        "fullName": "Roberto Silva",
                        "nationalID": "12345678-9",
                        "phone": "+56912345678"
                    },
                    "referenceID": "107LA-C",
                    "instructions": "Entregar en recepción del edificio",
                    "deliveryUnits": [
                        {
                            "lpn": "CODE-1C",
                            "items": [
                                {
                                    "description": "medicamentos",
                                    "quantity": 2,
                                    "sku": "MED-001"
                                }
                            ],
                            "volume": 0.5,
                            "weight": 3,
                            "price": 75,
                            "skills": [],
                            "evidences": []
                        }
                    ]
                },
                {
                    "contact": {
                        "email": "roberto.silva@email.com",
                        "fullName": "Roberto Silva",
                        "nationalID": "12345678-9",
                        "phone": "+56912345678"
                    },
                    "referenceID": "107LA-D",
                    "instructions": "Entregar en recepción del edificio",
                    "deliveryUnits": [
                        {
                            "lpn": "CODE-1D",
                            "items": [
                                {
                                    "description": "productos de limpieza",
                                    "quantity": 3,
                                    "sku": "LIMP-001"
                                }
                            ],
                            "volume": 3,
                            "weight": 15,
                            "price": 200,
                            "skills": [],
                            "evidences": []
                        },
                        {
                            "lpn": "CODE-1E",
                            "items": [
                                {
                                    "description": "ropa deportiva",
                                    "quantity": 2,
                                    "sku": "ROPA-DEP-001"
                                }
                            ],
                            "volume": 1,
                            "weight": 5,
                            "price": 120,
                            "skills": [],
                            "evidences": []
                        },
                        {
                            "lpn": "CODE-1F",
                            "items": [
                                {
                                    "description": "herramientas",
                                    "quantity": 1,
                                    "sku": "HERR-001"
                                }
                            ],
                            "volume": 2.5,
                            "weight": 8,
                            "price": 180,
                            "skills": [],
                            "evidences": []
                        }
                    ]
                }
            ]
        },
        {
            "type": "delivery",
            "addressInfo": {
                "addressLine1": "Calle 1007, la-florida",
                "addressLine2": "",
                "coordinates": {
                    "latitude": -33.5226641,
                    "longitude": -70.5996466
                },
                "politicalArea": {
                    "adminAreaLevel1": "",
                    "adminAreaLevel2": "",
                    "adminAreaLevel3": "",
                    "adminAreaLevel4": "",
                    "code": ""
                },
                "zipCode": ""
            },
            "nodeInfo": {
                "referenceID": "NODE-002"
            },
            "sequenceNumber": 2,
            "serviceTime": 300,
            "timeWindow": {
                "start": "09:00",
                "end": "17:00"
            },
            "unassignedReason": "",
            "orders": [
                {
                    "contact": {
                        "email": "maria.perez@email.com",
                        "fullName": "María Pérez",
                        "nationalID": "98765432-1",
                        "phone": "+56987654321"
                    },
                    "referenceID": "107LA-E",
                    "instructions": "Tocar timbre del apartamento 2B",
                    "deliveryUnits": [
                        {
                            "lpn": "CODE-1G",
                            "items": [
                                {
                                    "description": "libros",
                                    "quantity": 4,
                                    "sku": "LIB-001"
                                }
                            ],
                            "volume": 1.2,
                            "weight": 6,
                            "price": 90,
                            "skills": [],
                            "evidences": []
                        }
                    ]
                }
            ]
        },
        {
            "type": "delivery",
            "addressInfo": {
                "addressLine1": "Calle 1007, la-florida",
                "addressLine2": "",
                "coordinates": {
                    "latitude": -33.5226641,
                    "longitude": -70.5996466
                },
                "politicalArea": {
                    "adminAreaLevel1": "",
                    "adminAreaLevel2": "",
                    "adminAreaLevel3": "",
                    "adminAreaLevel4": "",
                    "code": ""
                },
                "zipCode": ""
            },
            "nodeInfo": {
                "referenceID": "NODE-003"
            },
            "sequenceNumber": 3,
            "serviceTime": 300,
            "timeWindow": {
                "start": "09:00",
                "end": "17:00"
            },
            "unassignedReason": "",
            "orders": [
                {
                    "contact": {
                        "email": "carlos.mendoza@email.com",
                        "fullName": "Carlos Mendoza",
                        "nationalID": "11223344-5",
                        "phone": "+56911223344"
                    },
                    "referenceID": "107LA-F",
                    "instructions": "Llamar antes de llegar",
                    "deliveryUnits": [
                        {
                            "lpn": "CODE-1H",
                            "items": [
                                {
                                    "description": "electrodomésticos",
                                    "quantity": 1,
                                    "sku": "ELEC-001"
                                }
                            ],
                            "volume": 4,
                            "weight": 25,
                            "price": 350,
                            "skills": [],
                            "evidences": []
                        }
                    ]
                }
            ]
        },
        {
            "type": "delivery",
            "addressInfo": {
                "addressLine1": "Calle 1016, Piso 17, la-florida",
                "addressLine2": "",
                "coordinates": {
                    "latitude": -33.5231166,
                    "longitude": -70.5830913
                },
                "politicalArea": {
                    "adminAreaLevel1": "",
                    "adminAreaLevel2": "",
                    "adminAreaLevel3": "",
                    "adminAreaLevel4": "",
                    "code": ""
                },
                "zipCode": ""
            },
            "nodeInfo": {
                "referenceID": "NODE-004"
            },
            "sequenceNumber": 4,
            "serviceTime": 300,
            "timeWindow": {
                "start": "09:00",
                "end": "17:00"
            },
            "unassignedReason": "",
            "orders": [
                {
                    "contact": {
                        "email": "roberto.silva2@email.com",
                        "fullName": "Roberto Silva",
                        "nationalID": "12345678-9",
                        "phone": "+56912345678"
                    },
                    "referenceID": "116LA",
                    "instructions": "Entregar en piso 17, apartamento A",
                    "deliveryUnits": [
                        {
                            "lpn": "CODE-2",
                            "items": [
                                {
                                    "description": "bebida 350ml",
                                    "quantity": 13,
                                    "sku": "BEB-350-001"
                                }
                            ],
                            "volume": 2,
                            "weight": 12,
                            "price": 100,
                            "skills": [],
                            "evidences": []
                        }
                    ]
                }
            ]
        },
        {
            "type": "delivery",
            "addressInfo": {
                "addressLine1": "Calle 1009, Piso 10, la-florida",
                "addressLine2": "",
                "coordinates": {
                    "latitude": -33.5301395,
                    "longitude": -70.5828204
                },
                "politicalArea": {
                    "adminAreaLevel1": "",
                    "adminAreaLevel2": "",
                    "adminAreaLevel3": "",
                    "adminAreaLevel4": "",
                    "code": ""
                },
                "zipCode": ""
            },
            "nodeInfo": {
                "referenceID": "NODE-005"
            },
            "sequenceNumber": 5,
            "serviceTime": 300,
            "timeWindow": {
                "start": "09:00",
                "end": "17:00"
            },
            "unassignedReason": "",
            "orders": [
                {
                    "contact": {
                        "email": "gabriela.torres@email.com",
                        "fullName": "Gabriela Torres",
                        "nationalID": "55667788-9",
                        "phone": "+56955667788"
                    },
                    "referenceID": "109LA",
                    "instructions": "Entregar en piso 10, apartamento B",
                    "deliveryUnits": [
                        {
                            "lpn": "CODE-3",
                            "items": [
                                {
                                    "description": "bebida 350ml",
                                    "quantity": 13,
                                    "sku": "BEB-350-001"
                                }
                            ],
                            "volume": 2,
                            "weight": 12,
                            "price": 100,
                            "skills": [],
                            "evidences": []
                        }
                    ]
                }
            ]
        },
        {
            "type": "delivery",
            "addressInfo": {
                "addressLine1": "Calle 1000, Piso 1, la-florida",
                "addressLine2": "",
                "coordinates": {
                    "latitude": -33.5304825,
                    "longitude": -70.5854977
                },
                "politicalArea": {
                    "adminAreaLevel1": "",
                    "adminAreaLevel2": "",
                    "adminAreaLevel3": "",
                    "adminAreaLevel4": "",
                    "code": ""
                },
                "zipCode": ""
            },
            "nodeInfo": {
                "referenceID": "NODE-006"
            },
            "sequenceNumber": 6,
            "serviceTime": 300,
            "timeWindow": {
                "start": "09:00",
                "end": "17:00"
            },
            "unassignedReason": "",
            "orders": [
                {
                    "contact": {
                        "email": "ignacio.jeria@email.com",
                        "fullName": "Ignacio Jeria",
                        "nationalID": "99887766-5",
                        "phone": "+56999887766"
                    },
                    "referenceID": "100LA",
                    "instructions": "Entregar en piso 1, apartamento C",
                    "deliveryUnits": [
                        {
                            "lpn": "CODE-4",
                            "items": [
                                {
                                    "description": "bebida 350ml",
                                    "quantity": 13,
                                    "sku": "BEB-350-001"
                                }
                            ],
                            "volume": 2,
                            "weight": 12,
                            "price": 100,
                            "skills": [],
                            "evidences": []
                        }
                    ]
                }
            ]
        },
        {
            "type": "delivery",
            "addressInfo": {
                "addressLine1": "Calle 1006, Piso 7, la-florida",
                "addressLine2": "",
                "coordinates": {
                    "latitude": -33.5414441,
                    "longitude": -70.5872566
                },
                "politicalArea": {
                    "adminAreaLevel1": "",
                    "adminAreaLevel2": "",
                    "adminAreaLevel3": "",
                    "adminAreaLevel4": "",
                    "code": ""
                },
                "zipCode": ""
            },
            "nodeInfo": {
                "referenceID": "NODE-007"
            },
            "sequenceNumber": 7,
            "serviceTime": 300,
            "timeWindow": {
                "start": "09:00",
                "end": "17:00"
            },
            "unassignedReason": "",
            "orders": [
                {
                    "contact": {
                        "email": "roberto.silva3@email.com",
                        "fullName": "Roberto Silva",
                        "nationalID": "12345678-9",
                        "phone": "+56912345678"
                    },
                    "referenceID": "106LA",
                    "instructions": "Entregar en piso 7, apartamento D",
                    "deliveryUnits": [
                        {
                            "lpn": "CODE-5",
                            "items": [
                                {
                                    "description": "bebida 350ml",
                                    "quantity": 13,
                                    "sku": "BEB-350-001"
                                }
                            ],
                            "volume": 2,
                            "weight": 12,
                            "price": 100,
                            "skills": [],
                            "evidences": []
                        }
                    ]
                }
            ]
        },
        {
            "type": "delivery",
            "addressInfo": {
                "addressLine1": "Calle 1013, Piso 14, la-florida",
                "addressLine2": "",
                "coordinates": {
                    "latitude": -33.5466664,
                    "longitude": -70.5596647
                },
                "politicalArea": {
                    "adminAreaLevel1": "",
                    "adminAreaLevel2": "",
                    "adminAreaLevel3": "",
                    "adminAreaLevel4": "",
                    "code": ""
                },
                "zipCode": ""
            },
            "nodeInfo": {
                "referenceID": "NODE-008"
            },
            "sequenceNumber": 8,
            "serviceTime": 300,
            "timeWindow": {
                "start": "09:00",
                "end": "17:00"
            },
            "unassignedReason": "",
            "orders": [
                {
                    "contact": {
                        "email": "ana.rodriguez@email.com",
                        "fullName": "Ana Rodriguez",
                        "nationalID": "44556677-8",
                        "phone": "+56944556677"
                    },
                    "referenceID": "113LA",
                    "instructions": "Entregar en piso 14, apartamento E",
                    "deliveryUnits": [
                        {
                            "lpn": "CODE-6",
                            "items": [
                                {
                                    "description": "bebida 350ml",
                                    "quantity": 13,
                                    "sku": "BEB-350-001"
                                }
                            ],
                            "volume": 2,
                            "weight": 12,
                            "price": 100,
                            "skills": [],
                            "evidences": []
                        }
                    ]
                }
            ]
        },
        {
            "type": "delivery",
            "addressInfo": {
                "addressLine1": "Calle 1015, Piso 16, la-florida",
                "addressLine2": "",
                "coordinates": {
                    "latitude": -33.5384557,
                    "longitude": -70.5767166
                },
                "politicalArea": {
                    "adminAreaLevel1": "",
                    "adminAreaLevel2": "",
                    "adminAreaLevel3": "",
                    "adminAreaLevel4": "",
                    "code": ""
                },
                "zipCode": ""
            },
            "nodeInfo": {
                "referenceID": "NODE-009"
            },
            "sequenceNumber": 9,
            "serviceTime": 300,
            "timeWindow": {
                "start": "09:00",
                "end": "17:00"
            },
            "unassignedReason": "",
            "orders": [
                {
                    "contact": {
                        "email": "lucia.herrera@email.com",
                        "fullName": "Lucia Herrera",
                        "nationalID": "33445566-7",
                        "phone": "+56933445566"
                    },
                    "referenceID": "115LA",
                    "instructions": "Entregar en piso 16, apartamento F",
                    "deliveryUnits": [
                        {
                            "lpn": "CODE-7",
                            "items": [
                                {
                                    "description": "bebida 350ml",
                                    "quantity": 13,
                                    "sku": "BEB-350-001"
                                }
                            ],
                            "volume": 2,
                            "weight": 12,
                            "price": 100,
                            "skills": [],
                            "evidences": []
                        }
                    ]
                }
            ]
        },
        {
            "type": "delivery",
            "addressInfo": {
                "addressLine1": "Calle 1018, Piso 19, la-florida",
                "addressLine2": "",
                "coordinates": {
                    "latitude": -33.5341326,
                    "longitude": -70.5560076
                },
                "politicalArea": {
                    "adminAreaLevel1": "",
                    "adminAreaLevel2": "",
                    "adminAreaLevel3": "",
                    "adminAreaLevel4": "",
                    "code": ""
                },
                "zipCode": ""
            },
            "nodeInfo": {
                "referenceID": "NODE-010"
            },
            "sequenceNumber": 10,
            "serviceTime": 300,
            "timeWindow": {
                "start": "09:00",
                "end": "17:00"
            },
            "unassignedReason": "",
            "orders": [
                {
                    "contact": {
                        "email": "fernando.castro@email.com",
                        "fullName": "Fernando Castro",
                        "nationalID": "22334455-6",
                        "phone": "+56922334455"
                    },
                    "referenceID": "118LA",
                    "instructions": "Entregar en piso 19, apartamento G",
                    "deliveryUnits": [
                        {
                            "lpn": "CODE-8",
                            "items": [
                                {
                                    "description": "bebida 350ml",
                                    "quantity": 13,
                                    "sku": "BEB-350-001"
                                }
                            ],
                            "volume": 2,
                            "weight": 12,
                            "price": 100,
                            "skills": [],
                            "evidences": []
                        }
                    ]
                }
            ]
        },
        {
            "type": "delivery",
            "addressInfo": {
                "addressLine1": "Calle 1002, Piso 3, la-florida",
                "addressLine2": "",
                "coordinates": {
                    "latitude": -33.5332085,
                    "longitude": -70.5516135
                },
                "politicalArea": {
                    "adminAreaLevel1": "",
                    "adminAreaLevel2": "",
                    "adminAreaLevel3": "",
                    "adminAreaLevel4": "",
                    "code": ""
                },
                "zipCode": ""
            },
            "nodeInfo": {
                "referenceID": "NODE-011"
            },
            "sequenceNumber": 11,
            "serviceTime": 300,
            "timeWindow": {
                "start": "09:00",
                "end": "17:00"
            },
            "unassignedReason": "",
            "orders": [
                {
                    "contact": {
                        "email": "maria.perez2@email.com",
                        "fullName": "Maria Perez",
                        "nationalID": "11223344-5",
                        "phone": "+56911223344"
                    },
                    "referenceID": "102LA",
                    "instructions": "Entregar en piso 3, apartamento H",
                    "deliveryUnits": [
                        {
                            "lpn": "CODE-9",
                            "items": [
                                {
                                    "description": "bebida 350ml",
                                    "quantity": 13,
                                    "sku": "BEB-350-001"
                                }
                            ],
                            "volume": 2,
                            "weight": 12,
                            "price": 100,
                            "skills": [],
                            "evidences": []
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