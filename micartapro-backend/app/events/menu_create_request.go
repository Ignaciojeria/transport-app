package events

import (
	"encoding/json"
	"strings"
	"time"

	cloudevents "github.com/cloudevents/sdk-go/v2"
)

/*
	STATION
*/

type Station string

const (
	StationKitchen Station = "KITCHEN"
	StationBar     Station = "BAR"
)

/*
	UNITS
*/

type UnitOfMeasure string

const (
	UnitEach        UnitOfMeasure = "EACH"
	UnitGram        UnitOfMeasure = "GRAM"
	UnitKilogram    UnitOfMeasure = "KILOGRAM"
	UnitMilliliter  UnitOfMeasure = "MILLILITER"
	UnitLiter       UnitOfMeasure = "LITER"
	UnitMeter       UnitOfMeasure = "METER"
	UnitSquareMeter UnitOfMeasure = "SQUARE_METER"
)

/*
	PRICING
*/

type PricingMode string

const (
	PricingUnit   PricingMode = "UNIT"
	PricingWeight PricingMode = "WEIGHT"
	PricingVolume PricingMode = "VOLUME"
	PricingLength PricingMode = "LENGTH"
	PricingArea   PricingMode = "AREA"
)

/*
	CURRENCY
*/

type Currency string

const (
	CurrencyUSD Currency = "USD"
	CurrencyCLP Currency = "CLP"
	CurrencyBRL Currency = "BRL"
)

type Pricing struct {
	Mode         PricingMode   `json:"mode"`
	Unit         UnitOfMeasure `json:"unit"`
	PricePerUnit float64       `json:"pricePerUnit"`
	BaseUnit     float64       `json:"baseUnit"`
	CostPerUnit  *float64      `json:"costPerUnit,omitempty"` // Costo por unidad (nil = no definido)
}

// HasCost indica si hay costo definido.
func (p Pricing) HasCost() bool {
	return p.CostPerUnit != nil
}

// EffectiveCostPerUnit retorna el costo por unidad; 0 si no está definido.
func (p Pricing) EffectiveCostPerUnit() float64 {
	if p.CostPerUnit == nil {
		return 0
	}
	return *p.CostPerUnit
}

/*
	MULTILINGUAL TEXT
*/

type MultilingualText struct {
	Base      string            `json:"base"`
	Languages map[string]string `json:"languages"`
}

// GetText retorna el texto en el idioma especificado, o el texto base si no está disponible
func (m MultilingualText) GetText(lang string) string {
	if m.Languages != nil {
		if text, ok := m.Languages[lang]; ok && text != "" {
			return text
		}
	}
	return m.Base
}

// GetBase retorna el texto base (idioma principal)
func (m MultilingualText) GetBase() string {
	return m.Base
}

/*
	FOOD ATTRIBUTES
*/

type FoodAttribute string

const (
	FoodGluten     FoodAttribute = "GLUTEN"
	FoodSeafood    FoodAttribute = "SEAFOOD"
	FoodNuts       FoodAttribute = "NUTS"
	FoodDairy      FoodAttribute = "DAIRY"
	FoodEggs       FoodAttribute = "EGGS"
	FoodSoy        FoodAttribute = "SOY"
	FoodVegan      FoodAttribute = "VEGAN"
	FoodVegetarian FoodAttribute = "VEGETARIAN"
	FoodSpicy      FoodAttribute = "SPICY"
	FoodAlcohol    FoodAttribute = "ALCOHOL"
)

/*
	MENU PRESENTATION STYLE
*/

type MenuPresentationStyle string

const (
	MenuStyleHero   MenuPresentationStyle = "HERO"
	MenuStyleModern MenuPresentationStyle = "MODERN"
)

/*
	DESCRIPTION WITH SELECTABLES
	description[]: cada bloque puede tener id y selectables para preferencias sin precio.
	sides[]: variantes reales con precio.
*/

// DescriptionSelectableOption es una opción seleccionable dentro de un bloque de descripción (preferencia sin precio).
type DescriptionSelectableOption struct {
	ID   string           `json:"id"`
	Name MultilingualText `json:"name"`
}

