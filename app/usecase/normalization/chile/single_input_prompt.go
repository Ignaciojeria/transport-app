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
	providerInput domain.AddressInfo) string

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

	return func(c context.Context, userInput, providerInput domain.AddressInfo) string {
		userText := userInput.AddressLine1
		userIndications := userInput.AddressLine2
		providerInputAddress := providerInput.AddressLine1
		providerLat := providerInput.Location.Lat()
		providerLon := providerInput.Location.Lon()

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

			**direccion ingresada por el usuario:** %s 
			**Indicaciones ingresadas por el usuario:** %s   
			**direccion sugerida por el proveedor:** %s  
			**coordenadas:** %.6f, %.6f  
			
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
				"providerAddress":"direccion sugerida por el proveedor",
				"addressLine1": "calle y numero normalizado",
				"addressLine2": "informacion adicional si aplica obtener desde direccion ingresada por el usuario",
				"district": "comuna normalizada",
				"province": "provincia correspondiente",
				"state": "region correspondiente",
				"latitude": %.6f,
				"longitude": %.6f
			}
			
			### **reglas de normalizacion:**
			1. **la comuna ingresada debe ser una comuna valida de chile.**
			2. **la provincia debe ser la correcta segun la comuna ingresada.**
			3. **la region debe coincidir con la provincia correspondiente.**
			4. **correccion ortografica minima solo en calles y nombres comunes.**
			5. **no se deben modificar los nombres oficiales de comunas provincias o regiones.**`,
			userText, userIndications, providerInputAddress, providerLat, providerLon,
			strings.Join(validRegions, ", "),
			strings.Join(validProvinces, "\n"),
			strings.Join(validCommunes, ", "),
			providerLat, providerLon)

		return prompt
	}, nil
}
