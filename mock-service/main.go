package main

import (
	"log"
	"math/rand"
	"os"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/basicauth"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/joho/godotenv"
)

// 1. Estructuras de datos (Los moldes que mapean las columnas de tus CSV)
type Establecimiento struct {
	ID       string `json:"id"`
	Fna      string `json:"fna"`
	Gna      string `json:"gna"`
	Nam      string `json:"nam"`
	Cui      string `json:"cui"`
	Cue      string `json:"cue"`
	Anx      string `json:"anx"`
	Ges      string `json:"ges"`
	NemMde   string `json:"nem_mde"`
	Tip      string `json:"tip"`
	Dep      string `json:"dep"`
	Dir      string `json:"dir"`
	Bar      string `json:"bar"`
	Com      string `json:"com"`
	Sag      string `json:"sag"`
	Geometry string `json:"geometry"`
}

type Alojamiento struct {
	ID        string `json:"id"`
	Nombre    string `json:"nombre"`
	Domicilio string `json:"domicilio"`
	Telefono  string `json:"telefono"`
	Email     string `json:"email"`
	Web       string `json:"web"`
	Tipo      string `json:"tipo"`
	Calle     string `json:"calle"`
	Altura    string `json:"altura"`
	Direccion string `json:"direccion"`
	Long      string `json:"long"`
	Lat       string `json:"lat"`
	Geometry  string `json:"geometry"`
}

func main() {
	// Cargar variables de entorno
	godotenv.Load("../.env")

	app := fiber.New()

	// 2. CONFIGURACIÓN DE CORS (Obligatorio para que conecte el frontend)
	app.Use(cors.New(cors.Config{
		AllowOrigins: "*",
		AllowHeaders: "Origin, Content-Type, Accept, Authorization",
		AllowMethods: "GET, POST, PUT, DELETE, OPTIONS",
	}))

	// 3. Muestra representativa hardcodeada de tus archivos CSV
	listaEscuelas := []Establecimiento{
		{ID: "1", Fna: "Escuela Nro 1", Gna: "EGB", Nam: "San Martín", Cui: "20123456789", Cue: "0600123", Anx: "0", Ges: "Estatal", NemMde: "EGB 1", Tip: "Escuela", Dep: "Morón", Dir: "Av. Rivadavia 1234", Bar: "Centro", Com: "Comuna 1", Sag: "123", Geometry: "POINT(-34.65 -58.62)"},
		{ID: "2", Fna: "Colegio Nacional", Gna: "Secundario", Nam: "Manuel Belgrano", Cui: "20987654321", Cue: "0600456", Anx: "0", Ges: "Privado", NemMde: "SEC 2", Tip: "Colegio", Dep: "Merlo", Dir: "Suipacha 450", Bar: "Merlo Norte", Com: "Comuna 2", Sag: "456", Geometry: "POINT(-34.66 -58.72)"},
	}

	listaHoteles := []Alojamiento{
		{ID: "1", Nombre: "Hotel Plaza", Domicilio: "Mitre 800", Telefono: "1144332211", Email: "info@hotelplaza.com", Web: "://hotelplaza.com", Tipo: "Hotel", Calle: "Mitre", Altura: "800", Direccion: "Mitre 800, Merlo", Long: "-58.723", Lat: "-34.661", Geometry: "POINT(-58.723 -34.661)"},
		{ID: "2", Nombre: "Hostel del Sur", Domicilio: "Colon 120", Telefono: "1155667788", Email: "contacto@hosteldelsur.com", Web: "://hosteldelsur.com", Tipo: "Hostel", Calle: "Colon", Altura: "120", Direccion: "Colon 120, Merlo", Long: "-58.725", Lat: "-34.663", Geometry: "POINT(-58.725 -34.663)"},
	}

	// 4. Endpoints Públicos
	app.Get("/health", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{"status": "ok"})
	})

	app.Get("/alojamientos", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"total": len(listaHoteles),
			"data":  listaHoteles,
		})
	})

	app.Get("/establecimientos_public", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"total": len(listaEscuelas),
			"data":  listaEscuelas})
	})

	// 5. Grupo de Rutas Protegidas (Basic Auth)
	auth := app.Group(
		"/",
		basicauth.New(basicauth.Config{
			Users: map[string]string{
				os.Getenv("BASIC_AUTH_USER"): os.Getenv("BASIC_AUTH_PASSWORD"),
			},
		}),
	)

	// Endpoint Protegido: Buscar Escuela por ID
	auth.Get("/establecimientos/:id", func(c *fiber.Ctx) error {
		idBuscado := c.Params("id")
		for _, esc := range listaEscuelas {
			if esc.ID == idBuscado {
				return c.JSON(esc)
			}
		}
		return c.Status(404).JSON(fiber.Map{"status": false, "error": "establecimiento no encontrado"})
	})

	// Endpoint Protegido: Mercado Pago Mock (Delay corregido de 100s a 1s)
	auth.Post("/mercadopago/pago", func(c *fiber.Ctx) error {
		type PagoRequest struct {
			CardNumber  string  `json:"card_number"`
			Amount      float64 `json:"amount"`
			Currency    string  `json:"currency"`
			Description string  `json:"description"`
		}

		var req PagoRequest
		if err := c.BodyParser(&req); err != nil {
			return c.Status(400).JSON(fiber.Map{"status": false, "error": "json inválido"})
		}

		if req.CardNumber != "4509953566233704" {
			return c.Status(400).JSON(fiber.Map{"status": false, "approved": false, "message": "tarjeta rechazada"})
		}

		if req.Amount > 50000 {
			return c.Status(400).JSON(fiber.Map{"status": false, "approved": false, "message": "saldo insuficiente"})
		}

		// Simula espera real de red para que el frontend pruebe spinners
		time.Sleep(1 * time.Second)

		return c.JSON(fiber.Map{
			"status":      true,
			"approved":    true,
			"payment_id":  "MP-MOCK-" + strconv.Itoa(rand.Intn(90000)+10000),
			"currency":    req.Currency,
			"amount":      req.Amount,
			"description": req.Description,
			"message":     "pago aprobado con éxito",
		})
	})

	// 6. Encendido del Servidor
	puerto := os.Getenv("MOCK_PORT")
	if puerto == "" {
		puerto = "3000"
	}

	log.Printf("Servidor Mock listo en el puerto %s", puerto)
	log.Fatal(app.Listen(":" + puerto))
}