// DescriptionSelectables define preferencias sin precio dentro de un bloque de descripción.
type DescriptionSelectables struct {
	Selection struct {
		Mode string `json:"mode"` // "SINGLE" o "MULTIPLE"
		Min  int    `json:"min"`
		Max  int    `json:"max"`
	} `json:"selection"`
	Options []DescriptionSelectableOption `json:"options"`
}

// DescriptionBlock es un elemento del array description con soporte opcional de selectables.
type DescriptionBlock struct {
	ID          string                 `json:"id,omitempty"`
	Base        string                 `json:"base"`
	Languages   map[string]string      `json:"languages"`
	Selectables *DescriptionSelectables `json:"selectables,omitempty"`
}

// GetText retorna el texto en el idioma especificado.
func (d DescriptionBlock) GetText(lang string) string {
	if d.Languages != nil {
		if text, ok := d.Languages[lang]; ok && text != "" {
			return text
		}
	}
	return d.Base
}

/*
	MENU MODELS
*/

type Side struct {
	ID             string           `json:"id"`
	Name           MultilingualText `json:"name"`
	FoodAttributes []FoodAttribute  `json:"foodAttributes,omitempty"`
	Pricing        Pricing          `json:"pricing"`
	PhotoUrl       string           `json:"photoUrl,omitempty"`
	Station        Station          `json:"station,omitempty"`
}

type MenuItem struct {
	ID             string            `json:"id"`
	Title          MultilingualText  `json:"title"`
	Description    []DescriptionBlock `json:"description,omitempty"` // Array: cada elemento es una dimensión; puede tener id y selectables para preferencias sin precio
	FoodAttributes []FoodAttribute   `json:"foodAttributes,omitempty"`
	Sides          []Side            `json:"sides,omitempty"`
	Pricing        Pricing           `json:"pricing"`
	PhotoUrl       string            `json:"photoUrl,omitempty"`
	Station        Station           `json:"station,omitempty"`
}

// GetDescriptionText une todos los elementos de Description en un solo texto para el idioma dado.
func (m MenuItem) GetDescriptionText(lang string) string {
	if len(m.Description) == 0 {
		return ""
	}
	parts := make([]string, 0, len(m.Description))
	for _, d := range m.Description {
		t := d.GetText(lang)
		if t != "" {
			parts = append(parts, t)
		}
	}
	return strings.Join(parts, " ")
}

type menuItemRaw struct {
	ID             string            `json:"id"`
	Title          MultilingualText  `json:"title"`
	Description    json.RawMessage   `json:"description"`
	FoodAttributes []FoodAttribute   `json:"foodAttributes,omitempty"`
	Sides          []Side            `json:"sides,omitempty"`
	Pricing        Pricing           `json:"pricing"`
	PhotoUrl       string            `json:"photoUrl,omitempty"`
	Station        Station           `json:"station,omitempty"`
}

// UnmarshalJSON acepta description como array de DescriptionBlock, array de MultilingualText (legacy), o objeto único (legacy).
func (m *MenuItem) UnmarshalJSON(data []byte) error {
	var raw menuItemRaw
	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}
	m.ID = raw.ID
	m.Title = raw.Title
	m.FoodAttributes = raw.FoodAttributes
	m.Sides = raw.Sides
	m.Pricing = raw.Pricing
	m.PhotoUrl = raw.PhotoUrl
	m.Station = raw.Station
	if len(raw.Description) == 0 {
		m.Description = nil
		return nil
	}
	var arr []DescriptionBlock
	if err := json.Unmarshal(raw.Description, &arr); err == nil {
		m.Description = arr
		return nil
	}
	var legacyArr []MultilingualText
	if err := json.Unmarshal(raw.Description, &legacyArr); err == nil {
		m.Description = make([]DescriptionBlock, len(legacyArr))
		for i, mt := range legacyArr {
			m.Description[i] = DescriptionBlock{Base: mt.Base, Languages: mt.Languages}
		}
		return nil
	}
	var single MultilingualText
	if err := json.Unmarshal(raw.Description, &single); err != nil {
		return err
	}
	m.Description = []DescriptionBlock{{Base: single.Base, Languages: single.Languages}}
	return nil
}

