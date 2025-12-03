import Link from 'next/link'
import Image from 'next/image'

export const metadata = {
  title: 'Términos y Condiciones - MiCartaPro',
  description: 'Términos y condiciones de uso de MiCartaPro. Conoce los términos que rigen el uso de nuestros servicios.',
}

export default function TermsPage() {
  return (
    <div className="min-h-screen bg-gradient-to-br from-blue-50 via-white to-indigo-50">
      {/* Navigation */}
      <nav className="border-b bg-white/80 backdrop-blur-sm sticky top-0 z-50">
        <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
          <div className="flex justify-between items-center h-16">
            <Link href="/" className="flex items-center space-x-2">
              <Image 
                src="/logo.png" 
                alt="MiCartaPro Logo" 
                width={240} 
                height={72}
                className="h-16 md:h-20 w-auto"
              />
            </Link>
            <Link 
              href="/"
              className="text-gray-600 hover:text-blue-600 transition-colors"
            >
              Volver al inicio
            </Link>
          </div>
        </div>
      </nav>

      {/* Content */}
      <div className="max-w-4xl mx-auto px-4 sm:px-6 lg:px-8 py-12 md:py-16">
        <div className="bg-white rounded-lg shadow-lg p-8 md:p-12">
          <h1 className="text-3xl md:text-4xl font-bold text-gray-900 mb-6">
            Términos y Condiciones de Uso
          </h1>
          
          <p className="text-gray-600 mb-8">
            <strong>Última actualización:</strong> {new Date().toLocaleDateString('es-ES', { year: 'numeric', month: 'long', day: 'numeric' })}
          </p>

          <div className="prose prose-lg max-w-none space-y-8 text-gray-700">
            <section>
              <h2 className="text-2xl font-semibold text-gray-900 mb-4">1. Aceptación de los Términos</h2>
              <p>
                Bienvenido a MiCartaPro. Estos Términos y Condiciones de Uso ("Términos") rigen tu acceso y uso de nuestros servicios de menú digital, incluyendo nuestro sitio web, aplicaciones y cualquier servicio relacionado (colectivamente, los "Servicios").
              </p>
              <p>
                Al acceder o utilizar nuestros Servicios, aceptas estar sujeto a estos Términos. Si no estás de acuerdo con alguna parte de estos Términos, no debes utilizar nuestros Servicios.
              </p>
            </section>

            <section>
              <h2 className="text-2xl font-semibold text-gray-900 mb-4">2. Descripción de los Servicios</h2>
              <p>
                MiCartaPro proporciona servicios de menú digital que permiten a los restaurantes crear, gestionar y compartir menús digitales con sus clientes. Nuestros servicios incluyen:
              </p>
              <ul className="list-disc pl-6 space-y-2">
                <li>Diseño y desarrollo de menús digitales personalizados</li>
                <li>Generación de códigos QR exclusivos</li>
                <li>Integración de carrito de compras</li>
                <li>Integración con WhatsApp para recepción de pedidos</li>
                <li>Diseño responsivo para todos los dispositivos</li>
              </ul>
            </section>

            <section>
              <h2 className="text-2xl font-semibold text-gray-900 mb-4">3. Elegibilidad</h2>
              <p>
                Para utilizar nuestros Servicios, debes:
              </p>
              <ul className="list-disc pl-6 space-y-2">
                <li>Tener al menos 18 años de edad o tener el consentimiento de un padre o tutor legal</li>
                <li>Tener la capacidad legal para celebrar contratos vinculantes</li>
                <li>No estar prohibido de utilizar nuestros Servicios bajo las leyes aplicables</li>
                <li>Proporcionar información precisa y completa al registrarte o utilizar nuestros Servicios</li>
              </ul>
            </section>

            <section>
              <h2 className="text-2xl font-semibold text-gray-900 mb-4">4. Uso de los Servicios</h2>
              
              <h3 className="text-xl font-semibold text-gray-900 mb-3 mt-4">4.1. Uso Permitido</h3>
              <p>Puedes utilizar nuestros Servicios únicamente para fines legales y de acuerdo con estos Términos. Aceptas:</p>
              <ul className="list-disc pl-6 space-y-2">
                <li>Utilizar los Servicios de manera responsable y ética</li>
                <li>Respetar los derechos de propiedad intelectual de otros</li>
                <li>Cumplir con todas las leyes y regulaciones aplicables</li>
                <li>No interferir con el funcionamiento de los Servicios</li>
              </ul>

              <h3 className="text-xl font-semibold text-gray-900 mb-3 mt-4">4.2. Uso Prohibido</h3>
              <p>Está estrictamente prohibido:</p>
              <ul className="list-disc pl-6 space-y-2">
                <li>Utilizar los Servicios para actividades ilegales o fraudulentas</li>
                <li>Intentar acceder no autorizado a nuestros sistemas o datos</li>
                <li>Transmitir virus, malware o código malicioso</li>
                <li>Realizar ingeniería inversa, descompilar o desensamblar nuestros Servicios</li>
                <li>Utilizar bots, scripts automatizados o métodos similares para acceder a los Servicios</li>
                <li>Copiar, modificar o distribuir nuestros Servicios sin autorización</li>
                <li>Utilizar los Servicios para enviar spam o comunicaciones no solicitadas</li>
                <li>Violar los derechos de privacidad o propiedad intelectual de otros</li>
              </ul>
            </section>

            <section>
              <h2 className="text-2xl font-semibold text-gray-900 mb-4">5. Cuentas y Registro</h2>
              <p>
                Para acceder a ciertas funcionalidades de nuestros Servicios, puede ser necesario crear una cuenta. Eres responsable de:
              </p>
              <ul className="list-disc pl-6 space-y-2">
                <li>Mantener la confidencialidad de tus credenciales de cuenta</li>
                <li>Todas las actividades que ocurran bajo tu cuenta</li>
                <li>Notificarnos inmediatamente de cualquier uso no autorizado de tu cuenta</li>
                <li>Proporcionar información precisa y actualizada</li>
              </ul>
            </section>

            <section>
              <h2 className="text-2xl font-semibold text-gray-900 mb-4">6. Contenido del Usuario</h2>
              <p>
                Al utilizar nuestros Servicios, puedes proporcionar contenido, incluyendo información de menús, imágenes, textos y otros materiales ("Contenido del Usuario"). Al proporcionar Contenido del Usuario, otorgas a MiCartaPro una licencia no exclusiva, mundial, libre de regalías y transferible para usar, reproducir, modificar y distribuir dicho contenido únicamente para proporcionar y mejorar nuestros Servicios.
              </p>
              <p>
                Eres responsable de asegurar que tu Contenido del Usuario:
              </p>
              <ul className="list-disc pl-6 space-y-2">
                <li>No infrinja los derechos de propiedad intelectual de terceros</li>
                <li>No contenga material difamatorio, obsceno o ilegal</li>
                <li>No viole los derechos de privacidad de otros</li>
                <li>Cumpla con todas las leyes y regulaciones aplicables</li>
              </ul>
            </section>

            <section>
              <h2 className="text-2xl font-semibold text-gray-900 mb-4">7. Propiedad Intelectual</h2>
              <p>
                Todos los derechos, títulos e intereses en y sobre nuestros Servicios, incluyendo pero no limitado a software, diseño, texto, gráficos, logos, iconos y compilaciones de datos, son propiedad de MiCartaPro o sus licenciantes y están protegidos por leyes de propiedad intelectual internacionales.
              </p>
              <p>
                No se te otorga ningún derecho, título o interés en nuestros Servicios excepto el derecho limitado de uso según estos Términos.
              </p>
            </section>

            <section>
              <h2 className="text-2xl font-semibold text-gray-900 mb-4">8. Precios y Pagos</h2>
              <p>
                Nuestros servicios están disponibles desde $150 USD con las siguientes condiciones:
              </p>
              <ul className="list-disc pl-6 space-y-2">
                <li><strong>Primer año:</strong> Gratis (oferta promocional con cupos limitados)</li>
                <li><strong>Renovación:</strong> $10 USD mensuales a partir del segundo año</li>
                <li>Los precios están sujetos a cambios con previo aviso</li>
                <li>Los pagos se procesan según los términos acordados en el contrato de servicio</li>
              </ul>
              <p className="mt-4">
                Todos los precios están expresados en dólares estadounidenses (USD) y pueden estar sujetos a impuestos según tu jurisdicción.
              </p>
            </section>

            <section>
              <h2 className="text-2xl font-semibold text-gray-900 mb-4">9. Disponibilidad del Servicio</h2>
              <p>
                Nos esforzamos por mantener nuestros Servicios disponibles de manera continua, pero no garantizamos que los Servicios estarán disponibles en todo momento, libres de interrupciones o errores. Podemos realizar mantenimiento programado o no programado que puede resultar en interrupciones temporales del servicio.
              </p>
            </section>

            <section>
              <h2 className="text-2xl font-semibold text-gray-900 mb-4">10. Limitación de Responsabilidad</h2>
              <p>
                EN LA MÁXIMA MEDIDA PERMITIDA POR LA LEY APLICABLE, MICARTAPRO Y SUS AFILIADOS NO SERÁN RESPONSABLES POR DAÑOS INDIRECTOS, INCIDENTALES, ESPECIALES, CONSECUENCIALES O PUNITIVOS, INCLUYENDO PERO NO LIMITADO A PÉRDIDA DE BENEFICIOS, DATOS O USO, RESULTANTES DEL USO O IMPOSIBILIDAD DE USAR NUESTROS SERVICIOS.
              </p>
              <p>
                Nuestra responsabilidad total hacia ti por cualquier reclamo relacionado con nuestros Servicios no excederá el monto que hayas pagado a MiCartaPro en los doce (12) meses anteriores al evento que dio lugar al reclamo.
              </p>
            </section>

            <section>
              <h2 className="text-2xl font-semibold text-gray-900 mb-4">11. Indemnización</h2>
              <p>
                Aceptas indemnizar, defender y eximir de responsabilidad a MiCartaPro, sus afiliados, directores, empleados y agentes de cualquier reclamo, demanda, pérdida, responsabilidad y gasto (incluyendo honorarios de abogados) que surjan de o estén relacionados con:
              </p>
              <ul className="list-disc pl-6 space-y-2">
                <li>Tu uso de los Servicios</li>
                <li>Tu violación de estos Términos</li>
                <li>Tu violación de cualquier ley o derecho de terceros</li>
                <li>Tu Contenido del Usuario</li>
              </ul>
            </section>

            <section>
              <h2 className="text-2xl font-semibold text-gray-900 mb-4">12. Terminación</h2>
              <p>
                Podemos terminar o suspender tu acceso a nuestros Servicios inmediatamente, sin previo aviso, por cualquier motivo, incluyendo pero no limitado a:
              </p>
              <ul className="list-disc pl-6 space-y-2">
                <li>Violación de estos Términos</li>
                <li>Uso fraudulento o ilegal de los Servicios</li>
                <li>No pago de tarifas aplicables</li>
                <li>Por razones operativas o comerciales</li>
              </ul>
              <p className="mt-4">
                También puedes terminar tu uso de los Servicios en cualquier momento. Al terminar, tu derecho a utilizar los Servicios cesará inmediatamente.
              </p>
            </section>

            <section>
              <h2 className="text-2xl font-semibold text-gray-900 mb-4">13. Servicios de Terceros</h2>
              <p>
                Nuestros Servicios pueden integrarse con servicios de terceros, incluyendo WhatsApp. Tu uso de estos servicios de terceros está sujeto a sus propios términos y condiciones y políticas de privacidad. No somos responsables de las prácticas de privacidad o el contenido de servicios de terceros.
              </p>
            </section>

            <section>
              <h2 className="text-2xl font-semibold text-gray-900 mb-4">14. Modificaciones de los Términos</h2>
              <p>
                Nos reservamos el derecho de modificar estos Términos en cualquier momento. Te notificaremos de cambios materiales publicando los Términos actualizados en esta página y actualizando la fecha de "Última actualización".
              </p>
              <p>
                Tu uso continuado de los Servicios después de que los Términos modificados entren en vigor constituye tu aceptación de los Términos modificados.
              </p>
            </section>

            <section>
              <h2 className="text-2xl font-semibold text-gray-900 mb-4">15. Ley Aplicable y Jurisdicción</h2>
              <p>
                Estos Términos se rigen e interpretan de acuerdo con las leyes aplicables, sin dar efecto a ningún principio de conflictos de leyes.
              </p>
              <p>
                Cualquier disputa que surja de o esté relacionada con estos Términos o nuestros Servicios será resuelta exclusivamente en los tribunales competentes, a menos que se acuerde lo contrario mediante arbitraje.
              </p>
            </section>

            <section>
              <h2 className="text-2xl font-semibold text-gray-900 mb-4">16. Disposiciones Generales</h2>
              
              <h3 className="text-xl font-semibold text-gray-900 mb-3 mt-4">16.1. Acuerdo Completo</h3>
              <p>Estos Términos constituyen el acuerdo completo entre tú y MiCartaPro respecto al uso de los Servicios.</p>

              <h3 className="text-xl font-semibold text-gray-900 mb-3 mt-4">16.2. Divisibilidad</h3>
              <p>Si alguna disposición de estos Términos se considera inválida o inaplicable, las disposiciones restantes permanecerán en pleno vigor y efecto.</p>

              <h3 className="text-xl font-semibold text-gray-900 mb-3 mt-4">16.3. Renuncia</h3>
              <p>Nuestra falta de ejercer o hacer valer cualquier derecho o disposición de estos Términos no constituirá una renuncia a tal derecho o disposición.</p>

              <h3 className="text-xl font-semibold text-gray-900 mb-3 mt-4">16.4. Cesión</h3>
              <p>No puedes ceder o transferir estos Términos o tus derechos bajo estos Términos sin nuestro consentimiento previo por escrito. Podemos ceder estos Términos sin restricciones.</p>
            </section>

            <section>
              <h2 className="text-2xl font-semibold text-gray-900 mb-4">17. Contacto</h2>
              <p>
                Si tienes preguntas sobre estos Términos y Condiciones, puedes contactarnos a través de:
              </p>
              <ul className="list-disc pl-6 space-y-2">
                <li><strong>WhatsApp:</strong> +56957857558</li>
                <li><strong>Email:</strong> legal@micartapro.com</li>
              </ul>
            </section>
          </div>

          <div className="mt-12 pt-8 border-t border-gray-200 flex flex-col sm:flex-row justify-between items-center gap-4">
            <Link 
              href="/"
              className="inline-flex items-center text-blue-600 hover:text-blue-700 font-medium"
            >
              ← Volver al inicio
            </Link>
            <Link 
              href="/privacy"
              className="inline-flex items-center text-blue-600 hover:text-blue-700 font-medium"
            >
              Ver Política de Privacidad →
            </Link>
          </div>
        </div>
      </div>
    </div>
  )
}

