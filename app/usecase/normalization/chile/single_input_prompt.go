package chile

import (
	"context"
	"embed"
	"encoding/json"
	"fmt"
	"strings"
	"transport-app/app/domain"
	"transport-app/app/shared/utils"

	ioc "github.com/Ignaciojeria/einar-ioc/v2"
)

//go:embed chile_hierarchy.json
var chileHierarchyJSON embed.FS

// ChileHierarchy representa la estructura administrativa de Chile.
type ChileHierarchy struct {
	Region     string `json:"region"`
	RegionCode string `json:"region_iso_3166_2"`
	Provincias []struct {
		Name    string `json:"name"`
		Comunas []struct {
			Name string `json:"name"`
			Code string `json:"code"`
		} `json:"comunas"`
	} `json:"provincias"`
}

var chileHierarchy []ChileHierarchy

type SingleInputPrompt func(
	c context.Context,
	userInput domain.AddressInfo,
) string

func init() {
	ioc.Registry(NewSingleInputPrompt)
}

func NewSingleInputPrompt() (SingleInputPrompt, error) {
	data, err := chileHierarchyJSON.ReadFile("chile_hierarchy.json")
	if err != nil {
		return nil, err
	}
	if err := json.Unmarshal(data, &chileHierarchy); err != nil {
		return nil, err
	}

	return func(c context.Context, input domain.AddressInfo) string {

		// Extraer listas dinámicamente respetando la estructura jerárquica oficial
		var validRegions, validProvinces, validCommunes []string
		provinceMap := make(map[string]string) // Mapa comuna -> provincia
		regionMap := make(map[string]string)   // Mapa provincia -> región

		for _, region := range chileHierarchy {
			validRegions = append(validRegions, utils.NormalizeText(fmt.Sprintf("%s (%s)", region.Region, region.RegionCode)))
			for _, province := range region.Provincias {
				regionMap[utils.NormalizeText(province.Name)] = utils.NormalizeText(region.Region)
				var communeList []string
				for _, commune := range province.Comunas {
					communeName := utils.NormalizeText(commune.Name)
					validCommunes = append(validCommunes, communeName)
					communeList = append(communeList, communeName)
					provinceMap[communeName] = utils.NormalizeText(province.Name)
				}
				validProvinces = append(validProvinces, fmt.Sprintf("%s -> %s", utils.NormalizeText(province.Name), strings.Join(communeList, ", ")))
			}
		}

		prompt := fmt.Sprintf(`normaliza la siguiente direccion en chile segun el formato estandar

			**region de entrada:** %s   
			**provincia de entrada:** %s  
			**comuna de entrada:** %s 
			
			---
			### jerarquia administrativa de chile (respetando nombres oficiales):
			
			#### regiones validas:
			%s

			#### provincias y sus comunas:
			%s

			#### comunas validas:
			%s

			---
			### **formato de salida esperado (json)**
			{
				"district": "comuna normalizada",
				"province": "provincia normalizada",
				"state": "region normalizada",
			}
			
			### **reglas de normalizacion:**
			1. **la comuna ingresada debe ser una comuna valida de chile.**
			2. **la provincia debe ser la correcta segun la comuna ingresada.**
			3. **la region debe coincidir con la provincia correspondiente.**
			4. **solo se deben normalizar nombres oficiales de comunas, provincias y regiones.**
			6. **no se deben modificar los nombres oficiales de comunas provincias o regiones.**
			7. **los nombres deben conservar la letra ñ si forman parte del nombre original. no reemplazar ni eliminar estos caracteres.**
			`,
			input.State.String(),
			input.Province.String(),
			input.District.String(),
			strings.Join(validRegions, ", "),
			strings.Join(validProvinces, "\n"),
			strings.Join(validCommunes, ", "))

		return prompt
	}, nil
}