type MenuCategory struct {
	Title MultilingualText `json:"title"`
	Items []MenuItem       `json:"items"`
}

/*
	BUSINESS
*/

type BusinessInfo struct {
	BusinessName  string   `json:"businessName"`
	Whatsapp      string   `json:"whatsapp"`
	BusinessHours []string `json:"businessHours"`
	Currency      Currency `json:"currency"`
}

// EffectiveCurrency retorna la moneda efectiva del negocio; si está vacía retorna CLP por defecto (compatibilidad con menús antiguos).
func (b BusinessInfo) EffectiveCurrency() Currency {
	if b.Currency == "" {
		return CurrencyCLP
	}
	return b.Currency
}

/*
	DELIVERY OPTIONS
*/

type DeliveryOptionType string

const (
	DeliveryOptionPickup   DeliveryOptionType = "PICKUP"
	DeliveryOptionDelivery DeliveryOptionType = "DELIVERY"
	DeliveryOptionDigital  DeliveryOptionType = "DIGITAL"
)

type TimeRequestType string

const (
	TimeRequestExact  TimeRequestType = "EXACT"
	TimeRequestWindow TimeRequestType = "WINDOW"
)

type TimeWindow struct {
	Start string `json:"start"`
	End   string `json:"end"`
}

type DeliveryOption struct {
	Type            DeliveryOptionType `json:"type"`
	RequireTime     bool               `json:"requireTime"`
	TimeRequestType TimeRequestType    `json:"timeRequestType,omitempty"`
	TimeWindows     []TimeWindow       `json:"timeWindows,omitempty"`
}

/*
	IMAGE GENERATION
*/

type CoverImageGenerationRequest struct {
	Prompt     string `json:"prompt"`
	ImageCount int    `json:"imageCount"`
	// Campos temporales para URLs pre-firmadas (no se persisten)
	UploadURL string `json:"uploadUrl,omitempty"`
	PublicURL string `json:"publicUrl,omitempty"`
}

type ImageGenerationRequest struct {
	MenuItemID  string `json:"menuItemId"`
	Prompt      string `json:"prompt"`
	AspectRatio string `json:"aspectRatio"`
	ImageCount  int    `json:"imageCount"`
	// Campos temporales para URLs pre-firmadas (no se persisten)
	UploadURL string `json:"uploadUrl,omitempty"`
	PublicURL string `json:"publicUrl,omitempty"`
}

/*
	IMAGE EDITION
*/

type CoverImageEditionRequest struct {
	Prompt            string `json:"prompt"`
	ReferenceImageUrl string `json:"referenceImageUrl"`
	ImageCount        int    `json:"imageCount"`
	// Campos temporales para URLs pre-firmadas (no se persisten)
	UploadURL string `json:"uploadUrl,omitempty"`
	PublicURL string `json:"publicUrl,omitempty"`
}

type ImageEditionRequest struct {
	MenuItemID        string `json:"menuItemId"`
	Prompt            string `json:"prompt"`
	ReferenceImageUrl string `json:"referenceImageUrl"`
	AspectRatio       string `json:"aspectRatio"`
	ImageCount        int    `json:"imageCount"`
	// Campos temporales para URLs pre-firmadas (no se persisten)
	UploadURL string `json:"uploadUrl,omitempty"`
	PublicURL string `json:"publicUrl,omitempty"`
}

/*
	EVENT
*/

type MenuCreateRequest struct {
	ID                          string                       `json:"id"`
	PresentationStyle           MenuPresentationStyle        `json:"presentationStyle"`
	CoverImage                  string                       `json:"coverImage"`
	FooterImage                 string                       `json:"footerImage"`
	BusinessInfo                BusinessInfo                 `json:"businessInfo"`
	Menu                        []MenuCategory               `json:"menu"`
	DeliveryOptions             []DeliveryOption             `json:"deliveryOptions,omitempty"`
	CoverImageGenerationRequest *CoverImageGenerationRequest `json:"coverImageGenerationRequest,omitempty"`
	ImageGenerationRequests     []ImageGenerationRequest     `json:"imageGenerationRequests,omitempty"`
	CoverImageEditionRequest    *CoverImageEditionRequest    `json:"coverImageEditionRequest,omitempty"`
	ImageEditionRequests        []ImageEditionRequest        `json:"imageEditionRequests,omitempty"`
}

