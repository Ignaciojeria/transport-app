import Link from 'next/link'
import Image from 'next/image'

export const metadata = {
  title: 'Política de Privacidad - MiCartaPro',
  description: 'Política de privacidad de MiCartaPro. Conoce cómo protegemos y gestionamos tus datos personales.',
}

export default function PrivacyPage() {
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
            Política de Privacidad
          </h1>
          
          <p className="text-gray-600 mb-8">
            <strong>Última actualización:</strong> {new Date().toLocaleDateString('es-ES', { year: 'numeric', month: 'long', day: 'numeric' })}
          </p>

          <div className="prose prose-lg max-w-none space-y-8 text-gray-700">
            <section>
              <h2 className="text-2xl font-semibold text-gray-900 mb-4">1. Introducción</h2>
              <p>
                MiCartaPro ("nosotros", "nuestro" o "la empresa") se compromete a proteger y respetar tu privacidad. Esta Política de Privacidad explica cómo recopilamos, usamos, divulgamos y protegemos tu información personal cuando utilizas nuestros servicios de menú digital.
              </p>
              <p>
                Al utilizar nuestros servicios, aceptas las prácticas descritas en esta política. Si no estás de acuerdo con esta política, por favor no utilices nuestros servicios.
              </p>
            </section>

            <section>
              <h2 className="text-2xl font-semibold text-gray-900 mb-4">2. Información que Recopilamos</h2>
              
              <h3 className="text-xl font-semibold text-gray-900 mb-3 mt-4">2.1. Información que Proporcionas Directamente</h3>
              <ul className="list-disc pl-6 space-y-2">
                <li><strong>Información de pedidos:</strong> Nombre del cliente, hora de retiro, y detalles completos del pedido cuando se envía a través de WhatsApp.</li>
                <li><strong>Información de contacto:</strong> Número de teléfono cuando contactas a través de WhatsApp para cotizaciones o consultas.</li>
                <li><strong>Información del restaurante:</strong> Datos proporcionados al solicitar nuestros servicios, incluyendo información de contacto y detalles del negocio.</li>
              </ul>

              <h3 className="text-xl font-semibold text-gray-900 mb-3 mt-4">2.2. Información Recopilada Automáticamente</h3>
              <ul className="list-disc pl-6 space-y-2">
                <li><strong>Información del carrito:</strong> Los items seleccionados se almacenan localmente en tu navegador (localStorage) para mejorar tu experiencia de usuario.</li>
                <li><strong>Información del dispositivo:</strong> Tipo de dispositivo, sistema operativo, y navegador utilizado para acceder a nuestros servicios.</li>
                <li><strong>Información de uso:</strong> Cómo interactúas con nuestros servicios, incluyendo páginas visitadas y tiempo de permanencia.</li>
              </ul>
            </section>

            <section>
              <h2 className="text-2xl font-semibold text-gray-900 mb-4">3. Base Legal para el Procesamiento (GDPR/LGPD)</h2>
              <p>Procesamos tu información personal basándonos en las siguientes bases legales:</p>
              <ul className="list-disc pl-6 space-y-2">
                <li><strong>Consentimiento:</strong> Cuando proporcionas explícitamente tu consentimiento para el procesamiento de tus datos.</li>
                <li><strong>Ejecución de contrato:</strong> Para procesar y gestionar pedidos de restaurantes que utilizan nuestros servicios.</li>
                <li><strong>Interés legítimo:</strong> Para mejorar nuestros servicios, prevenir fraudes y mantener la seguridad de nuestros sistemas.</li>
                <li><strong>Cumplimiento legal:</strong> Cuando es necesario para cumplir con obligaciones legales aplicables.</li>
              </ul>
            </section>

            <section>
              <h2 className="text-2xl font-semibold text-gray-900 mb-4">4. Uso de la Información</h2>
              <p>Utilizamos la información recopilada para los siguientes propósitos:</p>
              <ul className="list-disc pl-6 space-y-2">
                <li>Procesar y gestionar pedidos de restaurantes de manera eficiente</li>
                <li>Mejorar y personalizar la experiencia del usuario</li>
                <li>Responder a consultas, solicitudes de cotización y proporcionar soporte al cliente</li>
                <li>Mantener la funcionalidad del carrito de compras y recordar tus preferencias</li>
                <li>Desarrollar y mejorar nuestros servicios y funcionalidades</li>
                <li>Enviar comunicaciones relacionadas con nuestros servicios (con tu consentimiento)</li>
                <li>Cumplir con obligaciones legales y proteger nuestros derechos legales</li>
              </ul>
            </section>

            <section>
              <h2 className="text-2xl font-semibold text-gray-900 mb-4">5. Almacenamiento y Retención de Datos</h2>
              <p>
                <strong>Almacenamiento Local:</strong> Los datos del carrito se almacenan localmente en tu navegador (localStorage) y no se envían a nuestros servidores hasta que decidas enviar un pedido a través de WhatsApp. Puedes limpiar estos datos en cualquier momento desde la configuración de tu navegador.
              </p>
              <p>
                <strong>Retención de Datos:</strong> Conservamos tu información personal solo durante el tiempo necesario para cumplir con los propósitos descritos en esta política, a menos que la ley requiera o permita un período de retención más largo.
              </p>
              <p>
                <strong>Datos de Pedidos:</strong> Los pedidos se procesan directamente a través de WhatsApp y son gestionados por el restaurante correspondiente. No almacenamos información de pedidos en nuestros servidores.
              </p>
            </section>

            <section>
              <h2 className="text-2xl font-semibold text-gray-900 mb-4">6. Compartir y Divulgación de Información</h2>
              <p>No vendemos, alquilamos ni compartimos tu información personal con terceros para sus propios fines comerciales, excepto en las siguientes circunstancias:</p>
              
              <h3 className="text-xl font-semibold text-gray-900 mb-3 mt-4">6.1. Proveedores de Servicios</h3>
              <p>Podemos compartir información con proveedores de servicios de confianza que nos ayudan a operar nuestro negocio, incluyendo:</p>
              <ul className="list-disc pl-6 space-y-2">
                <li>Proveedores de hosting y servicios en la nube</li>
                <li>Servicios de análisis y monitoreo</li>
                <li>Proveedores de servicios de comunicación (WhatsApp)</li>
              </ul>

              <h3 className="text-xl font-semibold text-gray-900 mb-3 mt-4">6.2. Restaurantes</h3>
              <p>Cuando envías un pedido, la información se comparte con el restaurante correspondiente a través de WhatsApp para procesar tu pedido.</p>

              <h3 className="text-xl font-semibold text-gray-900 mb-3 mt-4">6.3. Requerimientos Legales</h3>
              <p>Podemos divulgar información cuando sea requerido por ley, orden judicial, o para proteger nuestros derechos legales, propiedad o seguridad, o la de nuestros usuarios.</p>

              <h3 className="text-xl font-semibold text-gray-900 mb-3 mt-4">6.4. Transferencias de Negocio</h3>
              <p>En caso de fusión, adquisición o venta de activos, tu información puede ser transferida como parte de esa transacción.</p>
            </section>

            <section>
              <h2 className="text-2xl font-semibold text-gray-900 mb-4">7. Transferencias Internacionales de Datos</h2>
              <p>
                Tus datos pueden ser procesados y almacenados en servidores ubicados fuera de tu país de residencia. Cuando transferimos datos internacionalmente, implementamos salvaguardas apropiadas para proteger tu información de acuerdo con esta política y las leyes aplicables de protección de datos.
              </p>
            </section>

            <section>
              <h2 className="text-2xl font-semibold text-gray-900 mb-4">8. Cookies y Tecnologías Similares</h2>
              <p>
                Utilizamos localStorage del navegador para almacenar temporalmente la información de tu carrito de compras. Esta información se elimina cuando limpias el almacenamiento de tu navegador o cuando utilizas la función de "Limpiar carrito" en nuestra aplicación.
              </p>
              <p>
                No utilizamos cookies de seguimiento de terceros ni tecnologías de rastreo invasivas.
              </p>
            </section>

            <section>
              <h2 className="text-2xl font-semibold text-gray-900 mb-4">9. Tus Derechos (GDPR/LGPD)</h2>
              <p>Dependiendo de tu ubicación, tienes los siguientes derechos respecto a tu información personal:</p>
              
              <ul className="list-disc pl-6 space-y-2">
                <li><strong>Derecho de Acceso:</strong> Puedes solicitar una copia de la información personal que tenemos sobre ti.</li>
                <li><strong>Derecho de Rectificación:</strong> Puedes solicitar que corrijamos cualquier información incorrecta o incompleta.</li>
                <li><strong>Derecho de Supresión:</strong> Puedes solicitar que eliminemos tu información personal en ciertas circunstancias.</li>
                <li><strong>Derecho a la Portabilidad:</strong> Puedes solicitar que transfiramos tu información a otro proveedor de servicios.</li>
                <li><strong>Derecho de Oposición:</strong> Puedes oponerte al procesamiento de tu información personal en ciertas circunstancias.</li>
                <li><strong>Derecho a Limitar el Procesamiento:</strong> Puedes solicitar que limitemos el procesamiento de tu información.</li>
                <li><strong>Derecho a Retirar el Consentimiento:</strong> Puedes retirar tu consentimiento en cualquier momento cuando el procesamiento se base en consentimiento.</li>
                <li><strong>Derecho a Presentar una Queja:</strong> Tienes derecho a presentar una queja ante la autoridad de protección de datos de tu jurisdicción.</li>
              </ul>
              
              <p className="mt-4">
                Para ejercer cualquiera de estos derechos, por favor contáctanos utilizando la información proporcionada en la sección de Contacto.
              </p>
            </section>

            <section>
              <h2 className="text-2xl font-semibold text-gray-900 mb-4">10. Seguridad de los Datos</h2>
              <p>
                Implementamos medidas de seguridad técnicas y organizativas razonables para proteger tu información personal contra acceso no autorizado, alteración, divulgación o destrucción. Estas medidas incluyen:
              </p>
              <ul className="list-disc pl-6 space-y-2">
                <li>Encriptación de datos en tránsito</li>
                <li>Almacenamiento seguro de información</li>
                <li>Acceso restringido a información personal</li>
                <li>Monitoreo regular de nuestros sistemas de seguridad</li>
              </ul>
              <p className="mt-4">
                Sin embargo, ningún método de transmisión por Internet o almacenamiento electrónico es 100% seguro. Aunque nos esforzamos por proteger tu información, no podemos garantizar su seguridad absoluta.
              </p>
            </section>

            <section>
              <h2 className="text-2xl font-semibold text-gray-900 mb-4">11. Privacidad de Menores</h2>
              <p>
                Nuestros servicios no están dirigidos a menores de 18 años. No recopilamos intencionalmente información personal de menores. Si descubrimos que hemos recopilado información de un menor sin el consentimiento de los padres, tomaremos medidas para eliminar esa información de nuestros sistemas.
              </p>
            </section>

            <section>
              <h2 className="text-2xl font-semibold text-gray-900 mb-4">12. Cambios a esta Política de Privacidad</h2>
              <p>
                Podemos actualizar esta Política de Privacidad ocasionalmente para reflejar cambios en nuestras prácticas o por razones legales, operativas o regulatorias. Te notificaremos de cualquier cambio material publicando la nueva política en esta página y actualizando la fecha de "Última actualización".
              </p>
              <p>
                Te recomendamos que revises esta política periódicamente para estar informado sobre cómo protegemos tu información.
              </p>
            </section>

            <section>
              <h2 className="text-2xl font-semibold text-gray-900 mb-4">13. Contacto</h2>
              <p>
                Si tienes preguntas, inquietudes o solicitudes relacionadas con esta Política de Privacidad o el procesamiento de tu información personal, puedes contactarnos a través de:
              </p>
              <ul className="list-disc pl-6 space-y-2">
                <li><strong>WhatsApp:</strong> +56957857558</li>
                <li><strong>Email:</strong> privacidad@micartapro.com</li>
              </ul>
              <p className="mt-4">
                Responderemos a tu solicitud dentro de un plazo razonable y de acuerdo con las leyes aplicables de protección de datos.
              </p>
            </section>

            <section>
              <h2 className="text-2xl font-semibold text-gray-900 mb-4">14. Jurisdicción y Ley Aplicable</h2>
              <p>
                Esta Política de Privacidad se rige por las leyes aplicables de protección de datos en tu jurisdicción, incluyendo pero no limitado al Reglamento General de Protección de Datos (GDPR) de la Unión Europea, la Ley General de Protección de Datos (LGPD) de Brasil, y otras leyes de protección de datos aplicables.
              </p>
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
              href="/terms"
              className="inline-flex items-center text-blue-600 hover:text-blue-700 font-medium"
            >
              Ver Términos y Condiciones →
            </Link>
          </div>
        </div>
      </div>
    </div>
  )
}