func (c MenuCreateRequest) ToCloudEvent(source string) cloudevents.Event {
	event := cloudevents.NewEvent()
	event.SetSubject("menu.create.request")
	event.SetType(EventMenuCreateRequested)
	event.SetSource(source)
	event.SetData(cloudevents.ApplicationJSON, c)
	event.SetTime(time.Now())
	return event
}

// EffectivePresentationStyle retorna el estilo de presentación efectivo; si está vacío retorna MODERN por defecto.
func (c *MenuCreateRequest) EffectivePresentationStyle() MenuPresentationStyle {
	if c == nil || c.PresentationStyle == "" {
		return MenuStyleModern
	}
	return c.PresentationStyle
}

// EnsurePresentationStyleDefault asigna MODERN a PresentationStyle si está vacío (útil tras deserializar JSON antiguo).
func (c *MenuCreateRequest) EnsurePresentationStyleDefault() {
	if c != nil && c.PresentationStyle == "" {
		c.PresentationStyle = MenuStyleModern
	}
}

// Clean limpia los campos temporales de solicitudes de generación/edición de imágenes
// después de que hayan sido procesadas. Estos campos no deben persistirse en storage/Supabase.
func (c *MenuCreateRequest) Clean() {
	c.CoverImageGenerationRequest = nil
	c.ImageGenerationRequests = nil
	c.CoverImageEditionRequest = nil
	c.ImageEditionRequests = nil
}

// NormalizeGCSURL corrige URLs de GCS mal formateadas
// Es idempotente: puede aplicarse múltiples veces sin causar efectos secundarios
// Es pública para poder usarse desde otros paquetes al asignar URLs
func NormalizeGCSURL(url string) string {
	if url == "" {
		return url
	}

	// Verificar si la URL ya está correctamente formateada
	if strings.HasPrefix(url, "https://storage.googleapis.com") || strings.HasPrefix(url, "http://storage.googleapis.com") {
		// Aún así, verificar duplicaciones de "https" o "http"
		url = strings.ReplaceAll(url, "httpshttps://", "https://")
		url = strings.ReplaceAll(url, "httphttp://", "http://")
		return url
	}

	// Corregir "https.storage.googleapis.com" → "https://storage.googleapis.com"
	// Solo aplicar si el patrón incorrecto está presente
	if strings.Contains(url, "https.storage.googleapis.com") {
		url = strings.ReplaceAll(url, "https.storage.googleapis.com", "https://storage.googleapis.com")
	}
	if strings.Contains(url, "http.storage.googleapis.com") {
		url = strings.ReplaceAll(url, "http.storage.googleapis.com", "http://storage.googleapis.com")
	}

	// Manejar casos donde se duplicó "https" (httpshttps://storage...)
	url = strings.ReplaceAll(url, "httpshttps://", "https://")
	url = strings.ReplaceAll(url, "httphttp://", "http://")

	return url
}

// NormalizeImageURLs normaliza todas las URLs de imágenes del menú antes de guardar en BD.
// Corrige formatos mal guardados como "httpshttps://" o "https.storage.googleapis.com".
// Es idempotente y seguro llamarlo múltiples veces.
func (c *MenuCreateRequest) NormalizeImageURLs() {
	if c == nil {
		return
	}

	c.CoverImage = NormalizeGCSURL(c.CoverImage)
	c.FooterImage = NormalizeGCSURL(c.FooterImage)

	for i := range c.Menu {
		for j := range c.Menu[i].Items {
			c.Menu[i].Items[j].PhotoUrl = NormalizeGCSURL(c.Menu[i].Items[j].PhotoUrl)
			for k := range c.Menu[i].Items[j].Sides {
				c.Menu[i].Items[j].Sides[k].PhotoUrl = NormalizeGCSURL(c.Menu[i].Items[j].Sides[k].PhotoUrl)
			}
		}
	}
}
